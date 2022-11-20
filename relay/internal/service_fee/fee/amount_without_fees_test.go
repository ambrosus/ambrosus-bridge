package fee

import (
	"testing"

	"github.com/shopspring/decimal"
)

func decimalFromString(s string) decimal.Decimal {
	i, _ := decimal.NewFromString(s)
	return i
}

func Test_possibleAmountWithoutFees(t *testing.T) {
	origPercentFromAmount := percentFromAmount
	defer func() { percentFromAmount = origPercentFromAmount }()

	type args struct {
		amount        decimal.Decimal
		tokenUsdPrice decimal.Decimal
		thisCoinPrice decimal.Decimal
		transferFee   decimal.Decimal
		minBridgeFee  decimal.Decimal
	}
	tests := []struct {
		name              string
		percentFromAmount map[uint64]int64
		args              args
		want              decimal.Decimal
		wantErr           bool
	}{

		// todo readable numbers
		{
			name: "$100K in AMB (should take 1%)",
			percentFromAmount: map[uint64]int64{
				0:       1.0 * 100, // 0..100_000$ => 1%
				100_000: 0.5 * 100, // 100_000...$ => 0.5%
			},
			args: args{
				amount:        decimalFromString("13848494031870104170397696"),
				tokenUsdPrice: decimal.NewFromFloat(0.00000000000000000000722114619610344),
				thisCoinPrice: decimal.NewFromFloat(0.00000000000000000000722114619610344),
				transferFee:   decimalFromString("380000000000000000000"),
				minBridgeFee:  decimal.NewFromFloat(5),
			},
			want:    decimalFromString("13711000229574360564746670.0511020376037148"),
			wantErr: false,
		},
		{
			name: "$20 in AMB (should take $5)",
			percentFromAmount: map[uint64]int64{
				0:       1.0 * 100, // 0..100_000$ => 1%
				100_000: 0.5 * 100, // 100_000...$ => 0.5%
			},
			args: args{
				amount:        decimalFromString("3000000000000000000000"),
				tokenUsdPrice: decimal.NewFromFloat(0.00000000000000000000722114619610344),
				thisCoinPrice: decimal.NewFromFloat(0.00000000000000000000722114619610344),
				transferFee:   decimalFromString("380000000000000000000"),
				minBridgeFee:  decimal.NewFromFloat(5),
			},
			want:    decimalFromString("1927589146623562276743.200106524153161"),
			wantErr: false,
		},
		{
			name: "$20 in AMB & percentFromAmount is 0 (should take $0)",
			percentFromAmount: map[uint64]int64{
				0: 0.0 * 100, // 0.. ...$ => 0%
			},
			args: args{
				amount:        decimalFromString("3000000000000000000000"),
				tokenUsdPrice: decimal.NewFromFloat(0.00000000000000000000722114619610344),
				thisCoinPrice: decimal.NewFromFloat(0.00000000000000000000722114619610344),
				transferFee:   decimalFromString("0"),
				minBridgeFee:  decimal.NewFromFloat(0),
			},
			want:    decimalFromString("3000000000000000000000"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			percentFromAmount = tt.percentFromAmount // override global `percentFromAmount`

			got, err := possibleAmountWithoutFees(tt.args.amount, tt.args.tokenUsdPrice, tt.args.transferFee, tt.args.thisCoinPrice, tt.args.minBridgeFee)
			if (err != nil) != tt.wantErr {
				t.Errorf("possibleAmountWithoutFees() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("possibleAmountWithoutFees() = %v, want %v", got, tt.want)
			}
		})
	}
}
