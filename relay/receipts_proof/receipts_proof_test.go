package test

import (
	"context"
	"relay/config"
	"relay/networks/amb"
	"relay/receipts_proof"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const url = "https://network.ambrosus-dev.io"

func TestCalcProof(t *testing.T) {

	ambBridge := amb.New(&config.Bridge{
		Url:             url,
		ContractAddress: common.HexToAddress("0x8b12C2C9C61Ae2081567a01091654B21a018db29"),
		PrivateKey:      nil,
		SafetyBlocks:    10,
	})

	log := getLog2(ambBridge, t)

	block, err := ambBridge.Client.BlockByHash(context.Background(), log.BlockHash)
	if err != nil {
		t.Fatal(err)
	}
	receipts, err := ambBridge.GetReceipts(log.BlockHash)
	if err != nil {
		t.Fatal(err)
	}

	proof, err := receipts_proof.CalcProof(receipts, log)
	if err != nil {
		t.Fatal(err)
	}

	receiptsRoot := receipts_proof.CheckProof(proof, log)
	if receiptsRoot != block.ReceiptHash() {
		t.Fatal("proof check failed")
	}
}

func getLog1(ambBridge *amb.Bridge, t *testing.T) *types.Log {
	transfers, err := ambBridge.Contract.FilterTransferEvent(nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if !transfers.Next() {
		if err := transfers.Error(); err != nil {
			t.Fatal(err)
		}
		t.Fatal("no transfers")
	}
	return &transfers.Event.Raw
}

func getLog2(ambBridge *amb.Bridge, t *testing.T) *types.Log {
	receipt, err := ambBridge.Client.TransactionReceipt(context.Background(), common.HexToHash("0xfe802d29486f3648fd082ecad8c8455aa751b18f36be7cff3cfd65245253233a"))
	if err != nil {
		t.Fatal(err)
	}
	return receipt.Logs[1]
}
