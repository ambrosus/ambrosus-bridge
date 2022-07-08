package eth

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/rs/zerolog"
)

const BridgeName = "ethereum"

type Bridge struct {
	nc.CommonBridge

	Ethash *ethash.Ethash
}

// New creates a new ethereum bridge.
func New(cfg *config.ETHConfig, baseLogger zerolog.Logger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}

	commonBridge.Logger = baseLogger.With().Str("bridge", BridgeName).Logger()

	return &Bridge{
		CommonBridge: commonBridge,
		Ethash:       ethash.New(cfg.EthashDir, cfg.EthashKeepPrevEpochs, cfg.EthashGenNextEpochs),
	}, nil
}

func (b *Bridge) Run() {
	b.Logger.Debug().Msg("Running ethereum bridge...")

	// go b.ensureDAGsExists()
	go b.UnlockTransfersLoop()
	go b.TriggerTransfersLoop()
	//b.SubmitTransfersLoop()
}
