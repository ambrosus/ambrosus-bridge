package contracts

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

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
		ParentHash  hexutil.Bytes `json:"parent_hash"`
		P2          hexutil.Bytes `json:"p2"`
		ReceiptHash hexutil.Bytes `json:"receipt_hash"`
		P3          hexutil.Bytes `json:"p3"`

		// seal part
		Step      hexutil.Bytes `json:"step"`
		Signature hexutil.Bytes `json:"signature"`

		Type       uint8 `json:"type_"`
		DeltaIndex int64 `json:"delta_index"`
	}
	tm := AuraBlockAura{
		t.P0Seal[:], t.P0Bare[:],
		t.ParentHash[:], t.P2, t.ReceiptHash[:], t.P3,
		t.Step, t.Signature,
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
		t.P1, t.ParentOrReceiptHash[:], t.P2, t.Difficulty[:], t.P3, t.Number[:], t.P4,
		t.P5, t.Nonce[:],
		t.P6,
		dslookup, wflookup,
	}
	return json.Marshal(&tm)
}
