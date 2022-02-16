package amb

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func Test_deltaVS(t *testing.T) {
	type args struct {
		a []common.Address
		b []common.Address
	}
	tests := []struct {
		name    string
		args    args
		want    *Delta
		wantErr bool
	}{
		{
			"Deleted validator",
			args{
				a: set("0", "1", "2"),
				b: set("0", "2"),
			},
			&Delta{Index: -2, Address: set("1")[0]},
			false,
		},
		{
			"Deleted validator from the end",
			args{
				a: set("0", "1", "2"),
				b: set("0", "1"),
			},
			&Delta{Index: -3, Address: set("2")[0]},
			false,
		},
		{
			"Added validator",
			args{
				a: set("0", "1", "2"),
				b: set("0", "1", "3", "2"),
			},
			&Delta{Index: 2, Address: set("3")[0]},
			false,
		},
		{
			"Deleted 2 validators (error)",
			args{
				a: set("0", "1", "2"),
				b: set("0"),
			},
			nil,
			true,
		},
		{
			"Validator sets are the same (error)",
			args{
				a: set("0", "1", "2"),
				b: set("0", "1", "2"),
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := deltaVS(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("deltaVS() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			assert.Equalf(t, tt.want, res, "deltaVS(%v, %v)", tt.args.a, tt.args.b)
		})
	}
}

// for fast creating test addresses
// set("1") 	 -> [0x0000000000000000000000000000000000000001]
// set("a", "b") -> [0x000000000000000000000000000000000000000a, 0x000000000000000000000000000000000000000b]
func set(addresses ...string) []common.Address {
	result := make([]common.Address, 0, len(addresses))
	for _, i := range addresses {
		addr := common.HexToAddress(fmt.Sprint("0x000000000000000000000000000000000000000", i))
		result = append(result, addr)
	}
	return result
}
