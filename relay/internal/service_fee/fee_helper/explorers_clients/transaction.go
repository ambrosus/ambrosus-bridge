package explorers_clients

import (
	"errors"
	"math/big"
)

var (
	ErrTxsNotFound = errors.New("transactions not found")
)

type Transaction struct {
	BlockNumber uint64
	Hash        string
	From        string
	To          string
	GasPrice    *big.Int
	GasUsed     uint64
	Value       *big.Int
}

func FilterTxsByFromToAddresses(txs []*Transaction, from string, to string) []*Transaction {
	var res []*Transaction
	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		if tx.From == from && tx.To == to {
			res = append(res, tx)
		}
	}
	return res
}

func TakeTxsUntilTxHash(txs []*Transaction, untilTxHash *string) (res []*Transaction, isReachedTheTxHash bool) {
	if untilTxHash != nil {
		for i, tx := range txs {
			if tx.Hash == *untilTxHash {
				return txs[:i], true
			}
		}
	}
	return txs, false
}
