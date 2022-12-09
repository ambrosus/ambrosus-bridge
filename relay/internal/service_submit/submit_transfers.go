package service_submit

import (
	"errors"
	"fmt"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/events"
	"github.com/rs/zerolog"
)

type SubmitTransfers struct {
	submitter Submitter
	receiver  Receiver
	logger    *zerolog.Logger
}

func NewSubmitTransfers(submitter Submitter) *SubmitTransfers {
	logger := submitter.GetLogger().With().
		Str("relayReceiver", submitter.Receiver().GetAuth().From.Hex()).
		Str("relaySubmitter", submitter.GetAuth().From.Hex()).
		Str("service", "SubmitTransfers").Logger()

	return &SubmitTransfers{
		submitter: submitter,
		receiver:  submitter.Receiver(),
		logger:    &logger,
	}
}

func (b *SubmitTransfers) Run() {
	cb.ShouldHavePk(b.receiver)

	for {
		// since we submit transfers to receiver, ensure that it is unpaused
		cb.EnsureContractUnpaused(b.receiver, b.logger)

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

	for i := uint64(1); ; i++ {
		nextEventId := lastEventId.Uint64() + i
		nextEvent, err := b.submitter.Events().GetTransfer(nextEventId)
		if errors.Is(err, events.ErrEventNotFound) { // no more old events
			return nil
		} else if err != nil {
			return fmt.Errorf("getEventById on id %v: %w", nextEventId, err)
		}

		b.logger.Info().Uint64("event_id", nextEventId).Msg("Send old event...")
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

	// main loop
	for {
		event, err := b.submitter.Events().WatchTransfer()
		if err != nil {
			return fmt.Errorf("watching transfers: %w", err)
		}
		b.logger.Info().Str("event_id", event.EventId.String()).Msg("Send event...")
		if err := b.processEvent(event); err != nil {
			return err
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
	if err := b.submitter.IsEventRemoved(&event.Raw); err != nil {
		return fmt.Errorf("IsEventRemoved: %w", err)
	}

	if err := b.submitter.SendEvent(event, safetyBlocks); err != nil {
		return fmt.Errorf("send event: %w", err)
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
