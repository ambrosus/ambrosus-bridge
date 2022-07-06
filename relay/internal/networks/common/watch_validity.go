package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
)

var ErrEmptyLockedTransfers = errors.New("empty locked transfers")

func (b *CommonBridge) ValidityWatchdog() {
	b.ShouldHavePk()

	for {
		b.EnsureContractUnpaused()

		if err := b.watchLockedTransfers(); err != nil {
			b.Logger.Error().Err(fmt.Errorf("ValidityWatchdog: %s", err)).Msg("ValidityWatchdog error")
		}
		time.Sleep(failSleepTIme)
	}
}

func (b *CommonBridge) checkOldLockedTransfers() error {
	b.Logger.Info().Msg("Checking old transfer submit events...")

	oldestLockedEventId, err := b.Contract.OldestLockedEventId(nil)
	if err != nil {
		return fmt.Errorf("get oldest locked event id: %w", err)
	}

	return b.checkOldLockedTransferFromId(oldestLockedEventId)
}

func (b *CommonBridge) checkOldLockedTransferFromId(oldestLockedEventId *big.Int) error {
	for i := int64(0); ; i++ {
		nextLockedEventId := new(big.Int).Add(oldestLockedEventId, big.NewInt(i))
		nextLockedTransfer, err := b.getLockedTransfers(nextLockedEventId, nil)
		if errors.Is(err, ErrEmptyLockedTransfers) {
			return nil
		} else if err != nil {
			return fmt.Errorf("GetLockedTransfers: %w", err)
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

	eventCh := make(chan *bindings.BridgeTransferSubmit)
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

			// Nodes may be out of sync and the event may be empty (but it's not confirmed info)
			// So we need to retry if the result is empty
			lockedTransfer, err := b.getLockedTransfers(event.EventId, &bind.CallOpts{
				BlockNumber: big.NewInt(int64(event.Raw.BlockNumber)),
			})
			if err != nil {
				return fmt.Errorf("GetLockedTransfers: %w", err)
			}
			if err := b.checkValidity(event.EventId, &lockedTransfer); err != nil {
				b.Logger.Error().Err(fmt.Errorf("checkValidity: %s", err)).Msg("checkValidity error")
			}
		}
	}
}

func (b *CommonBridge) checkValidity(lockedEventId *big.Int, lockedTransfer *bindings.CommonStructsLockedTransfers) error {
	sideEvent, err := b.SideBridge.GetEventById(lockedEventId)
	if err != nil && !errors.Is(err, networks.ErrEventNotFound) { // we'll handle the ErrEventNotFound later
		return fmt.Errorf("getEventById: %w", err)
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

	if err := b.ProcessTx("pause", func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.Contract.Pause(b.Auth)
	}); err != nil {
		return fmt.Errorf("pausing contract: %w", err)
	}

	b.Logger.Warn().Str("event_id", lockedEventId.String()).Msg("checkValidity: Contract has been paused")
	return nil

}

func (b *CommonBridge) getLockedTransfers(eventId *big.Int, opts *bind.CallOpts) (lockedTransfer bindings.CommonStructsLockedTransfers, err error) {
	err = retry.Do(
		func() error {
			lockedTransfer, err = b.Contract.GetLockedTransfers(opts, eventId)
			if err != nil {
				return err
			}

			// check
			if lockedTransfer.EndTimestamp.Uint64() == 0 {
				return ErrEmptyLockedTransfers
			}
			return nil
		},

		retry.Attempts(5),
		retry.Delay(time.Second*2),
		retry.LastErrorOnly(true),
	)
	return lockedTransfer, err
}
