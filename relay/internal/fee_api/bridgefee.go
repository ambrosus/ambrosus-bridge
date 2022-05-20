package fee_api

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/price"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/price_0x"
	"github.com/ethereum/go-ethereum/common"
)

var percentFromAmount = map[uint64]int64{
	0:       5 * 100, // 0..100_000$ => 5%
	100_000: 2 * 100, // 100_000...$ => 2%
}

func GetBridgeFee(bridge networks.BridgeFeeApi, tokenAddress common.Address, amount *big.Int) (*big.Int, error) {
	// get token symbol and decimals
	tokenSymbol, tokenDecimals, err := getTokenData(bridge, tokenAddress)
	if err != nil {
		return nil, fmt.Errorf("get token data: %w", err)
	}

	// get token price
	tokenToUsdtPrice, err := getTokenToUsdtPrice(tokenSymbol, tokenDecimals)
	if err != nil {
		return nil, fmt.Errorf("get token price: %w", err)
	}

	// get fee in usdt
	amountInUsdt := calcAmountInUsdt(amount, tokenToUsdtPrice, tokenDecimals)
	feePercent := getFeePercent(amountInUsdt)
	feeUsdt := calcBps(amountInUsdt, feePercent)

	// if fee < minBridgeFee then use the minBridgeFee
	if minBridgeFee := bridge.GetMinBridgeFee(); feeUsdt.Cmp(minBridgeFee) == -1 {
		feeUsdt = minBridgeFee
	}

	// convert usdt to native token
	nativeToUsdtPrice, err := bridge.CoinPrice()
	if err != nil {
		return nil, fmt.Errorf("get native to usdt price: %w", err)
	}

	// calc fee in native token
	nativeFee := calcNativeFee(feeUsdt, nativeToUsdtPrice, tokenDecimals)

	res, _ := nativeFee.Int(nil)
	return res, nil
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

func getTokenToUsdtPrice(tokenSymbol string, tokenDecimals uint8) (tokenToUsdtPrice float64, err error) {
	if tokenSymbol == "SAMB" {
		tokenToUsdtPrice, err = price.CoinToUsdt(price.Amb)
	} else {
		tokenToUsdtPrice, err = price_0x.CoinToUSDT(tokenSymbol, tokenDecimals)
	}
	return tokenToUsdtPrice, err
}

// tokenToUsdtPrice * (amount / 10^tokenDecimals)
func calcAmountInUsdt(amount *big.Int, tokenToUsdtPrice float64, tokenDecimals uint8) *big.Float {
	return new(big.Float).Mul(
		big.NewFloat(tokenToUsdtPrice),
		new(big.Float).Quo(
			new(big.Float).SetInt(amount),
			big.NewFloat(math.Pow10(int(tokenDecimals))),
		),
	)
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

// amount * bps / 10_000
func calcBps(amount *big.Float, bps int64) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).Mul(amount, big.NewFloat(float64(bps))),
		big.NewFloat(10_000),
	)
}

// usdtFeeWei / nativeToUsdtPrice * 10^tokenDecimals
func calcNativeFee(feeUsdt *big.Float, nativeToUsdtPrice float64, tokenDecimals uint8) *big.Float {
	nativeFee := new(big.Float)
	nativeFee.Quo(
		feeUsdt,
		big.NewFloat(nativeToUsdtPrice),
	)
	nativeFee.Mul(
		nativeFee,
		big.NewFloat(math.Pow10(int(tokenDecimals))),
	)
	return nativeFee
}
