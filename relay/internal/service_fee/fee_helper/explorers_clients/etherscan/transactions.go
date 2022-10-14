package etherscan

import (
	"errors"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper/explorers_clients"
	"github.com/nanmu42/etherscan-api"
)

const (
	maxTxListResponse = 10_000
)

// errors
var (
	ErrTxsNotFound = errors.New("etherscan server: No transactions found")
)

// That method wraps etherscan's `NormalTxByAddress` but returns our errors
func (e *Etherscan) normalTxByAddress(address string, startBlock *int, endBlock *int, page int, offset int, desc bool) (txs []etherscan.NormalTx, err error) {
	txs, err = e.client.NormalTxByAddress(address, startBlock, endBlock, page, offset, desc)
	if err != nil {
		if err.Error() == ErrTxsNotFound.Error() {
			return nil, explorers_clients.ErrTxsNotFound
		}
	}
	return
}

func (e *Etherscan) TxListByAddress(address string, txFilters explorers_clients.TxFilters) ([]*explorers_clients.Transaction, error) {
	var txs []*explorers_clients.Transaction

	var startBlock int = int(txFilters.FromBlock)
	for {
		pageTxs, err := e.normalTxByAddress(address, &startBlock, nil, 0, 0, true)
		if err != nil {
			return nil, err
		}
		startBlock = pageTxs[len(pageTxs)-1].BlockNumber

		ourTypeTx := toOurTxType(pageTxs)
		txsUntilTxHash, isReachedTheTxHash := explorers_clients.TakeTxsUntilTxHash(ourTypeTx, txFilters.UntilTxHash)
		txs = append(txs, txsUntilTxHash...)

		if len(pageTxs) != maxTxListResponse || isReachedTheTxHash {
			break
		}
	}

	txsWithoutDups := explorers_clients.RemoveTransactionsDups(txs)
	return txsWithoutDups, nil
}

func (e *Etherscan) TxListByFromToAddresses(from, to string, txFilters explorers_clients.TxFilters) ([]*explorers_clients.Transaction, error) {
	from, to = explorers_clients.ToLower(from), explorers_clients.ToLower(to)
	txs, err := e.TxListByAddress(from, txFilters)
	if err != nil {
		return nil, err
	}

	res := explorers_clients.FilterTxsByFromToAddresses(txs, from, to)
	return res, nil
}

func (e *Etherscan) TxListByFromListToAddresses(fromList []string, to string, txFilters explorers_clients.TxFilters) ([]*explorers_clients.Transaction, error) {
	fromList, to = explorers_clients.ToLower(fromList), explorers_clients.ToLower(to)
	txs, err := e.TxListByAddress(to, txFilters)
	if err != nil {
		return nil, err
	}

	res := explorers_clients.FilterTxsByFromToAddresses(txs, fromList, to)
	return res, nil
}

func toOurTxType(txs []etherscan.NormalTx) []*explorers_clients.Transaction {
	var mappedTxs []*explorers_clients.Transaction

	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		mappedTxs = append(mappedTxs, &explorers_clients.Transaction{
			BlockNumber: uint64(tx.BlockNumber),
			Hash:        tx.Hash,
			From:        tx.From,
			To:          tx.To,
			GasPrice:    tx.GasPrice.Int(),
			GasUsed:     uint64(tx.GasUsed),
			Input:       tx.Input,
		})
	}
	return mappedTxs
}
