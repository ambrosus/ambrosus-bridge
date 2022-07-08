package fee

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
)

func decimalFromString(s string) decimal.Decimal {
	i, _ := decimal.NewFromString(s)
	return i
}

func Test_getBridgeAndAmount(t *testing.T) {
	type args struct {
		reqAmount        decimal.Decimal
		isAmountWithFees bool
		tokenUsdPrice    decimal.Decimal
		thisCoinPrice    decimal.Decimal
		transferFee      decimal.Decimal
		minBridgeFee     decimal.Decimal
	}
	tests := []struct {
		name          string
		args          args
		wantBridgeFee decimal.Decimal
		wantAmount    decimal.Decimal
		wantErr       bool
	}{
		{
			name: "$100K in AMB (should take 5%)",
			args: args{
				reqAmount:        decimalFromString("13848494031870104170397696"),
				isAmountWithFees: true,
				tokenUsdPrice:    decimal.NewFromFloat(0.00000000000000000000722114619610344),
				thisCoinPrice:    decimal.NewFromFloat(0.00000000000000000000722114619610344),
				transferFee:      decimalFromString("380000000000000000000"),
				minBridgeFee:     decimal.NewFromFloat(5),
			},
			wantBridgeFee: decimalFromString("659433096755719246204219.6827408804224156"),
			wantAmount:    decimalFromString("13188661935114384924195179.3913578384840331"),
			wantErr:       false,
		},
		{
			name: "$20 in AMB (should take $5)",
			args: args{
				reqAmount:        decimalFromString("3000000000000000000000"),
				isAmountWithFees: true,
				tokenUsdPrice:    decimal.NewFromFloat(0.00000000000000000000722114619610344),
				thisCoinPrice:    decimal.NewFromFloat(0.00000000000000000000722114619610344),
				transferFee:      decimalFromString("380000000000000000000"),
				minBridgeFee:     decimal.NewFromFloat(5),
			},
			wantBridgeFee: decimalFromString("692410853376437723256.799893475846839"),
			wantAmount:    decimalFromString("1927589146623562276743.200106524153161"),
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bridgeFee, amount, err := getBridgeFeeAndAmount(tt.args.reqAmount, tt.args.tokenUsdPrice, tt.args.thisCoinPrice, tt.args.transferFee, tt.args.minBridgeFee, tt.args.isAmountWithFees)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBridgeFeeAndAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if bridgeFee.String() != tt.wantBridgeFee.String() {
				t.Errorf("getBridgeFeeAndAmount() bridgeFee = %v, wantBridgeFee %v", bridgeFee, tt.wantBridgeFee)
			}
			if amount.String() != tt.wantAmount.String() {
				t.Errorf("getBridgeFeeAndAmount() amount = %v, wantBridgeFee %v", amount, tt.wantAmount)
			}
		})
	}
}

// todo rewrite tests without this func

func getBridgeFeeAndAmount(
	reqAmount decimal.Decimal,
	tokenUsdPrice decimal.Decimal,
	thisCoinPrice decimal.Decimal,
	transferFee decimal.Decimal,
	minBridgeFee decimal.Decimal,
	isAmountWithFees bool,
) (decimal.Decimal, decimal.Decimal, error) {
	amount := reqAmount.Copy()

	// if amount contains fees, then we need change the amount to the possible amount without fees (when transfer *max* native coins)
	if isAmountWithFees {
		var err error
		amount, err = possibleAmountWithoutFees(amount, tokenUsdPrice, transferFee, thisCoinPrice, minBridgeFee)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}
	}

	// get bridge fee
	bridgeFee, err := getBridgeFee(thisCoinPrice, tokenUsdPrice, amount, minBridgeFee)
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("error when getting bridge fee: %w", err)
	}

	return bridgeFee, amount, nil
}
