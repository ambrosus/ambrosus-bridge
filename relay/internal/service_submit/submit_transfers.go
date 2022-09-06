package service_submit

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/logger/middlewares"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/rs/zerolog"
)

type SubmitTransfers struct {
	submitter Submitter
	receiver  Receiver
	logger    *zerolog.Logger
}

func NewSubmitTransfers(submitter Submitter) *SubmitTransfers {
	logger := submitter.GetLogger().With().
		Str("relay", submitter.GetAuth().From.Hex()).
		Str("service", "SubmitTransfers").Logger()

	return &SubmitTransfers{
		submitter: submitter,
		receiver:  submitter.Receiver(),
		logger:    &logger,
	}
}


func (b *SubmitTransfers) Run() {
	cb.ShouldHavePk(b.receiver)
	b.logger.WithLevel(middlewares.TgLevel).Msg("Relay has been started!")

	for {
		// since we submit transfers to receiver, ensure that it is unpaused
		cb.EnsureContractUnpaused(b.receiver)

		if err := b.watchTransfers(); err != nil {
			b.logger.Error().Err(err).Msg("")
		}
		time.Sleep(1 * time.Minute)
	}
}

func (b *SubmitTransfers) checkOldTransfers() error {
	b.logger.Info().Msg("Checking old events...")

	lastEventId, err := b.receiver.GetContract().InputEventId(nil)
	if err != nil {
		return fmt.Errorf("GetLastReceivedEventId: %w", err)
	}

	for i := int64(1); ; i++ {
		nextEventId := new(big.Int).Add(lastEventId, big.NewInt(i))
		nextEvent, err := cb.GetEventById(b.submitter.GetContract(), nextEventId)
		if errors.Is(err, networks.ErrEventNotFound) { // no more old events
			return nil
		} else if err != nil {
			return fmt.Errorf("getEventById on id %v: %w", nextEventId.String(), err)
		}

		b.logger.Info().Str("event_id", nextEventId.String()).Msg("Send old event...")
		if err := b.processEvent(nextEvent); err != nil {
			return err
		}
	}
}

func (b *SubmitTransfers) watchTransfers() error {
	if err := b.checkOldTransfers(); err != nil {
		return fmt.Errorf("checkOldTransfers: %w", err)
	}
	b.logger.Info().Msg("Listening new events...")

	// Subscribe to events
	eventCh := make(chan *bindings.BridgeTransfer)
	eventSub, err := b.submitter.GetWsContract().WatchTransfer(nil, eventCh, nil)
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
			b.logger.Info().Str("event_id", event.EventId.String()).Msg("Send event...")
			if err := b.processEvent(event); err != nil {
				return err
			}
		}
	}
}

func (b *SubmitTransfers) processEvent(event *bindings.BridgeTransfer) error {
	safetyBlocks, err := getMinSafetyBlocksNum(b.receiver.GetContract())
	if err != nil {
		return fmt.Errorf("GetMinSafetyBlocksNum: %w", err)
	}

	b.logger.Debug().Uint64("blockNum", event.Raw.BlockNumber+safetyBlocks).Msg("Waiting for block...")
	if err := cb.WaitForBlock(b.submitter.GetWsClient(), event.Raw.BlockNumber+safetyBlocks); err != nil {
		return fmt.Errorf("waitForBlock: %w", err)
	}

	// Check if the event has been removed.
	if err := isEventRemoved(b.submitter.GetContract(), event); err != nil {
		return fmt.Errorf("isEventRemoved: %w", err)
	}

	if err := b.submitter.SendEvent(event, safetyBlocks); err != nil {
		return fmt.Errorf("send event: %w", err)
	}

	return nil
}

func isEventRemoved(contract interfaces.BridgeContract, event *bindings.BridgeTransfer) error {
	newEvent, err := cb.GetEventById(contract, event.EventId)
	if err != nil {
		return err
	}
	if newEvent.Raw.BlockHash != event.Raw.BlockHash {
		return fmt.Errorf("looks like the event has been removed")
	}
	return nil
}

func getMinSafetyBlocksNum(contract interfaces.BridgeContract) (uint64, error) {
	safetyBlocks, err := contract.MinSafetyBlocks(nil)
	if err != nil {
		return 0, err
	}
	return safetyBlocks.Uint64(), nil
}
