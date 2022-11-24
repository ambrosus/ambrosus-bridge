package polygon

import (
	"context"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/limit_filtering"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients_rpc/rate_limiter"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
)

const BridgeName = "polygon"

type Bridge struct {
	nc.CommonBridge
}

// New creates a new polygon bridge.
func New(cfg *config.Network, sideBridgeName string, baseLogger zerolog.Logger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}

	rpcHTTPClient, err := rate_limiter.DialHTTP(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}

	limitFilteringClient := limit_filtering.NewClient(rpcHTTPClient, int64(29262814), int64(3499))
	commonBridge.Client = limitFilteringClient

	// Creating a new bridge contract instance.
	commonBridge.Contract, err = bindings.NewBridge(commonBridge.ContractAddress, commonBridge.Client)
	if err != nil {
		return nil, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	if cfg.WsURL != "" {
		rpcWSClient, err := rpc.DialWebsocket(context.Background(), cfg.WsURL, "")
		if err != nil {
			return nil, fmt.Errorf("dial ws: %w", err)
		}
		commonBridge.WsClient = limit_filtering.NewClient(rpcWSClient, int64(29262814), int64(3499))

		commonBridge.WsContract, err = bindings.NewBridge(commonBridge.ContractAddress, commonBridge.WsClient)
		if err != nil {
			return nil, fmt.Errorf("create contract ws: %w", err)
		}
	}

	commonBridge.Logger = baseLogger.With().Str("bridge", BridgeName).Str("sideBridge", sideBridgeName).Logger()

	return &Bridge{
		CommonBridge: commonBridge,
	}, nil
}
