package amb

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"relay/config"
	"relay/helpers"
	"relay/receipts_proof/mytrie"
	"testing"
)

var ambBridge = New(&config.Bridge{Url: "https://network.ambrosus.io"})

func TestHeader(t *testing.T) {
	h, err := ambBridge.HeaderByNumber(big.NewInt(16021709))
	if err != nil {
		t.Fatal(err)
	}

	bare, err := h.Rlp(false)
	if err != nil {
		t.Fatal(err)
	}

	seal, err := h.Rlp(true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%x\n", bare)
	fmt.Printf("%x\n", seal)
}

func TestEncoding(t *testing.T) {
	// todo HeaderByNumber return empty struct on fail but err == nil
	h, err := ambBridge.HeaderByNumber(big.NewInt(16021709))
	if err != nil {
		t.Fatal(err)
	}

	encodedBlockPoA := EncodeBlock(h, true)

	rlpBase := helpers.BytesConcat(encodedBlockPoA.P1, encodedBlockPoA.PrevHashOrReceiptRoot[:], encodedBlockPoA.P2, encodedBlockPoA.Timestamp, encodedBlockPoA.P3)

	// without seal
	rlpWithoutSeal := helpers.BytesConcat(encodedBlockPoA.P0Bare, rlpBase)
	hashWithoutSeal := common.BytesToHash(mytrie.Hash(rlpWithoutSeal))

	if hashWithoutSeal != h.Hash(false) {
		t.Fatalf("wrong bare hash")
	}

	// with seal
	rlpWithSeal := helpers.BytesConcat(encodedBlockPoA.P0Seal, rlpBase, encodedBlockPoA.S1, encodedBlockPoA.Signature, encodedBlockPoA.S2)
	hashWithSeal := common.BytesToHash(mytrie.Hash(rlpWithSeal))

	if hashWithSeal != h.Hash(true) {
		t.Fatalf("wrong seal hash")
	}

}
