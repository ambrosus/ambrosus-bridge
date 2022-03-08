package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Transfer event

type TransferEvent struct {
	EventId *big.Int
	Queue   []CommonStructsTransfer
	Raw     types.Log // Blockchain specific contextual infos
}

// EthTransfer represents a Transfer event raised by the Eth contract.
type EthTransfer struct {
	TransferEvent
}

// AmbTransfer represents a Transfer event raised by the Amb contract.
type AmbTransfer struct {
	TransferEvent
}

func (te *TransferEvent) Log() *types.Log {
	return &te.Raw
}

func (te *TransferEvent) ProofElements() [][]byte {
	return [][]byte{te.Raw.Address.Bytes(), te.EventId.Bytes(), te.Raw.Data}
}

// Validator Set change event

// VsInitiateChange represents a InitiateChange event raised by the Vs contract.
type VsInitiateChange struct {
	ParentHash [32]byte
	NewSet     []common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

func (vs *VsInitiateChange) Log() *types.Log {
	return &vs.Raw
}

func (vs *VsInitiateChange) ProofElements() [][]byte {
	return [][]byte{vs.Raw.Address.Bytes(), vs.Raw.Data}
}
