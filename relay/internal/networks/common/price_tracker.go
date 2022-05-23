package common

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

// PriceTrackerData is kinda memory DB for storing previous results and prev used event id
type PriceTrackerData struct {
	PrevUsedBlockNumber     uint64
	PrevSideUsedBlockNumber uint64
	TotalGasCost            uint64 // TODO: make it as big.Int
	WithdrawsCount          int
}

func (d *PriceTrackerData) save(blockNumber, sideBlockNumber uint64, totalGasCost uint64, withdrawsCount int) {
	d.PrevUsedBlockNumber = blockNumber
	d.PrevSideUsedBlockNumber = sideBlockNumber
	d.TotalGasCost += totalGasCost
	d.WithdrawsCount += withdrawsCount
}

func (b *CommonBridge) GasPerWithdraw(data *PriceTrackerData) (float64, error) {
	b.GasPerWithdrawLock.Lock()
	defer b.GasPerWithdrawLock.Unlock()

	// get the latest block numbers from both sides
	end, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		return 0, fmt.Errorf("get latest block number: %v", err)
	}
	endSide, err := b.SideBridge.(networks.BridgeFeeApi).GetLatestBlockNumber()
	if err != nil {
		return 0, fmt.Errorf("get side latest block number: %v", err)
	}

	//
	if data.PrevSideUsedBlockNumber == 0 {
		data.PrevUsedBlockNumber = end - 10000 // TODO: get block number from the first event за останій тиждень для кожної мережі
		data.PrevSideUsedBlockNumber = endSide - 10000
	}

	// get total gas cost
	totalGasCost, totalGas, err := b.SideBridge.(networks.BridgeFeeApi).UsedGas(data.PrevSideUsedBlockNumber, endSide)
	if err != nil {
		return 0, fmt.Errorf("get used gas: %w", err)
	}

	// get withdraws count
	withdrawsCount, err := b.withdrawCount(data.PrevUsedBlockNumber, end)
	if err != nil {
		return 0, fmt.Errorf("get withdraw count: %w", err)
	}

	b.Logger.Info().Msgf("withdraws count: %d, totalGas: %d, totalGasCost: %d", withdrawsCount, totalGas, totalGasCost)
	data.save(end, endSide, totalGasCost, withdrawsCount)
	return float64(data.TotalGasCost) / float64(data.WithdrawsCount), nil
}

func (b *CommonBridge) withdrawCount(startBlockNumber, endBlockNumber uint64) (int, error) {
	count := 0

	opts := &bind.FilterOpts{Start: startBlockNumber, End: &endBlockNumber}
	logs, err := b.Contract.FilterTransfer(opts, nil)
	if err != nil {
		return 0, fmt.Errorf("filter transfer: %w", err)
	}
	for logs.Next() {
		count += len(logs.Event.Queue)
	}

	return count, nil
}

func (b *CommonBridge) UsedGas(startBlockNumber, endBlockNumber uint64) (uint64, uint64, error) {
	// collect unique transaction hashes

	txs := map[common.Hash]interface{}{} // use as hashset, coz 1 tx can emit many events

	opts := &bind.FilterOpts{Start: startBlockNumber, End: &endBlockNumber}
	logsSubmit, err := b.Contract.FilterTransferSubmit(opts, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("filter transfer submit: %w", err)
	}
	for logsSubmit.Next() {
		txs[logsSubmit.Event.Raw.TxHash] = 0
	}
	logsUnlock, err := b.Contract.FilterTransferFinish(opts, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("filter transfer finish: %w", err)
	}
	for logsUnlock.Next() {
		txs[logsUnlock.Event.Raw.TxHash] = 0
	}

	// fetch transactions and sum gas

	totalGas := uint64(0)
	totalGasCost := uint64(0)

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
			atomic.AddUint64(&totalGas, tx.Gas())
			atomic.AddUint64(&totalGasCost, tx.Cost().Uint64())
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		return 0, 0, fmt.Errorf("calc used gas: %w", err)
	}

	return totalGasCost, totalGas, nil
}
