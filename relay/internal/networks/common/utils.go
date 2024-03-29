package common

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

func ShouldHavePk(b networks.Bridge) {
	if b.GetAuth() == nil {
		b.GetLogger().Fatal().Msg("Private key is required")
	}
}

func EnsureContractUnpaused(b networks.Bridge, logger *zerolog.Logger) {
	for {
		err := waitForUnpauseContract(b)
		if err == nil {
			return
		}

		logger.Error().Err(fmt.Errorf("EnsureContractUnpaused: %w", err)).Msg("")
		time.Sleep(time.Second * 30)
	}
}

func waitForUnpauseContract(b networks.Bridge) error {
	paused, err := b.GetContract().Paused(nil)
	if err != nil {
		return fmt.Errorf("Paused: %w", err)
	}
	if !paused {
		return nil
	}

	for {
		err := b.Events().WatchUnpaused()
		if err != nil {
			return fmt.Errorf("watching unpaused event: %w", err)
		}
		b.GetLogger().Info().Msg("Contracts is unpaused, continue working!")
		return nil
	}
}

// GetTxGasPrice returns gas price from raw response if transaction's type - DynamicFeeTx
// else - returns default gas price from transaction
func GetTxGasPrice(client ethclients.ClientInterface, tx *types.Transaction) (*big.Int, error) {
	if tx.Type() == types.DynamicFeeTxType {
		return client.TxGasPriceFromResponse(context.Background(), tx.Hash())
	}
	return tx.GasPrice(), nil
}

func WaitForBlock(wsClient ethclients.ClientInterface, targetBlockNum uint64) error {

	// todo maybe timeout (context)
	blockChannel := make(chan *types.Header)
	blockSub, err := wsClient.SubscribeNewHead(context.Background(), blockChannel)
	if err != nil {
		return fmt.Errorf("SubscribeNewHead: %w", err)
	}
	defer blockSub.Unsubscribe()

	currentBlockNum, err := wsClient.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("get last block num: %w", err)
	}

	for currentBlockNum < targetBlockNum {
		select {
		case err := <-blockSub.Err():
			return fmt.Errorf("listening new blocks: %w", err)

		case block := <-blockChannel:
			currentBlockNum = block.Number.Uint64()
		}
	}

	return nil
}

func WaitForNextBlock(wsClient ethclients.ClientInterface) error {
	latestBlock, err := wsClient.BlockNumber(context.Background())
	if err != nil {
		return fmt.Errorf("get latest block number: %w", err)
	}

	return WaitForBlock(wsClient, latestBlock+1)
}

func GetMultipliedEstimatedGasLimit(auth bind.TransactOpts, multiplier float64, txCallback networks.ContractCallFn) (*bind.TransactOpts, error) {
	var authChangedGasLimit = auth
	// Make tx without sending it for getting the gas limit.
	auth.NoSend = true
	tx, err := txCallback(&auth)
	if err != nil {
		return nil, err
	}

	changedGasLimit := uint64(float64(tx.Gas()) * multiplier)
	authChangedGasLimit.GasLimit = changedGasLimit
	return &authChangedGasLimit, nil
}
