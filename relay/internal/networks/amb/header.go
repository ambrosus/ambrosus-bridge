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

	SealFields []hexutil.Bytes `json:"sealFields"`
}

func (b *Bridge) HeaderByNumber(number *big.Int) (*Header, error) {
	body := &request{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{hexutil.EncodeBig(number), false},
	}
	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(body); err != nil {
		return nil, err
	}
	resp, err := http.Post(b.config.HttpURL, "application/json", payloadBuf)
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
	return &respData.Result, nil
}

func (h *Header) Rlp(withSeal bool) ([]byte, error) {
	headerAsSlice := []interface{}{
		h.ParentHash, h.UncleHash, h.Coinbase, h.Root,
		h.TxHash, h.ReceiptHash, h.Bloom, h.Difficulty.ToInt(),
		h.Number.ToInt(), h.GasLimit, h.GasUsed, h.Time, h.Extra,
	}

	if withSeal {
		headerAsSlice = append(headerAsSlice, h.Step(), h.Signature())
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

func (h *Header) Step() (r []byte) {
	_ = rlp.DecodeBytes(h.SealFields[0], &r)
	return
}

func (h *Header) Signature() (r []byte) {
	_ = rlp.DecodeBytes(h.SealFields[1], &r)
	return
}
