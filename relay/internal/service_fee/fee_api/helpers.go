package fee_api

import (
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

// amountUsd / priceUsd
func usd2Coin(amountUsd decimal.Decimal, priceUsd decimal.Decimal) decimal.Decimal {
	return amountUsd.Div(priceUsd)
}

// amountWei * priceUsd
func coin2Usd(amountWei decimal.Decimal, priceUsd decimal.Decimal) decimal.Decimal {
	return amountWei.Mul(priceUsd)
}

// through usd prices
func Coin2coin(amountWei, firstCoinPriceUsd, secondCoinPriceUsd decimal.Decimal) decimal.Decimal {
	firstCoinInUsd := coin2Usd(amountWei, firstCoinPriceUsd)
	return usd2Coin(firstCoinInUsd, secondCoinPriceUsd)
}

// amount * bps / 10_000
func calcBps(amount decimal.Decimal, bps int64) decimal.Decimal {
	return amount.Mul(decimal.NewFromInt(bps)).Div(decimal.NewFromInt(10_000))
}

func writeError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, err.Error())
}
