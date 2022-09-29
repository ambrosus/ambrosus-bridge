package service_pause_unpause_watchdog

import (
	"fmt"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/rs/zerolog"
)

const (
	MsgContractHasBeenPaused   = "Contract has been paused!"
	MsgContractHasBeenUnpaused = "Contract has been unpaused!"
)

type WatchPauseUnpauseBridgeContract struct {
	bridge networks.Bridge
	logger *zerolog.Logger

	lastPausedStatus bool
}

func NewWatchPauseUnpauseBridgeContract(bridge networks.Bridge) *WatchPauseUnpauseBridgeContract {
	logger := bridge.GetLogger().With().Str("service", "WatchPauseUnpauseBridgeContract").Logger()

	return &WatchPauseUnpauseBridgeContract{
		bridge:           bridge,
		logger:           &logger,
		lastPausedStatus: false, // if the contract is paused at the relay startup - log that
	}
}

func (b *WatchPauseUnpauseBridgeContract) Run() {
	for {
		if err := b.watchPauseUnpauseBridgeContract(); err != nil {
			b.logger.Error().Err(err).Msg("")
			time.Sleep(1 * time.Minute)
		}
	}
}

func (b *WatchPauseUnpauseBridgeContract) watchPauseUnpauseBridgeContract() error {
	paused, err := b.bridge.GetContract().Paused(nil)
	if err != nil {
		return fmt.Errorf("Paused: %w", err)
	}
	if paused != b.lastPausedStatus {
		var msg string
		if paused {
			msg = MsgContractHasBeenPaused
		} else {
			msg = MsgContractHasBeenUnpaused
		}

		b.logger.Warn().Msg(msg)
		b.lastPausedStatus = paused
	}

	if paused {
		if err := b.waitForUnpauseContract(); err != nil {
			return fmt.Errorf("waitForUnpauseContract: %w", err)
		}
	} else {
		if err := b.waitForPauseContract(); err != nil {
			return fmt.Errorf("waitForPauseContract: %w", err)
		}
	}

	return nil
}

func (b *WatchPauseUnpauseBridgeContract) waitForPauseContract() error {
	eventCh := make(chan *bindings.BridgePaused)
	eventSub, err := b.bridge.GetWsContract().WatchPaused(nil, eventCh)
	if err != nil {
		return fmt.Errorf("WatchPaused: %w", err)
	}
	defer eventSub.Unsubscribe()

	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching paused event: %w", err)
		case event := <-eventCh:
			if event.Raw.Removed {
				continue
			}

			b.logger.Warn().Msg(MsgContractHasBeenPaused)
			b.lastPausedStatus = true
			return nil
		}
	}
}

func (b *WatchPauseUnpauseBridgeContract) waitForUnpauseContract() error {
	eventCh := make(chan *bindings.BridgeUnpaused)
	eventSub, err := b.bridge.GetWsContract().WatchUnpaused(nil, eventCh)
	if err != nil {
		return fmt.Errorf("WatchUnpaused: %w", err)
	}
	defer eventSub.Unsubscribe()

	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching unpaused event: %w", err)
		case event := <-eventCh:
			if event.Raw.Removed {
				continue
			}

			b.logger.Warn().Msg(MsgContractHasBeenUnpaused)
			b.lastPausedStatus = false
			return nil
		}
	}
}
