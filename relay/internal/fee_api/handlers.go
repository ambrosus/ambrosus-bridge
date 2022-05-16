package fee_api

import (
	"encoding/json"
	"math/big"
	"net/http"
)

type Result struct {
	TransferFee *big.Int `json:"transfer_fee"`
	BridgeFee   *big.Int `json:"bridge_fee"`
	Signature   []byte   `json:"signature"`
}

/*
 * accepts a token address of this net, gets side token address from the contract (maybe accepts dev test or prod net and "amb" or "eth")
 * accepts an amount of tokens
 *
 * take the price from Uniswap
 * multiply that by the amount of tokens and by fee multiplier ("bridge fee")
 * get the average price for gas ("transfer fee")
 * sign that with time divided by some const
 *
 * respond the "bridge fee" and "transfer fee"
 */

func (p *FeeAPI) feesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get the token address from the query string
	tokenAddress := r.URL.Query().Get("token")
	if tokenAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(ErrTokenAddressNotPassed.Marshal())
		return
	}

	// get the bridge fee
	bridgeFee, err := p.GetBridgeFee(tokenAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewAppError(nil, "error when getting bridge fee", err.Error()).Marshal())
		return
	}

	// get the transfer fee
	transferFee, err := p.GetTransferFee()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewAppError(nil, "error when getting transfer fee", err.Error()).Marshal())
		return
	}

	// sign the price with private key
	// signature, err := signData(pk, tokenPrice, tokenAddress)
	signature, err := p.Sign(tokenAddress, transferFee, bridgeFee)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(NewAppError(nil, "error when signing data", err.Error()).Marshal())
		return
	}

	json.NewEncoder(w).Encode(Result{
		TransferFee: transferFee,
		BridgeFee:   bridgeFee,
		Signature:   signature,
	})
}
