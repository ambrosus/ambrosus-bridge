package explorers_clients

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/slices"
)

type TxFilters struct {
	FromBlock   uint64
	UntilTxHash *common.Hash
}

type Transaction struct {
	BlockNumber uint64
	Hash        string
	From        string
	To          string
	GasPrice    *big.Int
	GasUsed     uint64
	Input       string
}

func FilterTxsByFromBlock(txs []*Transaction, fromBlock uint64) []*Transaction {
	return FilterTxsByCallback(txs, func(tx *Transaction) bool { return tx.BlockNumber >= fromBlock })
}

func FilterTxsByFromToAddresses[T string | []string](txs []*Transaction, from T, to string) []*Transaction {
	var filterCallback func(tx *Transaction) bool
	switch from := any(from).(type) {
	case string:
		filterCallback = func(tx *Transaction) bool { return tx.From == from && tx.To == to }
	case []string:
		filterCallback = func(tx *Transaction) bool { return slices.Contains(from, tx.From) && tx.To == to }
	}

	return FilterTxsByCallback(txs, filterCallback)
}

func FilterTxsByCallback(txs []*Transaction, filterCallback func(tx *Transaction) bool) []*Transaction {
	res := []*Transaction{}

	for _, tx := range txs {
		if filterCallback(tx) {
			res = append(res, tx)
		}
	}

	return res
}

func TakeTxsUntilTxHash(txs []*Transaction, untilTxHash *common.Hash) (res []*Transaction, isReachedTheTxHash bool) {
	if untilTxHash != nil {
		for i, tx := range txs {
			if tx.Hash == untilTxHash.Hex() {
				return txs[:i], true
			}
		}
	}
	return txs, false
}

// RemoveTransactionsDups keep the order of the txs, unlike `helpers.Unique`
func RemoveTransactionsDups(m []*Transaction) []*Transaction {
	keys := make(map[string]bool)
	var list []*Transaction
	for _, entry := range m {
		if _, ok := keys[entry.Hash]; !ok {
			keys[entry.Hash] = true
			list = append(list, entry)
		}
	}
	return list
}

func ToLower[T string | []string](from T) (res T) {
	switch from := any(from).(type) {
	case string:
		res = any(strings.ToLower(from)).(T)

	case []string:
		var sl []string
		for _, v := range from {
			sl = append(sl, strings.ToLower(v))
		}
		res = any(sl).(T)
	}

	return res
}
