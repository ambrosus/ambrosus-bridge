package untrustless2

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/rs/zerolog"
)

type SubmitterUntrustless2 struct {
	networks.Bridge
	untrustlessReceiver service_submit.ReceiverUntrustless
	logger              *zerolog.Logger
}

func NewSubmitterUntrustless(bridge networks.Bridge, untrustlessReceiver service_submit.ReceiverUntrustless) (*SubmitterUntrustless2, error) {
	logger := bridge.GetLogger().With().Str("service", "SubmitterUntrustless2").Logger()

	return &SubmitterUntrustless2{
		Bridge:              bridge,
		untrustlessReceiver: untrustlessReceiver,
		logger:              &logger,
	}, nil
}

func (b *SubmitterUntrustless2) Receiver() service_submit.Receiver {
	return b.untrustlessReceiver
}

func (b *SubmitterUntrustless2) SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error {
	b.logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer untrustless...")
	err := b.untrustlessReceiver.SubmitTransferUntrustless(event)
	if err != nil {
		return fmt.Errorf("SubmitTransferUntrustless: %w", err)
	}

	metric.AddWithdrawalsCountMetric(b, len(event.Queue))
	return nil
}
