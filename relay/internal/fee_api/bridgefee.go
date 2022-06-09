package fee_api

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
)

var percentFromAmount = map[uint64]int64{
	0:       5 * 100, // 0..100_000$ => 5%
	100_000: 2 * 100, // 100_000...$ => 2%
}

func (p *FeeAPI) getBridgeFee(bridge BridgeFeeApi, nativeUsdPrice, tokenUsdPrice float64, amount *big.Int) (*big.Int, error) {
	// get fee in usd
	amountUsd := coin2Usd(amount, tokenUsdPrice)
	feeUsd := calcBps(amountUsd, getFeePercent(amountUsd))

	// if fee < minBridgeFee then use the minBridgeFee
	if minBridgeFee := bridge.GetMinBridgeFee(); feeUsd.Cmp(minBridgeFee) == -1 {
		feeUsd = minBridgeFee
	}

	// calc fee in native token
	feeNative := usd2Coin(feeUsd, nativeUsdPrice)
	return feeNative, nil
}

func getFeePercent(amountInUsdt *big.Float) (percent int64) {
	// use lower percent for higher amount
	for _, minUsdt := range helpers.SortedKeys(percentFromAmount) {
		percent_ := percentFromAmount[minUsdt]
		if amountInUsdt.Cmp(new(big.Float).SetUint64(minUsdt)) == -1 {
			break
		}
		percent = percent_
	}
	return percent
}
