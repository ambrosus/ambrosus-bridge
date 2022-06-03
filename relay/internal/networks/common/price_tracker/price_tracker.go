package price_tracker

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
)

// PriceTrackerData is kinda memory DB for storing previous results and prev used event id
type PriceTrackerData struct {
	// Bridge *nc.CommonBridge

	Submits   map[string]*contracts.BridgeTransferSubmit
	Unlocks   map[string]*contracts.BridgeTransferFinish
	Transfers map[string]*contracts.BridgeTransfer

	PrevUsedUnlockEventId   *big.Int // id when "GasPerWithdraw" was called last time (by the user from api)
	LastCaughtUnlockEventId *big.Int // id that was set by the watcher

	TotalGasCost   *big.Int
	WithdrawsCount int64
}

func (d *PriceTrackerData) Save(unlockEventId *big.Int, totalGasCost *big.Int, withdrawsCount int64) {
	d.PrevUsedUnlockEventId = unlockEventId

	d.TotalGasCost.Add(d.TotalGasCost, totalGasCost)
	d.WithdrawsCount += withdrawsCount
}
