package client

import (
	"context"
	"fmt"
	"net/http"
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
	s.logger.Info().Msg("Start sign operation")

	if err := s.startOperation(msg); err != nil {
		return nil, err
	}
	defer s.stopOperation()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	errCh := make(chan error, 10)
	var signature []byte

	go s.Tss.Sign(ctx, s.operation.InCh, s.operation.OutCh, errCh, msg, &signature)

	go func() { errCh <- s.receiver(conn) }()         // server -> Tss
	go func() { errCh <- s.transmitter(ctx, conn) }() // Tss -> server

	err = <-errCh
	return signature, err
}

func (s *Client) keygen() error {
	s.logger.Info().Msg("Start keygen operation")

	if err := s.startOperation(common.KeygenOperation); err != nil {
		return err
	}
	defer s.stopOperation()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := s.connect(ctx)
	if err != nil {
		return err
	}
	defer common.NormalClose(conn)

	errCh := make(chan error, 10)

	go s.Tss.Keygen(ctx, s.operation.InCh, s.operation.OutCh, errCh)

	go func() { errCh <- s.receiver(conn) }()         // server -> Tss
	go func() { errCh <- s.transmitter(ctx, conn) }() // Tss -> server

	err = <-errCh
	return err
}

func (s *Client) startOperation(msg []byte) error {
	s.Lock()
	defer s.Unlock()
	return s.operation.Start(msg)
}
func (s *Client) stopOperation() {
	s.operation.Stop()
}

func (s *Client) connect(ctx context.Context) (*websocket.Conn, error) {
	headers := make(http.Header)
	headers.Add(common.HeaderTssID, s.Tss.MyID())
	headers.Add(common.HeaderTssOperation, fmt.Sprintf("%x", s.operation.SignMsg)) // as hex

	s.logger.Debug().Msg("Connecting to server")

	conn, httpResp, err := websocket.DefaultDialer.DialContext(ctx, s.serverURL, headers)
	if err != nil {
		if httpResp != nil {
			return nil, fmt.Errorf("ws connect: %w. http resp: %v", err, httpResp.Status)
		} else {
			return nil, fmt.Errorf("ws connect: %w", err)
		}
		// todo here may be error about "This operation doesn't started by server", maybe sleep and retry here?
	}

	s.logger.Debug().Msg("Connected to server")

	return conn, nil
}

func (s *Client) receiver(conn *websocket.Conn) error {
	// breaks when connection closed
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}
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
