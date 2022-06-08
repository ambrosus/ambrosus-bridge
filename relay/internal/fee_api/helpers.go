package fee_api

import (
	"fmt"
	"math/big"
	"net/http"
)

// amountUsd / priceUsd
func usd2Coin(amountUsd *big.Float, priceUsd float64) *big.Int {
	resFloat := new(big.Float).Quo(
		amountUsd,
		big.NewFloat(priceUsd),
	)

	res, _ := resFloat.Int(nil)
	return res
}

// amountWei * priceUsd
func coin2Usd(amountWei *big.Int, priceUsd float64) *big.Float {
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

func writeError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, err.Error())
}
