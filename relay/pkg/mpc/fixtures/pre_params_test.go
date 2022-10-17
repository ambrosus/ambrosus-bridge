//go:build !ci

package fixtures

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/bnb-chain/tss-lib/ecdsa/keygen"
)

// Generate pre params fixtures
func TestGeneratePreParams(t *testing.T) {
	pp, err := keygen.GeneratePreParams(5 * time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	marshalled, err := json.Marshal(pp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Print(string(marshalled))
}
