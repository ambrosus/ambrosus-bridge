package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type Server struct {
	sync.Mutex
	Tss       *tss_wrap.Mpc
	operation common.Operation

	connections  map[string]*websocket.Conn
	connChangeCh chan byte // populates when client connect or disconnect; used for waitForConnections method

	logger *zerolog.Logger
}

// NewServer create and start new server
func NewServer(tss *tss_wrap.Mpc, logger *zerolog.Logger) *Server {
	s := &Server{
		Tss:          tss,
		connections:  make(map[string]*websocket.Conn),
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
// todo defer (stop operation; disconnect clients) after keygen/sign finished

func (s *Server) Sign(msg []byte) ([]byte, error) {
	s.logger.Info().Msg("Start sign operation")

	if err := s.startOperation(msg); err != nil {
		return nil, err
	}
	defer s.operation.Stop()

	s.waitForConnections()
	defer s.disconnectAll()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 10)
	var signature []byte

	go s.Tss.Sign(ctx, s.operation.InCh, s.operation.OutCh, errCh, msg, &signature)

	err := <-errCh
	if err != nil {
		s.logger.Error().Err(err).Msg("Sign failed")
	} else {
		s.logger.Info().Msg("Sign successfully finished")
	}
	return signature, err
}

func (s *Server) Keygen() error {
	s.logger.Info().Msg("Start keygen operation")

	if err := s.startOperation(common.KeygenOperation); err != nil {
		return err
	}
	defer s.operation.Stop()

	s.waitForConnections()
	defer s.disconnectAll()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 10)

	go s.Tss.Keygen(ctx, s.operation.InCh, s.operation.OutCh, errCh)

	err := <-errCh
	if err != nil {
		s.logger.Error().Err(err).Msg("Keygen failed")
	} else {
		s.logger.Info().Msg("Keygen successfully finished")
	}
	return err
}

func (s *Server) startOperation(msg []byte) error {
	s.Lock()
	defer s.Unlock()
	return s.operation.Start(msg)
}

func (s *Server) msgSender() {
	s.logger.Debug().Msg("Start msgSender")
	for {
		select {
		case msg := <-s.operation.OutCh:
			s.logger.Debug().Msg("Got message from Tss")
			err := s.sendMsg(msg)
			if err != nil {
				s.logger.Error().Err(err).Msg("Failed to send message")
				// todo repeat on err
			}
		}
	}
}

// sendMsg send message to own Tss or to another client(s)
func (s *Server) sendMsg(msg *tss_wrap.OutputMessage) error {
	if msg == nil || msg.SendToIds == nil {
		return fmt.Errorf("nil message")
	}

	for _, id := range msg.SendToIds {
		s.logger.Debug().Str("To", id).Msg("Send message to client")

		// send to own tss
		if id == s.Tss.MyID() {
			s.operation.InCh <- msg.Message
			s.logger.Debug().Str("To", id).Msg("Send message to myself successfully")
			continue
		}

		// send to another client
		conn, ok := s.connections[id]
		if !ok {
			s.logger.Warn().Msgf("Connection with id %s not found, ", id)
			return fmt.Errorf("connection %v not found", id)
			// todo maybe call waitForConnections on this error
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, msg.Message); err != nil {
			return fmt.Errorf("writeMessage: %w", err)
		}
		s.logger.Debug().Str("To", id).Msg("Send message to client sucessfully")
	}

	return nil
}
