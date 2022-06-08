package price_amb

import (
	"fmt"
	"testing"
)

func TestRequests(t *testing.T) {
	a, err := Get()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Amb: ", a)
}
