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
	fmt.Println("Eth to USDT: ", e)

}
func TestAmb(t *testing.T) {
	a, err := GetAmb()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Amb: ", a)
}
