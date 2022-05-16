package fee_api

import (
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

const signatureFeeTimestamp = 30 * 60 // 30 minutes

var percentFromAmount = map[uint64]int64{
	0:       5 * 100, // 0..100_000$ => 5%
	100_000: 2 * 100, // 100_000...$ => 2%
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
	var req reqParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := p.getFees(req)
	if err != nil {
		http.Error(w, string(err.Marshal()), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type reqParams struct {
	TokenAddress common.Address `json:"tokenAddress"`
	IsAmb        bool           `json:"isAmb"`
	Amount       *hexutil.Big   `json:"amount"`
}

type Result struct {
	BridgeFee   *big.Int      `json:"bridge_fee"`
	TransferFee *big.Int      `json:"transfer_fee"`
	Signature   hexutil.Bytes `json:"signature"`
}

func (p *FeeAPI) getFees(req reqParams) (*Result, *AppError) {
	bridge := p.amb
	if !req.IsAmb {
		bridge = p.side
	}

	// get the bridge fee
	bridgeFee, err := getBridgeFee(bridge, req.TokenAddress, (*big.Int)(req.Amount))
	if err != nil {
		return nil, NewAppError(nil, "error when getting bridge fee", err.Error())
	}

	// get the transfer fee
	transferFee, err := bridge.GetTransferFee()
	if err != nil {
		return nil, NewAppError(nil, "error when getting transfer fee", err.Error())
	}

	// sign the price with private key
	// signature, err := signData(pk, tokenPrice, tokenAddress)
	signature, err := p.Sign(req.TokenAddress, transferFee, bridgeFee)
	if err != nil {
		return nil, NewAppError(nil, "error when signing data", err.Error())
	}

	return &Result{
		BridgeFee:   bridgeFee,
		TransferFee: transferFee,
		Signature:   signature,
	}, nil
}

func getBridgeFee(bridge networks.BridgeFeeApi, tokenAddress common.Address, amount *big.Int) (*big.Int, error) {
	// get token price
	tokenToUsdtPrice := big.NewInt(0) // todo
	tokensInUsdt := new(big.Int).Mul(amount, tokenToUsdtPrice)

	// use lower percent for higher amount
	var percent int64
	for minUsdt, percent_ := range percentFromAmount {
		if tokensInUsdt.Uint64() < minUsdt {
			break
		}
		percent = percent_
	}

	// calc fee
	usdtFee := calcBps(tokensInUsdt, percent)

	// convert usdt to native token
	return bridge.UsdtToNative(usdtFee)
}

func calcBps(amount *big.Int, bps int64) *big.Int {
	// amount * bps / 10_000
	return new(big.Int).Div(new(big.Int).Mul(amount, big.NewInt(bps)), big.NewInt(10_000))
}
