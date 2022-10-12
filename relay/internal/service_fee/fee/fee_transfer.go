package fee

import (
	"github.com/shopspring/decimal"
)

// getTransferFee return amount in native coin of this network
// transfer fee for this bridge is
// submit/unlock txs in side relay wallet +
// trigger txs by this relay wallet
func (p *Fee) getTransferFee(bridge BridgeFeeApi, thisCoinPrice, sideCoinPrice decimal.Decimal) (decimal.Decimal, error) {
	feeI, err, _ := p.cache.Memoize("GetTransferFee"+bridge.GetName(), func() (interface{}, error) {
		this, side := bridge.GetTransferFee()
		return feeS{this, side}, nil
	})
	if err != nil {
		return decimal.Decimal{}, err // todo
	}
	fee := feeI.(feeS)
	feeUsd := coin2Usd(fee.this, thisCoinPrice).Add(coin2Usd(fee.side, sideCoinPrice))

	minTransferFeeUsd := bridge.GetMinTransferFee()
	if feeUsd.LessThan(minTransferFeeUsd) {
		feeUsd = minTransferFeeUsd
	}

	// convert it to native bridge currency
	feeThisNative := usd2Coin(feeUsd, thisCoinPrice)
	return feeThisNative, nil
}

type feeS struct {
	this, side decimal.Decimal
}
