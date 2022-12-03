package events

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
)

type Events interface {
	WatchPaused() error
	WatchUnpaused() error
	GetTransfer(eventID uint64) (*bindings.BridgeTransfer, error)
	WatchTransfer() (*bindings.BridgeTransfer, error)
	WatchTransferSubmit() (*bindings.BridgeTransferSubmit, error)
	WatchTransferFinish() (*bindings.BridgeTransferFinish, error)
	GetWithdrawsCount(fromEvent, toEvent uint64) (int, error)
}
