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

// Run MUST be called (as goroutine) for server working
func (s *Server) Run() {
	s.msgSender()
}

// todo if threshold < partyLen, do we need to provide current party or use full party? client doesn't know about current part of party

func (s *Server) Sign(ctx context.Context, msg []byte) ([]byte, error) {
	s.logger.Info().Msg("Start sign operation")

	var signature []byte
	err := s.doOperation(ctx, msg,
		func(ctx context.Context, errCh chan error) {
			s.Tss.Sign(ctx, s.operation.InCh, s.operation.OutCh, errCh, msg, &signature)
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
		func(ctx context.Context, errCh chan error) {
			s.Tss.Keygen(ctx, s.operation.InCh, s.operation.OutCh, errCh)
		},
		func() ([]byte, error) {
			addr, err := s.Tss.GetAddress()
			return addr.Bytes(), err
		},
	)

	return err
}

func (s *Server) doOperation(
	parentCtx context.Context,
	operation []byte,
	tssOperation func(ctx context.Context, errCh chan error),
	resultFunc func() ([]byte, error),
) error {
	if err := s.startOperation(operation); err != nil {
		return err
	}
	defer s.operation.Stop()

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	s.waitForConnections(ctx)
	// if users don't disconnect normally (at the end of function) they will receive this error
	defer s.disconnectAll(fmt.Errorf("some error happened"))

	errCh := make(chan error, 10)

	go tssOperation(ctx, errCh)

	err := <-errCh
	if err != nil {
		s.logger.Error().Err(err).Msg("Tss operation failed")
		return err
	}
	s.logger.Info().Msg("Tss operation successfully finished")

	s.resultsWaiter.Wait()
	if err = checkResults(s.results, resultFunc); err != nil {
		return err
	}

	// normal finish
	s.disconnectAll(nil)
	return nil
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
	for _, v := range results {
		if !bytes.Equal(v, ownResult) {
			return fmt.Errorf("results not equal")
		}
	}
	return nil
}
