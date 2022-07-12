package aura_proof

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func Test_deltaVS(t *testing.T) {
	type args struct {
		prev []common.Address
		curr []common.Address
	}
	tests := []struct {
		name    string
		args    args
		address common.Address
		index   uint16
		wantErr bool
	}{
		{
			"Deleted validator",
			args{
				prev: set("0", "1", "2"),
				curr: set("0", "2"),
			},
			addr("1"), 2,
			false,
		},
		{
			"Deleted validator from the end",
			args{
				prev: set("0", "1", "2"),
				curr: set("0", "1"),
			},
			addr("2"),
			3,
			false,
		},
		{
			"Added validator",
			args{
				prev: set("0", "1", "2"),
				curr: set("0", "1", "3", "2"),
			},
			addr("3"),
			0,
			false,
		},
		{
			"Deleted 2 validators (error)",
			args{
				prev: set("0", "1", "2"),
				curr: set("0"),
			},
			common.Address{}, 0,
			true,
		},
		{
			"Validator sets are the same (error)",
			args{
				prev: set("0", "1", "2"),
				curr: set("0", "1", "2"),
			},
			common.Address{}, 0,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address, index, err := deltaVS(tt.args.prev, tt.args.curr)
			if (err != nil) != tt.wantErr {
				t.Errorf("deltaVS() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.address, address, "deltaVS(%v, %v)", tt.args.prev, tt.args.curr)
			assert.Equalf(t, tt.index, index, "deltaVS(%v, %v)", tt.args.prev, tt.args.curr)
		})
	}
}

// addr("1") -> 0x0000000000000000000000000000000000000001
func addr(i string) common.Address {
	return common.HexToAddress(fmt.Sprint("0x000000000000000000000000000000000000000", i))
}

// for fast creating test addresses
// set("1") 	 -> [0x0000000000000000000000000000000000000001]
// set("a", "b") -> [0x000000000000000000000000000000000000000a, 0x000000000000000000000000000000000000000b]
func set(addresses ...string) []common.Address {
	result := make([]common.Address, 0, len(addresses))
	for _, i := range addresses {
		result = append(result, addr(i))
	}
	return result
}
