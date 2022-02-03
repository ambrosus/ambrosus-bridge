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

// EthTransferEvent represents a TransferEvent event raised by the Eth contract.
type EthTransferEvent struct {
	TransferEvent
}

// AmbTransferEvent represents a TransferEvent event raised by the Amb contract.
type AmbTransferEvent struct {
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
	P0Seal                []byte
	P0Bare                []byte
	P1                    []byte
	PrevHashOrReceiptRoot [32]byte
	P2                    []byte
	Step                  []byte
	S1                    []byte
	Signature             []byte
}

type ReceiptsProof [][]byte
