package common

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/helpers"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

const (
	eventsForGasCalc = 5
)

// PriceTrackerData is kinda memory DB for storing previous results and prev used event id
type PriceTrackerData struct {
	PrevUsedBlockNumber     uint64
	PrevSideUsedBlockNumber uint64
	TotalGasCost            *big.Int
	WithdrawsCount          int64
}

func (d *PriceTrackerData) save(blockNumber, sideBlockNumber uint64, totalGasCost *big.Int, withdrawsCount int64) {
	d.PrevUsedBlockNumber = blockNumber
	d.PrevSideUsedBlockNumber = sideBlockNumber
	d.TotalGasCost.Add(d.TotalGasCost, totalGasCost)
	d.WithdrawsCount += withdrawsCount
}

func (b *CommonBridge) initPriceTrackerData(data *PriceTrackerData) error {
	// init data if needed
	if data.TotalGasCost != nil {
		return nil
	}

	// get oldest locked event id from side net to get start block number
	oldestLockedEvendId, err := b.SideBridge.(networks.BridgeFeeApi).GetOldestLockedEventId()
	if err != nil {
		return fmt.Errorf("get side oldest locked event id: %v", err)
	}
	// there's no unlocked events, so we can't get start block number
	if oldestLockedEvendId.Cmp(big.NewInt(1)) == 0 {
		return nil
	}

	//
	oldestLockedEvendId.Sub(oldestLockedEvendId, big.NewInt(eventsForGasCalc+1)) // +1 cuz oldestLockedEventId is for locked event but we need unlocked event
	if oldestLockedEvendId.Cmp(big.NewInt(0)) < 0 {
		oldestLockedEvendId = big.NewInt(1) // force set to 1 if it's < 0
	}

	// get start block number for this net
	startBlockNumber, err := b.getStartBlockNumber(oldestLockedEvendId)
	if err != nil {
		return fmt.Errorf("get start block number: %v", err)
	}

	// get start block number for side net
	sideStartBlockNumber, err := b.getSideStartBlockNumber(oldestLockedEvendId)
	if err != nil {
		return fmt.Errorf("get side start block number: %v", err)
	}

	data.PrevUsedBlockNumber = startBlockNumber
	data.PrevSideUsedBlockNumber = sideStartBlockNumber
	data.TotalGasCost = big.NewInt(0)
	return nil
}

func (b *CommonBridge) GasPerWithdraw(data *PriceTrackerData) (*big.Int, error) {
	b.GasPerWithdrawLock.Lock()
	defer b.GasPerWithdrawLock.Unlock()

	// init data if needed
	if err := b.initPriceTrackerData(data); err != nil {
		return nil, fmt.Errorf("init price tracker data: %w", err)
	}

	// get the latest block number from side net
	endSide, err := b.SideBridge.(networks.BridgeFeeApi).GetLatestBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("get side latest block number: %v", err)
	}

	// get submits and unlocks for `UsedGas`
	eventUnlock, submits, unlocks, err := b.SideBridge.(networks.BridgeFeeApi).GetLastCorrectSubmitUnlockPair(
		data.PrevSideUsedBlockNumber,
		endSide,
	)
	if err != nil {
		return nil, fmt.Errorf("get last correct submit unlock pair: %w", err)
	}

	// if we didn't find the pair, then we don't need to calculate anything, just return previous results from data
	if eventUnlock != nil {
		// get total gas cost
		totalGasCost, totalGas, err := b.SideBridge.(networks.BridgeFeeApi).UsedGas(submits, unlocks)
		if err != nil {
			return nil, fmt.Errorf("get used gas: %w", err)
		}

		// get withdraws count, setting the `endBlockNumber` to the transfer event's block number
		eventTransfer, err := b.GetEventById(eventUnlock.EventId)
		if err != nil {
			return nil, fmt.Errorf("get event by id: %w", err)
		}
		withdrawsCount, err := b.withdrawCount(data.PrevUsedBlockNumber, eventTransfer.Raw.BlockNumber)
		if err != nil {
			return nil, fmt.Errorf("get withdraw count: %w", err)
		}

		b.Logger.Info().Msgf("withdraws count: %d, totalGas: %d, totalGasCost: %d", withdrawsCount, totalGas, totalGasCost)
		data.save(eventTransfer.Raw.BlockNumber, eventUnlock.Raw.BlockNumber, totalGasCost, withdrawsCount)
	}

	// if there's no transfers then return nil
	if data.WithdrawsCount == 0 {
		return nil, nil
	}

	return new(big.Int).Div(data.TotalGasCost, big.NewInt(data.WithdrawsCount)), nil
}

func (b *CommonBridge) GetLastCorrectSubmitUnlockPair(startBlockNumber, endBlockNumber uint64) (
	event *contracts.BridgeTransferFinish,
	submits []*contracts.BridgeTransferSubmit,
	unlocks []*contracts.BridgeTransferFinish,
	err error,
) {
	opts := &bind.FilterOpts{Start: startBlockNumber, End: &endBlockNumber}

	// get submit events
	logsSubmit, err := b.Contract.FilterTransferSubmit(opts, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("filter transfer submit: %w", err)
	}
	for logsSubmit.Next() {
		submits = append(submits, logsSubmit.Event)
	}

	// get unlock events
	logsUnlock, err := b.Contract.FilterTransferFinish(opts, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("filter transfer finish: %w", err)
	}
	for logsUnlock.Next() {
		unlocks = append(unlocks, logsUnlock.Event)
	}

	// intersect submits and unlocks
	submits, unlocks = helpers.IntersectionSubmitsUnlocks(submits, unlocks)
	if len(unlocks) != 0 {
		event = unlocks[len(unlocks)-1]
	}

	return event, submits, unlocks, nil
}

func (b *CommonBridge) withdrawCount(startBlockNumber, endBlockNumber uint64) (int64, error) {
	count := 0

	opts := &bind.FilterOpts{Start: startBlockNumber, End: &endBlockNumber}
	logs, err := b.Contract.FilterTransfer(opts, nil)
	if err != nil {
		return 0, fmt.Errorf("filter transfer: %w", err)
	}
	for logs.Next() {
		count += len(logs.Event.Queue)
	}

	return int64(count), nil
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

func (b *CommonBridge) getStartBlockNumber(eventId *big.Int) (uint64, error) {
	event, err := b.GetEventById(eventId)
	if err != nil {
		return 0, fmt.Errorf("get event by id: %w", err)
	}
	return event.Raw.BlockNumber, nil
}

func (b *CommonBridge) getSideStartBlockNumber(eventId *big.Int) (uint64, error) {
	event, err := b.SideBridge.(networks.BridgeFeeApi).GetTransferSubmitById(eventId)
	if err != nil {
		return 0, fmt.Errorf("get transfer submit by id: %w", err)
	}

	return event.Raw.BlockNumber, nil
}

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
