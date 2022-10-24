package common

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func ShouldHavePk(b networks.Bridge) {
	if b.GetAuth() == nil {
		b.GetLogger().Fatal().Msg("Private key is required")
	}
}

func EnsureContractUnpaused(b networks.Bridge) {
	for {
		err := waitForUnpauseContract(b)
		if err == nil {
			return
		}

		b.GetLogger().Error().Err(err).Msg("waitForUnpauseContract error")
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

	eventCh := make(chan *bindings.BridgeUnpaused)
	eventSub, err := b.GetWsContract().WatchUnpaused(nil, eventCh)
	if err != nil {
		return fmt.Errorf("WatchUnpaused: %w", err)
	}
	defer eventSub.Unsubscribe()

	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching unpaused event: %w", err)
		case event := <-eventCh:
			if event.Raw.Removed {
				continue
			}

			b.GetLogger().Info().Msg("Contracts is unpaused, continue working!")
			return nil
		}
	}
}

// GetEventById get `Transfer` event (emitted by this contract) by id.
func GetEventById(client interfaces.BridgeContract, eventId *big.Int) (*bindings.BridgeTransfer, error) {
	logs, err := client.FilterTransfer(nil, []*big.Int{eventId})
	if err != nil {
		return nil, fmt.Errorf("filter transfer: %w", err)
	}
	for logs.Next() {
		if !logs.Event.Raw.Removed {
			return logs.Event, nil
		}
	}
	return nil, networks.ErrEventNotFound
}

func EncodeTransferProof(client ethclients.ClientInterface, event *bindings.BridgeTransfer) (bindings.CommonStructsTransferProof, error) {
	proof, err := GetProof(client, event)
	if err != nil {
		return bindings.CommonStructsTransferProof{}, fmt.Errorf("GetProof: %w", err)
	}

	return bindings.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func GetProof(client ethclients.ClientInterface, event receipts_proof.ProofEvent) ([][]byte, error) {
	receipts, err := getReceipts(client, event.Log().BlockHash)
	if err != nil {
		return nil, fmt.Errorf("getReceipts: %w", err)
	}
	return receipts_proof.CalcProofEvent(receipts, event)
}

func getReceipts(client ethclients.ClientInterface, blockHash common.Hash) ([]*types.Receipt, error) {
	txsCount, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		return nil, fmt.Errorf("get transaction count: %w", err)
	}

	receipts := make([]*types.Receipt, txsCount)

	for i := uint(0); i < txsCount; i++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, i)
		if err != nil {
			return nil, fmt.Errorf("get transaction in block: %w", err)
		}
		receipts[i], err = client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return nil, fmt.Errorf("get transaction receipt: %w", err)
		}
	}

	return receipts, nil
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
