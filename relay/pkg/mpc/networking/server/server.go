package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type Server struct {
	sync.Mutex
	Tss       *tss_wrap.Mpc
	operation operation

	connections  map[string]*websocket.Conn
	connChangeCh chan byte // populates when client connect or disconnect; used for waitForConnections method

	logger *zerolog.Logger
}

type operation struct {
	started bool
	// will be sent to tss
	inCh chan []byte
	// received from tss or another client; will be sent to own tss or another client
	outCh   chan *tss_wrap.OutputMessage
	signMsg []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  65536,
	WriteBufferSize: 65536,
}

const chSize = 100

// NewServer create and start new server
func NewServer(tss *tss_wrap.Mpc, logger *zerolog.Logger) *Server {
	s := &Server{
		Tss:          tss,
		connections:  make(map[string]*websocket.Conn),
		connChangeCh: make(chan byte, 1000),
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
	if err := s.startOperation(msg); err != nil {
		return nil, err
	}
	s.logger.Info().Msg("Start sign operation")
	s.logger.Debug().Msg("Wait for connections")
	s.waitForConnections()
	s.logger.Debug().Msg("All connections established")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 10)
	var signature []byte

	go s.Tss.Sign(ctx, s.operation.inCh, s.operation.outCh, errCh, msg, &signature)

	err := <-errCh

	return signature, err
}

func (s *Server) Keygen() error {
	if err := s.startOperation(keygenOperation); err != nil {
		return err
	}
	s.logger.Info().Msg("Start keygen operation")
	s.logger.Debug().Msg("Wait for connections")
	s.waitForConnections()
	s.logger.Debug().Msg("All connections established")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 10)

	go s.Tss.Keygen(ctx, s.operation.inCh, s.operation.outCh, errCh)

	err := <-errCh
	return err
}

func (s *Server) startOperation(msg []byte) error {
	// todo don't sure if we need locks
	s.Lock()
	defer s.Unlock()

	if s.operation.started {
		return fmt.Errorf("already doing something")
	}
	// todo don't need to create channels every time;
	s.operation = operation{
		started: true,
		inCh:    make(chan []byte, chSize),
		outCh:   make(chan *tss_wrap.OutputMessage, chSize),
		signMsg: msg,
	}

	return nil
}

func (s *Server) msgSender() {
	for {
		select {
		case msg := <-s.operation.outCh:
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
	for _, id := range msg.SendToIds {

		// send to own tss
		if id == s.Tss.MyID() {
			s.operation.inCh <- msg.Message
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
	}

	return nil
}
