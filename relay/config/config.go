package config

import (
	"crypto/ecdsa"
	"encoding/hex"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type Bridge struct {
	Url               string
	HttpUrl           string // used for getting header in ambrosus
	ContractAddress   ethcommon.Address
	VSContractAddress ethcommon.Address
	PrivateKey        *ecdsa.PrivateKey
}

// todo load from json

var Config = map[string]*Bridge{
	"amb": {
		Url:               "wss://network.ambrosus-dev.io",
		HttpUrl:           "https://network.ambrosus-dev.io",
		ContractAddress:   ethcommon.HexToAddress(""),
		VSContractAddress: ethcommon.HexToAddress(""),
		PrivateKey:        ParsePK("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"),
	},
	"eth": {
		Url:             "wss://rinkeby.infura.io/ws/v3/01117e8ede8e4f36801a6a838b24f36c",
		ContractAddress: ethcommon.HexToAddress(""),
		PrivateKey:      ParsePK("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5"),
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
