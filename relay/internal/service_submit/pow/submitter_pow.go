package pow

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/pow/pow_proof"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/rs/zerolog"
)

type SubmitterPoW struct {
	networks.Bridge
	powReceiver service_submit.ReceiverPoW
	powEncoder  *pow_proof.PoWEncoder
	ethash      *ethash.Ethash
	logger      *zerolog.Logger
}

func NewSubmitterPoW(bridge networks.Bridge, powReceiver service_submit.ReceiverPoW, cfg *config.SubmitterPoW) (*SubmitterPoW, error) {
	logger := bridge.GetLogger().With().Str("service", "SubmitterPoW").Logger()
	ethashObj := ethash.New(cfg.EthashDir, cfg.EthashKeepPrevEpochs, cfg.EthashGenNextEpochs, logger)

	return &SubmitterPoW{
		Bridge:      bridge,
		powReceiver: powReceiver,
		powEncoder:  pow_proof.NewPoWEncoder(bridge, powReceiver, ethashObj),
		logger:      &logger,
	}, nil
}

func (b *SubmitterPoW) SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error {
	powProof, err := b.powEncoder.EncodePoWProof(event, safetyBlocks)
	if err != nil {
		return fmt.Errorf("encodePoWProof: %w", err)
	}

	for _, blockNum := range []uint64{event.Raw.BlockNumber, event.Raw.BlockNumber + safetyBlocks} {
		if err := b.checkEpochData(blockNum, event.EventId); err != nil {
			return fmt.Errorf("checkEpochData on block %v: %w", blockNum, err)
		}
	}

	b.logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer PoW...")
	err = b.powReceiver.SubmitTransferPoW(powProof)
	if err != nil {
		return fmt.Errorf("SubmitTransferPoW: %w", err)
	}
	return nil
}

func (b *SubmitterPoW) checkEpochData(blockNumber uint64, eventId *big.Int) error {
	epoch := blockNumber / 30000
	if isEpochSet, err := b.powReceiver.IsEpochSet(epoch); err != nil {
		return fmt.Errorf("IsEpochSet: %w", err)
	} else if isEpochSet {
		return nil
	}

	b.logger.Info().Str("event_id", eventId.String()).Msg("Submit epoch data...")
	epochData, err := b.ethash.GetEpochData(epoch)
	if err != nil {
		return fmt.Errorf("loadEpochDataFile: %w", err)
	}

	err = b.powReceiver.SubmitEpochData(epochData)
	if err != nil {
		return fmt.Errorf("SubmitEpochData: %w", err)
	}
	return nil
}

func (b *SubmitterPoW) ensureDAGsExists() {
	b.logger.Info().Msgf("Checking if DAG file exists...")

	// Getting last ethereum block number.
	blockNumber, err := b.GetClient().BlockNumber(context.Background())
	if err != nil {
		b.logger.Error().Err(err).Msgf("error getting last block number")
		return
	}
	b.ethash.GenDagForEpoch(blockNumber / 30000)
}
