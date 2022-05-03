package common

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
)

func (b *CommonBridge) SubmitTransfersLoop() {
	b.shouldHavePk()
	for {
		// since we submit transfers to SideBridge, ensure that it is unpaused
		b.SideBridge.EnsureContractUnpaused()

		if err := b.watchTransfers(); err != nil {
			b.Logger.Error().Err(err).Msg("watchTransfers error")
		}
		time.Sleep(failSleepTIme)
	}
}

func (b *CommonBridge) checkOldTransfers() error {
	b.Logger.Info().Msg("Checking old events...")

	lastEventId, err := b.SideBridge.GetLastEventId()
	if err != nil {
		return fmt.Errorf("GetLastEventId: %w", err)
	}

	for i := int64(1); ; i++ {
		nextEventId := new(big.Int).Add(lastEventId, big.NewInt(i))
		nextEvent, err := b.GetEventById(nextEventId)
		if errors.Is(err, networks.ErrEventNotFound) { // no more old events
			return nil
		} else if err != nil {
			return fmt.Errorf("GetEventById on id %v: %w", nextEventId.String(), err)
		}

		b.Logger.Info().Str("event_id", nextEventId.String()).Msg("Send old event...")
		if err := b.processEvent(nextEvent); err != nil {
			return err
		}
	}
}

func (b *CommonBridge) watchTransfers() error {
	if err := b.checkOldTransfers(); err != nil {
		return fmt.Errorf("checkOldTransfers: %w", err)
	}
	b.Logger.Info().Msg("Listening new events...")

	// Subscribe to events
	eventCh := make(chan *contracts.BridgeTransfer)
	eventSub, err := b.WsContract.WatchTransfer(nil, eventCh, nil)
	if err != nil {
		return fmt.Errorf("watchTransfer: %w", err)
	}
	defer eventSub.Unsubscribe()

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching transfers: %w", err)
		case event := <-eventCh:
			if event.Raw.Removed {
				continue
			}
			b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Send event...")
			if err := b.processEvent(event); err != nil {
				return err
			}
		}
	}
}

func (b *CommonBridge) processEvent(event *contracts.BridgeTransfer) error {
	safetyBlocks, err := b.SideBridge.GetMinSafetyBlocksNum(nil)
	if err != nil {
		return fmt.Errorf("GetMinSafetyBlocksNum: %w", err)
	}

	if err := b.WaitForBlock(event.Raw.BlockNumber + safetyBlocks); err != nil {
		return fmt.Errorf("WaitForBlock: %w", err)
	}

	// Check if the event has been removed.
	if err := b.IsEventRemoved(event); err != nil {
		return fmt.Errorf("isEventRemoved: %w", err)
	}

	if err := b.SendEvent(event, safetyBlocks); err != nil {
		return fmt.Errorf("send event: %w", err)
	}

	b.AddWithdrawalsCountMetric(len(event.Queue))
	return nil
}
