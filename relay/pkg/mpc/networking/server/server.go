package server

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/rs/zerolog"
)

type Server struct {
	sync.Mutex
	Tss       *tss_wrap.Mpc
	operation common.Operation

	connections  map[string]*common.Conn
	connChangeCh chan byte // populates when client connect or disconnect; used for waitForConnections method

	results       map[string][]byte
	resultsWaiter *sync.WaitGroup

	logger *zerolog.Logger
}

// NewServer create and start new server
func NewServer(tss *tss_wrap.Mpc, logger *zerolog.Logger) *Server {
	s := &Server{
		Tss:          tss,
		connections:  make(map[string]*common.Conn),
		connChangeCh: make(chan byte, 1000),
		operation:    common.NewOperation(),
		logger:       logger,
	}
	return s
}

// todo if threshold < partyLen, do we need to provide current party or use full party? client doesn't know about current part of party

func (s *Server) Sign(ctx context.Context, msg []byte) ([]byte, error) {
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

func (s *Server) Keygen(ctx context.Context) error {
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

func (s *Server) GetFullMsg() ([]byte, error) {
	// just to implement MpcSigner interface
	panic("can be called only on client")
}

func (s *Server) SetFullMsg(fullMsg []byte) {
	s.operation.FullMsg = fullMsg
}

func (s *Server) doOperation(
	ctx context.Context,
	operation []byte,
	tssOperation func(ctx context.Context) error,
	resultFunc func() ([]byte, error),
) error {
	if err := s.startOperation(operation); err != nil {
		return err
	}
	defer s.operation.Stop()

	if err := s.waitForConnections(ctx); err != nil {
		return fmt.Errorf("wait for connections: %w", err)
	}

	err := s.doOperation_(ctx, tssOperation, resultFunc)

	if err != nil {
		s.logger.Error().Err(err).Msg("Operation error")
		s.disconnectAll(fmt.Errorf("server error: %w", err))
		return err
	}

	s.logger.Info().Msg("Operation finished successfully")
	s.disconnectAll(nil)
	return nil
}

func (s *Server) doOperation_(
	ctx context.Context,
	tssOperation func(ctx context.Context) error,
	resultFunc func() ([]byte, error),
) error {

	errCh := make(chan common.OpError)

	// todo create channels here instead of operation struct
	// todo make tssOperation returns result
	go func() { errCh <- common.OpError{"tss", tssOperation(ctx)} }()
	go func() { errCh <- common.OpError{"res", s.receiver(s.operation.OutCh)} }()
	go func() { errCh <- common.OpError{"tra", s.transmitter(s.operation.OutCh)} }()

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
			}
			if err.Type == "res" {
				// receiver returns nil when all clients send results

				// todo wait for result from tss
				if err := checkResults(s.results, resultFunc); err != nil {
					return fmt.Errorf("check results: %w", err)
				}
				s.logger.Info().Msg("Results checked successfully")

				// close outCh so transmitter goroutine will finish (when all queued msgs will be sent)
				close(s.operation.OutCh)
			}
			if err.Type == "tra" {
				// transmitter will return nil when s.operation.OutCh channel closed (when all client send results)
				// at this point we received all results and sends all queued msgs, so finish protocol
				s.logger.Info().Msg("Transmitter finished successfully")
				return nil
			}

		}
	}
}

func (s *Server) startOperation(msg []byte) error {
	s.Lock()
	defer s.Unlock()

	s.results = make(map[string][]byte)
	s.resultsWaiter = new(sync.WaitGroup)
	s.resultsWaiter.Add(s.Tss.Threshold() - 1) // -1 because we don't need to wait for our own result

	return s.operation.Start(msg)
}

func checkResults(results map[string][]byte, resultFunc func() ([]byte, error)) error {
	ownResult, err := resultFunc()
	if err != nil {
		return fmt.Errorf("get result: %w", err)
	}
	for clientID, v := range results {
		if !bytes.Equal(v, ownResult) {
			return fmt.Errorf("client %v send different result (%v != %v)", clientID, v, ownResult)
		}
	}
	return nil
}
