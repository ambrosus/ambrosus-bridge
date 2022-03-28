package amb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof/mytrie"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type request struct {
	Jsonrpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}
type response struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  Header `json:"result"`
	Id      int    `json:"id"`
}

type Header struct {
	ParentHash  *common.Hash    `json:"parentHash"`
	UncleHash   *common.Hash    `json:"sha3Uncles"`
	Coinbase    *common.Address `json:"author"`
	Root        *common.Hash    `json:"stateRoot"`
	TxHash      *common.Hash    `json:"transactionsRoot"`
	ReceiptHash *common.Hash    `json:"receiptsRoot"`
	Bloom       *types.Bloom    `json:"logsBloom"`
	Difficulty  *hexutil.Big    `json:"difficulty"`
	Number      *hexutil.Big    `json:"number"`
	GasLimit    *hexutil.Uint64 `json:"gasLimit"`
	GasUsed     *hexutil.Uint64 `json:"gasUsed"`
	Time        *hexutil.Uint64 `json:"timestamp"`
	Extra       *hexutil.Bytes  `json:"extraData"`

	SealFields []string `json:"sealFields"`
	Step       uint64   `json:"step,string"`
	Signature  string   `json:"signature"`
}

func (b *Bridge) HeaderByNumber(number *big.Int) (header *Header, err error) {
	body := &request{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{hexutil.EncodeBig(number), true},
	}
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(body); err != nil {
		return nil, err
	}
	resp, err := http.Post(b.HttpUrl, "application/json", payloadBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData := new(response)
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	// Check if result is empty
	if respData.Result.Number == nil {
		return nil, fmt.Errorf("there is no header with number %d", number.Int64())
	}
	return &respData.Result, err
}

func (h *Header) Rlp(withSeal bool) ([]byte, error) {
	headerAsSlice := []interface{}{
		h.ParentHash, h.UncleHash, h.Coinbase, h.Root,
		h.TxHash, h.ReceiptHash, h.Bloom, h.Difficulty.ToInt(),
		h.Number.ToInt(), h.GasLimit, h.GasUsed, h.Time, h.Extra,
	}

	if withSeal {
		headerAsSlice = append(headerAsSlice,
			common.Hex2Bytes(fmt.Sprintf("%x", h.Step)), // int -> bytes
			common.Hex2Bytes(h.Signature),
		)
	}

	return rlp.EncodeToBytes(headerAsSlice)
}

func (h *Header) Hash(seal bool) common.Hash {
	rlp_, err := h.Rlp(seal)
	if err != nil {
		return common.Hash{}
	}
	return common.BytesToHash(mytrie.Hash(rlp_))
}
