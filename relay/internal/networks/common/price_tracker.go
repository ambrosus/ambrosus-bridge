package common

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

const (
	eventsForGasCalc = 5
)

// PriceTrackerData is kinda memory DB for storing previous results and prev used event id
type PriceTrackerData struct {
	Bridge *CommonBridge

	Submits   map[string]*contracts.BridgeTransferSubmit
	Unlocks   map[string]*contracts.BridgeTransferFinish
	Transfers map[string]*contracts.BridgeTransfer

	PrevUsedUnlockEventId   *big.Int // id when "GasPerWithdraw" was called last time (by the user from api)
	LastCaughtUnlockEventId *big.Int // id that was set by the watcher

	TotalGasCost   *big.Int
	WithdrawsCount int64
}

func (d *PriceTrackerData) Init() error {
	sideBridge := d.Bridge.SideBridge.(networks.BridgeFeeApi)
	// init data if needed
	if d.TotalGasCost != nil {
		return nil
	}

	// get oldest locked event id from side net to get
	oldestLockedEvendId, err := sideBridge.GetOldestLockedEventId()
	if err != nil {
		return fmt.Errorf("get side oldest locked event id: %v", err)
	}
	// there's no unlocked events, so we can't get
	if oldestLockedEvendId.Cmp(big.NewInt(1)) == 0 {
		return nil
	}

	//
	oldestUnlockedEventId := new(big.Int).Sub(oldestLockedEvendId, big.NewInt(1))
	oldestUnlockedStartEvendId := new(big.Int).Sub(oldestLockedEvendId, big.NewInt(eventsForGasCalc+1)) // +1 cuz oldestLockedEventId is for locked event but we need unlocked event
	if oldestUnlockedStartEvendId.Cmp(big.NewInt(0)) < 0 {
		oldestUnlockedStartEvendId = big.NewInt(1) // force set to 1 if it's < 0
	}

	var eventIds []*big.Int
	for i := oldestUnlockedStartEvendId; i.Cmp(oldestUnlockedEventId) <= 0; i = new(big.Int).Add(i, big.NewInt(1)) {
		eventIds = append(eventIds, i)
	}

	submits, err := sideBridge.GetTransferSubmitsByIds(eventIds)
	if err != nil {
		return fmt.Errorf("get transfer submits by ids: %v", err)
	}
	unlocks, err := sideBridge.GetTransferUnlocksByIds(eventIds)
	if err != nil {
		return fmt.Errorf("get transfer unlocks by ids: %v", err)
	}
	transfers, err := d.Bridge.GetEventsByIds(eventIds)
	if err != nil {
		return fmt.Errorf("get transfers by ids: %v", err)
	}

	d.Submits = make(map[string]*contracts.BridgeTransferSubmit)
	d.Unlocks = make(map[string]*contracts.BridgeTransferFinish)
	d.Transfers = make(map[string]*contracts.BridgeTransfer)
	for i, eventId := range eventIds {
		d.Submits[eventId.String()] = submits[i]
		d.Unlocks[eventId.String()] = unlocks[i]
		d.Transfers[eventId.String()] = transfers[i]
	}

	d.PrevUsedUnlockEventId = oldestUnlockedStartEvendId.Sub(oldestUnlockedStartEvendId, big.NewInt(1)) //
	d.LastCaughtUnlockEventId = oldestUnlockedEventId
	d.TotalGasCost = big.NewInt(0)
	return nil
}

func (d *PriceTrackerData) processEvent(event *contracts.BridgeTransferFinish) error {
	submit, err := d.Bridge.GetTransferSubmitById(event.EventId)
	if err != nil {
		return fmt.Errorf("get transfer submit event by id: %w", err)
	}
	transfer, err := d.Bridge.SideBridge.GetEventById(event.EventId)
	if err != nil {
		return fmt.Errorf("get transfer event by id: %w", err)
	}

	d.Unlocks[event.EventId.String()] = event
	d.Submits[event.EventId.String()] = submit
	d.Transfers[event.EventId.String()] = transfer
	d.LastCaughtUnlockEventId = event.EventId
	return nil
}

func (d *PriceTrackerData) save(unlockEventId *big.Int, totalGasCost *big.Int, withdrawsCount int64) {
	d.PrevUsedUnlockEventId = unlockEventId

	d.TotalGasCost.Add(d.TotalGasCost, totalGasCost)
	d.WithdrawsCount += withdrawsCount
}

// ------------------------------------------------------------------------

func (b *CommonBridge) WatchUnlocksLoop(sideData *PriceTrackerData) {
	b.shouldHavePk()
	for {
		b.EnsureContractUnpaused()

		if err := b.watchUnlocks(sideData); err != nil {
			b.Logger.Error().Err(err).Msg("price tracker watchUnlocks error")
		}
		time.Sleep(failSleepTIme)
	}
}

func (b *CommonBridge) checkOldUnlocks(sideData *PriceTrackerData) error {
	b.Logger.Info().Msg("Checking old unlock events...")

	for i := int64(1); ; i++ {
		nextEventId := new(big.Int).Add(sideData.LastCaughtUnlockEventId, big.NewInt(i))
		nextEvent, err := b.GetTransferUnlockById(nextEventId)
		if errors.Is(err, networks.ErrEventNotFound) { // no more old events
			return nil
		} else if err != nil {
			return fmt.Errorf("GetTransferUnlockById on id %v: %w", nextEventId.String(), err)
		}

		b.Logger.Info().Str("event_id", nextEventId.String()).Msg("Send old unlock event...")
		if err := sideData.processEvent(nextEvent); err != nil {
			return err
		}
	}
}

func (b *CommonBridge) watchUnlocks(sideData *PriceTrackerData) error {
	if err := b.checkOldUnlocks(sideData); err != nil {
		return fmt.Errorf("checkOldUnlocks: %w", err)
	}

	eventCh := make(chan *contracts.BridgeTransferFinish)
	eventSub, err := b.WsContract.WatchTransferFinish(nil, eventCh, nil)
	if err != nil {
		return fmt.Errorf("watchTransferFinish: %w", err)
	}
	defer eventSub.Unsubscribe()

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching unlock transfers: %w", err)
		case event := <-eventCh:
			if event.Raw.Removed {
				continue
			}

			b.Logger.Info().Str("event_id", event.EventId.String()).Msg("Found new TransferFinish event")
			if err := sideData.processEvent(event); err != nil {
				return err
			}
		}
	}
}

func (b *CommonBridge) GasPerWithdraw(data *PriceTrackerData) (*big.Int, error) {
	b.GasPerWithdrawLock.Lock()
	defer b.GasPerWithdrawLock.Unlock()

	// init data if needed
	if err := data.Init(); err != nil {
		return nil, fmt.Errorf("init price tracker data: %w", err)
	}

	// if we didn't find the pair, then we don't need to calculate anything, just return previous results from data
	if data.PrevUsedUnlockEventId.Cmp(data.LastCaughtUnlockEventId) != 0 {
		var submits []*contracts.BridgeTransferSubmit
		var unlocks []*contracts.BridgeTransferFinish
		var transfers []*contracts.BridgeTransfer
		start := new(big.Int).Add(data.PrevUsedUnlockEventId, big.NewInt(1))
		for i := start; i.Cmp(data.LastCaughtUnlockEventId) <= 0; i = new(big.Int).Add(i, big.NewInt(1)) {
			submits = append(submits, data.Submits[i.String()])
			unlocks = append(unlocks, data.Unlocks[i.String()])
			transfers = append(transfers, data.Transfers[i.String()])
		}

		// get total gas cost
		totalGasCost, totalGas, err := b.SideBridge.(networks.BridgeFeeApi).UsedGas(submits, unlocks)
		if err != nil {
			return nil, fmt.Errorf("get used gas: %w", err)
		}

		// get withdraws count
		withdrawsCount := b.withdrawCount(transfers)

		b.Logger.Info().Msgf("withdraws count: %d, totalGas: %d, totalGasCost: %d", withdrawsCount, totalGas, totalGasCost)
		data.save(data.Unlocks[data.LastCaughtUnlockEventId.String()].EventId, totalGasCost, withdrawsCount)
	}

	// if there's no transfers then return nil
	if data.WithdrawsCount == 0 {
		return nil, nil
	}

	b.Logger.Info().Msgf("withdraws count: %d, totalGasCost: %d", data.WithdrawsCount, data.TotalGasCost)
	return new(big.Int).Div(data.TotalGasCost, big.NewInt(data.WithdrawsCount)), nil
}

func (b *CommonBridge) withdrawCount(transfers []*contracts.BridgeTransfer) int64 {
	count := 0

	for _, transfer := range transfers {
		count += len(transfer.Queue)
	}

	return int64(count)
}

func (b *CommonBridge) UsedGas(logsSubmit []*contracts.BridgeTransferSubmit, logsUnlock []*contracts.BridgeTransferFinish) (*big.Int, *big.Int, error) {
	// collect unique transaction hashes

	txs := map[common.Hash]interface{}{} // use as hashset, coz 1 tx can emit many events

	for _, log := range logsSubmit {
		txs[log.Raw.TxHash] = 0
	}
	for _, log := range logsUnlock {
		txs[log.Raw.TxHash] = 0
	}

	// fetch transactions and sum gas

	totalGas := new(big.Int)
	totalGasCost := new(big.Int)

	lock := sync.Mutex{}
	sem := make(chan interface{}, 10) // max 20 simultaneous requests
	errGroup := new(errgroup.Group)
	for txHash := range txs {
		txHash := txHash
		errGroup.Go(func() error {
			sem <- 0
			defer func() { <-sem }()

			tx, _, err := b.Client.TransactionByHash(context.Background(), txHash)
			if err != nil {
				return fmt.Errorf("get transaction by hash: %w", err)
			}

			lock.Lock()
			totalGas.Add(totalGas, big.NewInt(int64(tx.Gas())))
			totalGasCost.Add(totalGasCost, tx.Cost())
			lock.Unlock()
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		return nil, nil, fmt.Errorf("calc used gas: %w", err)
	}

	return totalGasCost, totalGas, nil
}

// ---------getters---------

func (b *CommonBridge) GetOldestLockedEventId() (*big.Int, error) {
	return b.Contract.OldestLockedEventId(nil)
}

func (b *CommonBridge) GetTransferSubmitById(eventId *big.Int) (*contracts.BridgeTransferSubmit, error) {
	logSubmit, err := b.Contract.FilterTransferSubmit(nil, []*big.Int{eventId})
	if err != nil {
		return nil, fmt.Errorf("filter transfer submit: %w", err)
	}
	for logSubmit.Next() {
		if !logSubmit.Event.Raw.Removed {
			return logSubmit.Event, nil
		}
	}

	return nil, networks.ErrTransferSubmitNotFound
}

func (b *CommonBridge) GetTransferUnlockById(eventId *big.Int) (*contracts.BridgeTransferFinish, error) {
	logUnlock, err := b.Contract.FilterTransferFinish(nil, []*big.Int{eventId})
	if err != nil {
		return nil, fmt.Errorf("filter transfer finish: %w", err)
	}
	for logUnlock.Next() {
		if !logUnlock.Event.Raw.Removed {
			return logUnlock.Event, nil
		}
	}

	return nil, networks.ErrTransferFinishNotFound
}

func (b *CommonBridge) GetTransferSubmitsByIds(eventIds []*big.Int) (submits []*contracts.BridgeTransferSubmit, err error) {
	logSubmit, err := b.Contract.FilterTransferSubmit(nil, eventIds)
	if err != nil {
		return nil, fmt.Errorf("filter transfer submit: %w", err)
	}

	for logSubmit.Next() {
		if !logSubmit.Event.Raw.Removed {
			submits = append(submits, logSubmit.Event)
		}
	}
	return submits, nil
}
func (b *CommonBridge) GetTransferUnlocksByIds(eventIds []*big.Int) (unlocks []*contracts.BridgeTransferFinish, err error) {
	logUnlock, err := b.Contract.FilterTransferFinish(nil, eventIds)
	if err != nil {
		return nil, fmt.Errorf("filter transfer finish: %w", err)
	}

	for logUnlock.Next() {
		if !logUnlock.Event.Raw.Removed {
			unlocks = append(unlocks, logUnlock.Event)
		}
	}
	return unlocks, nil
}
