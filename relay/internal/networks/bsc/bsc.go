package eth

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/external_logger"
	"github.com/ethereum/go-ethereum/rpc"
)

const BridgeName = "binance"
const forkProtectionBlocks = 10 // fork protection

type Bridge struct {
	nc.CommonBridge
	Config     *config.ETHConfig
	sideBridge networks.BridgeReceiveClique
}

// New creates a new ethereum bridge.
func New(cfg *config.ETHConfig, externalLogger external_logger.ExternalLogger) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}
	commonBridge.Logger = logger.NewSubLogger(BridgeName, externalLogger)

	b := &Bridge{
		CommonBridge: commonBridge,
		Config:       cfg,
	}
	b.CommonBridge.Bridge = b
	return b, nil
}

func (b *Bridge) Run(sideBridge networks.BridgeReceiveClique) {
	b.sideBridge = sideBridge
	b.CommonBridge.SideBridge = sideBridge

	b.Logger.Debug().Msg("Running binance bridge...")

	go b.UnlockTransfersLoop()
	b.ListenTransfersLoop()
}

func (b *Bridge) SendEvent(event *contracts.BridgeTransfer) error {
	if err := b.WaitForBlock(event.Raw.BlockNumber + forkProtectionBlocks); err != nil {
		return fmt.Errorf("WaitForBlock: %w", err)
	}

	if err := b.IsEventRemoved(event); err != nil {
		return fmt.Errorf("isEventRemoved: %w", err)
	}

	cliqueProof, err := b.encodeCliqueProof(event)
	if err != nil {
		return fmt.Errorf("encodeCliqueProof: %w", err)
	}

	b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer Clique...")
	err = b.sideBridge.SubmitTransferClique(cliqueProof)
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
