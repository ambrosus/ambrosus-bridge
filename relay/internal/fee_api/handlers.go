package fee_api

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
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
		helpers.JSONError(w, NewAppError(nil, "error when decoding request body", err.Error()).Marshal(), http.StatusBadRequest)
		return
	}

	result, err := p.getFees(req)
	if err != nil {
		helpers.JSONError(w, err.Marshal(), http.StatusInternalServerError)
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
	BridgeFee   *hexutil.Big  `json:"bridgeFee"`
	TransferFee *hexutil.Big  `json:"transferFee"`
	Signature   hexutil.Bytes `json:"signature"`
}

func (p *FeeAPI) getFees(req reqParams) (*Result, *AppError) {
	bridge := p.amb
	if !req.IsAmb {
		bridge = p.side
	}

	// get the bridge fee
	bridgeFee, err := getBridgeFeeStub(bridge, req.TokenAddress, (*big.Int)(req.Amount)) // TODO: replace with `getBridgeFee`
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
	message := buildMessage(req.TokenAddress, transferFee, bridgeFee)
	signature, err := bridge.Sign(message)
	if err != nil {
		return nil, NewAppError(nil, "error when signing data", err.Error())
	}

	return &Result{
		BridgeFee:   (*hexutil.Big)(bridgeFee),
		TransferFee: (*hexutil.Big)(transferFee),
		Signature:   signature,
	}, nil
}

func getBridgeFeeStub(bridge networks.BridgeFeeApi, tokenAddress common.Address, amount *big.Int) (*big.Int, error) {
	return big.NewInt(53000000000000000), nil // 0.053 ether
}

func getBridgeFee(bridge networks.BridgeFeeApi, tokenAddress common.Address, amount *big.Int) (*big.Int, error) {
	// get token price
	tokenToUsdtPrice := big.NewInt(0) // todo
	tokensInUsdt := new(big.Int).Mul(amount, tokenToUsdtPrice)

	// use lower percent for higher amount
	var percent int64
	for _, minUsdt := range helpers.SortedKeys(percentFromAmount) {
		percent_ := percentFromAmount[minUsdt]
		if tokensInUsdt.Uint64() < uint64(percent_) {
			break
		}
		percent = percent_
	}

	// calc fee
	usdtFee := calcBps(tokensInUsdt, percent)

	// convert usdt to native token
	nativeToUsdtPrice, err := bridge.CoinPrice()
	if err != nil {
		return nil, err
	}

	_, _ = usdtFee, nativeToUsdtPrice
	nativeFee := new(big.Int) // todo nativeFee = usdtFee / nativeToUsdtPrice

	return nativeFee, nil
}

func calcBps(amount *big.Int, bps int64) *big.Int {
	// amount * bps / 10_000
	return new(big.Int).Div(new(big.Int).Mul(amount, big.NewInt(bps)), big.NewInt(10_000))
}

func buildMessage(tokenAddress common.Address, transferFee *big.Int, bridgeFee *big.Int) []byte {
	// tokenAddress + timestamp + transferFee + bridgeFee

	var data bytes.Buffer
	var b32 [32]byte // solidity fills uint256 with zero

	data.Write(tokenAddress.Bytes())

	timestamp := time.Now().Unix() / signatureFeeTimestamp
	data.Write(big.NewInt(timestamp).FillBytes(b32[:]))

	data.Write(transferFee.FillBytes(b32[:]))
	data.Write(bridgeFee.FillBytes(b32[:]))

	return accounts.TextHash(crypto.Keccak256(data.Bytes()))
}
