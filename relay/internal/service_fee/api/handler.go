package api

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type reqParams struct {
	TokenAddress common.Address `json:"tokenAddress"`
	IsAmb        bool           `json:"isAmb"`

	Amount           *hexutil.Big `json:"amount"`
	IsAmountWithFees bool         `json:"IsAmountWithFees"`
}

type response struct {
	BridgeFee   *hexutil.Big  `json:"bridgeFee"`
	TransferFee *hexutil.Big  `json:"transferFee"`
	Amount      *hexutil.Big  `json:"amount"`
	Signature   hexutil.Bytes `json:"signature"`
}

func (p *FeeAPI) feesHandler(w http.ResponseWriter, r *http.Request) {
	var req reqParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("error when decoding request body: %w", err))
		return
	}

	bridgeFee, transferFee, amount, signature, err := p.Service.GetFees(req.TokenAddress, (*big.Int)(req.Amount), req.IsAmb, req.IsAmountWithFees)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("error getting bridge_fee: %w", err))
		return
	}

	result := &response{
		BridgeFee:   (*hexutil.Big)(bridgeFee),
		TransferFee: (*hexutil.Big)(transferFee),
		Amount:      (*hexutil.Big)(amount),
		Signature:   signature,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func writeError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, err.Error())
}
