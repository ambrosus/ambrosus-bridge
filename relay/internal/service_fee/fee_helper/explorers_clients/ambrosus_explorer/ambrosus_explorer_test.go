package ambrosus_explorer

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGetTxs(t *testing.T) {
	explorer, err := NewAmbrosusExplorer("https://explorer-api.ambrosus-dev.io", nil)
	if err != nil {
		t.Fatal(err)
	}
	r, err := explorer.TxListByFromToAddresses("0x295C2707319ad4BecA6b5bb4086617fD6F240CfE", "0xf7E15b720867747a536137f4EFdAB4309225f8D6")
	if err != nil {
		t.Fatal(err)
	}

	jr, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(jr))
}
