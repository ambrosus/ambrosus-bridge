package service_pause_unpause_watchdog

import (
	"fmt"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/rs/zerolog"
)

const (
	MsgPaused         = "Contract has been paused!"
	MsgUnpaused       = "Contract has been unpaused!"
	MsgWithDateFormat = "%s\n%s"
)

func msgWithDate(msg string) string {
	return fmt.Sprintf(MsgWithDateFormat, msg, time.Now().Format(time.UnixDate))
}

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
	for {
		err := b.bridge.Events().WatchPaused()
		if err != nil {
			return fmt.Errorf("watching paused event: %w", err)
		}
		b.logger.Warn().Msg(msgWithDate(MsgPaused))
		b.lastPausedStatus = true
		return nil
	}
}

func (b *WatchPauseUnpause) waitForUnpause() error {
	for {
		err := b.bridge.Events().WatchUnpaused()
		if err != nil {
			return fmt.Errorf("watching unpaused event: %w", err)
		}
		b.logger.Warn().Msg(msgWithDate(MsgUnpaused))
		b.lastPausedStatus = false
		return nil
	}
}
