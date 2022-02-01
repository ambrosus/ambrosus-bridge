package receipts_proof

import (
	"context"
	"relay/config"
	"relay/networks/amb"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const url = "https://network.ambrosus-dev.io"

func TestCalcProof(t *testing.T) {

	ambBridge := amb.New(&config.Bridge{
		Url:             url,
		ContractAddress: common.HexToAddress("0xE3A1f4Af2c71957033BaD65771f59C4e797C0693"),
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

	proof, err := CalcProof(receipts, log)
	if err != nil {
		t.Fatal(err)
	}

	receiptsRoot := CheckProof(proof, log)
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
	receipt, err := ambBridge.Client.TransactionReceipt(context.Background(), common.HexToHash("0xeff82a20a691eb2d9fd3fe726ef09731b11aa12e56960544abb4612eb2c73ab3"))
	if err != nil {
		t.Fatal(err)
	}
	return receipt.Logs[0]
}
