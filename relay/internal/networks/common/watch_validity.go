package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
)

func (b *CommonBridge) WatchValidityLockedTransfersLoop() {
	// todo don't pause if contract already paused
	// todo watcher should be run as a separate instance
	// todo rename to watchdog?
	return
	for {
		if err := b.WatchValidityLockedTransfers(); err != nil {
			b.Logger.Error().Msgf("WatchValidityLockedTransfersLoop: %s", err)
		}
	}
}

func (b *CommonBridge) CheckOldLockedTransfers() error {
	b.Logger.Info().Msg("Checking old transfer submit events...")

	oldestLockedEventId, err := b.Contract.OldestLockedEventId(nil)
	if err != nil {
		return fmt.Errorf("get oldest locked event id: %w", err)
	}

	for i := int64(1); ; i++ {
		nextLockedEventId := new(big.Int).Add(oldestLockedEventId, big.NewInt(i))
		nextLockedTransfer, err := b.Contract.GetLockedTransfers(nil, nextLockedEventId)
		if err != nil {
			return fmt.Errorf("GetLockedTransfers: %w", err)
		}
		if nextLockedTransfer.EndTimestamp.Cmp(big.NewInt(0)) == 0 {
			return nil
		}

		nextSideEvent, err := b.SideBridge.GetEventById(oldestLockedEventId)
		if err != nil {
			// todo if event not found it may be наебка too
			return fmt.Errorf("GetEventById: %w", err)
		}

		b.Logger.Info().Str("event_id", nextLockedEventId.String()).Msg("Found old TransferSubmit event")
		if err := b.checkValidityLockedTransfers(nextSideEvent, &nextLockedTransfer); err != nil {
			return fmt.Errorf("checkValidityLockedTransfers: %w", err)
		}
	}
}

func (b *CommonBridge) WatchValidityLockedTransfers() error {
	if err := b.CheckOldLockedTransfers(); err != nil {
		return fmt.Errorf("CheckOldLockedTransfers: %w", err)
	}

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

			sideEvent, err := b.SideBridge.GetEventById(event.EventId)
			if err != nil {
				// todo if event not found it may be наебка too
				return fmt.Errorf("GetEventById: %w", err)
			}
			lockedTransfer, err := b.Contract.GetLockedTransfers(nil, event.EventId)
			if err != nil {
				return fmt.Errorf("GetLockedTransfers: %w", err)
			}
			if err := b.checkValidityLockedTransfers(sideEvent, &lockedTransfer); err != nil {
				b.Logger.Error().Msgf("checkValidityLockedTransfers: %s", err)
			}
		}
	}
}

func (b *CommonBridge) checkValidityLockedTransfers(sideEvent *contracts.BridgeTransfer, lockedTransfer *contracts.CommonStructsLockedTransfers) error {
	thisTransfers, err := json.MarshalIndent(lockedTransfer.Transfers, "", "  ")
	if err != nil {
		return fmt.Errorf("thisLockedTransfer marshal: %w", err)
	}
	sideTransfers, err := json.MarshalIndent(sideEvent.Queue, "", "  ")
	if err != nil {
		return fmt.Errorf("thisLockedTransfer marshal: %w", err)
	}

	if bytes.Equal(thisTransfers, sideTransfers) {
		b.Logger.Debug().Str("event_id", sideEvent.EventId.String()).Msg("Locked transfers are equal")
		return nil
	}

	b.Logger.Warn().Str("event_id", sideEvent.EventId.String()).Msgf(`
checkValidityLockedTransfers: Transfers mismatch in event %d\n
this network locked transfers: %s \n
side network transfer event: %s \n
Pausing contract...`, sideEvent.EventId, thisTransfers, sideTransfers)

	tx, txErr := b.Contract.Pause(b.Auth)
	if err := b.GetTransactionError(networks.GetTransactionErrorParams{Tx: tx, TxErr: txErr, MethodName: "pause"}); err != nil {
		return fmt.Errorf("pausing contract: %w", err)
	}

	b.Logger.Warn().Str("event_id", sideEvent.EventId.String()).Msg("checkValidityLockedTransfers: Contract has been paused")
	return nil

}
