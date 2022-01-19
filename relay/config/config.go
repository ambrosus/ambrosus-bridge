package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Network struct {
	Url          string
	SafetyBlocks int
}

type Bridge struct {
	ContractAddress ethcommon.Address
	PrivateKey      *ecdsa.PrivateKey
}

type BridgePair struct {
	Amb         *Bridge
	Side        *Bridge
	SideNetwork string
}

// todo load from json

var Networks = map[string]*Network{
	"amb": {
		Url:          "https://network.ambrosus.io",
		SafetyBlocks: 10,
	},
	"eth": {
		Url:          "https://rinkeby.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c",
		SafetyBlocks: 10,
	},
}

var Bridges = []BridgePair{
	{
		Amb: &Bridge{
			ContractAddress: ethcommon.HexToAddress(""),
			PrivateKey:      parsePK("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"),
		},
		Side: &Bridge{

			ContractAddress: ethcommon.HexToAddress(""),
			PrivateKey:      parsePK("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"),
		},
		SideNetwork: "eth",
	},
}

func parsePK(pk string) *ecdsa.PrivateKey {
	b, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	p, err := crypto.ToECDSA(b)
	if err != nil {
		panic(err)
	}
	return p
}
