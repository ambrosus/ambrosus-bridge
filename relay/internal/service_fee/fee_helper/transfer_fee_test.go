package fee_helper

import (
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper/explorers_clients"
	ec "github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper/explorers_clients"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestIsTrigger(t *testing.T) {
	assert.Equal(t, common.FromHex("746b5c42"), triggerMethodID)
	assert.True(t, isTrigger(&explorers_clients.Transaction{Input: "0x746b5c42"}))
	assert.False(t, isTrigger(&explorers_clients.Transaction{Input: "0x123456abcdef"}))
}

func Test_filterTxsWithTriggers(t *testing.T) {
	tests := []struct {
		name string
		txs  []*ec.Transaction
		want []*ec.Transaction
	}{
		{
			name: "Ok",
			txs: []*ec.Transaction{
				{Input: "0x123123"},
				{Input: "0xtesttest"},
				{Input: common.Bytes2Hex(triggerMethodID)},
			},
			want: []*ec.Transaction{
				{Input: common.Bytes2Hex(triggerMethodID)},
			},
		},
		{
			name: "Empty result",
			txs: []*ec.Transaction{
				{Input: "0x123123"},
				{Input: "0xtesttest"},
				{Input: "0xtesttest2"},
			},
			want: []*ec.Transaction{},
		},
		{
			name: "Empty input and output",
			txs:  []*ec.Transaction{},
			want: []*ec.Transaction{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, filterTxsWithTriggers(tt.txs), tt.want)
		})
	}
}
