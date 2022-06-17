package bsc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
)

func Test_splitVsChanges(t *testing.T) {
	f, err := ioutil.ReadFile("../../../posaproof.json")
	if err != nil {
		t.Fatal(err)
	}

	var proof c.CheckPoSAPoSAProof
	err = json.Unmarshal(f, &proof)
	if err != nil {
		t.Fatal(err)
	}

	changes := splitVsChanges(&proof, 1)
	if err != nil {
		t.Fatal(err)
	}

	for _, change := range changes {
		fmt.Println(len(change.Blocks))
	}
}
