package price

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test0x(t *testing.T) {
	e, err := TokenToUSD(&TokenInfo{Symbol: "ETH", Decimals: 18, Address: common.Address{}})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Eth to USD: ", e)
}

func Test0xBNB(t *testing.T) {
	e, err := TokenToUSD(&TokenInfo{Symbol: "WBNB", Decimals: 18, Address: common.Address{}})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("WBNB to USD: ", e)
}

func TestAmb(t *testing.T) {
	a, err := GetAmb()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Amb: ", a)
}

func TestKucoin(t *testing.T) {
	a, err := GetKucoin(&TokenInfo{Symbol: "USDT", Decimals: 18})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("USDT to USD: %.30f", a)
}
