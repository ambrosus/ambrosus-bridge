package common

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/events"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	common_ethclient "github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
)

type CommonBridge struct {
	networks.Bridge
	SideBridge networks.Bridge

	Client, WsClient ethclients.ClientInterface
	Contract         interfaces.BridgeContract
	EventsApi        events.Events
	ContractAddress  common.Address
	Auth             *bind.TransactOpts

	Logger zerolog.Logger
	Name   string
	Pk     *ecdsa.PrivateKey

	ContractCallLock *sync.Mutex
}

func New(cfg *config.Network, name string, eventsApi events.Events) (b CommonBridge, err error) {
	b.Name = name
	b.ContractAddress = common.HexToAddress(cfg.ContractAddr)
	b.EventsApi = eventsApi

	origin := GetAmbrosusOrigin()

	rpcHTTPClient, err := rpc.DialHTTP(cfg.HttpURL)
	if err != nil {
		return b, fmt.Errorf("dial http: %w", err)
	}
	rpcHTTPClient.SetHeader("Origin", origin)
	b.Client = common_ethclient.NewClient(rpcHTTPClient)

	// Creating a new bridge contract instance.
	b.Contract, err = bindings.NewBridge(b.ContractAddress, b.Client)
	if err != nil {
		return b, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	if cfg.WsURL != "" {
		rpcWSClient, err := rpc.DialWebsocket(context.Background(), cfg.WsURL, origin)
		if err != nil {
			return b, fmt.Errorf("dial ws: %w", err)
		}
		b.WsClient = common_ethclient.NewClient(rpcWSClient)
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
		metric.SetRelayBalanceMetric(&b)
	} else {
		b.Logger.Info().Msg("No private key provided")
	}

	b.ContractCallLock = &sync.Mutex{}
	return b, nil

}

func GetAmbrosusOrigin() string {
	var origin string
	stage := os.Getenv("STAGE")
	if stage == "prod" {
		origin = "https://ambrosus.io"
	} else if stage == "test" {
		origin = "https://ambrosus-test.io"
	} else if stage == "dev" {
		origin = "https://ambrosus-dev.io"
	}
	return origin
}

// interface `Bridge`

func (b *CommonBridge) Events() events.Events {
	return b.EventsApi
}

func (b *CommonBridge) GetClient() ethclients.ClientInterface {
	return b.Client
}

func (b *CommonBridge) GetWsClient() ethclients.ClientInterface {
	return b.WsClient
}

func (b *CommonBridge) GetContract() interfaces.BridgeContract {
	return b.Contract
}

func (b *CommonBridge) GetLogger() *zerolog.Logger {
	return &b.Logger
}

func (b *CommonBridge) GetName() string {
	return b.Name
}

func (b *CommonBridge) GetAuth() *bind.TransactOpts {
	return b.Auth
}

func (b *CommonBridge) GetContractAddress() common.Address {
	return b.ContractAddress
}

func (b *CommonBridge) GetRelayAddress() common.Address {
	return b.Auth.From
}

func (b *CommonBridge) IsEventRemoved(eventLog *types.Log) error {
	header, err := b.Client.HeaderByNumber(context.Background(), big.NewInt(int64(eventLog.BlockNumber)))
	if err != nil {
		return fmt.Errorf("HeaderByNumber: %w", err)
	}
	if header.Hash() != eventLog.BlockHash {
		cid, err := b.Client.ChainID(context.Background())
		print("chainId", cid, err)
		return fmt.Errorf("%s != %s (blockNum: %v)", header.Hash().Hex(), eventLog.BlockHash.Hex(), eventLog.BlockNumber)
	}
	return nil
}
