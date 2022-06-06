package fee_api

import (
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/price"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/price_0x"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kofalt/go-memoize"
)

var percentFromAmount = map[uint64]int64{
	0:       5 * 100, // 0..100_000$ => 5%
	100_000: 2 * 100, // 100_000...$ => 2%
}

func GetBridgeFee(bridge networks.BridgeFeeApi, nativeCoinPriceInUsd float64, cache *memoize.Memoizer, tokenAddress common.Address, amount *big.Int) (*big.Int, error) {
	// get token symbol and decimals
	tokenSymbol, tokenDecimals, err := getTokenData(bridge, tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("get token data: %w", err)
	}

	// get token price
	tokenToUsdtPrice, err := getTokenToUsdtPrice(tokenSymbol, tokenDecimals, cache)
	if err != nil {
		return nil, fmt.Errorf("get token price: %w", err)
	}

	// get fee in usdt
	amountInUsdt := Coin2Usd(amount, tokenToUsdtPrice, tokenDecimals)
	feePercent := getFeePercent(amountInUsdt)
	feeUsdt := calcBps(amountInUsdt, feePercent)

	// if fee < minBridgeFee then use the minBridgeFee
	if minBridgeFee := bridge.GetMinBridgeFee(); feeUsdt.Cmp(minBridgeFee) == -1 {
		feeUsdt = minBridgeFee
	}

	// calc fee in native token
	nativeFee := Usd2Coin(feeUsdt, nativeCoinPriceInUsd, 18)
	return nativeFee, nil
}

func getTokenData(bridge networks.BridgeFeeApi, tokenAddress common.Address) (string, uint8, error) {
	tokenContract, err := contracts.NewToken(tokenAddress, bridge.GetClient())
	if err != nil {
		return "", 0, fmt.Errorf("get token contract: %w", err)
	}
	tokenSymbol, err := tokenContract.Symbol(nil)
	if err != nil {
		return "", 0, fmt.Errorf("get token symbol: %w", err)
	}
	tokenDecimals, err := tokenContract.Decimals(nil)
	if err != nil {
		return "", 0, fmt.Errorf("get token decimals: %w", err)
	}
	return tokenSymbol, tokenDecimals, nil
}

func getTokenToUsdtPrice(tokenSymbol string, tokenDecimals uint8, cache *memoize.Memoizer) (tokenToUsdtPrice float64, err error) {
	var res interface{}
	if tokenSymbol == "SAMB" {
		res, err, _ = cache.Memoize("SAMB", func() (interface{}, error) {
			return price.CoinToUsdt(price.Amb)
		})
	} else {
		res, err, _ = cache.Memoize(tokenSymbol, func() (interface{}, error) {
			return price_0x.CoinToUSDT(tokenSymbol, tokenDecimals)
		})
	}
	return res.(float64), err
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
