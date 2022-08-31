package fee

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kofalt/go-memoize"
	"github.com/shopspring/decimal"
)

const (
	signatureFeeTimestamp = 30 * 60 // 30 minutes
	cacheExpiration       = time.Minute * 10
)

type BridgeFeeApi interface {
	networks.Bridge
	Sign(message []byte) ([]byte, error)
	GetTransferFee() (thisGas, sideGas decimal.Decimal)
	GetWrapperAddress() common.Address
	GetMinBridgeFee() decimal.Decimal   // GetMinBridgeFee returns the minimal bridge fee in usd
	GetMinTransferFee() decimal.Decimal // GetMinTransferFee returns the minimal transfer fee in usd
}

type Fee struct {
	amb, side   BridgeFeeApi
	priceGetter priceGetter
	cache       *memoize.Memoizer
}

func NewFee(amb, side BridgeFeeApi) *Fee {
	return &Fee{
		amb:         amb,
		side:        side,
		priceGetter: new(priceGetterS),
		cache:       memoize.NewMemoizer(cacheExpiration, time.Hour),
	}
}

func (p *Fee) GetFees(tokenAddress common.Address, reqAmount *big.Int, isAmb, isAmountWithFees bool) (
	bridgeFee, transferFee, amount *big.Int, signature []byte, err error) {

	bridge, sideBridge := p.getBridges(isAmb)

	// if token address is native, then it's "wrapWithdraw" and we need to work with "wrapperAddress"
	if tokenAddress == common.HexToAddress("0x0000000000000000000000000000000000000000") {
		tokenAddress = bridge.GetWrapperAddress()
	}

	// get fees
	bridgeFee, transferFee, amount, err = p.getFees(bridge, sideBridge, tokenAddress, decimal.NewFromBigInt(reqAmount, 0), isAmountWithFees)
	if err != nil {
		return
	}
	bridge.GetLogger().Debug().Msgf("bridgeFee: %s, transferFee: %s, amount: %s", bridgeFee.String(), transferFee.String(), amount.String())

	// sign fees with relay private key
	message := buildMessage(tokenAddress, transferFee, bridgeFee, amount)
	signature, err = bridge.Sign(message)
	if err != nil {
		err = fmt.Errorf("error when signing data: %w", err)
	}

	return bridgeFee, transferFee, amount, signature, err
}

func (p *Fee) getFees(bridge, sideBridge BridgeFeeApi, tokenAddress common.Address, amount decimal.Decimal, isAmountWithFees bool) (bridgeFeeBI, transferFeeBI, amountBI *big.Int, err error) {
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

	// if amount contains fees, then we need change the amount to the possible amount without fees (when transfer *max* native coins)
	if isAmountWithFees {
		amount, err = possibleAmountWithoutFees(amount, tokenUsdPrice, thisCoinPrice, transferFee, bridge.GetMinBridgeFee())
		if err != nil {
			return
		}
	}

	bridgeFee, err := getBridgeFee(thisCoinPrice, tokenUsdPrice, amount, bridge.GetMinBridgeFee())
	if err != nil {
		err = fmt.Errorf("error when getting bridge fee: %w", err)
		return
	}

	return bridgeFee.BigInt(), transferFee.BigInt(), amount.BigInt(), nil

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
