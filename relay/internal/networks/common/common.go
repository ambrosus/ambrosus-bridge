package common

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

type CommonBridge struct {
	networks.Bridge
	Client     *ethclient.Client
	WsClient   *ethclient.Client
	Contract   *contracts.Bridge
	WsContract *contracts.Bridge
	Auth       *bind.TransactOpts
	SideBridge networks.Bridge
	Logger     zerolog.Logger
	Name       string
}

func New(cfg config.Network, name string) (b CommonBridge, err error) {
	b.Name = name

	b.Client, err = ethclient.Dial(cfg.HttpURL)
	if err != nil {
		return b, fmt.Errorf("dial http: %w", err)
	}

	// Creating a new bridge contract instance.
	b.Contract, err = contracts.NewBridge(common.HexToAddress(cfg.ContractAddr), b.Client)
	if err != nil {
		return b, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	if cfg.WsURL != "" {
		b.WsClient, err = ethclient.Dial(cfg.WsURL)
		if err != nil {
			return b, fmt.Errorf("dial ws: %w", err)
		}

		b.WsContract, err = contracts.NewBridge(common.HexToAddress(cfg.ContractAddr), b.WsClient)
		if err != nil {
			return b, fmt.Errorf("create contract ws: %w", err)
		}
	}

	// create auth if privateKey provided
	if cfg.PrivateKey != "" {
		pk, err := parsePK(cfg.PrivateKey)
		if err != nil {
			return b, fmt.Errorf("parse private key: %w", err)
		}
		chainId, err := b.Client.ChainID(context.Background())
		if err != nil {
			return b, fmt.Errorf("chain id: %w", err)
		}
		b.Auth, err = bind.NewKeyedTransactorWithChainID(pk, chainId)
		if err != nil {
			return b, fmt.Errorf("new keyed transactor: %w", err)
		}

		// update metrics
		b.SetRelayBalanceMetric()
	}

	return b, nil

}

// GetLastEventId gets last contract event id.
func (b *CommonBridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// GetEventById gets contract event by id.
func (b *CommonBridge) GetEventById(eventId *big.Int) (*contracts.BridgeTransfer, error) {
	logs, err := b.Contract.FilterTransfer(nil, []*big.Int{eventId})
	if err != nil {
		return nil, fmt.Errorf("filter transfer: %w", err)
	}

	if logs.Next() {
		return logs.Event, nil
	}

	return nil, networks.ErrEventNotFound
}

func (b *CommonBridge) GetMinSafetyBlocksNum() (uint64, error) {
	safetyBlocks, err := b.Contract.MinSafetyBlocks(nil)
	if err != nil {
		return 0, err
	}
	return safetyBlocks.Uint64(), nil
}

func (b *CommonBridge) CheckOldEvents() error {
	b.Logger.Info().Msg("Checking old events...")

	lastEventId, err := b.SideBridge.GetLastEventId()
	if err != nil {
		return fmt.Errorf("GetLastEventId: %w", err)
	}

	for i := int64(1); ; i++ {
		nextEventId := new(big.Int).Add(lastEventId, big.NewInt(i))
		nextEvent, err := b.GetEventById(nextEventId)
		if errors.Is(err, networks.ErrEventNotFound) { // no more old events
			return nil
		} else if err != nil {
			return fmt.Errorf("GetEventById on id %v: %w", nextEventId.String(), err)
		}

		b.Logger.Info().Str("event_id", nextEventId.String()).Msg("Send old event...")
		if err := b.SendEvent(nextEvent); err != nil {
			return fmt.Errorf("send event: %w", err)
		}
	}
}

func (b *CommonBridge) Listen() error {
	if err := b.CheckOldEvents(); err != nil {
		return fmt.Errorf("CheckOldEvents: %w", err)
	}
	b.Logger.Info().Msg("Listening new events...")

	// Subscribe to events
	eventCh := make(chan *contracts.BridgeTransfer)
	eventSub, err := b.WsContract.WatchTransfer(nil, eventCh, nil)
	if err != nil {
		return fmt.Errorf("watchTransfer: %w", err)
	}
	defer eventSub.Unsubscribe()

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching transfers: %w", err)
		case event := <-eventCh:
			b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Send event...")

			if err := b.SendEvent(event); err != nil {
				return fmt.Errorf("send event: %w", err)
			}
		}
	}
}

func (b *CommonBridge) ProcessTx(params networks.GetTxErrParams) error {
	if err := b.Bridge.GetTxErr(params); err != nil {
		return err
	}

	receipt, err := bind.WaitMined(context.Background(), b.Client, params.Tx)
	if err != nil {
		return fmt.Errorf("wait mined: %w", err)
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		err = b.GetFailureReason(params.Tx)
		if err != nil {
			return fmt.Errorf("GetFailureReason: %w", helpers.ParseError(err))
		}
	}

	return nil
}
