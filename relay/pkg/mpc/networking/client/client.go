package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
	ec "github.com/ethereum/go-ethereum/common"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/networking/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
	"github.com/rs/zerolog"
)

type Client struct {
	sync.Mutex
	logger *zerolog.Logger

	Tss       *tss_wrap.Mpc
	operation []byte

	serverURL   string
	accessToken string
}

func NewClient(tss *tss_wrap.Mpc, serverURL string, accessToken string, logger *zerolog.Logger) *Client {
	return &Client{
		Tss:         tss,
		serverURL:   serverURL,
		accessToken: accessToken,
		logger:      logger,
	}
}

func (s *Client) Sign(ctx context.Context, party []string, msg []byte) ([]byte, error) {
	s.logger.Info().Msg("Start sign operation")

	signature, err := s.doOperation(ctx, msg,
		func(ctx context.Context, inCh <-chan []byte, outCh chan<- *tss_wrap.Message) ([]byte, error) {
			return s.Tss.Sign(ctx, party, inCh, outCh, msg)
		},
	)

	return signature, err
}

func (s *Client) Keygen(ctx context.Context, party []string, optionalPreParams ...keygen.LocalPreParams) error {
	s.logger.Info().Msg("Start keygen operation")

	_, err := s.doOperation(ctx, common.KeygenOperation,
		func(ctx context.Context, inCh <-chan []byte, outCh chan<- *tss_wrap.Message) ([]byte, error) {
			err := s.Tss.Keygen(ctx, party, inCh, outCh, optionalPreParams...)
			if err != nil {
				return nil, err
			}
			addr, err := s.Tss.GetAddress()
			return addr.Bytes(), err
		},
	)

	return err
}

func (s *Client) Reshare(ctx context.Context, partyIDsOld, partyIDsNew []string, thresholdNew int) error {
	s.logger.Info().Msg("Start reshare operation")

	_, err := s.doOperation(ctx, common.ReshareOperation,
		func(ctx context.Context, inCh <-chan []byte, outCh chan<- *tss_wrap.Message) ([]byte, error) {
			err := s.Tss.Reshare(ctx, partyIDsOld, partyIDsNew, thresholdNew, inCh, outCh)
			if err != nil {
				return nil, err
			}
			addr, err := s.Tss.GetAddress()
			return addr.Bytes(), err
		},
	)

	return err
}

func (s *Client) SetFullMsg(fullMsg []byte) {
	// just to implement MpcSigner interface
}

func (s *Client) GetFullMsg() ([]byte, error) {
	// todo ctx
	serverURL := strings.Replace(s.serverURL, "ws", "http", 1)
	resp, err := http.Get(serverURL)
	if resp == nil {
		return nil, fmt.Errorf("resp is nil, err: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body to buffer: %w", err)
	}

	return b, nil
}

func (s *Client) GetTssAddress() (ec.Address, error) {
	return s.Tss.GetAddress()
}

func (s *Client) doOperation(ctx context.Context, operation []byte, tssOperation common.OperationFunc) ([]byte, error) {
	if err := s.startOperation(operation); err != nil {
		return nil, err
	}
	defer s.stopOperation()

	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}

	result, err := s.doOperation_(ctx, conn, tssOperation)

	if err != nil {
		s.logger.Error().Err(err).Msg("Operation error")
		conn.Close(fmt.Errorf("client error: %w", err))
		return nil, err
	}

	s.logger.Info().Msg("Operation finished successfully")
	conn.NormalClose()
	return result, nil
}

func (s *Client) doOperation_(
	ctx context.Context,
	conn *common.Conn,
	tssOperation common.OperationFunc,
) (ownResult []byte, err error) {

	inCh := make(chan []byte, 10)
	outCh := make(chan *tss_wrap.Message, 10)
	errCh := make(chan common.OpError, 3)

	go func() {
		ownResult, err = tssOperation(ctx, inCh, outCh)
		errCh <- common.OpError{"tss", err}
	}()
	go func() { errCh <- common.OpError{"res", s.receiver(conn, inCh)} }()     // server -> Tss
	go func() { errCh <- common.OpError{"tra", s.transmitter(conn, outCh)} }() // Tss -> server

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case err := <-errCh:
			if err.Err != nil {
				return nil, fmt.Errorf("%s error: %w", err.Type, err.Err)
			}

			// if err is nil, it means that some goroutine successfully finished

			if err.Type == "tss" {
				s.logger.Info().Msg("Tss operation finished successfully")

				// tss operation finished successfully, close connection close outCh so transmitter goroutine will finish
				close(outCh)
			}

			if err.Type == "tra" {
				// transmitter will return nil when s.operation.OutCh channel closed (at the end of tssOperation)

				// all messages sent, now can send result (signature or address) to server
				resultMsg := append(common.ResultPrefix, ownResult...)
				if err := conn.Write(resultMsg); err != nil {
					return nil, fmt.Errorf("write result message: %w", err)
				}

				s.logger.Info().Msg("Result sent")

			}

			if err.Type == "res" {
				// receiver returns nil when connection closed normally (after server receive all results)

				// server closed connection, all done
				s.logger.Info().Msg("Receiver finished successfully")

				return ownResult, nil
			}

		}
	}

}

func (s *Client) startOperation(msg []byte) error {
	s.Lock()
	defer s.Unlock()
	if s.operation != nil {
		return fmt.Errorf("operation already started")
	}

	s.operation = msg
	return nil
}

func (s *Client) stopOperation() {
	s.Lock()
	defer s.Unlock()
	s.operation = nil
}
