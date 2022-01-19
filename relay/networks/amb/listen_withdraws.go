package amb

import (
	"context"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"relay/helpers"
	"relay/networks/common"
)

const SafetyBlocks = 9

type Listener struct {
	client *ethclient.Client
}

func (l *Listener) Run() {

	// todo listen
}

func (l *Listener) onEvent(event *types.Log) (*common.Withdraw, error) {
	block, err := l.client.BlockByHash(context.Background(), event.BlockHash)
	if err != nil {
		return nil, err
	}
	header, err := HeaderByNumber(block.Number())
	if err != nil {
		return nil, err
	}
	receipts := types.Receipts(l.getReceipts(block.Hash()))

	withdraw := &common.Withdraw{
		Network:       "",  // event first topic
		EventId:       nil, // event second topic
		Blocks:        make([]*common.AmbBlock, SafetyBlocks+1),
		Events:        nil,
		ReceiptsProof: helpers.CalcProof(&receipts, event.Data),
	}

	withdraw.Blocks[0] = EncodeBlock(header, true)

	for i := int64(1); i <= SafetyBlocks; i++ {
		h, err := HeaderByNumber(add(block.Number(), i))
		if err != nil {
			return nil, err
		}
		withdraw.Blocks[i] = EncodeBlock(h, false)
	}

	return withdraw, nil
}

func (l *Listener) getReceipts(blockHash ethcommon.Hash) []*types.Receipt {
	txsCount, err := l.client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		panic(err)
	}

	receipts := make([]*types.Receipt, 0, txsCount)

	for i := uint(0); i < txsCount; i++ {
		tx, err := l.client.TransactionInBlock(context.Background(), blockHash, i)
		if err != nil {
			panic(err)
		}
		receipt, err := l.client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			panic(err)
		}
		receipts = append(receipts, receipt)
	}
	return receipts
}

func add(a *big.Int, b int64) *big.Int {
	res := &big.Int{}
	return res.Add(a, big.NewInt(b))
}
