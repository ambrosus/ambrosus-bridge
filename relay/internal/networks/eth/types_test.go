package eth

import (
	"math/big"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
)

func TestDump(t *testing.T) {
	b, err := New(&config.ETHConfig{
		Network: config.Network{HttpURL: "https://rinkeby.infura.io/v3/01117e8ede8e4f36801a6a838b24f36c",
			ContractAddr: "0x8fC5BcE484C6B937aCf6B96268cc9318c65255cD"},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = b.GetEventById(big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
}
