package price_0x

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {
	e, err := CoinToUSDT("ETH", 18)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Eth to USDT: ", e)

}
