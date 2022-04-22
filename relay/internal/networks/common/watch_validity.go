package common

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
)

func (b *CommonBridge) WatchValidityLockedTransfersLoop() {
	// todo don't pause if contract already paused
	// todo check old locked transfers and not only oldest
	// todo watcher should be run as a separate instance
	// todo rename to watchdog?
	return
	for {
		b.EnsureContractUnpaused()

		if err := b.WatchValidityLockedTransfers(); err != nil {
			b.Logger.Error().Msgf("WatchValidityLockedTransfersLoop: %s", err)
		}
	}
}

func (b *CommonBridge) WatchValidityLockedTransfers() error {
	eventCh := make(chan *contracts.BridgeTransferSubmit)
	eventSub, err := b.WsContract.WatchTransferSubmit(nil, eventCh, nil)
	if err != nil {
		return fmt.Errorf("WatchTransferSubmit: %w", err)
	}

	defer eventSub.Unsubscribe()

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching submit transfers: %w", err)
		case event := <-eventCh:
			b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Found new TransferSubmit event")

			if err := b.checkValidityLockedTransfers(event); err != nil {
				b.Logger.Error().Msgf("checkValidityLockedTransfers: %s", err)
			}
		}
	}
}

func (b *CommonBridge) checkValidityLockedTransfers(event *contracts.BridgeTransferSubmit) error {
	sideEvent, err := b.SideBridge.GetEventById(event.EventId)
	if err != nil {
		// todo if event not found it may be наебка too
		return fmt.Errorf("GetEventById: %w", err)
	}
	lockedTransfer, err := b.Contract.GetLockedTransfers(nil, event.EventId)
	if err != nil {
		return fmt.Errorf("GetLockedTransfers: %w", err)
	}

	thisTransfers, err := json.MarshalIndent(lockedTransfer.Transfers, "", "  ")
	if err != nil {
		return fmt.Errorf("thisLockedTransfer marshal: %w", err)
	}
	sideTransfers, err := json.MarshalIndent(sideEvent.Queue, "", "  ")
	if err != nil {
		return fmt.Errorf("thisLockedTransfer marshal: %w", err)
	}

	if bytes.Equal(thisTransfers, sideTransfers) {
		b.Logger.Debug().Str("event_id", event.EventId.String()).Msg("Locked transfers are equal")
		return nil
	}

	b.Logger.Warn().Str("event_id", event.EventId.String()).Msgf(`
checkValidityLockedTransfers: Transfers mismatch in event %d\n
this network locked transfers: %s \n
side network transfer event: %s \n
Pausing contract...`, event.EventId, thisTransfers, sideTransfers)

	tx, txErr := b.Contract.Pause(b.Auth)
	if err := b.GetTransactionError(networks.GetTransactionErrorParams{Tx: tx, TxErr: txErr, MethodName: "pause"}); err != nil {
		return fmt.Errorf("pausing contract: %w", err)
	}

	b.Logger.Warn().Str("event_id", event.EventId.String()).Msg("checkValidityLockedTransfers: Contract has been paused")
	return nil

}
