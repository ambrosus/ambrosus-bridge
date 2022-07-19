package common

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
)

type CommonBridge struct {
	networks.Bridge
	SideBridge networks.Bridge

	Client, WsClient     ethclients.ClientInterface
	Contract, WsContract interfaces.BridgeContract
	Auth                 *bind.TransactOpts

	Logger zerolog.Logger
	Name   string
	Pk     *ecdsa.PrivateKey

	ContractCallLock *sync.Mutex
}

func New(cfg *config.Network, name string) (b CommonBridge, err error) {
	b.Name = name

	origin := GetAmbrosusOrigin()

	rpcHTTPClient, err := rpc.DialHTTP(cfg.HttpURL)
	if err != nil {
		return b, fmt.Errorf("dial http: %w", err)
	}
	rpcHTTPClient.SetHeader("Origin", origin)
	b.Client = ethclient.NewClient(rpcHTTPClient)

	// Creating a new bridge contract instance.
	b.Contract, err = bindings.NewBridge(common.HexToAddress(cfg.ContractAddr), b.Client)
	if err != nil {
		return b, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	if cfg.WsURL != "" {
		rpcWSClient, err := rpc.DialWebsocket(context.Background(), cfg.WsURL, origin)
		if err != nil {
			return b, fmt.Errorf("dial ws: %w", err)
		}
		b.WsClient = ethclient.NewClient(rpcWSClient)

		b.WsContract, err = bindings.NewBridge(common.HexToAddress(cfg.ContractAddr), b.WsClient)
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

func (b *CommonBridge) GetClient() ethclients.ClientInterface {
	return b.Client
}

func (b *CommonBridge) GetWsClient() ethclients.ClientInterface {
	return b.WsClient
}

func (b *CommonBridge) GetContract() interfaces.BridgeContract {
	return b.Contract
}

func (b *CommonBridge) GetWsContract() interfaces.BridgeContract {
	return b.WsContract
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
