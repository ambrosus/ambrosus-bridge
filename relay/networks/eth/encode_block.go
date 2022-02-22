package eth

import (
	"bytes"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/helpers"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func EncodeBlock(header *types.Header, isEventBlock bool) (*contracts.CheckPoWBlockPoW, error) {
	// split rlp encoded header (bytes) by
	// - receiptHash (for event block) / parentHash (for safety block)
	// - Difficulty (for PoW)

	rlpHeader, err := rlp.EncodeToBytes(header)
	if err != nil {
		return nil, err
	}

	splitEls := make([][]byte, 2)
	if isEventBlock {
		splitEls[0] = header.ReceiptHash.Bytes()
	} else {
		splitEls[0] = header.ParentHash.Bytes()
	}

	splitEls[1] = header.Difficulty.Bytes()

	splitted, err := helpers.BytesSplit(rlpHeader, splitEls)
	if err != nil {
		return nil, err
	}

	return &contracts.CheckPoWBlockPoW{
		P1:                    splitted[0],
		PrevHashOrReceiptRoot: helpers.BytesToBytes32(splitEls[0]),
		P2:                    splitted[1],
		Difficulty:            splitEls[1],
		P3:                    splitted[2],
	}, nil

}

func EncodeHeaderWithoutNonceToRLP(header *types.Header) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := rlp.Encode(buffer, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
		header.BaseFee,
	})

	return buffer.Bytes(), err
}
