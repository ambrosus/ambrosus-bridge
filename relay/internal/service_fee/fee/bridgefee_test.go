package fee

import (
	"testing"

	"github.com/shopspring/decimal"
)

func Test_getBridgeFee(t *testing.T) {
	type args struct {
		nativeUsdPrice decimal.Decimal
		tokenUsdPrice  decimal.Decimal
		amount         decimal.Decimal
		minBridgeFee   decimal.Decimal
	}
	tests := []struct {
		name    string
		args    args
		want    decimal.Decimal
		wantErr bool
	}{
		{
			name: "$100K in AMB (should take 5%)",
			args: args{
				nativeUsdPrice: decimal.NewFromFloat(0.00000000000000000000722114619610344),
				tokenUsdPrice:  decimal.NewFromFloat(0.00000000000000000000722114619610344),
				amount:         decimalFromString("13188661935114384924195179.3913578384840331"),
				minBridgeFee:   decimal.NewFromFloat(5),
			},
			want:    decimalFromString("659433096755719246204219.6827408804224156"),
			wantErr: false,
		},
		{
			name: "$20 in AMB (should take $5)",
			args: args{
				nativeUsdPrice: decimal.NewFromFloat(0.00000000000000000000722114619610344),
				amount:         decimalFromString("3000000000000000000000"),
				tokenUsdPrice:  decimal.NewFromFloat(0.00000000000000000000722114619610344),
				minBridgeFee:   decimal.NewFromFloat(5),
			},
			want:    decimalFromString("692410853376437723256.799893475846839"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBridgeFee(tt.args.nativeUsdPrice, tt.args.tokenUsdPrice, tt.args.amount, tt.args.minBridgeFee)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBridgeFee() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("getBridgeFee() = %v, want %v", got, tt.want)
			}
		})
	}
}
