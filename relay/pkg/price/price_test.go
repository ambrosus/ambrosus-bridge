package price

import (
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func Test0x(t *testing.T) {
	os.Setenv("0X_API_KEY", "c7a70fdd-3474-425b-9fc3-2e40a7275ca3") // test api key
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
