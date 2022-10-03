package client

import (
	"context"
	"fmt"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var keygenOperation = []byte("keygen")

type Client struct {
	sync.Mutex // locked while signing
	Logger     *zerolog.Logger

	tss       *tss_wrap.Mpc
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

func (s *Client) Sign(msg []byte) ([]byte, error) {
	if err := s.startOperation(msg); err != nil {
		return nil, err
	}
	// sync operation, wait for sign
	// todo if threshold < partyLen, do we need to provide current party or use full party?
	// todo client doesn't know about current part of party
	sign, err := s.tss.Sign(s.operation.inCh, s.operation.outCh, msg)
	return sign, err

	// todo stop Sign function if websocket error happened by context

	// todo disconnect clients
}

func (s *Client) Keygen() error {
	if err := s.startOperation(keygenOperation); err != nil {
		return err
	}

	// sync operation, wait
	err := s.tss.Keygen(s.operation.inCh, s.operation.outCh)
	return err
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

func (s *Client) sendMsg(ctx context.Context, operation []byte) error {
	conn, _, err := websocket.DefaultDialer.Dial(s.serverURL, nil)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, []byte(s.tss.MyID())); err != nil {
		return fmt.Errorf("write 1 message: %w", err)
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, operation); err != nil {
		return fmt.Errorf("write 2 message: %w", err)
	}

	// protocol begins

	// todo goroutines for receiver and transmitter

	// server -> tss
	go func() {
		// todo how to stop this (if error somewhere else happened)?
		for {
			// todo is cycle break on disconnect?
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				return fmt.Errorf("read protocol message: %w", err)
			}

			s.operation.inCh <- msgBytes
		}
	}()

	// tss -> server
	go func() {
		for {
			select {
			case msg := <-s.operation.outCh:
				if err := conn.WriteMessage(websocket.BinaryMessage, msg.Message); err != nil {
					// todo retry
					return fmt.Errorf("write protocol message: %w", err)
				}
			case <-ctx.Done():
				return fmt.Errorf("context done: %w", ctx.Err())
			}
		}
	}()

}
