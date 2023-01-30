package bsc

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/events"
	common_ethclient "github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients_rpc/rate_limiter"
	"github.com/rs/zerolog"
)

const BridgeName = "binance"

type Bridge struct {
	nc.CommonBridge
}

// New creates a new ethereum bridge.
func New(cfg *config.Network, sideBridgeName string, eventsApi events.Events, baseLogger zerolog.Logger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg, BridgeName, eventsApi)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}

	commonBridge.Logger = baseLogger.With().Str("bridge", BridgeName).Str("sideBridge", sideBridgeName).Logger()

	// ///////////////////
	origin := nc.GetAmbrosusOrigin()

	rpcHTTPClient, err := rate_limiter.DialHTTP(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}
	rpcHTTPClient.SetHeader("Origin", origin)

	commonBridge.Client = common_ethclient.NewClient(rpcHTTPClient)

	// Creating a new bridge contract instance.
	commonBridge.Contract, err = bindings.NewBridge(commonBridge.ContractAddress, commonBridge.Client)
	if err != nil {
		return nil, fmt.Errorf("create contract http: %w", err)
	}

	return &Bridge{
		CommonBridge: commonBridge,
	}, nil
}
