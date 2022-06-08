package fee_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type reqParams struct {
	TokenAddress common.Address `json:"tokenAddress"`
	IsAmb        bool           `json:"isAmb"`

	Amount           *hexutil.Big `json:"amount"`
	IsAmountWithFees bool         `json:"IsAmountWithFees"`
}

type result struct {
	BridgeFee   *hexutil.Big  `json:"bridgeFee"`
	TransferFee *hexutil.Big  `json:"transferFee"`
	Signature   hexutil.Bytes `json:"signature"`
}

func (p *FeeAPI) feesHandler(w http.ResponseWriter, r *http.Request) {
	var req reqParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Errorf("error when decoding request body: %w", err))
		return
	}

	result, err := p.getFees(req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("error getting bridge_fee: %w", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (p *FeeAPI) getFees(req reqParams) (*result, error) {
	var bridge, sideBridge BridgeFeeApi
	if req.IsAmb {
		bridge, sideBridge = p.amb, p.side
	} else {
		bridge, sideBridge = p.side, p.amb
	}

	// if token address is native, then it's "wrapWithdraw" and we need to work with "wrapperAddress"
	if req.TokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		req.TokenAddress = bridge.GetWrapperAddress()
	}

	// get coin prices of this and side bridges for reusing it below
	thisCoinPrice, err := p.getTokenPrice(bridge, common.Address{})
	if err != nil {
		return nil, fmt.Errorf("error when getting this bridge coin price: %w", err)
	}
	sideCoinPrice, err := p.getTokenPrice(sideBridge, common.Address{})
	if err != nil {
		return nil, fmt.Errorf("error when getting side bridge coin price: %w", err)
	}

	// get fees
	bridgeFee, err := p.getBridgeFee(bridge, thisCoinPrice, req.TokenAddress, (*big.Int)(req.Amount))
	if err != nil {
		return nil, fmt.Errorf("error when getting bridge fee: %w", err)
	}
	transferFee, err := p.getTransferFee(bridge, thisCoinPrice, sideCoinPrice)
	if err != nil {
		return nil, fmt.Errorf("error when getting transfer fee: %w", err)
	}

	amount := (*big.Int)(req.Amount)
	// if amount contains fees, then we need to remove them (when transfer *max* native coins)
	if req.IsAmountWithFees {
		amount.Sub(amount, new(big.Int).Add(bridgeFee, transferFee))
		if amount.Cmp(big.NewInt(0)) <= 0 {
			return nil, fmt.Errorf("amount is too small")
		}
	}

	// sign the price with private key
	message := buildMessage(req.TokenAddress, transferFee, bridgeFee, (*big.Int)(req.Amount))
	signature, err := bridge.Sign(message)
	if err != nil {
		err = fmt.Errorf("error when signing data: %w", err)
	}

	return &result{
		BridgeFee:   (*hexutil.Big)(bridgeFee),
		TransferFee: (*hexutil.Big)(transferFee),
		Signature:   signature,
	}, err
}

func buildMessage(tokenAddress common.Address, transferFee, bridgeFee, amount *big.Int) []byte {
	var data bytes.Buffer
	var b32 [32]byte // solidity fills uint256 with zero

	data.Write(tokenAddress.Bytes())

	timestamp := time.Now().Unix() / signatureFeeTimestamp
	data.Write(big.NewInt(timestamp).FillBytes(b32[:]))

	data.Write(transferFee.FillBytes(b32[:]))
	data.Write(bridgeFee.FillBytes(b32[:]))
	data.Write(amount.FillBytes(b32[:]))

	return accounts.TextHash(crypto.Keccak256(data.Bytes()))
}

func (p *FeeAPI) getTransferFee(bridge BridgeFeeApi, thisCoinPrice, sideCoinPrice float64) (*big.Int, error) {
	feeSideNativeI, err, _ := p.cache.Memoize("GetTransferFee", func() (interface{}, error) {
		return bridge.GetTransferFee(), nil
	})
	if err != nil {
		return nil, err // todo
	}
	feeSideNative := feeSideNativeI.(*big.Int)

	// convert it to native bridge currency
	feeUsd := coin2Usd(feeSideNative, sideCoinPrice)
	feeThisNative := usd2Coin(feeUsd, thisCoinPrice)
	return feeThisNative, nil
}
