package enc2sol

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"relay/enc2sol/mytrie"
)

const SafetyBlocks = 9

var bigOne = big.NewInt(1)

type EncodedBlock [][]byte

func Encode(client *ethclient.Client, event *types.Log) (*types.Block, EncodedBlock, []EncodedBlock, [][]byte) {
	block, err := client.BlockByHash(context.Background(), event.BlockHash)
	if err != nil {
		panic(err)
	}

	receipts := types.Receipts(findReceipts(client, block.Hash()))
	proof := getReceiptProof(&receipts, event.Data)

	encodedBlock := EncodeBlock(block.Header(), true)
	otherBlocks := EncodeOtherBlock(client, block.Number())

	return block, encodedBlock, otherBlocks, proof
}

func EncodeOtherBlock(client *ethclient.Client, blockNum *big.Int) []EncodedBlock {
	var result []EncodedBlock

	for i := int64(0); i < SafetyBlocks; i++ {
		block, err := client.BlockByNumber(context.Background(), blockNum.Add(blockNum, bigOne))
		if err != nil {
			panic(err)
		}
		encodedBlock := EncodeBlock(block.Header(), false)
		result = append(result, encodedBlock)
	}
	return result
}

func EncodeBlock(header *types.Header, main bool) (result EncodedBlock) {
	splitEls := make([][]byte, 2)
	splitEls[1], _ = hex.DecodeString(fmt.Sprintf("%x", header.Time))
	if main {
		splitEls[0] = header.ReceiptHash.Bytes()
	} else {
		splitEls[0] = header.ParentHash.Bytes()
	}
	return encodeBlock(header, splitEls)
}

func encodeBlock(header *types.Header, splitEls [][]byte) (result EncodedBlock) {
	rlpHeader, _ := rlp.EncodeToBytes(header)

	for _, se := range splitEls {
		r := bytes.SplitN(rlpHeader, se, 2)
		result = append(result, r[0], se)
		rlpHeader = r[1]
	}
	result = append(result, rlpHeader)
	return
}

func getReceiptProof(receipts *types.Receipts, eventDataToSearch []byte) [][]byte {
	hasher := mytrie.NewStackTrie()
	types.DeriveSha(receipts, hasher)
	return CalcProof(hasher, eventDataToSearch)

}

func findReceipts(client *ethclient.Client, blockHash common.Hash) []*types.Receipt {
	txsCount, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		panic(err)
	}

	receipts := make([]*types.Receipt, 0, txsCount)

	for i := uint(0); i < txsCount; i++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, i)
		if err != nil {
			panic(err)
		}
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			panic(err)
		}
		receipts = append(receipts, receipt)
	}
	return receipts
}
