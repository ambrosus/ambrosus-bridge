package contracts

import (
	"encoding/json"
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

func (t *CommonStructsTransfer) MarshalJSON() ([]byte, error) {
	type Transfer struct {
		TokenAddress common.Address `json:"tokenAddress"`
		ToAddress    common.Address `json:"toAddress"`
		Amount       *hexutil.Big   `json:"amount"`
	}
	tm := Transfer{t.TokenAddress, t.ToAddress, (*hexutil.Big)(t.Amount)}
	return json.Marshal(&tm)
}

// todo maybe create `type ReceiptProof []hexutil.Bytes`

func (t *CommonStructsTransferProof) MarshalJSON() ([]byte, error) {
	type TransferProof struct {
		ReceiptProof []hexutil.Bytes         `json:"receipt_proof"`
		EventId      *hexutil.Big            `json:"event_id"`
		Transfers    []CommonStructsTransfer `json:"transfers"`
	}
	rp := make([]hexutil.Bytes, len(t.ReceiptProof))
	for i, v := range t.ReceiptProof {
		rp[i] = v
	}
	tm := TransferProof{rp, (*hexutil.Big)(t.EventId), t.Transfers}
	return json.Marshal(&tm)
}

func (t *CheckAuraBlockAura) MarshalJSON() ([]byte, error) {
	type AuraBlockAura struct {
		P0Seal hexutil.Bytes `json:"p0_seal"`
		P0Bare hexutil.Bytes `json:"p0_bare"`

		// common (for bare and seal headers) part
		P1          hexutil.Bytes `json:"p1"`
		ParentHash  hexutil.Bytes `json:"parent_hash"`
		P2          hexutil.Bytes `json:"p2"`
		ReceiptHash hexutil.Bytes `json:"receipt_hash"`
		P3          hexutil.Bytes `json:"p3"`

		// seal part
		S1        hexutil.Bytes `json:"s1"`
		Step      hexutil.Bytes `json:"step"`
		S2        hexutil.Bytes `json:"s2"`
		Signature hexutil.Bytes `json:"signature"`

		Type       uint8 `json:"type_"`
		DeltaIndex int64 `json:"delta_index"`
	}
	tm := AuraBlockAura{
		t.P0Seal, t.P0Bare,
		t.P1, t.ParentHash[:], t.P2, t.ReceiptHash[:], t.P3,
		t.S1, t.Step, t.S2, t.Signature,
		t.Type, t.DeltaIndex,
	}
	return json.Marshal(&tm)
}

func (t *CheckAuraValidatorSetProof) MarshalJSON() ([]byte, error) {
	type ValidatorSetProof struct {
		ReceiptProof []hexutil.Bytes `json:"receipt_proof"`
		DeltaAddress common.Address  `json:"delta_address"`
		DeltaIndex   int64           `json:"delta_index"`
	}
	rp := make([]hexutil.Bytes, len(t.ReceiptProof))
	for i, v := range t.ReceiptProof {
		rp[i] = v
	}
	tm := ValidatorSetProof{rp, t.DeltaAddress, t.DeltaIndex}
	return json.Marshal(&tm)
}
