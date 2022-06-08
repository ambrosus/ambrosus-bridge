package helpers

import (
	"math/big"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/stretchr/testify/assert"
)

func TestIntersectionSubmitsUnlocks(t *testing.T) {
	submits := generateSubmits(2, 8)
	unlocks := generateUnlocks(1, 7)

	resSubmits, resUnlocks := IntersectionSubmitsUnlocks(submits, unlocks)

	assert.Equal(t, resSubmits, generateSubmits(2, 7))
	assert.Equal(t, resUnlocks, generateUnlocks(2, 7))
}

func generateSubmits(start, end int) []*bindings.BridgeTransferSubmit {
	var res []*bindings.BridgeTransferSubmit
	for i := start; i <= end; i++ {
		res = append(res, &bindings.BridgeTransferSubmit{EventId: big.NewInt(int64(i))})
	}
	return res
}

func generateUnlocks(start, end int) []*bindings.BridgeTransferFinish {
	var res []*bindings.BridgeTransferFinish
	for i := start; i <= end; i++ {
		res = append(res, &bindings.BridgeTransferFinish{EventId: big.NewInt(int64(i))})
	}
	return res
}
