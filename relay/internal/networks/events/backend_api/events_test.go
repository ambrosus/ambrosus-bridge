package backend_api

import (
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetTransferEventReal(t *testing.T) {
	api := NewEventsApi("localhost:8080", "amb", "eth", &zerolog.Logger{})
	resp, err := api.GetTransfer(1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestGetTransferEvent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
"blockNumber": "0x16865a8",
"blockHash": "0x196d7c2a3334eb5bb18fee6c0483343bb224a5a2884d01c0d3441213691757d8",
"transactionIndex": "0x65",
"removed": false,
"address": "0x92fa52d3043725D00Eab422440C4e9ef3ba180d3",
"data": "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000007a477aa8ed4884509387dba81ba6f2b7c97597e2000000000000000000000000bb5bd958bfdd217283b4723b63186069c27d5f9c00000000000000000000000000000000000000000000001866a485eb4118342b",
"topics": [
	"0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2",
	"0x000000000000000000000000000000000000000000000000000000000000008f"
],
"transactionHash": "0x5d219c7ce6b985a6537ded3e1fddfa420cd8def15bbe83ec65dc73fd71d9637a",
"logIndex": "0x13c"
}`))
	}))
	expectedAmount, _ := new(big.Int).SetString("450118041502534349867", 10)
	expectedEvent := &bindings.BridgeTransfer{
		EventId: big.NewInt(143),
		Queue: []bindings.CommonStructsTransfer{
			{
				TokenAddress: common.HexToAddress("0x7A477aA8ED4884509387Dba81BA6F2B7C97597e2"),
				ToAddress:    common.HexToAddress("0xbB5bd958bFDD217283B4723B63186069c27D5F9C"),
				Amount:       expectedAmount,
			},
		},
		Raw: types.Log{
			Address: common.HexToAddress("0x92fa52d3043725D00Eab422440C4e9ef3ba180d3"),
			Topics: []common.Hash{
				common.HexToHash("0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2"),
				common.HexToHash("0x000000000000000000000000000000000000000000000000000000000000008f"),
			},
			Data:        common.FromHex("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000007a477aa8ed4884509387dba81ba6f2b7c97597e2000000000000000000000000bb5bd958bfdd217283b4723b63186069c27d5f9c00000000000000000000000000000000000000000000001866a485eb4118342b"),
			BlockNumber: uint64(23618984),
			TxHash:      common.HexToHash("0x5d219c7ce6b985a6537ded3e1fddfa420cd8def15bbe83ec65dc73fd71d9637a"),
			TxIndex:     101,
			BlockHash:   common.HexToHash("0x196d7c2a3334eb5bb18fee6c0483343bb224a5a2884d01c0d3441213691757d8"),
			Index:       316,
			Removed:     false,
		},
	}

	url := strings.TrimPrefix(ts.URL, "http://")
	api := NewEventsApi(url, "bsc", "amb", &zerolog.Logger{})
	resp, err := api.GetTransfer(143)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedEvent, resp)
}
