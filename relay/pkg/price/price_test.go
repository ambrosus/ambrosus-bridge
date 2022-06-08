package price

import (
	"fmt"
	"testing"
)

func TestCoinToUSD(t *testing.T) {
	e, err := CoinToUSD(EthUrl, "ETH", 18)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Eth to USDT: ", e)

}
func TestAmb(t *testing.T) {
	a, err := AmbToUSD()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Amb: ", a)
}
