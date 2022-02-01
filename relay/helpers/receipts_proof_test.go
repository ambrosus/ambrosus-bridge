package helpers

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

const url = "https://rinkeby.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c"

func TestCalcProof(t *testing.T) {
	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x4e4bf7bb5f732326af2425ffe02359a0f9049c1367ecc7ca2cc84237315093bc"))
	if err != nil {
		t.Fatal(err)
	}

	block, err := client.BlockByHash(context.Background(), receipt.Logs[0].BlockHash)
	if err != nil {
		t.Fatal(err)
	}

	receipts := types.Receipts(FindReceipts(client, block.Hash()))
	proof := CalcProof(&receipts, receipt.Logs[0].Data)

	if !CheckProof(receipt.Logs[0].Data, proof, block.ReceiptHash()) {
		t.Fatal("proof check failed")
	}
}
