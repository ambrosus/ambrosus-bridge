package common

import (
	"fmt"
	"math/big"
	"time"
)

func (b *CommonBridge) TriggerTransfersLoop() {
	b.shouldHavePk()
	for {
		b.EnsureContractUnpaused()

		if err := b.checkWithdraws(); err != nil {
			b.Logger.Error().Err(err).Msg("watchWithdraws error")
		}
		time.Sleep(failSleepTIme)
	}
}

func (b *CommonBridge) checkWithdraws() error {
	b.Logger.Info().Msg("Checking old events...")

	timeFrameSeconds, lastTimeFrame, err := fetchTimeParams(b)
	if err != nil {
		return fmt.Errorf("fetchTimeParams error: %w", err)
	}

	// when we should trigger transfers from lastTimeFrame
	triggerAt := calcTriggerAt(lastTimeFrame, timeFrameSeconds)

	if time.Until(triggerAt) > 0 {
		b.Logger.Info().Msg("Sleep until next time frame")
		time.Sleep(time.Until(triggerAt)) // sleep to the moment where we should trigger transfers
		return nil                        // return so we can get actual `lastTimeFrame` value in next iteration
	}

	queue, err := b.Contract.Queue(nil, big.NewInt(0))
	if err != nil {
		return fmt.Errorf("Contract.Queue error: %w", err)
	} else if queue.Amount.Int64() == 0 {
		b.Logger.Info().Msg("Queue empty, skipping...")

		// if lastTimeFrame has no transfers we should sleep at least until current time frame end
		currentTimeFrame := time.Now().Unix() / timeFrameSeconds
		triggerAt = calcTriggerAt(currentTimeFrame, timeFrameSeconds)
		time.Sleep(time.Until(triggerAt)) // sleep to the moment where we should trigger transfers
		return nil
	}

	return b.triggerTransfers()
}

func (b *CommonBridge) triggerTransfers() error {
	b.Logger.Info().Msg("Triggering transfers...")

	// todo trigger transfers
	return nil
}

func calcTriggerAt(timeFrameId, timeFrameSeconds int64) time.Time {
	// end of timeFrameId + 1/4 of timeFrameSeconds
	triggerAtUnix := (timeFrameId+1)*timeFrameSeconds + timeFrameSeconds*1/4
	return time.Unix(triggerAtUnix, 0)
}

func fetchTimeParams(b *CommonBridge) (timeFrame, lastTimeFrame int64, err error) {
	timeFrameSeconds, err := b.Contract.TimeframeSeconds(nil)
	if err != nil {
		return 0, 0, fmt.Errorf("TimeframeSeconds: %w", err)
	}

	lastTimeframe, err := b.Contract.LastTimeframe(nil)
	if err != nil {
		return 0, 0, fmt.Errorf("LastTimeframe: %w", err)
	}

	return timeFrameSeconds.Int64(), lastTimeframe.Int64(), nil
}
