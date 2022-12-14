package bsc

import (
	"context"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/events"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/limit_filtering"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients_rpc/rate_limiter"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/mitchellh/mapstructure"
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

	var specificSettings *config.BSCSpecificSettings
	if err = mapstructure.Decode(cfg.SpecificSettings, &specificSettings); err != nil {
		return nil, fmt.Errorf("failed to cast `specificSettings`: %w", err)
	}

	rpcHTTPClient, err := rate_limiter.DialHTTP(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}
	rpcHTTPClient.SetHeader("Origin", origin)

	limitFilteringClient := limit_filtering.NewClient(rpcHTTPClient, specificSettings.FilterLogsFromBlock, specificSettings.FilterLogsLimitBlocks)
	commonBridge.Client = limitFilteringClient

	// Creating a new bridge contract instance.
	commonBridge.Contract, err = bindings.NewBridge(commonBridge.ContractAddress, commonBridge.Client)
	if err != nil {
		return nil, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	if cfg.WsURL != "" {
		rpcWSClient, err := rpc.DialWebsocket(context.Background(), cfg.WsURL, origin)
		if err != nil {
			return nil, fmt.Errorf("dial ws: %w", err)
		}
		commonBridge.WsClient = limit_filtering.NewClient(rpcWSClient, specificSettings.FilterLogsFromBlock, specificSettings.FilterLogsLimitBlocks)
	}

	return &Bridge{
		CommonBridge: commonBridge,
	}, nil
}
