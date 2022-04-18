package common

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (b *CommonBridge) WatchValidityLockedTransfersLoop() {
	for {
		if err := b.WatchValidityLockedTransfers(); err != nil {
			b.Logger.Error().Msgf("WatchValidityLockedTransfersLoop: %s", err)
		}
	}
}

func (b *CommonBridge) WatchValidityLockedTransfers() error {
	watchOpts := &bind.WatchOpts{Context: context.Background()}
	eventChannel := make(chan *contracts.BridgeTransferSubmit)
	eventSub, err := b.WsContract.WatchTransferSubmit(watchOpts, eventChannel, nil)
	if err != nil {
		return fmt.Errorf("WatchTransferSubmit: %w", err)
	}

	defer eventSub.Unsubscribe()

	if err := b.checkValidityLockedTransfers(&contracts.BridgeTransferSubmit{EventId: big.NewInt(1)}); err != nil {
		b.Logger.Error().Msgf("checkValidityLockedTransfers: %s", err)
		return nil
	}
	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching submit transfers: %w", err)
		case event := <-eventChannel:
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
		return fmt.Errorf("GetEventById: %w", err)
	}
	lockedTransfer, err := b.Contract.GetLockedTransfers(nil, event.EventId)
	if err != nil {
		return fmt.Errorf("GetLockedTransfers: %w", err)
	}
	thisLockedTransfers := lockedTransfer.Transfers
	sideLockedTransfers := sideEvent.Queue

	// Check the length equality.
	thisLockedTransfersLen := len(thisLockedTransfers)
	sideLockedTransfersLen := len(sideLockedTransfers)
	if thisLockedTransfersLen != sideLockedTransfersLen {
		b.Logger.Warn().Msgf(
			"checkValidityLockedTransfers: length of locked transfers of this network (%d) != of side network (%d). Pause contract...",
			thisLockedTransfersLen,
			sideLockedTransfersLen,
		)

		tx, txErr := b.Contract.Pause(b.Auth)
		if err := b.GetTransactionError(networks.GetTransactionErrorParams{Tx: tx, TxErr: txErr, MethodName: "pause"}); err != nil {
			return fmt.Errorf("Pause: %w", err)
		}

		b.Logger.Warn().Str("event_id", event.EventId.String()).Msg("checkValidityLockedTransfers: Contract has been paused")
		return nil
	}

	// Check the equality of locked transfers.
	for i := range thisLockedTransfers {
		if !reflect.DeepEqual(thisLockedTransfers[i], sideLockedTransfers[i]) {
			// Marshal transfers for logging.
			thisLockedTransferJson, err := json.MarshalIndent(thisLockedTransfers[i], "", "  ")
			if err != nil {
				return fmt.Errorf("thisLockedTransfer marshal: %w", err)
			}
			sideLockedTransferJson, err := json.MarshalIndent(sideLockedTransfers[i], "", "  ")
			if err != nil {
				return fmt.Errorf("sideLockedTransfer marshal: %w", err)
			}
			b.Logger.Warn().Str("event_id", event.EventId.String()).Msgf(`
checkValidityLockedTransfers: Queue mismatch in event %d by index %d!

this network locked transfer: %s

side network locked transfer: %s

Pause contract...`, event.EventId, i, thisLockedTransferJson, sideLockedTransferJson)

			tx, txErr := b.Contract.Pause(b.Auth)
			if err := b.GetTransactionError(networks.GetTransactionErrorParams{Tx: tx, TxErr: txErr, MethodName: "pause"}); err != nil {
				return fmt.Errorf("Pause: %w", err)
			}

			b.Logger.Warn().Str("event_id", event.EventId.String()).Msg("checkValidityLockedTransfers: Contract has been paused")
			return nil
		}
	}

	return nil
}
