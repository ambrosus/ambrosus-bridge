package common

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/sync/errgroup"
)

// failSleepTIme is how many seconds to sleep between iterations in infinity loops
const failSleepTIme = time.Second * 30

func (b *CommonBridge) EnsureContractUnpaused() {
	for {
		err := b.waitForUnpauseContract()
		if err == nil {
			return
		}

		b.Logger.Error().Err(err).Msg("waitForUnpauseContract error")
		time.Sleep(failSleepTIme)
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
		case event := <-eventCh:
			if event.Raw.Removed {
				continue
			}

			b.Logger.Info().Msg("Contracts is unpaused, continue working!")
			return nil
		}
	}
}

func (b *CommonBridge) WaitForBlock(targetBlockNum uint64) error {
	b.Logger.Debug().Uint64("blockNum", targetBlockNum).Msg("Waiting for block...")

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

func (b *CommonBridge) GetProof(event receipts_proof.ProofEvent) ([][]byte, error) {
	receipts, err := b.GetReceipts(event.Log().BlockHash)
	if err != nil {
		return nil, fmt.Errorf("GetReceipts: %w", err)
	}
	return receipts_proof.CalcProofEvent(receipts, event)
}

func (b *CommonBridge) IsEventRemoved(event *contracts.BridgeTransfer) error {
	b.Logger.Debug().Str("event_id", event.EventId.String()).Msg("Checking if the event has been removed...")

	newEvent, err := b.GetEventById(event.EventId)
	if err != nil {
		return err
	}
	if newEvent.Raw.BlockHash != event.Raw.BlockHash {
		return fmt.Errorf("looks like the event has been removed")
	}
	return nil
}

func (b *CommonBridge) getFailureReason(tx *types.Transaction) error {
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

func (b *CommonBridge) GetLastProcessedBlockNum(currEventId *big.Int) (uint64, error) {
	prevEventId := new(big.Int).Sub(currEventId, big.NewInt(1))
	prevEvent, err := b.GetEventById(prevEventId)
	if err != nil {
		return 0, fmt.Errorf("GetEventById: %w", err)
	}
	if prevEventId.Uint64() == 0 {
		return prevEvent.Raw.BlockNumber, nil
	}

	// todo specify block when prevEvent submitted in side network for 100$ correct `minSafetyBlocks` value
	minSafetyBlocks, err := b.SideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return 0, fmt.Errorf("get block by hash: %w", err)
	}

	return prevEvent.Raw.BlockNumber + minSafetyBlocks, nil
}

func (b *CommonBridge) shouldHavePk() {
	if b.Auth == nil {
		b.Logger.Fatal().Msg("Private key is required")
	}
}
