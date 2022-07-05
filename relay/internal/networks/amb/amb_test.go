package amb

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	ambBridge, err := New(&config.AMBConfig{
		Network: config.Network{HttpURL: "https://network.ambrosus.io"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	h, err := ambBridge.ParityClient.ParityHeaderByNumber(context.Background(), big.NewInt(13000000))
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

func TestGasPrice(t *testing.T) {
	ambBridge, err := New(&config.AMBConfig{
		Network: config.Network{
			HttpURL:      "https://network.ambrosus-dev.io",
			ContractAddr: "0x617F296c197266305904063CEFB07C9E3295D743",
			PrivateKey:   "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5",
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ethBridge, err := New(&config.AMBConfig{
		Network: config.Network{
			HttpURL:      "https://ropsten.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c",
			ContractAddr: "0xAd6557e9793F119e4d8601Eb5cB1b79b26d89fDb",
			PrivateKey:   "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5",
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	ambBridge.SideBridge = ethBridge
	ethBridge.SideBridge = ambBridge

	//todo
	//r, err := ambBridge.CommonBridge.GasPerWithdraw(ambBridge.PriceTrackerData)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Logf("amb->eth %d\n", r)
	//
	//r, err = ethBridge.CommonBridge.GasPerWithdraw(ethBridge.PriceTrackerData)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Logf("eth->amb %d", r)

}
