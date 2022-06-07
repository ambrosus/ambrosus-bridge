package helpers

import (
	"math/big"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/stretchr/testify/assert"
)

func TestIntersectionSubmitsUnlocks(t *testing.T) {
	submits := generateSubmits(2, 8)
	unlocks := generateUnlocks(1, 7)

	resSubmits, resUnlocks := IntersectionSubmitsUnlocks(submits, unlocks)

	assert.Equal(t, resSubmits, generateSubmits(2, 7))
	assert.Equal(t, resUnlocks, generateUnlocks(2, 7))
}

func generateSubmits(start, end int) []*contracts.BridgeTransferSubmit {
	var res []*contracts.BridgeTransferSubmit
	for i := start; i <= end; i++ {
		res = append(res, &contracts.BridgeTransferSubmit{EventId: big.NewInt(int64(i))})
	}
	return res
}

func generateUnlocks(start, end int) []*contracts.BridgeTransferFinish {
	var res []*contracts.BridgeTransferFinish
	for i := start; i <= end; i++ {
		res = append(res, &contracts.BridgeTransferFinish{EventId: big.NewInt(int64(i))})
	}
	return res
}
