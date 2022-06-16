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
	"github.com/shopspring/decimal"
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
	Amount      *hexutil.Big  `json:"amount"`
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
	// get token price
	tokenUsdPrice, err := p.getTokenPrice(bridge, req.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("get token price: %w", err)
	}
	bridge.GetLogger().Debug().Msgf("thisCoinPrice: %s, sideCoinPrice: %s, tokenUsdPrice: %s", thisCoinPrice.String(), sideCoinPrice.String(), tokenUsdPrice.String())

	// get transfer fee
	transferFee, err := p.getTransferFee(bridge, thisCoinPrice, sideCoinPrice)
	if err != nil {
		return nil, fmt.Errorf("error when getting transfer fee: %w", err)
	}
	bridge.GetLogger().Debug().Msgf("transferFee: %s", transferFee.String())

	bridgeFee, amount, err := getBridgeFeeAndAmount(
		decimal.NewFromBigInt((*big.Int)(req.Amount), 0),
		req.IsAmountWithFees,
		tokenUsdPrice,
		thisCoinPrice,
		transferFee,
		bridge.GetMinBridgeFee(),
	)
	if err != nil {
		return nil, err
	}

	// make the fees as big int (throw away the decimal part)
	bridgeFeeBigInt := bridgeFee.BigInt()
	transferFeeBigInt := transferFee.BigInt()
	amountBigInt := amount.BigInt()
	bridge.GetLogger().Debug().Msgf("bridgeFeeBigInt: %s, transferFeeBigInt: %s, amountBigInt: %s", bridgeFeeBigInt.String(), transferFeeBigInt.String(), amountBigInt.String())
	bridge.GetLogger().Debug().Msgf("bridgeFeeHex: %s, transferFeeHex: %s, amountHex: %s", (*hexutil.Big)(bridgeFeeBigInt).String(), (*hexutil.Big)(transferFeeBigInt).String(), (*hexutil.Big)(amountBigInt).String())

	// sign the price with private key
	message := buildMessage(req.TokenAddress, transferFeeBigInt, bridgeFeeBigInt, amountBigInt)
	signature, err := bridge.Sign(message)
	if err != nil {
		err = fmt.Errorf("error when signing data: %w", err)
	}

	return &result{
		BridgeFee:   (*hexutil.Big)(bridgeFeeBigInt),
		TransferFee: (*hexutil.Big)(transferFeeBigInt),
		Amount:      (*hexutil.Big)(amountBigInt),
		Signature:   signature,
	}, err
}

func getBridgeFeeAndAmount(
	reqAmount decimal.Decimal,
	isAmountWithFees bool,
	tokenUsdPrice decimal.Decimal,
	thisCoinPrice decimal.Decimal,
	transferFee decimal.Decimal,
	minBridgeFee decimal.Decimal,
) (decimal.Decimal, decimal.Decimal, error) {
	amount := reqAmount.Copy()

	// if amount contains fees, then we need change the amount to the possible amount without fees (when transfer *max* native coins)
	if isAmountWithFees {
		amount = possibleAmountWithoutFees(amount, tokenUsdPrice, transferFee, thisCoinPrice, minBridgeFee)

		if amount.Cmp(decimal.New(0, 0)) <= 0 {
			return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("amount is too small")
		}

	}

	// get bridge fee
	bridgeFee, err := getBridgeFee(thisCoinPrice, tokenUsdPrice, amount, minBridgeFee)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("error when getting bridge fee: %w", err)
	}

	return bridgeFee, amount, nil
}

func possibleAmountWithoutFees(
	amount,
	tokenUsdPrice,
	transferFee,
	thisCoinPrice,
	minBridgeFee decimal.Decimal,
) decimal.Decimal {
	transferFeeUsd := coin2Usd(transferFee, thisCoinPrice)

	amountUsd := coin2Usd(amount, tokenUsdPrice)
	feePercent := getFeePercent(amountUsd)

	// if fee < minBridgeFee then use the minBridgeFee
	if calcBps(amountUsd, feePercent).Cmp(minBridgeFee) == -1 {
		// amountUsd - (transferFeeUsd + minBridgeFee)
		amountUsd = amountUsd.Sub(minBridgeFee.Add(transferFeeUsd))
		return usd2Coin(amountUsd, tokenUsdPrice)
	}

	// (amountUsd - transferFeeUsd) / %
	newAmountUsd := amountUsd.Div(decimal.NewFromFloat(float64(feePercent+10_000) / 10_000))

	// if fee < minBridgeFee then use the minBridgeFee
	if calcBps(newAmountUsd, getFeePercent(newAmountUsd)).Cmp(minBridgeFee) == -1 {
		// amountUsd - (transferFeeUsd + minBridgeFee)
		newAmountUsd = amountUsd.Sub(minBridgeFee.Add(transferFeeUsd))
		return usd2Coin(newAmountUsd, tokenUsdPrice)
	}
	// if fee percent of new amount if different from the old one, then recalculate with the new one
	if newFeePercent := getFeePercent(newAmountUsd); newFeePercent != feePercent {
		newAmountUsd = amountUsd.Div(decimal.NewFromFloat(float64(newFeePercent+10_000) / 10_000))
	}

	newAmountUsd = newAmountUsd.Sub(transferFeeUsd)
	return usd2Coin(newAmountUsd, tokenUsdPrice)
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

func (p *FeeAPI) getTransferFee(bridge BridgeFeeApi, thisCoinPrice, sideCoinPrice decimal.Decimal) (decimal.Decimal, error) {
	feeSideNativeI, err, _ := p.cache.Memoize("GetTransferFee", func() (interface{}, error) {
		return bridge.GetTransferFee(), nil
	})
	if err != nil {
		return decimal.Decimal{}, err // todo
	}
	feeSideNative := feeSideNativeI.(*big.Int)

	// convert it to native bridge currency
	feeUsd := coin2Usd(decimal.NewFromBigInt(feeSideNative, 0), sideCoinPrice)
	feeThisNative := usd2Coin(feeUsd, thisCoinPrice)
	return feeThisNative, nil
}
