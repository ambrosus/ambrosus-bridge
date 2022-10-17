package untrustless_mpc

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/rs/zerolog"
)

type SubmitterUntrustlessMpc struct {
	networks.Bridge
	untrustlessReceiver service_submit.ReceiverUntrustlessMpc
	logger              *zerolog.Logger

	isServer bool
}

func NewSubmitterUntrustlessMpc(bridge networks.Bridge, untrustlessReceiver service_submit.ReceiverUntrustlessMpc, isServer bool) (*SubmitterUntrustlessMpc, error) {
	logger := bridge.GetLogger().With().Str("service", "SubmitterUntrustlessMpc").Logger()

	return &SubmitterUntrustlessMpc{
		Bridge:              bridge,
		untrustlessReceiver: untrustlessReceiver,
		logger:              &logger,
		isServer:            isServer,
	}, nil
}

func (b *SubmitterUntrustlessMpc) Receiver() service_submit.Receiver {
	return b.untrustlessReceiver
}

func (b *SubmitterUntrustlessMpc) SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error {
	// check is it server or client
	if b.isServer {
		b.logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer untrustless mpc as server...")
		err := b.untrustlessReceiver.SubmitTransferUntrustlessMpcServer(event)
		if err != nil {
			return fmt.Errorf("SubmitTransferUntrustlessMpcServer: %w", err)
		}
		metric.AddWithdrawalsCountMetric(b, len(event.Queue))
	} else {
		b.logger.Info().Str("event_id", event.EventId.String()).Msg("Submit transfer untrustless mpc as client...")
		err := b.untrustlessReceiver.SubmitTransferUntrustlessMpcClient(event)
		if err != nil {
			return fmt.Errorf("SubmitTransferUntrustlessMpcClient: %w", err)
		}
	}

	return nil
}
