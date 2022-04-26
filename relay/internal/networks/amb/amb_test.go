package amb

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

var (
	parentHashPrefix = []byte{0xA0}
	stepPrefix       = []byte{0x84}
	signaturePrefix  = []byte{0xB8, 0x41}
)

func TestHeader(t *testing.T) {
	ambBridge, err := New(&config.AMBConfig{
		Network: config.Network{HttpURL: "https://network.ambrosus.io"},
	}, nil)
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
	ambBridge, err := New(&config.AMBConfig{
		Network: config.Network{HttpURL: "https://network.ambrosus.io"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	h, err := ambBridge.HeaderByNumber(big.NewInt(16021709))
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

func TestGasPrice(t *testing.T) {
	ambBridge, err := New(&config.AMBConfig{
		Network: config.Network{
			HttpURL:      "https://network.ambrosus-dev.io",
			ContractAddr: "0x4D6c49Dd98A25AD7bDAA64aA9ad7e0E9221bC2Ce",
			PrivateKey:   "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5",
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ethBridge, err := New(&config.AMBConfig{
		Network: config.Network{
			HttpURL:      "https://ropsten.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c",
			ContractAddr: "0xF634554393053F77F3C5093B52E3064A49aaB851",
			PrivateKey:   "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5",
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ambBridge.SideBridge = ethBridge
	ethBridge.SideBridge = ambBridge

	r, err := ethBridge.CommonBridge.GasPerWithdraw(big.NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("amb->eth", r)

	r, err = ambBridge.CommonBridge.GasPerWithdraw(big.NewInt(2))
	if err != nil {
		t.Fatal(err)
	}
	t.Log("eth->amb", r)

}
