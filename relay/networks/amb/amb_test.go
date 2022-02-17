package amb

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/helpers"
	"github.com/ambrosus/ambrosus-bridge/relay/receipts_proof/mytrie"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	ambBridge, err := New(&config.Bridge{Url: "wss://network.ambrosus.io", HttpUrl: "https://network.ambrosus.io"})
	if err != nil {
		t.Fatal(err)
	}

	h, err := ambBridge.HeaderByNumber(big.NewInt(13000000))
	if err != nil {
		t.Fatal(err)
	}

	bareRlp, err := h.Rlp(false)
	if err != nil {
		t.Fatal(err)
	}

	sealRlp, err := h.Rlp(true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%x\n", bareRlp)
	fmt.Printf("%x\n", sealRlp)

	assert.Equal(t, common.HexToHash("0x63deebcabaa73e872ee762e7b1dc12b849a5032d536362d5428a6992f5b5e848"), h.Hash(false), "bare (for signature)")
	assert.Equal(t, common.HexToHash("0xbd002f9a7e73ba2a1a33e90ad196301727e6d1aacd3e5a2c594b0f455f967d9f"), h.Hash(true), "seal (for parent_hash)")

}

func TestEncoding(t *testing.T) {
	ambBridge, err := New(&config.Bridge{Url: "wss://network.ambrosus.io", HttpUrl: "https://network.ambrosus.io"})
	if err != nil {
		t.Fatal(err)
	}

	// todo HeaderByNumber return empty struct on fail but err == nil
	h, err := ambBridge.HeaderByNumber(big.NewInt(16021709))
	if err != nil {
		t.Fatal(err)
	}

	block, err := EncodeBlock(h)
	if err != nil {
		t.Fatal(err)
	}

	rlpCommon := helpers.BytesConcat(block.P1, block.ParentHash[:], block.P2, block.ReceiptHash[:], block.P3)

	// without seal
	rlpWithoutSeal := helpers.BytesConcat(block.P0Bare, rlpCommon)
	hashWithoutSeal := common.BytesToHash(mytrie.Hash(rlpWithoutSeal))

	if hashWithoutSeal != h.Hash(false) {
		t.Fatalf("wrong bare hash")
	}

	// with seal
	rlpWithSeal := helpers.BytesConcat(block.P0Seal, rlpCommon, block.S1, block.Step, block.S2, block.Signature)
	hashWithSeal := common.BytesToHash(mytrie.Hash(rlpWithSeal))

	if hashWithSeal != h.Hash(true) {
		t.Fatalf("wrong seal hash")
	}

}
