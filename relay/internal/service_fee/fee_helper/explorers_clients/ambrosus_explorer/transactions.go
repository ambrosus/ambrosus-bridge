package ambrosus_explorer

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper/explorers_clients"
	"github.com/ethereum/go-ethereum/common"
)

const (
	maxTxListResponse        = 5_000
	txByAddressRequestFormat = "accounts/%s/transactions"
)

// errors
var (
	ErrTxsNotFound = errors.New("No results found")
)

type txListByAddressResponse struct {
	Data       []txResponse `json:"data"`
	Meta       meta         `json:"meta"`
	Pagination pagination   `json:"pagination"`
}
type txResponse struct {
	BlockNumber uint64   `json:"blockNumber"`
	Hash        string   `json:"hash"`
	From        string   `json:"from"`
	To          string   `json:"to"`
	GasPrice    *big.Int `json:"gasPrice"`
	GasUsed     uint64   `json:"gasUsed"`
	Input       string   `json:"input"`
}

type meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"` // when request is unsuccessful
}
type pagination struct {
	HasNext     bool `json:"hasNext"`
	HasPrevious bool `json:"hasPrevious"`
	Previous    int  `json:"previous"`
}

func (e *AmbrosusExplorer) txListByAddressWithMeta(address string, limit int, page int) (*txListByAddressResponse, error) {
	u := e.apiUrl.JoinPath(fmt.Sprintf(txByAddressRequestFormat, address))
	q := u.Query()
	q.Set("limit", fmt.Sprint(limit))
	q.Set("page", fmt.Sprint(page))
	u.RawQuery = q.Encode()

	resp, err := e.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var respData txListByAddressResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if respData.Meta.Message == ErrTxsNotFound.Error() {
		return nil, explorers_clients.ErrTxsNotFound
	}
	if respData.Meta.Message != "" {
		return nil, fmt.Errorf(respData.Meta.Message)
	}

	return &respData, nil
}

func (e *AmbrosusExplorer) TxListByAddress(address string, untilTxHash *common.Hash) ([]*explorers_clients.Transaction, error) {
	var txs []*explorers_clients.Transaction

	var currentPage = 1
	for {
		pageTxs, err := e.txListByAddressWithMeta(address, maxTxListResponse, currentPage)
		if err != nil {
			return nil, err
		}
		currentPage++

		ourTypeTx := toOurTxType(pageTxs.Data)
		txsUntilTxHash, isReachedTheTxHash := explorers_clients.TakeTxsUntilTxHash(ourTypeTx, untilTxHash)
		txs = append(txs, txsUntilTxHash...)

		// TODO: when "limit" is big enough (>1000), "hasNext" is false, but there's the next page, should check the next page like etherscan client
		if !pageTxs.Pagination.HasNext || isReachedTheTxHash {
			break
		}
	}

	txsWithoutDups := explorers_clients.RemoveTransactionsDups(txs)
	return txsWithoutDups, nil

}

func (e *AmbrosusExplorer) TxListByFromToAddresses(from, to string, untilTxHash *common.Hash) ([]*explorers_clients.Transaction, error) {
	txs, err := e.TxListByAddress(from, untilTxHash)
	if err != nil {
		return nil, err
	}

	res := explorers_clients.FilterTxsByFromToAddresses(txs, from, to)
	return res, nil
}

func toOurTxType(txs []txResponse) []*explorers_clients.Transaction {
	var mappedTxs []*explorers_clients.Transaction

	for i := 0; i < len(txs); i++ {
		tx := txs[i]
		mappedTxs = append(mappedTxs, &explorers_clients.Transaction{
			BlockNumber: tx.BlockNumber,
			Hash:        tx.Hash,
			From:        tx.From,
			To:          tx.To,
			GasPrice:    tx.GasPrice,
			GasUsed:     tx.GasUsed,
			Input:       tx.Input,
		})
	}
	return mappedTxs
}