package eth

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/events"
	"github.com/rs/zerolog"
)

const BridgeName = "ethereum"

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

	return &Bridge{
		CommonBridge: commonBridge,
	}, nil
}
