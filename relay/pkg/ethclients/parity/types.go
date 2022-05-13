package parity

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

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
	return common.BytesToHash(crypto.Keccak256(rlp_))
}

func (h *Header) Step() (r []byte) {
	_ = rlp.DecodeBytes(h.SealFields[0], &r)
	return
}

func (h *Header) Signature() (r []byte) {
	_ = rlp.DecodeBytes(h.SealFields[1], &r)
	return
}
