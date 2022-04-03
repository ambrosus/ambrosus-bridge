package ethash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEpoch(t *testing.T) {
	ethash := New("./test")
	_, err := ethash.getCache(1)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ethash.getDag(1)
	if err != nil {
		t.Fatal(err)
	}
	edata, err := ethash.GetEpochData(1)
	if err != nil {
		t.Fatal(err)
	}

	a := fmt.Sprintf("%x", edata.MerkleNodes[0])
	if a != "e27b1e1d38a0f0cbfa44991921fd28964253c34a5602235017f6495a6f570718" {
		t.Fatal(a)
	}

}

func TestBytesToUint32Slice(t *testing.T) {
	a := []byte{0x10, 0x20, 0x30, 0x40}
	b := bytesToUint32Slice(a)
	assert.Equal(t, []uint32{1076895760}, b)
}
