package common

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"golang.org/x/sync/errgroup"
)

func (b *CommonBridge) EnsureContractUnpaused() {
	for {
		err := b.waitForUnpauseContract()
		if err == nil {
			return
		}

		b.Logger.Error().Err(err).Msg("waitForUnpauseContract error")
		time.Sleep(10 * time.Second)
	}
}

func (b *CommonBridge) waitForUnpauseContract() error {
	paused, err := b.Contract.Paused(nil)
	if err != nil {
		return fmt.Errorf("Paused: %w", err)
	}
	if !paused {
		return nil
	}

	eventCh := make(chan *contracts.BridgeUnpaused)
	eventSub, err := b.WsContract.WatchUnpaused(nil, eventCh)
	if err != nil {
		return fmt.Errorf("WatchUnpaused: %w", err)
	}
	defer eventSub.Unsubscribe()

	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching unpaused event: %w", err)
		case <-eventCh:
			b.Logger.Info().Msg("Contracts is unpaused, continue working!")
			return nil
		}
	}
}

func (b *CommonBridge) WaitForBlock(targetBlockNum uint64) error {
	// todo maybe timeout (context)
	blockChannel := make(chan *types.Header)
	blockSub, err := b.WsClient.SubscribeNewHead(context.Background(), blockChannel)
	if err != nil {
		return fmt.Errorf("SubscribeNewHead: %w", err)
	}
	defer blockSub.Unsubscribe()

	currentBlockNum, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("get last block num: %w", err)
	}

	for currentBlockNum < targetBlockNum {
		select {
		case err := <-blockSub.Err():
			return fmt.Errorf("listening new blocks: %w", err)

		case block := <-blockChannel:
			currentBlockNum = block.Number.Uint64()
		}
	}

	return nil
}

func (b *CommonBridge) GetReceipts(blockHash common.Hash) ([]*types.Receipt, error) {
	txsCount, err := b.Client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		return nil, fmt.Errorf("get transaction count: %w", err)
	}

	receipts := make([]*types.Receipt, txsCount)

	errGroup := new(errgroup.Group)
	for i := uint(0); i < txsCount; i++ {
		i := i // https://golang.org/doc/faq#closures_and_goroutines ¯\_(ツ)_/¯
		errGroup.Go(func() error {
			tx, err := b.Client.TransactionInBlock(context.Background(), blockHash, i)
			if err != nil {
				return fmt.Errorf("get transaction in block: %w", err)
			}
			receipt, err := b.Client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				return fmt.Errorf("get transaction receipt: %w", err)
			}

			receipts[i] = receipt
			return nil
		})
	}

	return receipts, errGroup.Wait()
}

func (b *CommonBridge) GetFailureReason(tx *types.Transaction) error {
	_, err := b.Client.CallContract(context.Background(), ethereum.CallMsg{
		From:     b.Auth.From,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}, nil)

	return err
}

func (b *CommonBridge) SetUsedGasMetric(usedGas uint64) {
	metric.UsedGas.WithLabelValues(b.Name).Observe(float64(usedGas))
}

func (b *CommonBridge) SetRelayBalanceMetric() {
	balance, err := b.getBalanceGWei(b.Auth.From)
	if err != nil {
		b.Logger.Error().Err(err).Msg("error when getting contract balance in GWei")
		return
	}

	metric.RelayBalance.WithLabelValues(b.Name).Set(balance)
}

func (b *CommonBridge) getBalanceGWei(address common.Address) (float64, error) {
	balance, err := b.Client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return 0, err
	}
	balanceGWei := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(params.GWei))
	balanceFloat64, _ := balanceGWei.Float64()
	return balanceFloat64, nil
}

func parsePK(pk string) (*ecdsa.PrivateKey, error) {
	b, err := hex.DecodeString(pk)
	if err != nil {
		return nil, err
	}
	return crypto.ToECDSA(b)
}
