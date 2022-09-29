package service_pause_unpause_watchdog

import (
	"fmt"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/rs/zerolog"
)

const (
	MsgPaused   = "Contract has been paused!"
	MsgUnpaused = "Contract has been unpaused!"
)

type WatchPauseUnpause struct {
	bridge networks.Bridge
	logger *zerolog.Logger

	lastPausedStatus bool
}

func NewWatchPauseUnpause(bridge networks.Bridge) *WatchPauseUnpause {
	logger := bridge.GetLogger().With().Str("service", "WatchPauseUnpause").Logger()

	return &WatchPauseUnpause{
		bridge:           bridge,
		logger:           &logger,
		lastPausedStatus: false, // if the contract is paused at the relay startup - log that
	}
}

func (b *WatchPauseUnpause) Run() {
	for {
		if err := b.watchPauseUnpause(); err != nil {
			b.logger.Error().Err(err).Msg("")
			time.Sleep(1 * time.Minute)
		}
	}
}

func (b *WatchPauseUnpause) watchPauseUnpause() error {
	paused, err := b.bridge.GetContract().Paused(nil)
	if err != nil {
		return fmt.Errorf("Paused: %w", err)
	}
	if paused != b.lastPausedStatus {
		var msg string
		if paused {
			msg = MsgPaused
		} else {
			msg = MsgUnpaused
		}

		b.logger.Warn().Msg(msg)
		b.lastPausedStatus = paused
	}

	if paused {
		if err := b.waitForUnpause(); err != nil {
			return fmt.Errorf("waitForUnpause: %w", err)
		}
	} else {
		if err := b.waitForPause(); err != nil {
			return fmt.Errorf("waitForPause: %w", err)
		}
	}

	return nil
}

func (b *WatchPauseUnpause) waitForPause() error {
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

			b.logger.Warn().Msg(MsgPaused)
			b.lastPausedStatus = true
			return nil
		}
	}
}

func (b *WatchPauseUnpause) waitForUnpause() error {
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

			b.logger.Warn().Msg(MsgUnpaused)
			b.lastPausedStatus = false
			return nil
		}
	}
}
