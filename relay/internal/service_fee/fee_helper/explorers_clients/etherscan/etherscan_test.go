package etherscan

import (
	"testing"

	"github.com/nanmu42/etherscan-api"
	"github.com/stretchr/testify/assert"
)

func TestEtherscan(t *testing.T) {
	e := Etherscan{etherscan.New(etherscan.Mainnet, "DY4Z86MQ2D9E24C6HB98PTA79EKJ5TQIFX")}
	res, err := e.TxListByFromToAddressesUntilTxHash(
		"0xFAE075e12116FBfE65c58e1Ef0E6CA959cA37ded",
		"0x0de2669e8a7a6f6cc0cbd3cf2d1eead89e243208",
		"0xa1e0662e972ef31f563df927f95d9cfe23cbc8d2672b7e5d9f1cd6543bb08f35",
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range res {
		t.Log(v.Hash)
		t.Log(v.BlockNumber)
	}
	
	assert.Equal(t, "0x54e359093e5ae4587640e2f168980a7ae7db9dcf8b1d67857b1dbeddd3799b26", res[0].Hash)
}

func TestEtherscan2(t *testing.T) {
	e := Etherscan{etherscan.New(etherscan.Mainnet, "DY4Z86MQ2D9E24C6HB98PTA79EKJ5TQIFX")}
	res, err := e.TxListByFromToAddresses(
		"0xFAE075e12116FBfE65c58e1Ef0E6CA959cA37ded",
		"0x0de2669e8a7a6f6cc0cbd3cf2d1eead89e243208",
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range res {
		t.Log(v.Hash)
		t.Log(v.BlockNumber)
	}
	
	// assert.Equal(t, "0x54e359093e5ae4587640e2f168980a7ae7db9dcf8b1d67857b1dbeddd3799b26", res[0].Hash)
}
