package fee_api

import (
	"math/big"
)

// amountUsd / priceUsd
func Usd2Coin(amountUsd *big.Float, priceUsd float64) *big.Int {
	resFloat := new(big.Float).Quo(
		amountUsd,
		big.NewFloat(priceUsd),
	)

	res, _ := resFloat.Int(nil)
	return res
}

// amountWei * priceUsd
func Coin2Usd(amountWei *big.Int, priceUsd float64) *big.Float {
	return new(big.Float).Mul(
		new(big.Float).SetInt(amountWei),
		big.NewFloat(priceUsd),
	)
}

// amount * bps / 10_000
func calcBps(amount *big.Float, bps int64) *big.Float {
	return new(big.Float).Quo(
		new(big.Float).Mul(amount, big.NewFloat(float64(bps))),
		big.NewFloat(10_000),
	)
}
