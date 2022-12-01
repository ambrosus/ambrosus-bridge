package events

import (
	"encoding/json"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/common"
)

type TransferEvent struct {
	EventId helpers.BigInt
	Queue   []TransferStruct
}

func (e *TransferEvent) Unmarshall(data []byte) error {
	if err := json.Unmarshal(data, e); err != nil {
		return fmt.Errorf("unmarshal transfer event: %w", err)
	}
	return nil
}

// technically speaking, `TransferEvent` struct can be used for BridgeTransfer, BridgeTransferSubmit and BridgeTransferFinish :)

func (e *TransferEvent) toBridgeTransfer() *bindings.BridgeTransfer {
	queue := make([]bindings.CommonStructsTransfer, len(e.Queue))
	for i, item := range e.Queue {
		queue[i] = item.ToBinding()
	}
	return &bindings.BridgeTransfer{
		EventId: e.EventId.BigInt,
		Queue:   queue,
	}
}

func (e *TransferEvent) toBridgeTransferFinish() *bindings.BridgeTransferFinish {
	return &bindings.BridgeTransferFinish{
		EventId: e.EventId.BigInt,
	}
}
func (e *TransferEvent) toBridgeTransferSubmit() *bindings.BridgeTransferSubmit {
	return &bindings.BridgeTransferSubmit{
		EventId: e.EventId.BigInt,
	}
}

type TransferStruct struct {
	TokenAddress common.Address
	ToAddress    common.Address
	Amount       helpers.BigInt
}

func (e *TransferStruct) ToBinding() bindings.CommonStructsTransfer {
	return bindings.CommonStructsTransfer{
		TokenAddress: e.TokenAddress,
		ToAddress:    e.ToAddress,
		Amount:       e.Amount.BigInt,
	}
}
