package config

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Bridge struct {
	Url               string
	HttpUrl           string // used for getting header in ambrosus
	ContractAddress   ethcommon.Address
	VSContractAddress ethcommon.Address
	PrivateKey        *ecdsa.PrivateKey
	SafetyBlocks      uint64
	ChainID           *big.Int
}

// todo load from json

var Config = map[string]*Bridge{
	"amb": {
		Url:               "wss://network.ambrosus-dev.io",
		HttpUrl:           "https://network.ambrosus-dev.io",
		ContractAddress:   ethcommon.HexToAddress(""),
		VSContractAddress: ethcommon.HexToAddress(""),
		PrivateKey:        ParsePK("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"),
		SafetyBlocks:      10,
	},
	"eth": {
		Url:             "wss://rinkeby.infura.io/ws/v3/01117e8ede8e4f36801a6a838b24f36c",
		ContractAddress: ethcommon.HexToAddress("0x74586F2646faf275D1c90b4F347BF9C0c5E9E3a1"),
		PrivateKey:      ParsePK("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"),
		SafetyBlocks:    10,
		ChainID:         big.NewInt(1),
	},
	"oelocal": {
		Url:               "ws://127.0.0.1:8546",
		HttpUrl:           "http://127.0.0.1:8545",
		ContractAddress:   ethcommon.HexToAddress(""),
		VSContractAddress: ethcommon.HexToAddress(""),
		PrivateKey:        ParsePK("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"),
		SafetyBlocks:      1,
		ChainID:           big.NewInt(1),
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
