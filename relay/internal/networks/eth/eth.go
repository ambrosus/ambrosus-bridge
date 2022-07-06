package eth

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
)

const BridgeName = "ethereum"

type Bridge struct {
	nc.CommonBridge

	Ethash *ethash.Ethash
}

// New creates a new ethereum bridge.
func New(cfg *config.ETHConfig, externalLogger logger.Hook) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}
	commonBridge.Logger = logger.NewSubLogger(BridgeName, externalLogger)

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
