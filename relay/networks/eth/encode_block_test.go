package eth

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const url = "https://mainnet.infura.io/v3/73b4b0b7af6c459e97f8350277274156"

func Test(t *testing.T) {
	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), big.NewInt(11111))
	if err != nil {
		t.Fatal(err)
	}

	encodedBlock, err := EncodeBlock(block.Header(), true)
	if err != nil {
		t.Fatal(err)
	}

	if common.BytesToHash(encodedBlock.PrevHashOrReceiptRoot[:]) != block.ReceiptHash() {
		t.Fatal("receiptsHash from encoded block != original")
	}
}
