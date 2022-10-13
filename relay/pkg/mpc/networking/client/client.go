package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/rs/zerolog"
)

type Client struct {
	sync.Mutex
	logger *zerolog.Logger

	Tss       *tss_wrap.Mpc
	operation common.Operation

	serverURL string
}

func NewClient(tss *tss_wrap.Mpc, serverURL string, logger *zerolog.Logger) *Client {
	s := &Client{
		Tss:       tss,
		serverURL: serverURL,
		logger:    logger,
		operation: common.NewOperation(),
	}
	return s
}

func (s *Client) Sign(ctx context.Context, msg []byte) ([]byte, error) {
	for {
		sig, err := s.sign(ctx, msg)
		if err == nil {
			return sig, nil
		}
		s.logger.Error().Err(err).Msg("sign error")
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		time.Sleep(time.Second)
	}
}

func (s *Client) Keygen(ctx context.Context) error {
	for {
		err := s.keygen(ctx)
		if err == nil {
			return nil
		}
		s.logger.Error().Err(err).Msg("keygen error")
		if ctx.Err() != nil {
			return ctx.Err()
		}
		time.Sleep(time.Second)
	}
}

func (s *Client) SetFullMsg(fullMsg []byte) {
	// just to implement MpcSigner interface
}

func (s *Client) GetFullMsg() ([]byte, error) {
	resp, err := http.Get(s.serverURL + common.EndpointFullMsg)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("read body to buffer: %w", err)
	}

	return buf.Bytes(), nil
}

func (s *Client) sign(ctx context.Context, msg []byte) ([]byte, error) {
	s.logger.Info().Msg("Start sign operation")

	var signature []byte
	err := s.doOperation(ctx, msg,
		func(ctx context.Context) (err error) {
			signature, err = s.Tss.SignSync(ctx, s.operation.InCh, s.operation.OutCh, msg)
			return err
		},
		func() ([]byte, error) {
			return signature, nil
		},
	)

	return signature, err
}

func (s *Client) keygen(ctx context.Context) error {
	s.logger.Info().Msg("Start keygen operation")

	err := s.doOperation(ctx, common.KeygenOperation,
		func(ctx context.Context) error {
			return s.Tss.KeygenSync(ctx, s.operation.InCh, s.operation.OutCh)
		},
		func() ([]byte, error) {
			addr, err := s.Tss.GetAddress()
			return addr.Bytes(), err
		},
	)

	return err
}

func (s *Client) doOperation(
	ctx context.Context,
	operation []byte,
	tssOperation func(ctx context.Context) error,
	resultFunc func() ([]byte, error),
) error {
	if err := s.startOperation(operation); err != nil {
		return err
	}
	defer s.stopOperation()

	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}

	err = s.doOperation_(ctx, conn, tssOperation, resultFunc)

	if err != nil {
		s.logger.Error().Err(err).Msg("Operation error")
		conn.Close(fmt.Errorf("client error: %w", err))
		return err
	}

	s.logger.Info().Msg("Operation finished successfully")
	conn.NormalClose()
	return nil
}

func (s *Client) doOperation_(
	ctx context.Context,
	conn *common.Conn,
	tssOperation func(ctx context.Context) error,
	resultFunc func() ([]byte, error),
) error {
	errCh := make(chan common.OpError)

	// todo same that in server

	go func() { errCh <- common.OpError{"tss", tssOperation(ctx)} }()
	go func() { errCh <- common.OpError{"res", s.receiver(conn)} }()    // server -> Tss
	go func() { errCh <- common.OpError{"tra", s.transmitter(conn)} }() // Tss -> server

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-errCh:
			if err.Err != nil {
				return fmt.Errorf("%s error: %w", err.Type, err.Err)
			}

			// if err is nil, it means that some goroutine successfully finished

			if err.Type == "tss" {
				s.logger.Info().Msg("Tss operation finished successfully")

				// tss operation finished successfully, close connection close outCh so transmitter goroutine will finish
				close(s.operation.OutCh)
			}

			if err.Type == "tra" {
				// transmitter will return nil when s.operation.OutCh channel closed (at the end of tssOperation)

				// all messages sent, now can send result (signature or address) to server
				if err := sendResult(conn, resultFunc); err != nil {
					return fmt.Errorf("send result: %w", err)
				}
				s.logger.Info().Msg("Result sent")

			}

			if err.Type == "res" {
				// receiver returns nil when connection closed normally (after server receive all results)

				// server closed connection, all done
				s.logger.Info().Msg("Receiver finished successfully")

				return nil
			}

		}
	}

}

func (s *Client) startOperation(msg []byte) error {
	s.Lock()
	defer s.Unlock()
	return s.operation.Start(msg)
}
func (s *Client) stopOperation() {
	s.operation.Stop()
}
