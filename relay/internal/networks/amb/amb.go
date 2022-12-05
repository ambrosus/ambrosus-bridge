package amb

import (
	"context"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/events"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/parity"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
)

const BridgeName = "ambrosus"

type Bridge struct {
	nc.CommonBridge
	ParityClient *parity.Client
}

// New creates a new ambrosus bridge.
func New(cfg *config.Network, sideBridgeName string, eventsApi events.Events, baseLogger zerolog.Logger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg, BridgeName, eventsApi)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}

	commonBridge.Logger = baseLogger.With().Str("bridge", BridgeName).Str("sideBridge", sideBridgeName).Logger()

	// ///////////////////
	origin := nc.GetAmbrosusOrigin()

	rpcHTTPClient, err := rpc.DialHTTP(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}
	rpcHTTPClient.SetHeader("Origin", origin)

	parityClient := parity.NewClient(rpcHTTPClient)
	commonBridge.Client = parityClient

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
		commonBridge.WsClient = parity.NewClient(rpcWSClient)
	}

	return &Bridge{
		CommonBridge: commonBridge,
	}, nil
}
