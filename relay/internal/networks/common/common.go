package common

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

type CommonBridge struct {
	networks.Bridge
	Client     ethclients.ClientInterface
	WsClient   ethclients.ClientInterface
	Contract   *contracts.Bridge
	WsContract *contracts.Bridge
	Auth       *bind.TransactOpts
	SideBridge networks.Bridge
	Logger     zerolog.Logger
	Name       string
	Pk         *ecdsa.PrivateKey

	ContractCallLock *sync.Mutex
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
		pk, err := helpers.ParsePK(cfg.PrivateKey)
		if err != nil {
			return b, err
		}
		b.Pk = pk
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
	} else {
		b.Logger.Info().Msg("No private key provided")
	}

	b.ContractCallLock = &sync.Mutex{}
	return b, nil

}

// GetLastReceivedEventId get last event id submitted in this contract.
func (b *CommonBridge) GetLastReceivedEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// GetEventById get `Transfer` event (emitted by this contract) by id.
func (b *CommonBridge) GetEventById(eventId *big.Int) (*contracts.BridgeTransfer, error) {
	logs, err := b.Contract.FilterTransfer(nil, []*big.Int{eventId})
	if err != nil {
		return nil, fmt.Errorf("filter transfer: %w", err)
	}
	for logs.Next() {
		if !logs.Event.Raw.Removed {
			return logs.Event, nil
		}
	}
	return nil, networks.ErrEventNotFound
}

// GetEventsByIds gets contract events by ids.
func (b *CommonBridge) GetEventsByIds(eventIds []*big.Int) (transfers []*contracts.BridgeTransfer, err error) {
	logTransfer, err := b.Contract.FilterTransfer(nil, eventIds)
	if err != nil {
		return nil, fmt.Errorf("filter transfer: %w", err)
	}

	for logTransfer.Next() {
		if !logTransfer.Event.Raw.Removed {
			transfers = append(transfers, logTransfer.Event)
		}
	}
	return transfers, nil
}

func (b *CommonBridge) GetMinSafetyBlocksNum() (uint64, error) {
	safetyBlocks, err := b.Contract.MinSafetyBlocks(nil)
	if err != nil {
		return 0, err
	}
	return safetyBlocks.Uint64(), nil
}

func (b *CommonBridge) GetName() string {
	return b.Name
}

func (b *CommonBridge) GetClient() ethclients.ClientInterface {
	return b.Client
}

func (b *CommonBridge) GetContract() *contracts.Bridge {
	return b.Contract
}

func (b *CommonBridge) GetWsContract() *contracts.Bridge {
	return b.WsContract
}

func (b *CommonBridge) GetLogger() *zerolog.Logger {
	return &b.Logger
}
