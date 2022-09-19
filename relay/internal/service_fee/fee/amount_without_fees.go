package fee

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func possibleAmountWithoutFees(amount, tokenUsdPrice, thisCoinPrice, transferFee, minBridgeFee decimal.Decimal) (decimal.Decimal, error) {
	transferFeeUsd := coin2Usd(transferFee, thisCoinPrice)
	amountUsd := coin2Usd(amount, tokenUsdPrice)
	amountUsd = possibleAmountWithoutFeesUsd(amountUsd, transferFeeUsd, minBridgeFee)
	amount = usd2Coin(amountUsd, tokenUsdPrice)
	if amount.Cmp(decimal.New(0, 0)) <= 0 {
		return amount, fmt.Errorf("amount is too small")
	}
	return amount, nil
}

func possibleAmountWithoutFeesUsd(amountWithBridgeFee, transferFee, minBridgeFee decimal.Decimal) decimal.Decimal {
	feePercentWithBridgeFee := getFeePercent(amountWithBridgeFee)
	if shouldUseMinBridgeFee(amountWithBridgeFee, feePercentWithBridgeFee, minBridgeFee) {
		return getAmountWithoutFees(amountWithBridgeFee, minBridgeFee, transferFee)
	}

	amountWithoutBridgeFee := getAmountWithoutBridgeFeeRatio(amountWithBridgeFee, feePercentWithBridgeFee)
	feePercentWithoutBridgeFee := getFeePercent(amountWithoutBridgeFee)
	if shouldUseMinBridgeFee(amountWithoutBridgeFee, feePercentWithoutBridgeFee, minBridgeFee) {
		return getAmountWithoutFees(amountWithBridgeFee, minBridgeFee, transferFee)
	}

	// if fee percent of amount without bridge fee is different from the old one, then recalculate with the new one
	if feePercentWithoutBridgeFee != feePercentWithBridgeFee {
		amountWithoutBridgeFee = getAmountWithoutBridgeFeeRatio(amountWithBridgeFee, feePercentWithoutBridgeFee)
	}
	return getAmountWithoutTransferFee(amountWithoutBridgeFee, transferFee)
}

func shouldUseMinBridgeFee(amount decimal.Decimal, feePercent int64, minBridgeFee decimal.Decimal) bool {
	return calcBps(amount, feePercent).Cmp(minBridgeFee) < 0
}

func getAmountWithoutFees(amount decimal.Decimal, bridgeFee decimal.Decimal, transferFee decimal.Decimal) decimal.Decimal {
	return amount.Sub(bridgeFee).Sub(transferFee)
}

// "ratio" because 10/1.5=6.6667, not 5
func getAmountWithoutBridgeFeeRatio(amount decimal.Decimal, feePercent int64) decimal.Decimal {
	return amount.Div(decimal.NewFromInt(feePercent).Div(decimal.NewFromInt(10_000)).Add(decimal.NewFromInt(1)))
}

func getAmountWithoutTransferFee(amount, transferFee decimal.Decimal) decimal.Decimal {
	return amount.Sub(transferFee)
}
