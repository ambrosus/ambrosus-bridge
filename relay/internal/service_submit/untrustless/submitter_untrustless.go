package untrustless

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/rs/zerolog"
)

type SubmitterUntrustless struct {
	networks.Bridge
	untrustlessReceiver service_submit.ReceiverUntrustless
	logger              *zerolog.Logger
}

func NewSubmitterUntrustless(bridge networks.Bridge, untrustlessReceiver service_submit.ReceiverUntrustless) (*SubmitterUntrustless, error) {
	logger := bridge.GetLogger().With().Str("service", "SubmitterUntrustless").Logger()

	return &SubmitterUntrustless{
		Bridge:              bridge,
		untrustlessReceiver: untrustlessReceiver,
		logger:              &logger,
	}, nil
}

func (b *SubmitterUntrustless) SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error {
	// check is already confirmed by relay
	if isEventAlreadyConfirmed, err := b.untrustlessReceiver.IsEventAlreadyConfirmed(event); err != nil {
		return fmt.Errorf("is event already confirmed: %w", err)
	} else if isEventAlreadyConfirmed {
		b.logger.Info().Str("event_id", event.EventId.String()).Msg("Event is already confirmed, so skip...")
		return nil
	}

	b.logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer untrustless...")
	err := b.untrustlessReceiver.SubmitTransferUntrustless(event)
	if err != nil {
		return fmt.Errorf("SubmitTransferUntrustless: %w", err)
	}
	return nil
}
