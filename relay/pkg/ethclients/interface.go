package ethclients

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

type ClientInterface interface {
	ethereum.ChainReader
	ethereum.TransactionReader
	ethereum.ChainStateReader
	ethereum.ChainSyncReader
	ethereum.ContractCaller
	ethereum.GasEstimator
	ethereum.GasPricer
	ethereum.LogFilterer
	ethereum.PendingStateReader
	ethereum.PendingContractCaller
	bind.ContractBackend
	ChainID(ctx context.Context) (*big.Int, error)
	Close()
	BlockNumber(ctx context.Context) (uint64, error)
	TxGasPriceFromResponse(ctx context.Context, txHash common.Hash) (*big.Int, error)
}
