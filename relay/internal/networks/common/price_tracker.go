package common

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
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

func (b *CommonBridge) GasPerWithdraw(data *PriceTrackerData) (*big.Int, error) {
	b.GasPerWithdrawLock.Lock()
	defer b.GasPerWithdrawLock.Unlock()

	// get the latest block numbers from both sides
	end, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get latest block number: %v", err)
	}
	endSide, err := b.SideBridge.(networks.BridgeFeeApi).GetLatestBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("get side latest block number: %v", err)
	}

	//
	if data.PrevSideUsedBlockNumber == 0 {
		data.PrevUsedBlockNumber = end - 10000 // TODO: get block number from the first event за останій тиждень для кожної мережі
		data.PrevSideUsedBlockNumber = endSide - 10000
		data.TotalGasCost = big.NewInt(0)
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

	// if there's no transfers then return default transfer fee
	if data.WithdrawsCount == 0 {
		return b.DefaultTransferFeeWei, nil
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

	// make zip of the events and iterate in reverse order while we'll find the correct pair
	z := helpers.ZipDiff(submits, unlocks)
	for i := len(z) - 1; i >= 0; i-- {
		submitId := new(big.Int)
		unlockId := new(big.Int)

		// needed to compare big.Ints later
		if z[i].First != nil {
			submitId = z[i].First.EventId
		}
		// it's unreal, but let it be
		if z[i].Second != nil {
			unlockId = z[i].Second.EventId
		}

		// if ids are equal, then we found the correct pair
		if submitId.Cmp(unlockId) == 0 {
			event = z[i].Second     // also can be got from `unlocks[-1]`, but "Explicit is better than implicit"
			submits = submits[:i+1] // cut unneeded events
			unlocks = unlocks[:i+1] // cut unneeded events
			break
		}

		// if it's the end, but we didn't find the pair
		if i == 0 {
			submits, unlocks = nil, nil
		}
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
