package price_amb

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {
	a, err := Get(Amb)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Amb: ", a)

	e, err := Get(Eth)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Eth: ", e)

	b, err := Get(Bnb)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Bnb: ", b)
}
