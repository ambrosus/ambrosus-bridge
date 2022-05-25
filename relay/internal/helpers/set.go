package helpers

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
)

func IntersectionSubmitsUnlocks(submits []*contracts.BridgeTransferSubmit, unlocks []*contracts.BridgeTransferFinish) (
	[]*contracts.BridgeTransferSubmit,
	[]*contracts.BridgeTransferFinish,
) {
	var eventIds []*big.Int
	hash := make(map[string]bool)

	for _, item := range submits {
		hash[item.EventId.String()] = true
	}

	for _, item := range unlocks {
		if _, ok := hash[item.EventId.String()]; ok {
			eventIds = append(eventIds, item.EventId)
		}
	}

	// build submits and unlocks from intersection
	var resSubmits []*contracts.BridgeTransferSubmit
	var resUnlocks []*contracts.BridgeTransferFinish
	startSubmits := 0
	startUnlocks := 0

	// find the start point
	for i := 0; i < len(eventIds); i++ {
		if submits[i].EventId.Cmp(eventIds[0]) == 0 {
			startSubmits = i
			break
		}
	}
	for i := 0; i < len(eventIds); i++ {
		if unlocks[i].EventId.Cmp(eventIds[0]) == 0 {
			startUnlocks = i
			break
		}
	}

	resSubmits = submits[startSubmits : len(eventIds)+startSubmits]
	resUnlocks = unlocks[startUnlocks : len(eventIds)+startUnlocks]
	return resSubmits, resUnlocks
}
