package fee

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func possibleAmountWithoutFees(amount, tokenUsdPrice, transferFee, thisCoinPrice, minBridgeFee decimal.Decimal) (decimal.Decimal, error) {
	transferFeeUsd := coin2Usd(transferFee, thisCoinPrice)
	amountUsd := coin2Usd(amount, tokenUsdPrice)
	amountUsd = possibleAmountWithoutFeesUsd(amountUsd, transferFeeUsd, minBridgeFee)
	amount = usd2Coin(amountUsd, tokenUsdPrice)
	if amount.Cmp(decimal.New(0, 0)) <= 0 {
		return amount, fmt.Errorf("amount is too small")
	}
	return amount, nil
}

func possibleAmountWithoutFeesUsd(amountUsd1, transferFeeUsd, minBridgeFee decimal.Decimal) decimal.Decimal {
	// todo why float64(feePercent1+10_000) / 10_000) ??
	// todo add comments, i understand nothing :(

	// part 1
	feePercent1 := getFeePercent(amountUsd1)

	// if fee < minBridgeFee then use the minBridgeFee
	if calcBps(amountUsd1, feePercent1).Cmp(minBridgeFee) == -1 {
		// amountUsd1 - (transferFeeUsd + minBridgeFee)
		return amountUsd1.Sub(minBridgeFee).Sub(transferFeeUsd)
	}

	// part 2

	// (amountUsd1 - transferFeeUsd) / %
	amountUsd2 := amountUsd1.Div(decimal.NewFromFloat(float64(feePercent1)/10_000 + 1))
	feePercent2 := getFeePercent(amountUsd2)

	// if fee < minBridgeFee then use the minBridgeFee
	if calcBps(amountUsd2, feePercent2).Cmp(minBridgeFee) == -1 {
		// amountUsd1 - (transferFeeUsd + minBridgeFee)
		return amountUsd1.Sub(minBridgeFee).Sub(transferFeeUsd)
	}

	// part 3

	// if fee percent of new amount if different from the old one, then recalculate with the new one
	if feePercent2 != feePercent1 {
		amountUsd2 = amountUsd1.Div(decimal.NewFromFloat(float64(feePercent2)/10_000 + 1))
	}

	return amountUsd2.Sub(transferFeeUsd)
}
