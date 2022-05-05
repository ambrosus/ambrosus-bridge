package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	nc "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ethereum/go-ethereum/rpc"
)

const BridgeName = "ethereum"

type Bridge struct {
	nc.CommonBridge
	sideBridge networks.BridgeReceiveEthash
	ethash     *ethash.Ethash
}

// New creates a new ethereum bridge.
func New(cfg *config.ETHConfig, externalLogger logger.Hook) (*Bridge, error) {
	commonBridge, err := nc.New(cfg.Network, BridgeName)
	if err != nil {
		return nil, fmt.Errorf("create commonBridge: %w", err)
	}
	commonBridge.Logger = logger.NewSubLogger(BridgeName, externalLogger)

	b := &Bridge{
		CommonBridge: commonBridge,
		ethash:       ethash.New(cfg.EthashDir, cfg.EthashKeepPrevEpochs, cfg.EthashGenNextEpochs),
	}
	b.CommonBridge.Bridge = b
	return b, nil
}

func (b *Bridge) SetSideBridge(sideBridge networks.BridgeReceiveEthash) {
	b.sideBridge = sideBridge
	b.CommonBridge.SideBridge = sideBridge
}

func (b *Bridge) Run() {
	b.Logger.Debug().Msg("Running ethereum bridge...")

	go b.ensureDAGsExists()
	go b.UnlockTransfersLoop()
	b.SubmitTransfersLoop()
}

func (b *Bridge) SendEvent(event *contracts.BridgeTransfer, safetyBlocks uint64) error {
	powProof, err := b.encodePoWProof(event, safetyBlocks)
	if err != nil {
		return fmt.Errorf("encodePoWProof: %w", err)
	}

	for _, blockNum := range []uint64{event.Raw.BlockNumber, event.Raw.BlockNumber + safetyBlocks} {
		if err := b.checkEpochData(blockNum, event.EventId); err != nil {
			return fmt.Errorf("checkEpochData on block %v: %w", blockNum, err)
		}
	}

	b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer PoW...")
	err = b.sideBridge.SubmitTransferPoW(powProof)
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

func (b *Bridge) checkEpochData(blockNumber uint64, eventId *big.Int) error {
	epoch := blockNumber / 30000
	if isEpochSet, err := b.sideBridge.IsEpochSet(epoch); err != nil {
		return fmt.Errorf("IsEpochSet: %w", err)
	} else if isEpochSet {
		return nil
	}

	b.Logger.Info().Str("event_id", eventId.String()).Msg("Submit epoch data...")
	epochData, err := b.ethash.GetEpochData(epoch)
	if err != nil {
		return fmt.Errorf("loadEpochDataFile: %w", err)
	}

	err = b.sideBridge.SubmitEpochData(epochData)
	if err != nil {
		return fmt.Errorf("SubmitEpochData: %w", err)
	}
	return nil
}

func (b *Bridge) ensureDAGsExists() {
	b.Logger.Info().Msgf("Checking if DAG file exists...")

	// Getting last ethereum block number.
	blockNumber, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		b.Logger.Error().Msgf("error getting last block number: %s", err.Error())
		return
	}
	b.ethash.GenDagForEpoch(blockNumber / 30000)
}
