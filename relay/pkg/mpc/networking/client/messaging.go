package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/gorilla/websocket"
)

func (s *Client) connect(ctx context.Context) (*common.Conn, error) {
	headers := make(http.Header)
	headers.Add(common.HeaderTssID, s.Tss.MyID())
	headers.Add(common.HeaderTssOperation, fmt.Sprintf("%x", s.operation)) // as hex
	headers.Add(common.HeaderAccessToken, s.accessToken)

	s.logger.Debug().Str("url", s.serverURL).Msg("Connecting to server")

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

func (s *Client) receiver(conn *common.Conn, inCh chan<- []byte) error {
	// breaks when connection closed
	for {
		msgBytes, err := conn.Read()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil
			}
			return fmt.Errorf("read message: %w", err)
		}

		inCh <- msgBytes
	}
}

func (s *Client) transmitter(conn *common.Conn, outCh <-chan *tss_wrap.Message) error {
	for msg := range outCh {
		msgBytes, err := msg.Marshall()
		if err != nil {
			return fmt.Errorf("marshal message: %w", err)
		}
		if err := conn.Write(msgBytes); err != nil {
			return fmt.Errorf("write message: %w", err)
		}
		s.logger.Debug().Msg("Send message to server successfully")
	}
	return nil
}
