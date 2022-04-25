package amb

import (
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/ethereum/go-ethereum/common"
)

const BridgeName = "ambrosus"

type Bridge struct {
	nc.CommonBridge
	VSContract *contracts.Vs
	Config     *config.AMBConfig
	sideBridge networks.BridgeReceiveAura
}

// New creates a new ambrosus bridge.
func New(cfg *config.AMBConfig, externalLogger external_logger.ExternalLogger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}
	commonBridge.Logger = logger.NewSubLogger(BridgeName, externalLogger)

	// Creating a new ambrosus VS contract instance.
	vsContract, err := contracts.NewVs(common.HexToAddress(cfg.VSContractAddr), commonBridge.Client)
	if err != nil {
		return nil, fmt.Errorf("create vs contract: %w", err)
	}

	b := &Bridge{
		CommonBridge: commonBridge,
		VSContract:   vsContract,
		Config:       cfg,
	}
	b.CommonBridge.Bridge = b
	return b, nil
}

func (b *Bridge) Run(sideBridge networks.BridgeReceiveAura) {
	b.sideBridge = sideBridge
	b.CommonBridge.SideBridge = sideBridge

	go b.UnlockOldestTransfersLoop()
	go b.WatchValidityLockedTransfersLoop()

	b.Logger.Info().Msg("Ambrosus bridge runned!")

	for {
		b.EnsureContractUnpaused()

		if err := b.Listen(); err != nil {
			b.Logger.Error().Err(err).Msg("Listen error")
		}
	}
}

func (b *Bridge) SendEvent(event *contracts.BridgeTransfer) error {
	b.Logger.Debug().Str("event_id", event.EventId.String()).Msg("Waiting for safety blocks...")

	// Wait for safety blocks.
	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return fmt.Errorf("GetMinSafetyBlocksNum: %w", err)
	}

	if err := b.WaitForBlock(event.Raw.BlockNumber + safetyBlocks); err != nil {
		return fmt.Errorf("WaitForBlock: %w", err)
	}

	b.Logger.Debug().Str("event_id", event.EventId.String()).Msg("Checking if the event has been removed...")

	// Check if the event has been removed.
	if err := b.isEventRemoved(event); err != nil {
		return fmt.Errorf("isEventRemoved: %w", err)
	}

	auraProof, err := b.getBlocksAndEvents(event, safetyBlocks)
	if err != nil {
		return fmt.Errorf("getBlocksAndEvents: %w", err)
	}

	b.Logger.Debug().Str("event_id", event.EventId.String()).Msg("Submit transfer Aura...")

	err = b.sideBridge.SubmitTransferAura(auraProof)
	if err != nil {
		return fmt.Errorf("SubmitTransferAura: %w", err)
	}
	return nil
}

func (b *Bridge) GetTransactionError(params networks.GetTransactionErrorParams, txParams ...interface{}) error {
	if params.TxErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message
		// so, make low-level call method to get the full error message
		err := b.getFailureReasonViaCall(params.MethodName, txParams...)
		if err != nil {
			return fmt.Errorf("getFailureReasonViaCall: %w", err)
		}
		return params.TxErr
	}
	return nil
}
func (b *Bridge) ProcessTx(params networks.GetTransactionErrorParams, txParams ...interface{}) error {
	if err := b.GetTransactionError(params, txParams...); err != nil {
		return err
	}

	err := b.waitForTxMined(params.Tx)
	if err != nil {
		return fmt.Errorf("waitForTxMined: %w", err)
	}
	return nil
}

func (b *Bridge) isEventRemoved(event *contracts.BridgeTransfer) error {
	block, err := b.HeaderByNumber(big.NewInt(int64(event.Raw.BlockNumber)))
	if err != nil {
		return fmt.Errorf("HeaderByNumber: %w", err)
	}

	if block.Hash(true) != event.Raw.BlockHash {
		return fmt.Errorf("block hash != event's block hash")
	}
	return nil
}
