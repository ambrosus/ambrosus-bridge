package ambrosus_explorer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AmbrosusExplorer_TxListByAddress(t *testing.T) {
	e, err := NewAmbrosusExplorer(MainnetAPIBase, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := e.TxListByAddress(
		"0x0f071e1785e3E115360E04c9C8D53e958E6f85FE",
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}

func Test_AmbrosusExplorer_TxListByFromToAddresses(t *testing.T) {
	e, err := NewAmbrosusExplorer(MainnetAPIBase, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := e.TxListByFromToAddresses(
		// "0x0f071e1785e3E115360E04c9C8D53e958E6f85FE",
		// "0x19caBC1E34Ab0CC5C62DaA1394f6022B38b75c78",
		"",
		"0xf9cF8775CD6697c1bB52fdbD9625C8f5Fae6c62C",
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}

func Test_AmbrosusExplorer_TxListByFromToAddressesUntilTxHash(t *testing.T) {
	e, err := NewAmbrosusExplorer(MainnetAPIBase, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := e.TxListByFromToAddressesUntilTxHash(
		"0x0f071e1785e3E115360E04c9C8D53e958E6f85FE",
		"0x19caBC1E34Ab0CC5C62DaA1394f6022B38b75c78",
		"0xf9fd4889f5a0ef71c716bb815897b4303a7e81cdb85a7af2fb09dda04ba8514f",
	)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range res {
		t.Log(v.Hash)
	}
	
	assert.Equal(t, "0xce40b2d98099dce9219619110f5be954b2383ef97395c9b6cf4098a16d374b3e", res[0].Hash)
}
