package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
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
	operation []byte

	serverURL  string
	httpClient *http.Client
}

func NewClient(tss *tss_wrap.Mpc, serverURL string, httpClient *http.Client, logger *zerolog.Logger) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	s := &Client{
		Tss:        tss,
		serverURL:  serverURL,
		httpClient: httpClient,
		logger:     logger,
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
	resp, err := s.httpClient.Get("http://" + s.serverURL + common.EndpointFullMsg)
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

	signature, err := s.doOperation(ctx, msg,
		func(ctx context.Context, inCh <-chan []byte, outCh chan<- *tss_wrap.Message) ([]byte, error) {
			return s.Tss.SignSync(ctx, inCh, outCh, msg)
		},
	)

	return signature, err
}

func (s *Client) keygen(ctx context.Context) error {
	s.logger.Info().Msg("Start keygen operation")

	_, err := s.doOperation(ctx, common.KeygenOperation,
		func(ctx context.Context, inCh <-chan []byte, outCh chan<- *tss_wrap.Message) ([]byte, error) {
			err := s.Tss.KeygenSync(ctx, inCh, outCh)
			if err != nil {
				return nil, err
			}
			addr, err := s.Tss.GetAddress()
			return addr.Bytes(), err
		},
	)

	return err
}

func (s *Client) doOperation(ctx context.Context, operation []byte, tssOperation common.OperationFunc) ([]byte, error) {
	if err := s.startOperation(operation); err != nil {
		return nil, err
	}
	defer s.stopOperation()

	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	result, err := s.doOperation_(ctx, conn, tssOperation)

	if err != nil {
		s.logger.Error().Err(err).Msg("Operation error")
		conn.Close(fmt.Errorf("client error: %w", err))
		return nil, err
	}

	s.logger.Info().Msg("Operation finished successfully")
	conn.NormalClose()
	return result, nil
}

func (s *Client) doOperation_(
	ctx context.Context,
	conn *common.Conn,
	tssOperation common.OperationFunc,
) (ownResult []byte, err error) {

	inCh := make(chan []byte, 10)
	outCh := make(chan *tss_wrap.Message, 10)
	errCh := make(chan common.OpError, 3)

	go func() {
		ownResult, err = tssOperation(ctx, inCh, outCh)
		errCh <- common.OpError{"tss", err}
	}()
	go func() { errCh <- common.OpError{"res", s.receiver(conn, inCh)} }()     // server -> Tss
	go func() { errCh <- common.OpError{"tra", s.transmitter(conn, outCh)} }() // Tss -> server

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case err := <-errCh:
			if err.Err != nil {
				return nil, fmt.Errorf("%s error: %w", err.Type, err.Err)
			}

			// if err is nil, it means that some goroutine successfully finished

			if err.Type == "tss" {
				s.logger.Info().Msg("Tss operation finished successfully")

				// tss operation finished successfully, close connection close outCh so transmitter goroutine will finish
				close(outCh)
			}

			if err.Type == "tra" {
				// transmitter will return nil when s.operation.OutCh channel closed (at the end of tssOperation)

				// all messages sent, now can send result (signature or address) to server
				resultMsg := append(common.ResultPrefix, ownResult...)
				if err := conn.Write(resultMsg); err != nil {
					return nil, fmt.Errorf("write result message: %w", err)
				}

				s.logger.Info().Msg("Result sent")

			}

			if err.Type == "res" {
				// receiver returns nil when connection closed normally (after server receive all results)

				// server closed connection, all done
				s.logger.Info().Msg("Receiver finished successfully")

				return ownResult, nil
			}

		}
	}

}

func (s *Client) startOperation(msg []byte) error {
	s.Lock()
	defer s.Unlock()
	if s.operation != nil {
		return fmt.Errorf("operation already started")
	}

	s.operation = msg
	return nil
}

func (s *Client) stopOperation() {
	s.Lock()
	defer s.Unlock()
	s.operation = nil
}
