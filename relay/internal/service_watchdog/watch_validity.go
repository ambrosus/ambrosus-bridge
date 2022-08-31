package service_watchdog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

var errEmptyLockedTransfers = errors.New("empty locked transfers")

type WatchTransfersValidity struct {
	bridge       networks.Bridge
	eventEmitter interfaces.BridgeContract
	logger       *zerolog.Logger
}

func NewWatchTransfersValidity(bridge networks.Bridge, eventEmitter interfaces.BridgeContract) *WatchTransfersValidity {
	logger := bridge.GetLogger().With().Str("service", "WatchTransfersValidity").Logger()

	return &WatchTransfersValidity{
		bridge:       bridge,
		eventEmitter: eventEmitter,
		logger:       &logger,
	}
}

func (b *WatchTransfersValidity) Run() {
	cb.ShouldHavePk(b.bridge)

	for {
		cb.EnsureContractUnpaused(b.bridge)

		if err := b.watchLockedTransfers(); err != nil {
			b.logger.Error().Err(err).Msg("")
		}
		time.Sleep(1 * time.Minute)
	}
}

func (b *WatchTransfersValidity) checkOldLockedTransfers() error {
	b.logger.Info().Msg("Checking old transfer submit events...")

	oldestLockedEventId, err := b.bridge.GetContract().OldestLockedEventId(nil)
	if err != nil {
		return fmt.Errorf("get oldest locked event id: %w", err)
	}

	return b.CheckOldLockedTransferFromId(oldestLockedEventId)
}

func (b *WatchTransfersValidity) CheckOldLockedTransferFromId(oldestLockedEventId *big.Int) error {
	for i := int64(0); ; i++ {
		nextLockedEventId := new(big.Int).Add(oldestLockedEventId, big.NewInt(i))
		nextLockedTransfer, err := b.getLockedTransfers(nextLockedEventId, nil)
		if errors.Is(err, errEmptyLockedTransfers) {
			return nil
		} else if err != nil {
			return fmt.Errorf("GetLockedTransfers: %w", err)
		}

		b.logger.Info().Str("event_id", nextLockedEventId.String()).Msg("Found old TransferSubmit event")
		if err := b.checkValidity(nextLockedEventId, &nextLockedTransfer); err != nil {
			return fmt.Errorf("checkValidity: %w", err)
		}
	}
}

func (b *WatchTransfersValidity) watchLockedTransfers() error {
	if err := b.checkOldLockedTransfers(); err != nil {
		return fmt.Errorf("checkOldLockedTransfers: %w", err)
	}

	eventCh := make(chan *bindings.BridgeTransferSubmit)
	eventSub, err := b.bridge.GetWsContract().WatchTransferSubmit(nil, eventCh, nil)
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

			b.logger.Info().Str("event_id", event.EventId.String()).Msg("Found new TransferSubmit event")

			// Nodes may be out of sync and the event may be empty (but it's not confirmed info)
			// So we need to retry if the result is empty
			lockedTransfer, err := b.getLockedTransfers(event.EventId, &bind.CallOpts{
				BlockNumber: big.NewInt(int64(event.Raw.BlockNumber)),
			})
			if err != nil {
				return fmt.Errorf("GetLockedTransfers: %w", err)
			}
			if err := b.checkValidity(event.EventId, &lockedTransfer); err != nil {
				b.logger.Error().Err(fmt.Errorf("checkValidity: %s", err)).Msg("checkValidity error")
			}
		}
	}
}

func (b *WatchTransfersValidity) checkValidity(lockedEventId *big.Int, lockedTransfer *bindings.CommonStructsLockedTransfers) error {
	sideEvent, err := common.GetEventById(b.eventEmitter, lockedEventId)
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
		b.logger.Debug().Str("event_id", lockedEventId.String()).Msg("Locked transfers are equal")
		return nil
	}

	b.logger.Warn().Str("event_id", lockedEventId.String()).Msgf(`
checkValidity: Transfers mismatch in event %d\n
this network locked transfers: %s \n
side network transfer event: %s \n
Pausing contract...`, lockedEventId, thisTransfers, sideTransfers)

	if err := b.bridge.ProcessTx("pause", b.bridge.GetAuth(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.bridge.GetContract().Pause(opts)
	}); err != nil {
		return fmt.Errorf("pausing contract: %w", err)
	}

	b.logger.Warn().Str("event_id", lockedEventId.String()).Msg("checkValidity: Contract has been paused")
	return nil

}

func (b *WatchTransfersValidity) getLockedTransfers(eventId *big.Int, opts *bind.CallOpts) (lockedTransfer bindings.CommonStructsLockedTransfers, err error) {
	err = retry.Do(
		func() error {
			lockedTransfer, err = b.bridge.GetContract().GetLockedTransfers(opts, eventId)
			if err != nil {
				return err
			}

			// check
			if lockedTransfer.EndTimestamp.Uint64() == 0 {
				return errEmptyLockedTransfers
			}
			return nil
		},

		retry.Attempts(5),
		retry.Delay(time.Second*2),
		retry.DelayType(retry.FixedDelay),
		retry.LastErrorOnly(true),
	)
	return lockedTransfer, err
}
