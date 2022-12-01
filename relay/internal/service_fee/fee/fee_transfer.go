package fee

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// getTransferFee return amount in native coin of this network
// transfer fee for this bridge is
// submit/unlock txs in side relay wallet +
// trigger txs by this relay wallet
func (p *Fee) getTransferFee(bridge BridgeFeeApi, thisCoinPrice, sideCoinPrice decimal.Decimal) (decimal.Decimal, error) {
	feeThis, feeSide, err := bridge.GetTransferFee()
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("GetTransferFee: %w", err)
	}
	feeUsd := coin2Usd(feeThis, thisCoinPrice).Add(coin2Usd(feeSide, sideCoinPrice))

	minTransferFeeUsd := bridge.GetMinTransferFee()
	if feeUsd.LessThan(minTransferFeeUsd) {
		feeUsd = minTransferFeeUsd
	}

	// convert it to native bridge currency
	feeThisNative := usd2Coin(feeUsd, thisCoinPrice)
	return feeThisNative, nil
}
