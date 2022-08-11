package fee_api

import (
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/shopspring/decimal"
)

var EtherDecimal = decimal.NewFromFloat(params.Ether)

func toDecimal[T float64 | int](val T) decimal.Decimal {
	var res decimal.Decimal

	switch v := any(val).(type) {
	case int:
		res = decimal.NewFromInt(int64(v))
	case float64:
		res = decimal.NewFromFloat(v)
	}

	return res
}

func toWei[T float64 | int](val T) decimal.Decimal {
	res := toDecimal(val)
	res = res.Mul(EtherDecimal)
	return res
}

func Test_coin2coin(t *testing.T) {
	// numbers can be so small and just "rounded" to zero without that line
	decimal.DivisionPrecision = 30

	type args struct {
		amountWei          decimal.Decimal
		firstCoinPriceUsd  decimal.Decimal
		secondCoinPriceUsd decimal.Decimal
	}

	tests := []struct {
		name string
		args args
		want decimal.Decimal
	}{
		{
			name: "1 ETH to AMB",
			args: args{
				amountWei:          EtherDecimal,
				firstCoinPriceUsd:  toDecimal(2000),
				secondCoinPriceUsd: toDecimal(0.005),
			},
			want: toWei(400_000),
		},
		{
			name: "2 ETH to AMB",
			args: args{
				amountWei:          EtherDecimal.Mul(toDecimal(2)),
				firstCoinPriceUsd:  toDecimal(2000),
				secondCoinPriceUsd: toDecimal(0.005),
			},
			want: toWei(800_000),
		},
		{
			name: "1 AMB to ETH",
			args: args{
				amountWei:          EtherDecimal,
				firstCoinPriceUsd:  toDecimal(0.005),
				secondCoinPriceUsd: toDecimal(2000),
			},
			want: toWei(0.0000025),
		},
		{
			name: "2 AMB to ETH",
			args: args{
				amountWei:          EtherDecimal.Mul(toDecimal(2)),
				firstCoinPriceUsd:  toDecimal(0.005),
				secondCoinPriceUsd: toDecimal(2000),
			},
			want: toWei(0.000005),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Coin2coin(tt.args.amountWei, tt.args.firstCoinPriceUsd, tt.args.secondCoinPriceUsd); got.Cmp(tt.want) != 0 {
				t.Errorf("coin2coin() = %v, want %v", got, tt.want)
			}
		})
	}
}
