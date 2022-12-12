package bindings

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
