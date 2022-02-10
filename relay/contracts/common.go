package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

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

// CommonStructsTransfer is an auto generated low-level Go binding around an user-defined struct.
type CommonStructsTransfer struct {
	TokenAddress common.Address
	ToAddress    common.Address
	Amount       *big.Int
}

// CheckPoWBlockPoW is an auto generated low-level Go binding around an user-defined struct.
type CheckPoWBlockPoW struct {
	P1                    []byte
	PrevHashOrReceiptRoot [32]byte
	P2                    []byte
	Difficulty            []byte
	P3                    []byte
}

// CheckPoABlockPoA is an auto generated low-level Go binding around an user-defined struct.
type CheckPoABlockPoA struct {
	P0Seal []byte // header prefix when encoded with seal
	P0Bare []byte // header prefix when encoded without seal

	// common (for bare and seal headers) part
	P1          []byte // bytes after header prefix and before ParentHash (de facto ParentHash prefix)
	ParentHash  [32]byte
	P2          []byte // bytes between ParentHash and ReceiptHash
	ReceiptHash [32]byte
	P3          []byte // bytes after ReceiptHash and before seal part

	// seal part
	S1        []byte // step prefix
	Step      []byte
	S2        []byte // signature prefix
	Signature []byte
}

type ReceiptsProof [][]byte
