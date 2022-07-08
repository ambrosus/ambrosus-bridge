package bsc

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/rs/zerolog"
)

const BridgeName = "binance"

type Bridge struct {
	nc.CommonBridge
}

// New creates a new ethereum bridge.
func New(cfg *config.BSCConfig, baseLogger zerolog.Logger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}

	commonBridge.Logger = baseLogger.With().Str("bridge", BridgeName).Logger()

	return &Bridge{
		CommonBridge: commonBridge,
	}, nil
}

//
//func (b *Bridge) Run() {
//	b.Logger.Debug().Msg("Running binance bridge...")
//
//	go b.UnlockTransfersLoop()
//	b.SubmitTransfersLoop()
//}
