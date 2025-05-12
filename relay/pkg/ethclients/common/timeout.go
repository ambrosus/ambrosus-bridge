package common

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const timeout = 2 * time.Minute

// Add context with timeout to other methods

func (ec *Client) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.BlockByHash(ctx, hash)
}
func (ec *Client) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.BlockByNumber(ctx, number)
}

func (ec *Client) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.HeaderByHash(ctx, hash)
}
func (ec *Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.HeaderByNumber(ctx, number)
}
func (ec *Client) SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.SubscribeNewHead(ctx, ch)
}

func (ec *Client) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.TransactionCount(ctx, blockHash)
}
func (ec *Client) TransactionInBlock(ctx context.Context, blockHash common.Hash, index uint) (*types.Transaction, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.TransactionInBlock(ctx, blockHash, index)
}
func (ec *Client) TransactionByHash(ctx context.Context, txHash common.Hash) (tx *types.Transaction, isPending bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.TransactionByHash(ctx, txHash)
}
func (ec *Client) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.TransactionReceipt(ctx, txHash)
}

func (ec *Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.BalanceAt(ctx, account, blockNumber)
}
func (ec *Client) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.StorageAt(ctx, account, key, blockNumber)
}
func (ec *Client) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return ec.Client.CodeAt(ctx, account, blockNumber)
}
func (ec *Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.NonceAt(ctx, account, blockNumber)
}

func (ec *Client) SyncProgress(ctx context.Context) (*ethereum.SyncProgress, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.SyncProgress(ctx)
}

func (ec *Client) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.CallContract(ctx, call, blockNumber)
}

func (ec *Client) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.EstimateGas(ctx, call)
}
func (ec *Client) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.SuggestGasPrice(ctx)
}
func (ec *Client) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.SuggestGasTipCap(ctx)
}

func (ec *Client) FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.FilterLogs(ctx, q)
}
func (ec *Client) SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.SubscribeFilterLogs(ctx, q, ch)
}

func (ec *Client) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.SendTransaction(ctx, tx)
}

func (ec *Client) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.PendingBalanceAt(ctx, account)
}
func (ec *Client) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.PendingStorageAt(ctx, account, key)
}
func (ec *Client) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.PendingNonceAt(ctx, account)
}
func (ec *Client) PendingTransactionCount(ctx context.Context) (uint, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.PendingTransactionCount(ctx)
}
func (ec *Client) PendingCodeAt(ctx context.Context, contract common.Address) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.PendingCodeAt(ctx, contract)
}
func (ec *Client) PendingCallContract(ctx context.Context, call ethereum.CallMsg) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.PendingCallContract(ctx, call)
}

func (ec *Client) ChainID(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.ChainID(ctx)
}
func (ec *Client) BlockNumber(ctx context.Context) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return ec.Client.BlockNumber(ctx)
}
