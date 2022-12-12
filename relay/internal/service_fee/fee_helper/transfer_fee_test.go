//go:build !ci

package fee_helper

import (
	"testing"
)

func Test_getTransferFee(t *testing.T) {
	thisGas, sideGas, err := getTransferFee(
		"http://backoffice-api.ambrosus-test.io/relay/fee?networkThis=amb&networkSide=eth")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("thisGas:", thisGas)
	t.Log("sideGas:", sideGas)
}
