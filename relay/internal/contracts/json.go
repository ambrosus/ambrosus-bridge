package contracts

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// we can't use hexutil.Big instead of big.Int in that structs fields because of dirty solution
// in geth abi bindings code, that has a separate check for handling exactly big.Int type ¯\_(ツ)_/¯

func (t *CommonStructsTransfer) MarshalJSON() ([]byte, error) {
	type Transfer struct {
		TokenAddress common.Address
		ToAddress    common.Address
		Amount       *hexutil.Big
	}
	tm := Transfer{t.TokenAddress, t.ToAddress, (*hexutil.Big)(t.Amount)}
	return json.Marshal(&tm)
}

// todo maybe create `type ReceiptProof []hexutil.Bytes`

func (t *CommonStructsTransferProof) MarshalJSON() ([]byte, error) {
	type TransferProof struct {
		ReceiptProof []hexutil.Bytes
		EventId      *hexutil.Big
		Transfers    []CommonStructsTransfer
	}
	rp := make([]hexutil.Bytes, len(t.ReceiptProof))
	for _, i := range t.ReceiptProof {
		rp = append(rp, i)
	}
	tm := TransferProof{rp, (*hexutil.Big)(t.EventId), t.Transfers}
	return json.Marshal(&tm)
}
