package parity

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	parityClient, err := Dial("https://network.ambrosus.io")
	if err != nil {
		t.Fatal(err)
	}

	h, err := parityClient.ParityHeaderByNumber(context.Background(), big.NewInt(13000000))
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
