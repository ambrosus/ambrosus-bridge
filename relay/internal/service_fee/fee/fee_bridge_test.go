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
			name: "< $100K in AMB (should take 1%)",
			args: args{
				nativeUsdPrice: decimal.NewFromFloat(1),   // output in $
				tokenUsdPrice:  decimal.NewFromFloat(0.4), // 1 amb cost 0.4 $
				amount:         decimalFromString("2500"), // 1k$; fee is 1% = 10$
				minBridgeFee:   decimal.NewFromFloat(5),
			},
			want:    decimalFromString("10"),
			wantErr: false,
		},
		{
			name: "$100K in AMB (should take 0.5%)",
			args: args{
				nativeUsdPrice: decimal.NewFromFloat(1),     // output in $
				tokenUsdPrice:  decimal.NewFromFloat(0.4),   // 1 amb cost 0.4 $
				amount:         decimalFromString("500000"), // 200k$; fee is 0.5% = 1k$
				minBridgeFee:   decimal.NewFromFloat(5),
			},
			want:    decimalFromString("1000"),
			wantErr: false,
		},
		{
			name: "$100 in AMB (1% fee less then $5 => should take $5)",
			args: args{
				nativeUsdPrice: decimal.NewFromFloat(1),   // output in $
				tokenUsdPrice:  decimal.NewFromFloat(0.4), // 1 amb cost 0.4 $
				amount:         decimalFromString("250"),  // 100$; fee is 1$, but minBridgeFee is 5$
				minBridgeFee:   decimal.NewFromFloat(5),
			},
			want:    decimalFromString("5"),
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
