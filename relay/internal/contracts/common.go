package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

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
	P0Seal hexutil.Bytes // header prefix when encoded with seal
	P0Bare hexutil.Bytes // header prefix when encoded without seal

	// common (for bare and seal headers) part
	P1          hexutil.Bytes // bytes after header prefix and before ParentHash (de facto ParentHash prefix)
	ParentHash  common.Hash
	P2          hexutil.Bytes // bytes between ParentHash and ReceiptHash
	ReceiptHash common.Hash
	P3          hexutil.Bytes // bytes after ReceiptHash and before seal part

	// seal part
	S1        hexutil.Bytes // step prefix
	Step      hexutil.Bytes
	S2        hexutil.Bytes // signature prefix
	Signature hexutil.Bytes

	Type       uint8
	DeltaIndex int64
}

// CheckAuraValidatorSetProof is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraValidatorSetProof struct {
	ReceiptProof [][]byte
	DeltaAddress common.Address
	DeltaIndex   int64
}

// CheckAuraAuraProof is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraAuraProof struct {
	Blocks    []CheckAuraBlockAura
	Transfer  CommonStructsTransferProof
	VsChanges []CheckAuraValidatorSetProof
}

// CheckPoWBlockPoW is an auto generated low-level Go binding around an user-defined struct.
type CheckPoWBlockPoW struct {
	P0WithNonce    hexutil.Bytes
	P0WithoutNonce hexutil.Bytes

	P1                  hexutil.Bytes
	ParentOrReceiptHash common.Hash
	P2                  hexutil.Bytes
	Difficulty          hexutil.Bytes
	P3                  hexutil.Bytes
	Number              hexutil.Bytes
	P4                  hexutil.Bytes // end when extra end

	P5    hexutil.Bytes // after extra
	Nonce hexutil.Bytes

	P6 hexutil.Bytes

	DataSetLookup    []*big.Int
	WitnessForLookup []*big.Int
}

// CheckPoWPoWProof is an auto generated low-level Go binding around an user-defined struct.
type CheckPoWPoWProof struct {
	Blocks   []CheckPoWBlockPoW
	Transfer CommonStructsTransferProof
}

// CommonStructsConstructorArgs is an auto generated low-level Go binding around an user-defined struct.
type CommonStructsConstructorArgs struct {
	SideBridgeAddress  common.Address
	RelayAddress       common.Address
	TokenThisAddresses []common.Address
	TokenSideAddresses []common.Address
	Fee                *big.Int
	FeeRecipient       common.Address
	TimeframeSeconds   *big.Int
	LockTime           *big.Int
	MinSafetyBlocks    *big.Int
}
