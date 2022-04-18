package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

type CommonBridge struct {
	networks.Bridge
	Client     *ethclient.Client
	WsClient   *ethclient.Client
	Contract   *contracts.Bridge
	WsContract *contracts.Bridge
	Auth       *bind.TransactOpts
	SideBridge networks.Bridge
	Logger     zerolog.Logger
}

// GetLastEventId gets last contract event id.
func (b *CommonBridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// GetEventById gets contract event by id.
func (b *CommonBridge) GetEventById(eventId *big.Int) (*contracts.BridgeTransfer, error) {
	opts := &bind.FilterOpts{Context: context.Background()}

	logs, err := b.Contract.FilterTransfer(opts, []*big.Int{eventId})
	if err != nil {
		return nil, fmt.Errorf("filter transfer: %w", err)
	}

	if logs.Next() {
		return logs.Event, nil
	}

	return nil, networks.ErrEventNotFound
}

func (b *CommonBridge) GetMinSafetyBlocksNum() (uint64, error) {
	safetyBlocks, err := b.Contract.MinSafetyBlocks(nil)
	if err != nil {
		return 0, err
	}
	return safetyBlocks.Uint64(), nil
}

func (b *CommonBridge) CheckOldEvents() error {
	b.Logger.Info().Msg("Checking old events...")

	lastEventId, err := b.SideBridge.GetLastEventId()
	if err != nil {
		return fmt.Errorf("GetLastEventId: %w", err)
	}

	i := big.NewInt(1)
	for {
		nextEventId := big.NewInt(0).Add(lastEventId, i)
		nextEvent, err := b.GetEventById(nextEventId)
		if err != nil {
			if errors.Is(err, networks.ErrEventNotFound) {
				// no more old events
				return nil
			}
			return fmt.Errorf("GetEventById on id %v: %w", nextEventId.String(), err)
		}

		b.Logger.Info().Str("event_id", nextEventId.String()).Msg("Send old event...")

		if err := b.SendEvent(nextEvent); err != nil {
			return fmt.Errorf("send event: %w", err)
		}

		i = big.NewInt(0).Add(i, big.NewInt(1))
	}
}

func (b *CommonBridge) Listen() error {
	if err := b.CheckOldEvents(); err != nil {
		return fmt.Errorf("CheckOldEvents: %w", err)
	}

	b.Logger.Info().Msg("Listening new events...")

	// Subscribe to events
	watchOpts := &bind.WatchOpts{Context: context.Background()}
	eventChannel := make(chan *contracts.BridgeTransfer)
	eventSub, err := b.WsContract.WatchTransfer(watchOpts, eventChannel, nil)
	if err != nil {
		return fmt.Errorf("watchTransfer: %w", err)
	}

	defer eventSub.Unsubscribe()

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching transfers: %w", err)
		case event := <-eventChannel:
			b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Send event...")

			if err := b.SendEvent(event); err != nil {
				return fmt.Errorf("send event: %w", err)
			}
		}
	}
}

func (b *CommonBridge) UnlockOldestTransfersLoop() {
	for {
		if err := b.UnlockOldestTransfers(); err != nil {
			b.Logger.Error().Msgf("UnlockOldestTransfersLoop: %s", err)
		}
	}
}

func (b *CommonBridge) UnlockOldestTransfers() error {
	// Get oldest transfer timestamp.
	oldestLockedEventId, err := b.Contract.OldestLockedEventId(nil)
	if err != nil {
		return fmt.Errorf("get oldest locked event id: %w", err)
	}
	lockedTransferTime, err := b.Contract.LockedTransfers(nil, oldestLockedEventId)
	if err != nil {
		return fmt.Errorf("get locked transfer time %v: %w", oldestLockedEventId, err)
	}
	if lockedTransferTime.Cmp(big.NewInt(0)) == 0 {
		lockTime, err := b.Contract.LockTime(nil)
		if err != nil {
			return fmt.Errorf("get lock time: %w", err)
		}

		b.Logger.Info().Str("event_id", oldestLockedEventId.String()).Msgf(
			"UnlockOldestTransfers: there are no locked transfers with that id. Sleep %v seconds...",
			lockTime.Uint64(),
		)
		time.Sleep(time.Duration(lockTime.Uint64()) * time.Second)
		return nil
	}

	// Get the latest block.
	latestBlock, err := b.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("get latest block: %w", err)
	}

	// Check if the unlocking is allowed and get the sleep time.
	sleepTime := lockedTransferTime.Int64() - int64(latestBlock.Time())
	if sleepTime > 0 {
		b.Logger.Info().Str("event_id", oldestLockedEventId.String()).Msgf(
			"UnlockOldestTransfers: sleep %v seconds...",
			sleepTime,
		)
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}

	// Unlock the oldest transfer.
	b.Logger.Info().Str("event_id", oldestLockedEventId.String()).Msg("UnlockOldestTransfers: unlocking...")
	err = b.unlockTransfers(oldestLockedEventId)
	if err != nil {
		return fmt.Errorf("unlock locked transfer %v: %w", oldestLockedEventId, err)
	}
	b.Logger.Info().Str("event_id", oldestLockedEventId.String()).Msg("UnlockOldestTransfers: unlocked")
	return nil
}

func (b *CommonBridge) unlockTransfers(eventId *big.Int) error {
	tx, txErr := b.Contract.UnlockTransfers(b.Auth, eventId)
	return b.GetTransactionError(
		networks.GetTransactionErrorParams{Tx: tx, TxErr: txErr, MethodName: "unlockTransfers"},
		eventId,
	)
}

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
