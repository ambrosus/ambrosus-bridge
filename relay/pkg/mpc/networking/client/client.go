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

func (s *Client) sign(ctx context.Context, msg []byte) ([]byte, error) {
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

func (s *Client) keygen(ctx context.Context) error {
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

func (s *Client) doOperation(
	parentCtx context.Context,
	operation []byte,
	tssOperation func(ctx context.Context, errCh chan error),
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
	defer conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseProtocolError, "some error happened"))

	tssErrCh := make(chan error, 10)
	wsErrCh := make(chan error, 10)

	go tssOperation(ctx, tssErrCh)

	go func() { wsErrCh <- s.receiver(conn) }()         // server -> Tss
	go func() { wsErrCh <- s.transmitter(ctx, conn) }() // Tss -> server

	// wait for nil in tssErrCh
	// any value in wsErrCh at this point means error (even nil, coz it means connection closed)
	select {
	case err = <-tssErrCh:
		if err != nil {
			return err
		}
	case err = <-wsErrCh:
		// ws shouldn't return anything at this point
		return fmt.Errorf("ws error: %w", err)
	}
	// tss operation finished!

	s.logger.Debug().Msg("Tss operation finished")

	// sending result (signature or address) to server
	if err = sendResult(conn, resultFunc); err != nil {
		return err
	}
	s.logger.Debug().Msg("Result sent")

	// server close connection normally when all clients send they result;
	// wsErrCh got nil (in receiver goroutine) when connection closed normally
	err = <-wsErrCh
	if err != nil {
		return fmt.Errorf("ws error: %w", err)
	}

	common.NormalClose(conn)
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

func sendResult(conn *websocket.Conn, resultFunc func() ([]byte, error)) error {
	result, err := resultFunc()
	if err != nil {
		return fmt.Errorf("get result: %w", err)
	}
	resultMsg := append(common.ResultPrefix, result...)
	// todo fix concurrent write to conn (this goroutine and transmitter)
	if err := conn.WriteMessage(websocket.BinaryMessage, resultMsg); err != nil {
		return fmt.Errorf("write result message: %w", err)
	}
	return nil
}
