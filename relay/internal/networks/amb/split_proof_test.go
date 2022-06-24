package amb

import (
	"math/big"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/stretchr/testify/assert"
)

func Test_setCorrectVsChanges(t *testing.T) {
	type args struct {
		vsChanges    []bindings.CheckAuraValidatorSetProof
		blockCounter int
	}
	tests := []struct {
		name string
		args args
		want []bindings.CheckAuraValidatorSetProof
	}{
		{
			"block counter is 0",
			args{
				vsChanges: []bindings.CheckAuraValidatorSetProof{
					{EventBlock: big.NewInt(0)},
					{EventBlock: big.NewInt(1)},
					{EventBlock: big.NewInt(3)},
					{EventBlock: big.NewInt(5)},
					{EventBlock: big.NewInt(7)},
					{EventBlock: big.NewInt(8)},
				},
				blockCounter: 0,
			},
			[]bindings.CheckAuraValidatorSetProof{
				{EventBlock: big.NewInt(0)},
				{EventBlock: big.NewInt(1)},
				{EventBlock: big.NewInt(3)},
				{EventBlock: big.NewInt(5)},
				{EventBlock: big.NewInt(7)},
				{EventBlock: big.NewInt(8)},
			},
		},
		{
			"block counter is 100",
			args{
				vsChanges: []bindings.CheckAuraValidatorSetProof{
					{EventBlock: big.NewInt(100)},
					{EventBlock: big.NewInt(101)},
					{EventBlock: big.NewInt(103)},
					{EventBlock: big.NewInt(105)},
					{EventBlock: big.NewInt(107)},
					{EventBlock: big.NewInt(108)},
				},
				blockCounter: 100,
			},
			[]bindings.CheckAuraValidatorSetProof{
				{EventBlock: big.NewInt(0)},
				{EventBlock: big.NewInt(1)},
				{EventBlock: big.NewInt(3)},
				{EventBlock: big.NewInt(5)},
				{EventBlock: big.NewInt(7)},
				{EventBlock: big.NewInt(8)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := setCorrectVsChangesEventBlock(tt.args.vsChanges, tt.args.blockCounter)

			// need for cases when big.Int.Sub(1, 1) is not equal big.NewInt(0) ("abs" in first case is "{}", in second is "nil")
			for i := 0; i < len(res); i++ {
				if res[i].EventBlock.Cmp(big.NewInt(0)) == 0 {
					res[i].EventBlock = big.NewInt(0)
				}
			}

			assert.EqualValuesf(t, tt.want, res, "setCorrectVsChangesEventBlock(%v, %v)", tt.args.vsChanges, tt.args.blockCounter)
		})
	}
}
