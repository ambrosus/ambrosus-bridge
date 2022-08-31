package etherscan

import (
	"errors"
	"strings"

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

func (e *Etherscan) TxListByAddress(address string) (explorers_clients.Transactions, error) {
	var txs []etherscan.NormalTx

	firstPageTxs, err := e.normalTxByAddress(address, nil, nil, 0, 0, true)
	if err != nil {
		return nil, err
	}
	txs = append(txs, firstPageTxs...)

	for {
		nextPageTxs, err := e.normalTxByAddress(address, &txs[len(txs)-1].BlockNumber, nil, 0, 0, true)
		if err != nil {
			return nil, err
		}
		txs = append(txs, nextPageTxs...)

		if len(nextPageTxs) != maxTxListResponse {
			break
		}
	}

	mappedTxs := mapTxs(txs)
	txsWithoutDups := explorers_clients.RemoveTransactionsDups(mappedTxs)

	return txsWithoutDups, nil
}

func (e *Etherscan) TxListByAddressUntilTxHash(address string, untilTxHash string) (explorers_clients.Transactions, error) {
	var txs []etherscan.NormalTx

	firstPageTxs, err := e.normalTxByAddress(address, nil, nil, 0, 0, true)
	if err != nil {
		return nil, err
	}
	txsUntilTxHash, isReachedTheTxHash := getTxsUntilTxHash(firstPageTxs, untilTxHash)
	txs = append(txs, txsUntilTxHash...)

	if !isReachedTheTxHash {
		var nextPageTxs = firstPageTxs
		for {
			nextPageTxs, err = e.normalTxByAddress(address, &txs[len(txs)-1].BlockNumber, nil, 0, 0, true)
			if err != nil {
				return nil, err
			}
			txsUntilTxHash, isReachedTheTxHash := getTxsUntilTxHash(nextPageTxs, untilTxHash)
			txs = append(txs, txsUntilTxHash...)

			if isReachedTheTxHash {
				break
			}
			if len(nextPageTxs) != maxTxListResponse {
				return nil, nil
			}
		}
	}

	mappedTxs := mapTxs(txs)
	txsWithoutDups := explorers_clients.RemoveTransactionsDups(mappedTxs)

	return txsWithoutDups, nil
}

func getTxsUntilTxHash(txs []etherscan.NormalTx, untilTxHash string) (res []etherscan.NormalTx, isReachedTheTxHash bool) {
	for _, tx := range txs {
		if tx.Hash == untilTxHash {
			return res, true
		}
		res = append(res, tx)
	}

	return res, false
}

func (e *Etherscan) TxListByFromToAddresses(from, to string) (explorers_clients.Transactions, error) {
	from, to = strings.ToLower(from), strings.ToLower(to)
	txs, err := e.TxListByAddress(from)
	if err != nil {
		return nil, err
	}

	res := explorers_clients.FilterTxsByFromToAddresses(txs, from, to)
	return res, nil
}

func (e *Etherscan) TxListByFromToAddressesUntilTxHash(from, to string, untilTxHash string) (explorers_clients.Transactions, error) {
	from, to = strings.ToLower(from), strings.ToLower(to)
	txs, err := e.TxListByAddressUntilTxHash(from, untilTxHash)
	if err != nil {
		return nil, err
	}

	res := explorers_clients.FilterTxsByFromToAddresses(txs, from, to)
	return res, nil
}

func mapTxs(txs []etherscan.NormalTx) explorers_clients.Transactions {
	var mappedTxs explorers_clients.Transactions

	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		mappedTxs = append(mappedTxs, &explorers_clients.Transaction{
			BlockNumber: uint64(tx.BlockNumber),
			Hash:        tx.Hash,
			From:        tx.From,
			To:          tx.To,
			GasPrice:    tx.GasPrice.Int(),
			GasUsed:     uint64(tx.GasUsed),
			Value:       tx.Value.Int(),
		})
	}
	return mappedTxs
}
