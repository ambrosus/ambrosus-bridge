package posa

import (
	"context"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit/posa/posa_proof"
	"github.com/rs/zerolog"
)

type SubmitterPoSA struct {
	networks.Bridge
	posaReceiver service_submit.ReceiverPoSA
	posaEncoder  *posa_proof.PoSAEncoder
	logger       *zerolog.Logger
}

func NewSubmitterPoSA(bridge networks.Bridge, posaReceiver service_submit.ReceiverPoSA) (*SubmitterPoSA, error) {
	chainId, err := bridge.GetClient().ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("chain id: %w", err)
	}

	return &SubmitterPoSA{
		Bridge:       bridge,
		posaReceiver: posaReceiver,
		posaEncoder:  posa_proof.NewPoSAEncoder(bridge, posaReceiver, chainId),
		logger:       bridge.GetLogger(), // todo maybe sublogger?
	}, nil
}

func (b *SubmitterPoSA) SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error {
	posaProof, err := b.posaEncoder.EncodePoSAProof(event, safetyBlocks)
	if err != nil {
		return fmt.Errorf("encodePoSAProof: %w", err)
	}

	b.logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer PoSA...")
	err = b.posaReceiver.SubmitTransferPoSA(posaProof)
	if err != nil {
		return fmt.Errorf("SubmitTransferPoW: %w", err)
	}
	return nil
}
