package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/gorilla/websocket"
)

func (s *Client) connect(ctx context.Context) (*common.Conn, error) {
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

	return &common.Conn{Conn: conn}, nil
}

func (s *Client) receiver(conn *common.Conn) error {
	// breaks when connection closed
	for {
		msgBytes, err := conn.Read()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}
			return fmt.Errorf("read message: %w", err)
		}

		s.operation.InCh <- msgBytes
	}
}

func (s *Client) transmitter(conn *common.Conn) error {
	for msg := range s.operation.OutCh {
		msgBytes, err := msg.Marshall()
		if err != nil {
			return fmt.Errorf("marshal message: %w", err)
		}
		if err := conn.Write(msgBytes); err != nil {
			return fmt.Errorf("write message: %w", err)
		}
	}
	return nil
}

func sendResult(conn *common.Conn, resultFunc func() ([]byte, error)) error {
	result, err := resultFunc()
	if err != nil {
		return fmt.Errorf("get result: %w", err)
	}
	resultMsg := append(common.ResultPrefix, result...)
	// todo fix concurrent write to conn (this goroutine and transmitter)
	if err := conn.Write(resultMsg); err != nil {
		return fmt.Errorf("write result message: %w", err)
	}
	return nil
}
