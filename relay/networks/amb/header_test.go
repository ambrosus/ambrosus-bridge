package amb

import (
	"fmt"
	"math/big"
	"testing"
)

func TestHeaderHash(t *testing.T) {
	number := big.NewInt(16021709)
	h, _ := HeaderByNumber(number)

	seal, _ := h.SealRlp()
	fmt.Printf("%x\n", h.Rlp(true))
	fmt.Printf("%x\n", h.Rlp(false))
	fmt.Printf("%x\n", seal)
}
