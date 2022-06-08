package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sortedKeys(t *testing.T) {
	type args struct {
		m map[uint64]string
	}
	tests := []struct {
		name string
		args args
		want []uint64
	}{
		{
			name: "Check order",
			args: args{map[uint64]string{
				1002: "",
				1000: "",
				1001: "",
			}},
			want: []uint64{
				1000,
				1001,
				1002,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, SortedKeys(tt.args.m), "SortedKeys(%v)", tt.args.m)
		})
	}
}
