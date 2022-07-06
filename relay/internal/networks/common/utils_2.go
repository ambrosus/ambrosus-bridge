package common

import (
	"context"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/sync/errgroup"
)

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
