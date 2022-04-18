package common

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
}

func New(cfg config.Network) (*CommonBridge, error) {
	client, err := ethclient.Dial(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}

	// Creating a new bridge contract instance.
	contract, err := contracts.NewBridge(common.HexToAddress(cfg.ContractAddr), client)
	if err != nil {
		return nil, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	var wsClient *ethclient.Client
	var wsContract *contracts.Bridge
	if cfg.WsURL != "" {
		wsClient, err = ethclient.Dial(cfg.WsURL)
		if err != nil {
			return nil, fmt.Errorf("dial ws: %w", err)
		}

		wsContract, err = contracts.NewBridge(common.HexToAddress(cfg.ContractAddr), wsClient)
		if err != nil {
			return nil, fmt.Errorf("create contract ws: %w", err)
		}
	}

	// create auth if privateKey provided
	var auth *bind.TransactOpts
	if cfg.PrivateKey != "" {
		pk, err := parsePK(cfg.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
		chainId, err := client.ChainID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("chain id: %w", err)
		}
		auth, err = bind.NewKeyedTransactorWithChainID(pk, chainId)
		if err != nil {
			return nil, fmt.Errorf("new keyed transactor: %w", err)
		}
	}

	return &CommonBridge{
		Client:     client,
		WsClient:   wsClient,
		Contract:   contract,
		WsContract: wsContract,
		Auth:       auth,
	}, nil

}

// GetLastEventId gets last contract event id.
func (b *CommonBridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// GetEventById gets contract event by id.
func (b *CommonBridge) GetEventById(eventId *big.Int) (*contracts.BridgeTransfer, error) {
	opts := &bind.FilterOpts{Context: context.Background()}

	logs, err := b.Contract.FilterTransfer(opts, []*big.Int{eventId})
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

	i := big.NewInt(1)
	for {
		nextEventId := big.NewInt(0).Add(lastEventId, i)
		nextEvent, err := b.GetEventById(nextEventId)
		if err != nil {
			if errors.Is(err, networks.ErrEventNotFound) {
				// no more old events
				return nil
			}
			return fmt.Errorf("GetEventById on id %v: %w", nextEventId.String(), err)
		}

		b.Logger.Info().Str("event_id", nextEventId.String()).Msg("Send old event...")

		if err := b.SendEvent(nextEvent); err != nil {
			return fmt.Errorf("send event: %w", err)
		}

		i = big.NewInt(0).Add(i, big.NewInt(1))
	}
}

func (b *CommonBridge) Listen() error {
	if err := b.CheckOldEvents(); err != nil {
		return fmt.Errorf("CheckOldEvents: %w", err)
	}

	b.Logger.Info().Msg("Listening new events...")

	// Subscribe to events
	watchOpts := &bind.WatchOpts{Context: context.Background()}
	eventChannel := make(chan *contracts.BridgeTransfer)
	eventSub, err := b.WsContract.WatchTransfer(watchOpts, eventChannel, nil)
	if err != nil {
		return fmt.Errorf("watchTransfer: %w", err)
	}

	defer eventSub.Unsubscribe()

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching transfers: %w", err)
		case event := <-eventChannel:
			b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Send event...")

			if err := b.SendEvent(event); err != nil {
				return fmt.Errorf("send event: %w", err)
			}
		}
	}
}

func parsePK(pk string) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(pk)
	if err != nil {
		return nil, err
	}
	return crypto.ToECDSA(b)
}
