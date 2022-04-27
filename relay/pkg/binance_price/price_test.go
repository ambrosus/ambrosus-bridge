package binance_price

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {
	e, err := CoinToAmb(Eth)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Eth to Amb: ", e)

	b, err := CoinToAmb(Bnb)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Bnb to Amb: ", b)
}
