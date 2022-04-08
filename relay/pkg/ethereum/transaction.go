package ethereum

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/errgroup"
)

func GetReceipts(client *ethclient.Client, blockHash common.Hash) ([]*types.Receipt, error) {
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

	err = errGroup.Wait()
	if err != nil {
		return nil, err
	}

	return receipts, nil
}

func GetFailureReason(client *ethclient.Client, auth *bind.TransactOpts, tx *types.Transaction) error {
	_, err := client.CallContract(context.Background(), ethereum.CallMsg{
		From:     auth.From,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}, nil)

	return err
}
