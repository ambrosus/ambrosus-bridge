package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var ErrTimeuot = errors.New("timeout error")

func (b *Bridge) GetTransactionReceipt(txHash common.Hash) (*types.Receipt, error) {
	timeout := 2 * time.Minute
	receipts := make(chan *types.Receipt)

	go func(chan *types.Receipt) {
		for {
			receipt, _ := b.Client.TransactionReceipt(context.Background(), txHash)

			if receipt != nil {
				receipts <- receipt
			}
		}
	}(receipts)

	select {
	case receipt := <-receipts:
		return receipt, nil
	case <-time.After(timeout):
		return nil, ErrTimeuot
	}
}

func (b *Bridge) prepareTransaction(from common.Address, privateKey *ecdsa.PrivateKey, valueInWei *big.Int) (*bind.TransactOpts, error) {
	nonce, err := b.Client.PendingNonceAt(context.Background(), from)
	if err != nil {
		return nil, err
	}

	gasPrice, err := b.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := b.getAuth()
	if err != nil {
		return nil, err
	}

	auth.From = from
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = valueInWei
	auth.GasPrice = gasPrice

	return auth, nil
}
