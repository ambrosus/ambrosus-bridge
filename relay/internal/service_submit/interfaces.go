package service_submit

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
)

type Submitter interface {
	networks.Bridge
	SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error
	Receiver() Receiver
}

type Receiver interface {
	networks.Bridge
}

type ReceiverUntrustless interface {
	Receiver
	SubmitTransferUntrustless(*bindings.BridgeTransfer) error
	IsEventAlreadyConfirmed(event *bindings.BridgeTransfer) (bool, error)
}

type ReceiverUntrustless2 interface {
	Receiver
	SubmitTransferUntrustless(*bindings.BridgeTransfer) error
}

type ReceiverUntrustlessMpc interface {
	Receiver
	SubmitTransferUntrustlessMpcServer(*bindings.BridgeTransfer) error
	SubmitTransferUntrustlessMpcClient(*bindings.BridgeTransfer) error
}
