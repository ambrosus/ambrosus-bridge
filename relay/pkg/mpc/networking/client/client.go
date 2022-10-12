package client

import (
	"context"
	"sync"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/gorilla/websocket"
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

// just to implement MpcSigner interface
func (s *Client) SetFullMsg(fullMsg []byte) {}

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
	parentCtx context.Context,
	operation []byte,
	tssOperation func(ctx context.Context) error,
	resultFunc func() ([]byte, error),
) error {
	if err := s.startOperation(operation); err != nil {
		return err
	}
	defer s.stopOperation()

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	// if connection don't close normally (at the end of function) disconnect with error
	defer conn.ErrorClose(websocket.CloseProtocolError, "some error happened")

	errCh := make(chan common.OpError)

	go func() { errCh <- common.OpError{"tss", tssOperation(ctx)} }()
	go func() { errCh <- common.OpError{"res", s.receiver(conn)} }()    // server -> Tss
	go func() { errCh <- common.OpError{"tra", s.transmitter(conn)} }() // Tss -> server

	// wait for nil error in tss operation goroutine; any other error is error!
	if err := (<-errCh).Check("tss"); err != nil {
		return err
	}
	// tss operation finished successfully

	// close outCh so transmitter goroutine will finish (when all queued msgs will be sent)
	close(s.operation.OutCh)
	if err := (<-errCh).Check("tra"); err != nil {
		return err
	}

	// now can send result to server
	// sending result (signature or address) to server
	if err = sendResult(conn, resultFunc); err != nil {
		return err
	}
	s.logger.Debug().Msg("Result sent")

	// server close connection normally when all clients send they result;
	// receiver returns nil when connection closed normally
	if err := (<-errCh).Check("res"); err != nil {
		return err
	}

	conn.NormalClose()
	return nil
}

func (s *Client) startOperation(msg []byte) error {
	s.Lock()
	defer s.Unlock()
	return s.operation.Start(msg)
}
func (s *Client) stopOperation() {
	s.operation.Stop()
}
