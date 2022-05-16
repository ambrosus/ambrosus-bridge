package price

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {
	a, err := CoinToUsdt(Amb)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Amb: ", a)

	e, err := CoinToUsdt(Eth)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Eth: ", e)

	b, err := CoinToUsdt(Bnb)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Bnb: ", b)
}
