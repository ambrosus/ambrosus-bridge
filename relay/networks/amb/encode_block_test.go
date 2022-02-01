package amb

import (
	"math/big"
	"relay/helpers"
	"relay/receipts_proof/mytrie"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test(t *testing.T) {
	h, err := HeaderByNumber(big.NewInt(16021709))
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
