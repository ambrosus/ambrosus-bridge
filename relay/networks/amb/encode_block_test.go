package amb

import (
	"math/big"
	"relay/helpers/mytrie"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test(t *testing.T) {
	h, err := HeaderByNumber(big.NewInt(16021709))
	if err != nil {
		t.Fatal(err)
	}

	encodedBlock := EncodeBlock(h, true)
	rlpEncodedBlock := append(encodedBlock.P1, encodedBlock.PrevHashOrReceiptRoot[:]...)
	rlpEncodedBlock = append(rlpEncodedBlock, encodedBlock.P2...)
	rlpEncodedBlock = append(rlpEncodedBlock, encodedBlock.Timestamp...)
	rlpEncodedBlock = append(rlpEncodedBlock, encodedBlock.P3...)

	// with seal
	rlpEncodedBlockWithSeal := append(encodedBlock.P0Seal, rlpEncodedBlock...)
	rlpEncodedBlockWithSeal = append(rlpEncodedBlockWithSeal, encodedBlock.S1...)
	rlpEncodedBlockWithSeal = append(rlpEncodedBlockWithSeal, encodedBlock.Signature...)
	rlpEncodedBlockWithSeal = append(rlpEncodedBlockWithSeal, encodedBlock.S2...)
	encodedBlockHashWithSeal := common.BytesToHash(mytrie.Hash(rlpEncodedBlockWithSeal))
	headerHashWithSeal := h.Hash(true)
	if encodedBlockHashWithSeal != headerHashWithSeal {
		t.Fatalf("Header hash with seal %s != encodedBlock hash with seal %s", headerHashWithSeal, encodedBlockHashWithSeal)
	}

	// without seal
	rlpEncodedBlockWithoutSeal := append(encodedBlock.P0Bare, rlpEncodedBlock...)
	encodedBlockHashWithoutSeal := common.BytesToHash(mytrie.Hash(rlpEncodedBlockWithoutSeal))
	headerHashWithoutSeal := h.Hash(false)
	if encodedBlockHashWithoutSeal != headerHashWithoutSeal {
		t.Fatalf("Header hash without seal %s != encodedBlock hash without seal %s", headerHashWithSeal, encodedBlockHashWithSeal)
	}
}
