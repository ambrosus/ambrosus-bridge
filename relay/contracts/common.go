package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransferProof interface {
	EventId() *big.Int
}

// CommonStructsTransfer is an auto generated low-level Go binding around an user-defined struct.
type CommonStructsTransfer struct {
	TokenAddress common.Address
	ToAddress    common.Address
	Amount       *big.Int
}

// CommonStructsTransferProof is an auto generated low-level Go binding around an user-defined struct.
type CommonStructsTransferProof struct {
	ReceiptProof [][]byte
	EventId      *big.Int
	Transfers    []CommonStructsTransfer
}

// CheckAuraBlockAura is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraBlockAura struct {
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

	Type int64
}

// CheckAuraValidatorSetProof is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraValidatorSetProof struct {
	ReceiptProof [][]byte
	DeltaAddress common.Address
	DeltaIndex   int64
}

// CheckAuraAuraProof is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraAuraProof struct {
	Blocks    []*CheckAuraBlockAura
	Transfer  *CommonStructsTransferProof
	VsChanges []*CheckAuraValidatorSetProof
}

// CheckPoWBlockPoW is an auto generated low-level Go binding around an user-defined struct.
type CheckPoWBlockPoW struct {
	P0WithNonce    []byte
	P0WithoutNonce []byte

	P1                  []byte
	ParentOrReceiptHash [32]byte
	P2                  []byte
	Difficulty          []byte
	P3                  []byte
	Number              []byte
	P4                  []byte // end when extra end

	P5    []byte // after extra
	Nonce []byte

	P6 []byte

	DataSetLookUp    []*big.Int
	WitnessForLookUp []*big.Int
}

// CheckPoWPoWProof is an auto generated low-level Go binding around an user-defined struct.
type CheckPoWPoWProof struct {
	Blocks   []*CheckPoWBlockPoW
	Transfer *CommonStructsTransferProof
}

func (p *CheckAuraAuraProof) EventId() *big.Int {
	return p.Transfer.EventId
}

func (p *CheckPoWPoWProof) EventId() *big.Int {
	return p.Transfer.EventId
}
