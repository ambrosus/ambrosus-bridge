package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
)

func (b *CommonBridge) ValidityWatchdog(sideBridge networks.Bridge) {
	b.SideBridge = sideBridge

	for {
		b.EnsureContractUnpaused()

		if err := b.watchLockedTransfers(); err != nil {
			b.Logger.Error().Msgf("ValidityWatchdog: %s", err)
		}
	}
}

func (b *CommonBridge) checkOldLockedTransfers() error {
	b.Logger.Info().Msg("Checking old transfer submit events...")

	oldestLockedEventId, err := b.Contract.OldestLockedEventId(nil)
	if err != nil {
		return fmt.Errorf("get oldest locked event id: %w", err)
	}

	for i := int64(0); ; i++ {
		nextLockedEventId := new(big.Int).Add(oldestLockedEventId, big.NewInt(i))
		nextLockedTransfer, err := b.Contract.GetLockedTransfers(nil, nextLockedEventId)
		if err != nil {
			return fmt.Errorf("GetLockedTransfers: %w", err)
		}
		if nextLockedTransfer.EndTimestamp.Uint64() == 0 {
			return nil
		}

		b.Logger.Info().Str("event_id", nextLockedEventId.String()).Msg("Found old TransferSubmit event")
		if err := b.checkValidity(nextLockedEventId, &nextLockedTransfer); err != nil {
			return fmt.Errorf("checkValidity: %w", err)
		}
	}
}

func (b *CommonBridge) watchLockedTransfers() error {
	if err := b.checkOldLockedTransfers(); err != nil {
		return fmt.Errorf("checkOldLockedTransfers: %w", err)
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
			if event.Raw.Removed {
				continue
			}

			b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Found new TransferSubmit event")

			lockedTransfer, err := b.Contract.GetLockedTransfers(nil, event.EventId)
			if err != nil {
				return fmt.Errorf("GetLockedTransfers: %w", err)
			}
			if err := b.checkValidity(event.EventId, &lockedTransfer); err != nil {
				b.Logger.Error().Msgf("checkValidity: %s", err)
			}
		}
	}
}

func (b *CommonBridge) checkValidity(lockedEventId *big.Int, lockedTransfer *contracts.CommonStructsLockedTransfers) error {
	sideEvent, err := b.SideBridge.GetEventById(lockedEventId)
	if err != nil && !errors.Is(err, networks.ErrEventNotFound) { // we'll handle the ErrEventNotFound later
		return fmt.Errorf("GetEventById: %w", err)
	}

	thisTransfers, err := json.MarshalIndent(lockedTransfer.Transfers, "", "  ")
	if err != nil {
		return fmt.Errorf("thisLockedTransfer marshal: %w", err)
	}
	sideTransfers := []byte("event not found")
	if sideEvent != nil {
		sideTransfers, err = json.MarshalIndent(sideEvent.Queue, "", "  ")
		if err != nil {
			return fmt.Errorf("thisLockedTransfer marshal: %w", err)
		}
	}

	if bytes.Equal(thisTransfers, sideTransfers) {
		b.Logger.Debug().Str("event_id", lockedEventId.String()).Msg("Locked transfers are equal")
		return nil
	}

	b.Logger.Warn().Str("event_id", lockedEventId.String()).Msgf(`
checkValidity: Transfers mismatch in event %d\n
this network locked transfers: %s \n
side network transfer event: %s \n
Pausing contract...`, lockedEventId, thisTransfers, sideTransfers)

	tx, txErr := b.Contract.Pause(b.Auth)
	if err := b.ProcessTx(networks.GetTxErrParams{Tx: tx, TxErr: txErr, MethodName: "pause"}); err != nil {
		return fmt.Errorf("pausing contract: %w", err)
	}

	b.Logger.Warn().Str("event_id", lockedEventId.String()).Msg("checkValidity: Contract has been paused")
	return nil

}
