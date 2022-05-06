package helpers

import (
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestRlpOne(t *testing.T) {
	testRlp(t, 257)
}

func TestRlpMany(t *testing.T) {
	for i := 2; i < 1000; i++ {
		testRlp(t, i)
	}
}

func testRlp(t *testing.T, i int) {
	r, _ := rlp.EncodeToBytes(make([]byte, i))
	actual := RlpPrefix(i)
	expected := r[:len(actual)]

	assert.Equal(t, expected, actual, "", i, r)

}
