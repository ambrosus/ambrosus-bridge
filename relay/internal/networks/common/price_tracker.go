package common

import (
	"context"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

func (b *CommonBridge) GasPerWithdraw(afterEventId *big.Int) (float64, error) {
	totalGasCost, totalGas, err := b.usedGas(afterEventId)
	if err != nil {
		return 0, fmt.Errorf("get used gas: %w", err)
	}

	withdrawsCount, err := b.SideBridge.WithdrawCount(afterEventId)
	if err != nil {
		return 0, fmt.Errorf("get withdraw count: %w", err)
	}

	b.Logger.Info().Msgf("withdraws count: %d, totalGas: %d, totalGasCost: %d", withdrawsCount, totalGas, totalGasCost)
	return float64(totalGasCost) / float64(withdrawsCount), nil
}

func (b *CommonBridge) WithdrawCount(afterEventId *big.Int) (int, error) {
	event, err := b.GetEventById(afterEventId)
	if err != nil {
		return 0, fmt.Errorf("failed to get event: %w", err)
	}

	count := 0

	opts := &bind.FilterOpts{Start: event.Raw.BlockNumber}
	logs, err := b.Contract.FilterTransfer(opts, nil)
	if err != nil {
		return 0, fmt.Errorf("filter transfer: %w", err)
	}
	for logs.Next() {
		count += len(logs.Event.Queue)
	}

	return count, nil
}

func (b *CommonBridge) usedGas(afterEventId *big.Int) (uint64, uint64, error) {
	event, err := b.GetEventById(afterEventId)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get event: %w", err)
	}

	// collect unique transaction hashes

	txs := map[common.Hash]interface{}{} // use as hashset, coz 1 tx can emit many events

	opts := &bind.FilterOpts{Start: event.Raw.BlockNumber}
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
