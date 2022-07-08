package amb

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/parity"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
)

const BridgeName = "ambrosus"

type Bridge struct {
	nc.CommonBridge
	ParityClient *parity.Client
}

// New creates a new ambrosus bridge.
func New(cfg *config.Network, baseLogger zerolog.Logger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}

	commonBridge.Logger = baseLogger.With().Str("bridge", BridgeName).Logger()

	// ///////////////////

	parityClient, err := parity.Dial(cfg.HttpURL)
	if err != nil {
		return nil, fmt.Errorf("dial http: %w", err)
	}
	commonBridge.Client = parityClient

	// Creating a new bridge contract instance.
	commonBridge.Contract, err = bindings.NewBridge(common.HexToAddress(cfg.ContractAddr), commonBridge.Client)
	if err != nil {
		return nil, fmt.Errorf("create contract http: %w", err)
	}

	// Create websocket instances if wsUrl provided
	if cfg.WsURL != "" {
		commonBridge.WsClient, err = parity.Dial(cfg.WsURL)
		if err != nil {
			return nil, fmt.Errorf("dial ws: %w", err)
		}

		commonBridge.WsContract, err = bindings.NewBridge(common.HexToAddress(cfg.ContractAddr), commonBridge.WsClient)
		if err != nil {
			return nil, fmt.Errorf("create contract ws: %w", err)
		}
	}

	return &Bridge{
		CommonBridge: commonBridge,
	}, nil
}
