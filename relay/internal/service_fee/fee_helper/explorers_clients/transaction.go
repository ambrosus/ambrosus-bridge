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

type Transactions []*Transaction

func RemoveTransactionsDups(m Transactions) Transactions {
	keys := make(map[string]bool)
	list := Transactions{}
	for _, entry := range m {
		if _, ok := keys[entry.Hash]; !ok {
			keys[entry.Hash] = true
			list = append(list, entry)
		}
	}
	return list
}

func FilterTxsByFromToAddresses(txs Transactions, from string, to string) Transactions {
	var res Transactions
	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		if tx.From == from && tx.To == to {
			res = append(res, tx)
		}
	}
	return res
}
