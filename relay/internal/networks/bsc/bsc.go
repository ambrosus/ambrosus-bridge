package bsc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/bsc/posa_proof"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/rs/zerolog"
)

const BridgeName = "binance"

type Bridge struct {
	nc.CommonBridge
	sideBridge networks.BridgeReceivePoSA
	chainId    *big.Int // cache chainId, cos it used many times in encode_block

	posaEncoder *posa_proof.PoSAEncoder
}

// New creates a new ethereum bridge.
func New(cfg *config.BSCConfig, baseLogger zerolog.Logger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}

	commonBridge.Logger = baseLogger.With().Str("bridge", BridgeName).Logger()

	b := &Bridge{
		CommonBridge: commonBridge,
	}

	b.chainId, err = b.Client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}

	b.CommonBridge.Bridge = b
	return b, nil
}

func (b *Bridge) SetSideBridge(sideBridge networks.BridgeReceivePoSA) {
	b.sideBridge = sideBridge
	b.CommonBridge.SideBridge = sideBridge

	// todo refactor: move to separate service along with fetching chainId
	b.posaEncoder = posa_proof.NewPoSAEncoder(b, sideBridge, b.chainId)
}

func (b *Bridge) Run() {
	b.Logger.Debug().Msg("Running binance bridge...")

	go b.UnlockTransfersLoop()
	b.SubmitTransfersLoop()
}

func (b *Bridge) SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error {
	posaProof, err := b.posaEncoder.EncodePoSAProof(event, safetyBlocks)
	if err != nil {
		return fmt.Errorf("encodePoSAProof: %w", err)
	}

	b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer PoSA...")
	err = b.sideBridge.SubmitTransferPoSA(posaProof)
	if err != nil {
		return fmt.Errorf("SubmitTransferPoW: %w", err)
	}
	return nil
}
