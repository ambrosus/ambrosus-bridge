package bsc

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/ethereum/go-ethereum/rpc"
)

const BridgeName = "binance"

type Bridge struct {
	nc.CommonBridge
	sideBridge networks.BridgeReceivePoSA
	chainId    *big.Int // cache chainId, cos it used many times in encode_block
}

// New creates a new ethereum bridge.
func New(cfg *config.BSCConfig, externalLogger external_logger.ExternalLogger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}
	commonBridge.Logger = logger.NewSubLogger(BridgeName, externalLogger)

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
}

func (b *Bridge) Run() {
	b.Logger.Debug().Msg("Running binance bridge...")

	go b.UnlockTransfersLoop()
	b.SubmitTransfersLoop()
}

func (b *Bridge) SendEvent(event *contracts.BridgeTransfer, safetyBlocks uint64) error {
	posaProof, err := b.encodePoSAProof(event)
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

func (b *Bridge) GetTxErr(params networks.GetTxErrParams) error {
	if params.TxErr != nil {
		if params.TxErr.Error() == "execution reverted" {
			dataErr := params.TxErr.(rpc.DataError)
			return fmt.Errorf("contract runtime error: %s", dataErr.ErrorData())
		}
		return params.TxErr
	}
	return nil
}
