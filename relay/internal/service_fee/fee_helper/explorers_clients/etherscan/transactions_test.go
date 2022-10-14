package etherscan

import (
	"fmt"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper/explorers_clients"
	"github.com/stretchr/testify/assert"
)

func TestGetTxs(t *testing.T) {
	explorer, err := NewEtherscan("https://api-ropsten.etherscan.io/api?apikey=DY4Z86MQ2D9E24C6HB98PTA79EKJ5TQIFX", nil)
	if err != nil {
		t.Fatal(err)
	}
	r, err := explorer.TxListByFromToAddresses("0x295C2707319ad4BecA6b5bb4086617fD6F240CfE", "0x5570DF279dA47Fa46aBa4c1e27b55c8e1d6648A1", explorers_clients.TxFilters{0, nil})
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(len(r))
}

func TestGetTxsByFromList(t *testing.T) {
	explorer, err := NewEtherscan("https://api-ropsten.etherscan.io/api?apikey=DY4Z86MQ2D9E24C6HB98PTA79EKJ5TQIFX", nil)
	if err != nil {
		t.Fatal(err)
	}

	fromList := []string{"0xd693a3cc5686e74ca2e72e8120a2f2013b8ee66e", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe"}
	r, err := explorer.TxListByFromListToAddresses(fromList, "0x5570DF279dA47Fa46aBa4c1e27b55c8e1d6648A1", explorers_clients.TxFilters{0, nil})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(r))

	for _, tx := range r {
		fmt.Println(tx.BlockNumber)
		assert.Contains(t, fromList, tx.From)
	}
}

func TestGetTxsByFromListAndFromBlock(t *testing.T) {
	explorer, err := NewEtherscan("https://api-ropsten.etherscan.io/api?apikey=DY4Z86MQ2D9E24C6HB98PTA79EKJ5TQIFX", nil)
	if err != nil {
		t.Fatal(err)
	}

	fromBlock := uint64(12_900_000)
	fromList := []string{"0xd693a3cc5686e74ca2e72e8120a2f2013b8ee66e", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe"}
	r, err := explorer.TxListByFromListToAddresses(fromList, "0x5570DF279dA47Fa46aBa4c1e27b55c8e1d6648A1", explorers_clients.TxFilters{uint64(fromBlock), nil})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(r))

	for _, tx := range r {
		assert.Contains(t, fromList, tx.From)
		assert.GreaterOrEqual(t, tx.BlockNumber, fromBlock)
	}
}
