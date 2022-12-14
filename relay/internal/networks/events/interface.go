package events

import (
	"errors"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
)

var ErrEventNotFound = errors.New("error event not found")

type Events interface {
	WatchPaused() error
	WatchUnpaused() error
	GetTransfer(eventID uint64) (*bindings.BridgeTransfer, error)
	WatchTransfer() (*bindings.BridgeTransfer, error)
	WatchTransferSubmit() (*bindings.BridgeTransferSubmit, error)
	WatchTransferFinish() (*bindings.BridgeTransferFinish, error)
	GetWithdrawsCount(fromEvent, toEvent uint64) (int, error)
}
