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
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	go b.UnlockTransfersLoop()

	b.Logger.Info().Msg("Ambrosus bridge runned!")

	b.ListenTransfersLoop()
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

func (b *Bridge) GetTxErr(params networks.GetTxErrParams) error {
	if params.TxErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message at eth_estimateGas, so
		// do eth_call method to get the full error message
		err := b.Contract.Raw().Call(&bind.CallOpts{From: b.Auth.From}, nil, params.MethodName, params.TxParams...)
		if err != nil {
			return fmt.Errorf("getFailureReasonViaCall: %w", helpers.ParseError(err))
		}
		return params.TxErr
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
