package fee_helper

import (
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper/explorers_clients"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestIsTrigger(t *testing.T) {
	assert.Equal(t, common.FromHex("746b5c42"), triggerMethodID)
	assert.Equal(t, true, isTrigger(explorers_clients.Transaction{Input: "0x746b5c42"}))
	assert.Equal(t, false, isTrigger(explorers_clients.Transaction{Input: "0x123456abcdef"}))
}
