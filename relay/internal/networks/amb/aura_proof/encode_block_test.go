package aura_proof

import (
	"context"
	"math/big"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/parity"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	parentHashPrefix = []byte{0xA0}
	stepPrefix       = []byte{0x84}
	signaturePrefix  = []byte{0xB8, 0x41}
)

func TestEncoding(t *testing.T) {
	parityClient, err := parity.Dial("https://network.ambrosus.io")
	if err != nil {
		t.Fatal(err)
	}
	h, err := parityClient.ParityHeaderByNumber(context.Background(), big.NewInt(16021709))
	if err != nil {
		t.Fatal(err)
	}

	block, err := EncodeBlock(h)
	if err != nil {
		t.Fatal(err)
	}

	rlpCommon := helpers.BytesConcat(parentHashPrefix, block.ParentHash[:], block.P2, block.ReceiptHash[:], block.P3)

	// without seal
	rlpWithoutSeal := helpers.BytesConcat(block.P0Bare[:], rlpCommon)
	hashWithoutSeal := common.BytesToHash(crypto.Keccak256(rlpWithoutSeal))

	if hashWithoutSeal != h.Hash(false) {
		t.Fatalf("wrong bare hash")
	}

	// with seal
	rlpWithSeal := helpers.BytesConcat(block.P0Seal[:], rlpCommon, stepPrefix, block.Step[:], signaturePrefix, block.Signature)
	hashWithSeal := common.BytesToHash(crypto.Keccak256(rlpWithSeal))

	if hashWithSeal != h.Hash(true) {
		t.Fatalf("wrong seal hash")
	}

}
