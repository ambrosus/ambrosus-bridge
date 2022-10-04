package client

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

var keygenOperation = []byte("keygen")

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, s.serverURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ws connect: %w", err)
	}
	defer conn.Close()

	errCh := make(chan error, 10)
	var signature []byte

	go s.Tss.Sign(ctx, s.operation.InCh, s.operation.OutCh, errCh, msg, &signature)
	go s.sendMsgs(ctx, conn, errCh)

	err = <-errCh
	return signature, err
}

func (s *Client) keygen() error {
	if err := s.startOperation(keygenOperation); err != nil {
		return err
	}
	defer s.stopOperation()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, s.serverURL, nil)
	if err != nil {
		return fmt.Errorf("ws connect: %w", err)
	}
	defer conn.Close()

	errCh := make(chan error, 10)

	go s.Tss.Keygen(ctx, s.operation.InCh, s.operation.OutCh, errCh)
	go s.sendMsgs(ctx, conn, errCh)

	err = <-errCh
	return err
}

func (s *Client) startOperation(msg []byte) error {
	s.Lock()
	defer s.Unlock()

	if s.operation.Started {
		return fmt.Errorf("already doing something")
	}
	s.operation = common.NewOperation(msg)

	return nil
}
func (s *Client) stopOperation() {
	s.operation.Stop()
}

func (s *Client) sendMsgs(ctx context.Context, conn *websocket.Conn, errCh chan error) {
	if err := s.sendStartMsgs(conn); err != nil {
		errCh <- err
		return
	}

	// protocol begins

	// server -> Tss
	go func() {
		errCh <- s.receiver(conn)
	}()

	// Tss -> server
	go func() {
		errCh <- s.transmitter(ctx, conn)
	}()
}

func (s *Client) sendStartMsgs(conn *websocket.Conn) error {
	// send ID
	if err := conn.WriteMessage(websocket.BinaryMessage, []byte(s.Tss.MyID())); err != nil {
		return fmt.Errorf("write 1 message: %w", err)
	}
	// send operation (keygen or sign msg)
	if err := conn.WriteMessage(websocket.BinaryMessage, s.operation.SignMsg); err != nil {
		return fmt.Errorf("write 2 message: %w", err)
	}
	return nil
}

func (s *Client) receiver(conn *websocket.Conn) error {
	// breaks when connection closed
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read message: %w", err)
		}

		s.operation.InCh <- msgBytes
	}
}

func (s *Client) transmitter(ctx context.Context, conn *websocket.Conn) error {
	for {
		select {
		case msg := <-s.operation.OutCh:
			msgBytes, err := msg.Marshall()
			if err != nil {
				return fmt.Errorf("marshal message: %w", err)
			}
			if err := conn.WriteMessage(websocket.BinaryMessage, msgBytes); err != nil {
				return fmt.Errorf("write message: %w", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
