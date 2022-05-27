package fee_api

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const signatureFeeTimestamp = 30 * 60 // 30 minutes

type reqParams struct {
	TokenAddress common.Address `json:"tokenAddress"`
	IsAmb        bool           `json:"isAmb"`
	Amount       *hexutil.Big   `json:"amount"`
}

type result struct {
	BridgeFee   *hexutil.Big  `json:"bridgeFee"`
	TransferFee *hexutil.Big  `json:"transferFee"`
	Signature   hexutil.Bytes `json:"signature"`
}

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

func (p *FeeAPI) getFees(req reqParams) (*result, *AppError) {
	bridge := p.amb
	if !req.IsAmb {
		bridge = p.side
	}

	// if token address is native, then it's "wrapWithdraw" and we need to work with "wrapperAddress"
	var tokenAddress = req.TokenAddress
	if req.TokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		var err error
		tokenAddress, err = bridge.GetWrapperAddress()
		if err != nil {
			return nil, NewAppError(nil, "error when getting wrapper address", err.Error())
		}
	}

	// get the bridge fee
	bridgeFee, err := GetBridgeFee(bridge, tokenAddress, (*big.Int)(req.Amount))
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
	message := buildMessage(tokenAddress, transferFee, bridgeFee)
	signature, err := bridge.Sign(message)
	if err != nil {
		return nil, NewAppError(nil, "error when signing data", err.Error())
	}

	return &result{
		BridgeFee:   (*hexutil.Big)(bridgeFee),
		TransferFee: (*hexutil.Big)(transferFee),
		Signature:   signature,
	}, nil
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
