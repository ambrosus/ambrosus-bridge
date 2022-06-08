package bindings

import (
	"github.com/ethereum/go-ethereum/core/types"
)

// Transfer event

func (te *BridgeTransfer) Log() *types.Log {
	return &te.Raw
}

func (te *BridgeTransfer) ProofElements() [][]byte {
	return [][]byte{te.Raw.Address.Bytes(), te.EventId.Bytes(), te.Raw.Data}
}

// Validator Set change event

func (vs *VsInitiateChange) Log() *types.Log {
	return &vs.Raw
}

func (vs *VsInitiateChange) ProofElements() [][]byte {
	return [][]byte{vs.Raw.Address.Bytes(), vs.Raw.Data}
}
