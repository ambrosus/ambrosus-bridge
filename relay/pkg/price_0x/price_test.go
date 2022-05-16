package price_0x

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/params"
)

func TestRequests(t *testing.T) {
	e, err := CoinToUSDT("ETH", params.Ether)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Eth to USDT: ", e)

}
