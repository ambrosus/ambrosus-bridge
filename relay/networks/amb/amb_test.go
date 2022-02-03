package amb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/big"
	"relay/config"
	"relay/helpers"
	"relay/receipts_proof/mytrie"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var ambBridge = New(&config.Bridge{Url: "https://network.ambrosus.io"})

func TestHeader(t *testing.T) {
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
	// todo HeaderByNumber return empty struct on fail but err == nil
	h, err := ambBridge.HeaderByNumber(big.NewInt(16021709))
	if err != nil {
		t.Fatal(err)
	}

	encodedBlockPoA, err := EncodeBlock(h, true)
	if err != nil {
		t.Fatal(err)
	}

	rlpBase := helpers.BytesConcat(encodedBlockPoA.P1, encodedBlockPoA.PrevHashOrReceiptRoot[:], encodedBlockPoA.P2)

	// without seal
	rlpWithoutSeal := helpers.BytesConcat(encodedBlockPoA.P0Bare, rlpBase)
	hashWithoutSeal := common.BytesToHash(mytrie.Hash(rlpWithoutSeal))

	if hashWithoutSeal != h.Hash(false) {
		t.Fatalf("wrong bare hash")
	}

	// with seal
	rlpWithSeal := helpers.BytesConcat(encodedBlockPoA.P0Seal, rlpBase, encodedBlockPoA.S1, encodedBlockPoA.Step, encodedBlockPoA.S2, encodedBlockPoA.Signature)
	hashWithSeal := common.BytesToHash(mytrie.Hash(rlpWithSeal))

	if hashWithSeal != h.Hash(true) {
		t.Fatalf("wrong seal hash")
	}

}
