package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Bridge struct {
	Url             string
	ContractAddress ethcommon.Address
	PrivateKey      *ecdsa.PrivateKey
	SafetyBlocks    uint64
	ChainID         *big.Int
}

// todo load from json

var Config = map[string]*Bridge{
	"amb": {
		Url:             "https://network.ambrosus-dev.io",
		ContractAddress: ethcommon.HexToAddress(""),
		PrivateKey:      ParsePK("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"),
		SafetyBlocks:    10,
		ChainID:         big.NewInt(30741),
	},
	"eth": {
		Url:             "https://rinkeby.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c",
		ContractAddress: ethcommon.HexToAddress(""),
		PrivateKey:      ParsePK("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"),
		SafetyBlocks:    10,
		ChainID:         big.NewInt(1),
	},
}

func ParsePK(pk string) *ecdsa.PrivateKey {
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
