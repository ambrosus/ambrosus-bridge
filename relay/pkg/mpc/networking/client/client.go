package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var keygenOperation = []byte("keygen")

type Client struct {
	sync.Mutex // locked while signing
	logger     *zerolog.Logger

	Tss       *tss_wrap.Mpc
	operation operation

	serverURL string
}

func NewClient(tss *tss_wrap.Mpc, serverURL string, logger *zerolog.Logger) *Client {
	s := &Client{
		tss:       tss,
		serverURL: serverURL,
		Logger:    logger,
	}
	return s
}

type operation struct {
	started bool
	// will be sent to tss
	inCh chan []byte
	// received from tss or another client; will be sent to own tss or another client
	outCh   chan *tss_wrap.OutputMessage
	signMsg []byte
}

const chSize = 10

func NewClient(tss *tss_wrap.Mpc, serverURL string, logger *zerolog.Logger) *Client {
	s := &Client{
		Tss:       tss,
		serverURL: serverURL,
		logger:    logger,
	}
	return s
}

func (s *Client) Sign(msg []byte) ([]byte, error) {
	for {
		sig, err := s.sign(msg)
		if err == nil {
			return sig, nil
		}
		s.logger.Error().Err(err).Msg("sign error")
		time.Sleep(time.Second)
	}
}

func (s *Client) Keygen() error {
	for {
		err := s.keygen()
		if err == nil {
			return nil
		}
		s.logger.Error().Err(err).Msg("keygen error")
		time.Sleep(time.Second)
	}
}

func (s *Client) sign(msg []byte) ([]byte, error) {
	if err := s.startOperation(msg); err != nil {
		return nil, err
	}
	defer s.stopOperation()

	conn, _, err := websocket.DefaultDialer.Dial(s.serverURL, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 10)
	var signature []byte

	go s.sendMsgs(ctx, conn, errCh)
	go s.Tss.Sign(ctx, s.operation.inCh, s.operation.outCh, errCh, msg, &signature)

	err = <-errCh
	return signature, err

	// todo close channels (hope this stop gourotines)
	// todo disconnect
}

func (s *Client) keygen() error {
	if err := s.startOperation(keygenOperation); err != nil {
		return err
	}
	defer s.stopOperation()

	conn, _, err := websocket.DefaultDialer.Dial(s.serverURL, nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 10)

	go s.sendMsgs(ctx, conn, errCh)
	go s.Tss.Keygen(ctx, s.operation.inCh, s.operation.outCh, errCh)

	err = <-errCh
	return err

	// todo close channels (hope this stop gourotines)
	// todo disconnect clients
}

func (s *Client) startOperation(msg []byte) error {
	// todo don't sure if we need this
	s.Lock()
	defer s.Unlock()

	if s.operation.started {
		return fmt.Errorf("already doing something")
	}
	s.operation = operation{
		started: true,
		inCh:    make(chan []byte, chSize),
		outCh:   make(chan *tss_wrap.OutputMessage, chSize),
		signMsg: msg,
	}

	return nil
}
func (s *Client) stopOperation() {
	s.operation.started = false
	close(s.operation.inCh)
	close(s.operation.outCh)
}

func (s *Client) sendMsgs(ctx context.Context, conn *websocket.Conn, errCh chan error) {

	if err := conn.WriteMessage(websocket.BinaryMessage, []byte(s.Tss.MyID())); err != nil {
		errCh <- fmt.Errorf("write 1 message: %w", err)
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, s.operation.signMsg); err != nil {
		errCh <- fmt.Errorf("write 2 message: %w", err)
		return
	}

	// protocol begins

	// server -> Tss
	go func() {
		// breaks when connection closed
		for {
			// todo is cycle break on disconnect?
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				errCh <- fmt.Errorf("read message: %w", err)
				return
			}

			s.operation.inCh <- msgBytes
		}
	}()

	// Tss -> server
	go func() {
		for {
			select {
			case msg := <-s.operation.outCh:
				msgBytes, err := msg.Marshall()
				if err != nil {
					errCh <- fmt.Errorf("marshal message: %w", err)
					return
				}
				if err := conn.WriteMessage(websocket.BinaryMessage, msgBytes); err != nil {
					// todo retry
					errCh <- fmt.Errorf("write message: %w", err)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}
