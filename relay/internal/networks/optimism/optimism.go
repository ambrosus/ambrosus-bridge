package optimism

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/rs/zerolog"
)

const BridgeName = "optimism"

type Bridge struct {
	nc.CommonBridge
}

// New creates a new optimism bridge.
func New(cfg *config.Network, sideBridgeName string, baseLogger zerolog.Logger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}

	commonBridge.Logger = baseLogger.With().Str("bridge", BridgeName).Str("sideBridge", sideBridgeName).Logger()

	return &Bridge{
		CommonBridge: commonBridge,
	}, nil
}
