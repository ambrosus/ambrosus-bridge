package enc2sol

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"testing"
)

const url = "https://rinkeby.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c"

func Test(t *testing.T) {
	client, err := ethclient.Dial(url)
	if err != nil {
		t.Fatal(err)
	}

	receipt, err := client.TransactionReceipt(context.Background(), common.HexToHash("0x4e4bf7bb5f732326af2425ffe02359a0f9049c1367ecc7ca2cc84237315093bc"))
	if err != nil {
		t.Fatal(err)
	}

	block, encodedBlock, otherBlocks, proof := Encode(client, receipt.Logs[0])

	if common.BytesToHash(encodedBlock[1]) != block.ReceiptHash() {
		t.Fatal("receiptsHash from encoded block != original")
	}

	if !CheckProof(receipt.Logs[0].Data, proof, block.ReceiptHash()) {
		t.Fatal("proof check failed")
	}

	fmt.Printf("const rlpBlocks = %v;\n", formatBlocks(encodedBlock, otherBlocks))
	fmt.Printf("const proof = %v;\n", formatList(proof))

}

func formatBlocks(mainBlock EncodedBlock, otherBlocks []EncodedBlock) string {
	blocks := append([]EncodedBlock{mainBlock}, otherBlocks...)

	r := "[\n"
	for _, e := range blocks {
		r += formatList(e) + ", \n"
	}
	return r + "]"
}

func formatList(l [][]byte) string {
	r := "["
	for _, e := range l {
		r += fmt.Sprintf("\"0x%x\", ", e)
	}
	return r + "]"
}
