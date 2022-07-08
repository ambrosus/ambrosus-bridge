package fee

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kofalt/go-memoize"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

const (
	signatureFeeTimestamp = 30 * 60 // 30 minutes
	cacheExpiration       = time.Minute * 10
)

type BridgeFeeApi interface {
	networks.Bridge
	Sign(message []byte) ([]byte, error)
	GetTransferFee() *big.Int
	GetWrapperAddress() common.Address
	GetMinBridgeFee() decimal.Decimal // GetMinBridgeFee returns the minimal bridge fee that can be used
	GetDefaultTransferFee() decimal.Decimal
}

type Fee struct {
	amb, side BridgeFeeApi

	cache  *memoize.Memoizer
	Logger *zerolog.Logger
}

func NewFee(amb, side BridgeFeeApi, logger zerolog.Logger) *Fee {
	return &Fee{
		amb:    amb,
		side:   side,
		cache:  memoize.NewMemoizer(cacheExpiration, time.Hour),
		Logger: &logger,
	}
}

func (p *Fee) GetFees(tokenAddress common.Address, reqAmount *big.Int, isAmb, isAmountWithFees bool) (
	bridgeFeeBigInt, transferFeeBigInt, amountBigInt *big.Int, signature []byte, err error) {

	bridge, sideBridge := p.getBridges(isAmb)

	// if token address is native, then it's "wrapWithdraw" and we need to work with "wrapperAddress"
	if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		tokenAddress = bridge.GetWrapperAddress()
	}

	// get fees
	bridgeFee, transferFee, amount, err := p.getFees(bridge, sideBridge, tokenAddress, decimal.NewFromBigInt(reqAmount, 0), isAmountWithFees)
	if err != nil {
		err = fmt.Errorf("error when signing data: %w", err)
	}
	bridge.GetLogger().Debug().Msgf("bridgeFee: %s, transferFee: %s, amount: %s", bridgeFee.String(), transferFee.String(), amount.String())

	// sign fees with re;ay private key
	message := buildMessage(tokenAddress, transferFeeBigInt, bridgeFeeBigInt, amountBigInt)
	signature, err = bridge.Sign(message)
	if err != nil {
		err = fmt.Errorf("error when signing data: %w", err)
	}

	return bridgeFeeBigInt, transferFeeBigInt, amountBigInt, signature, err
}

func (p *Fee) getFees(bridge, sideBridge BridgeFeeApi, tokenAddress common.Address, amount decimal.Decimal, isAmountWithFees bool) (bridgeFeeBigInt, transferFeeBigInt, amountBigInt *big.Int, err error) {
	// get coin prices of this and side bridges
	thisCoinPrice, sideCoinPrice, tokenUsdPrice, err := p.getPrices(bridge, sideBridge, tokenAddress)
	if err != nil {
		return
	}

	// get transfer fee
	transferFee, err := p.getTransferFee(bridge, thisCoinPrice, sideCoinPrice)
	if err != nil {
		return
	}

	bridge.GetLogger().Debug().Msgf("thisCoinPrice: %s, sideCoinPrice: %s, tokenUsdPrice: %s, transferFee: %s",
		thisCoinPrice.String(), sideCoinPrice.String(), tokenUsdPrice.String(), transferFee.String())

	bridgeFee, amount, err := getBridgeFeeAndAmount(
		amount,
		isAmountWithFees,
		tokenUsdPrice,
		thisCoinPrice,
		transferFee,
		bridge.GetMinBridgeFee(),
	)
	if err != nil {
		return
	}

	return bridgeFee.BigInt(), transferFee.BigInt(), amount.BigInt(), nil
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

func (p *Fee) getTransferFee(bridge BridgeFeeApi, thisCoinPrice, sideCoinPrice decimal.Decimal) (decimal.Decimal, error) {
	feeSideNativeI, err, _ := p.cache.Memoize("GetTransferFee"+bridge.GetName(), func() (interface{}, error) {
		return bridge.GetTransferFee(), nil
	})
	if err != nil {
		return decimal.Decimal{}, err // todo
	}
	feeSideNative := feeSideNativeI.(*big.Int)

	if feeSideNative == nil {
		return bridge.GetDefaultTransferFee(), nil
	}

	// convert it to native bridge currency
	feeUsd := coin2Usd(decimal.NewFromBigInt(feeSideNative, 0), sideCoinPrice)
	feeThisNative := usd2Coin(feeUsd, thisCoinPrice)
	return feeThisNative, nil
}

func (p *Fee) getBridges(isAmb bool) (bridge, sideBridge BridgeFeeApi) {
	if isAmb {
		return p.amb, p.side
	}
	return p.side, p.amb
}
