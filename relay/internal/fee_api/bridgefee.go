package fee_api

import (
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/common"
)

var percentFromAmount = map[uint64]int64{
	0:       5 * 100, // 0..100_000$ => 5%
	100_000: 2 * 100, // 100_000...$ => 2%
}

func (p *FeeAPI) getBridgeFee(bridge networks.BridgeFeeApi, nativeUsdPrice float64, tokenAddress common.Address, amount *big.Int) (*big.Int, error) {
	// get token price
	tokenUsdPrice, err := p.getTokenPrice(bridge, tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("get token price: %w", err)
	}

	// get fee in usd
	amountUsd := Coin2Usd(amount, tokenUsdPrice)
	feeUsd := calcBps(amountUsd, getFeePercent(amountUsd))

	// if fee < minBridgeFee then use the minBridgeFee
	if minBridgeFee := bridge.GetMinBridgeFee(); feeUsd.Cmp(minBridgeFee) == -1 {
		feeUsd = minBridgeFee
	}

	// calc fee in native token
	feeNative := Usd2Coin(feeUsd, nativeUsdPrice)
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
