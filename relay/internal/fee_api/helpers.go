package fee_api

import (
	"math"
	"math/big"
)

// coinsInUsdt / coinPriceInUsd * 1e18
func Usd2Coin(coinAmountInUsd *big.Float, coinPriceInUsd float64, coinDecimals uint8) *big.Int {
	resFloat := new(big.Float)
	resFloat.Quo(
		coinAmountInUsd,
		big.NewFloat(coinPriceInUsd),
	)
	resFloat.Mul(
		resFloat,
		big.NewFloat(math.Pow10(int(coinDecimals))),
	)

	res, _ := resFloat.Int(nil)
	return res
}

// coinPriceInUsd * (amount / 10^coinDecimals)
func Coin2Usd(coinsWei *big.Int, coinPriceInUsd float64, coinDecimals uint8) *big.Float {
	return new(big.Float).Mul(
		big.NewFloat(coinPriceInUsd),
		new(big.Float).Quo(
			new(big.Float).SetInt(coinsWei),
			big.NewFloat(math.Pow10(int(coinDecimals))),
		),
	)
}

// amount * bps / 10_000
func calcBps(amount *big.Float, bps int64) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).Mul(amount, big.NewFloat(float64(bps))),
		big.NewFloat(10_000),
	)
}
