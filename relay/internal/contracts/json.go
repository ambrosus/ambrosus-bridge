package contracts

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// common

func (t *CommonStructsTransfer) MarshalJSON() ([]byte, error) {
	type Transfer struct {
		TokenAddress common.Address `json:"tokenAddress"`
		ToAddress    common.Address `json:"toAddress"`
		Amount       *hexutil.Big   `json:"amount"`
	}
	tm := Transfer{t.TokenAddress, t.ToAddress, (*hexutil.Big)(t.Amount)}
	return json.Marshal(&tm)
}

func (t *CommonStructsTransferProof) MarshalJSON() ([]byte, error) {
	type TransferProof struct {
		ReceiptProof []hexutil.Bytes         `json:"receiptProof"`
		EventId      *hexutil.Big            `json:"eventId"`
		Transfers    []CommonStructsTransfer `json:"transfers"`
	}
	rp := make([]hexutil.Bytes, len(t.ReceiptProof))
	for i, v := range t.ReceiptProof {
		rp[i] = v
	}
	tm := TransferProof{rp, (*hexutil.Big)(t.EventId), t.Transfers}
	return json.Marshal(&tm)
}

// AURA

func (t *CheckAuraAuraProof) MarshalJSON() ([]byte, error) {
	type CheckAuraAuraProof struct {
		Blocks             []CheckAuraBlockAura         `json:"blocks"`
		Transfer           CommonStructsTransferProof   `json:"transfer"`
		VsChanges          []CheckAuraValidatorSetProof `json:"vsChanges"`
		TransferEventBlock uint64                       `json:"transferEventBlock"`
	}
	tm := CheckAuraAuraProof{t.Blocks, t.Transfer, t.VsChanges, t.TransferEventBlock}
	return json.Marshal(&tm)
}

func (t *CheckAuraBlockAura) MarshalJSON() ([]byte, error) {
	type AuraBlockAura struct {
		P0Seal hexutil.Bytes `json:"p0Seal"`
		P0Bare hexutil.Bytes `json:"p0Bare"`

		// common (for bare and seal headers) part
		ParentHash  hexutil.Bytes `json:"parentHash"`
		P2          hexutil.Bytes `json:"p2"`
		ReceiptHash hexutil.Bytes `json:"receiptHash"`
		P3          hexutil.Bytes `json:"p3"`

		// seal part
		Step      hexutil.Bytes `json:"step"`
		Signature hexutil.Bytes `json:"signature"`

		FinalizedVs uint64 `json:"finalizedVs"`
	}
	tm := AuraBlockAura{
		t.P0Seal[:], t.P0Bare[:],
		t.ParentHash[:], t.P2, t.ReceiptHash[:], t.P3,
		t.Step[:], t.Signature,
		t.FinalizedVs,
	}
	return json.Marshal(&tm)
}

func (t *CheckAuraValidatorSetProof) MarshalJSON() ([]byte, error) {
	type ValidatorSetProof struct {
		ReceiptProof []hexutil.Bytes `json:"receiptProof"`
		DeltaAddress common.Address  `json:"deltaAddress"`
		DeltaIndex   int64           `json:"deltaIndex"`
	}
	rp := make([]hexutil.Bytes, len(t.ReceiptProof))
	for i, v := range t.ReceiptProof {
		rp[i] = v
	}
	tm := ValidatorSetProof{rp, t.DeltaAddress, t.DeltaIndex}
	return json.Marshal(&tm)
}

// POW

func (t *CheckPoWPoWProof) MarshalJSON() ([]byte, error) {
	type CheckPoWPoWProof struct {
		Blocks   []CheckPoWBlockPoW         `json:"blocks"`
		Transfer CommonStructsTransferProof `json:"transfer"`
	}
	tm := CheckPoWPoWProof{t.Blocks, t.Transfer}
	return json.Marshal(&tm)

}

func (t *CheckPoWBlockPoW) MarshalJSON() ([]byte, error) {
	type PoWBlockPoW struct {
		P0WithNonce    hexutil.Bytes `json:"p0WithNonce"`
		P0WithoutNonce hexutil.Bytes `json:"p0WithoutNonce"`

		P1                  hexutil.Bytes `json:"p1"`
		ParentOrReceiptHash hexutil.Bytes `json:"parentOrReceiptHash"`
		P2                  hexutil.Bytes `json:"p2"`
		Difficulty          hexutil.Bytes `json:"difficulty"`
		P3                  hexutil.Bytes `json:"p3"`
		Number              hexutil.Bytes `json:"number"`
		P4                  hexutil.Bytes `json:"p4"` // end when extra end

		P5    hexutil.Bytes `json:"p5"` // after extra
		Nonce hexutil.Bytes `json:"nonce"`

		P6 hexutil.Bytes `json:"p6"`

		DataSetLookup    []*hexutil.Big `json:"dataSetLookup"`
		WitnessForLookup []*hexutil.Big `json:"witnessForLookup"`
	}
	dslookup := make([]*hexutil.Big, len(t.DataSetLookup))
	wflookup := make([]*hexutil.Big, len(t.WitnessForLookup))
	for i, v := range t.DataSetLookup {
		dslookup[i] = (*hexutil.Big)(v)
	}
	for i, v := range t.WitnessForLookup {
		wflookup[i] = (*hexutil.Big)(v)
	}

	tm := PoWBlockPoW{
		t.P0WithNonce[:], t.P0WithoutNonce[:],
		t.P1, t.ParentOrReceiptHash[:], t.P2,
		t.Difficulty, t.P3, t.Number, t.P4,
		t.P5, t.Nonce, t.P6,
		dslookup, wflookup,
	}
	return json.Marshal(&tm)
}

// PoSA

func (t *CheckPoSAPoSAProof) MarshalJSON() ([]byte, error) {
	type PoSAPoSAProof struct {
		Blocks             []CheckPoSABlockPoSA       `json:"blocks"`
		Transfer           CommonStructsTransferProof `json:"transfer"`
		TransferEventBlock uint64                     `json:"transferEventBlock"`
	}
	tm := PoSAPoSAProof{t.Blocks, t.Transfer, t.TransferEventBlock}
	return json.Marshal(&tm)

}

func (t *CheckPoSABlockPoSA) MarshalJSON() ([]byte, error) {
	type PoSABlockPoSA struct {
		P0Signed   hexutil.Bytes `json:"p0Signed"`
		P0Unsigned hexutil.Bytes `json:"p0Unsigned"`

		ParentHash  hexutil.Bytes `json:"parentHash"`
		P1          hexutil.Bytes `json:"p1"`
		ReceiptHash hexutil.Bytes `json:"receiptHash"`
		P2          hexutil.Bytes `json:"p2"`
		Number      hexutil.Bytes `json:"number"`
		P3          hexutil.Bytes `json:"p3"`

		P4Signed   hexutil.Bytes `json:"p4Signed"`
		P4Unsigned hexutil.Bytes `json:"p4Unsigned"`
		ExtraData  hexutil.Bytes `json:"extraData"`

		P5 hexutil.Bytes `json:"p5"`
	}

	tm := PoSABlockPoSA{
		t.P0Signed[:], t.P0Unsigned[:],
		t.ParentHash[:], t.P1, t.ReceiptHash[:], t.P2, t.Number, t.P3,
		t.P4Signed, t.P4Unsigned, t.ExtraData,
		t.P5,
	}
	return json.Marshal(&tm)
}
