package service_trigger

import (
	"fmt"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

type TriggerTransfers struct {
	bridge networks.Bridge
	logger *zerolog.Logger
}

func NewTriggerTransfers(bridge networks.Bridge) *TriggerTransfers {
	logger := bridge.GetLogger().With().Str("service", "TriggerTransfers").Logger()

	return &TriggerTransfers{
		bridge: bridge,
		logger: &logger,
	}
}

func (b *TriggerTransfers) Run() {
	cb.ShouldHavePk(b.bridge)
	for {
		cb.EnsureContractUnpaused(b.bridge)

		if err := b.checkTriggerTransfers(); err != nil {
			b.logger.Error().Err(err).Msg("")
		}
		time.Sleep(1 * time.Minute)
	}
}

func (b *TriggerTransfers) checkTriggerTransfers() error {
	b.logger.Info().Msg("checkTriggerTransfers... ")

	timeFrameSeconds, lastTimeFrame, err := fetchTimeParams(b.bridge.GetContract())
	if err != nil {
		return fmt.Errorf("fetchTimeParams error: %w", err)
	}

	// when we should trigger transfers from lastTimeFrame
	triggerAt := calcTriggerAt(lastTimeFrame, timeFrameSeconds)

	remained := time.Until(triggerAt)
	if remained > 0 {
		b.logger.Info().Msgf("Sleep until next time frame (%s)", remained.String())
		time.Sleep(remained) // sleep to the moment where we should trigger transfers
		return nil           // return so we can get actual `lastTimeFrame` value in next iteration
	}

	isQueueEmpty, err := b.bridge.GetContract().IsQueueEmpty(nil)
	if err != nil {
		return fmt.Errorf("IsQueueEmpty: %w", err)
	} else if isQueueEmpty {
		// if lastTimeFrame has no transfers we should sleep at least until current time frame end
		currentTimeFrame := time.Now().Unix() / timeFrameSeconds
		triggerAt = calcTriggerAt(currentTimeFrame, timeFrameSeconds)

		remained := time.Until(triggerAt)
		b.logger.Info().Msgf("Queue empty, skipping... (sleep for %s)", remained.String())
		time.Sleep(remained) // sleep to the moment where we should trigger transfers
		return nil
	}

	if err := b.triggerTransfers(); err != nil {
		return fmt.Errorf("triggerTransfers error: %w", err)
	}
	return nil
}

func (b *TriggerTransfers) triggerTransfers() error {
	b.logger.Info().Msg("Triggering transfers...")

	return b.bridge.ProcessTx("triggerTransfers", b.bridge.GetAuth(), func(opts *bind.TransactOpts) (*types.Transaction, error) {
		return b.bridge.GetContract().TriggerTransfers(opts)
	})
}

func calcTriggerAt(timeFrameId, timeFrameSeconds int64) time.Time {
	// end of timeFrameId + 1/4 of timeFrameSeconds
	triggerAtUnix := (timeFrameId+1)*timeFrameSeconds + timeFrameSeconds*1/4
	return time.Unix(triggerAtUnix, 0)
}

func fetchTimeParams(contract interfaces.BridgeContract) (timeFrame, lastTimeFrame int64, err error) {
	timeFrameSeconds, err := contract.TimeframeSeconds(nil)
	if err != nil {
		return 0, 0, fmt.Errorf("TimeframeSeconds: %w", err)
	}

	lastTimeframe, err := contract.LastTimeframe(nil)
	if err != nil {
		return 0, 0, fmt.Errorf("LastTimeframe: %w", err)
	}

	return timeFrameSeconds.Int64(), lastTimeframe.Int64(), nil
}
