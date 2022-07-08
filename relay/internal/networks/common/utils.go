package common

import (
	"context"
	"fmt"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/sync/errgroup"
)

// failSleepTIme is how many seconds to sleep between iterations in infinity loops
const failSleepTIme = time.Second * 30

func EncodeTransferProof(client ethclients.ClientInterface, event *bindings.BridgeTransfer) (*bindings.CommonStructsTransferProof, error) {
	proof, err := GetProof(client, event)
	if err != nil {
		return nil, fmt.Errorf("GetProof: %w", err)
	}

	return &bindings.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func GetProof(client ethclients.ClientInterface, event receipts_proof.ProofEvent) ([][]byte, error) {
	receipts, err := getReceipts(client, event.Log().BlockHash)
	if err != nil {
		return nil, fmt.Errorf("getReceipts: %w", err)
	}
	return receipts_proof.CalcProofEvent(receipts, event)
}

func getReceipts(client ethclients.ClientInterface, blockHash common.Hash) ([]*types.Receipt, error) {
	txsCount, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		return nil, fmt.Errorf("get transaction count: %w", err)
	}

	receipts := make([]*types.Receipt, txsCount)

	errGroup := new(errgroup.Group)
	for i := uint(0); i < txsCount; i++ {
		i := i // https://golang.org/doc/faq#closures_and_goroutines ¯\_(ツ)_/¯
		errGroup.Go(func() error {
			tx, err := client.TransactionInBlock(context.Background(), blockHash, i)
			if err != nil {
				return fmt.Errorf("get transaction in block: %w", err)
			}
			receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				return fmt.Errorf("get transaction receipt: %w", err)
			}

			receipts[i] = receipt
			return nil
		})
	}

	return receipts, errGroup.Wait()
}

func (b *CommonBridge) waitForUnpauseContract() error {
	paused, err := b.Contract.Paused(nil)
	if err != nil {
		return fmt.Errorf("Paused: %w", err)
	}
	if !paused {
		return nil
	}

	eventCh := make(chan *bindings.BridgeUnpaused)
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
