package amb

import (
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
				a: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					common.HexToAddress("0x0000000000000000000000000000000000000001"),
					common.HexToAddress("0x0000000000000000000000000000000000000002"),
				},
				b: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					common.HexToAddress("0x0000000000000000000000000000000000000002"),
				},
			},
			&Delta{Index: 1, Address: common.HexToAddress("0x0000000000000000000000000000000000000001")},
			false,
		},
		{
			"Deleted validator from the end",
			args{
				a: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					common.HexToAddress("0x0000000000000000000000000000000000000001"),
					common.HexToAddress("0x0000000000000000000000000000000000000002"),
				},
				b: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					common.HexToAddress("0x0000000000000000000000000000000000000001"),
				},
			},
			&Delta{Index: 2, Address: common.HexToAddress("0x0000000000000000000000000000000000000002")},
			false,
		},
		{
			"Added validator",
			args{
				a: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					common.HexToAddress("0x0000000000000000000000000000000000000001"),
					common.HexToAddress("0x0000000000000000000000000000000000000002"),
				},
				b: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					common.HexToAddress("0x0000000000000000000000000000000000000001"),
					common.HexToAddress("0x0000000000000000000000000000000000000003"),
					common.HexToAddress("0x0000000000000000000000000000000000000002"),
				},
			},
			&Delta{Index: 2, Address: common.HexToAddress("0x0000000000000000000000000000000000000003")},
			false,
		},
		{
			"Deleted 2 validators (error)",
			args{
				a: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					common.HexToAddress("0x0000000000000000000000000000000000000001"),
					common.HexToAddress("0x0000000000000000000000000000000000000002"),
				},
				b: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
				},
			},
			nil,
			true,
		},
		{
			"Validator sets are the same (error)",
			args{
				a: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					common.HexToAddress("0x0000000000000000000000000000000000000001"),
					common.HexToAddress("0x0000000000000000000000000000000000000002"),
				},
				b: []common.Address{
					common.HexToAddress("0x0000000000000000000000000000000000000000"),
					common.HexToAddress("0x0000000000000000000000000000000000000001"),
					common.HexToAddress("0x0000000000000000000000000000000000000002"),
				},
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
