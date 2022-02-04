package amb

import (
	"bytes"
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"net/http"
	"relay/receipts_proof/mytrie"
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
	Signature  string   `json:"signature"`
}

func (b *Bridge) HeaderByNumber(number *big.Int) (*Header, error) {
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
	resp, err := http.Post(b.config.Url, "application/json", payloadBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData := new(response)
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
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

		for _, seal := range h.SealFields {
			sealRlpDecoded, err := decodeRlpHex(seal)
			if err != nil {
				return nil, err
			}
			headerAsSlice = append(headerAsSlice, *sealRlpDecoded)
		}
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

func decodeRlpHex(rlpHex string) (*[]byte, error) {
	prefixedRlpBytes, err := hexutil.Decode(rlpHex)
	if err != nil {
		return nil, err
	}

	rlpBytes := new([]byte)
	err = rlp.DecodeBytes(prefixedRlpBytes, rlpBytes)
	return rlpBytes, err
}
