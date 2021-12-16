package enc2sol

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"relay/enc2sol/mytrie"
)

// todo get header from
// https://network.ambrosus.io
//{
//  "jsonrpc": "2.0",
//  "id": 1,
//  "method": "eth_getBlockByNumber",
//  "params": ["0xf478cd", true]
//}

type Header struct {
	ParentHash  common.Hash    `json:"parentHash"`
	UncleHash   common.Hash    `json:"sha3Uncles"`
	Coinbase    common.Address `json:"author"`
	Root        common.Hash    `json:"stateRoot"`
	TxHash      common.Hash    `json:"transactionsRoot"`
	ReceiptHash common.Hash    `json:"receiptsRoot"`
	Bloom       types.Bloom    `json:"logsBloom"`
	Difficulty  *big.Int       `json:"difficulty"`
	Number      *big.Int       `json:"number"`
	GasLimit    uint64         `json:"gasLimit"`
	GasUsed     uint64         `json:"gasUsed"`
	Time        uint64         `json:"timestamp"`
	Extra       []byte         `json:"extraData"`

	SealFields []string `json:"sealFields"`
	Signature  string   `json:"signature"`
}

func (h *Header) Hash() common.Hash {
	// todo without seal and sign
	return common.BytesToHash(mytrie.Hash(rlpEnc(h)))
}

func (h *Header) SealHash() common.Hash {
	// todo without sign
	return common.BytesToHash(mytrie.Hash(rlpEnc(h)))
}
