package eth

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"relay/helpers"
)

func EncodeBlock(header *types.Header, isEventBlock bool) (result BlockInput) {
	// split rlp encoded header (bytes) by
	// - receiptHash (for event block) / parentHash (for safety block)
	// - Difficulty (for PoW)

	rlpHeader, _ := rlp.EncodeToBytes(header)

	splitEls := make([][]byte, 2)
	if isEventBlock {
		splitEls[0] = header.ReceiptHash.Bytes()
	} else {
		splitEls[0] = header.ParentHash.Bytes()
	}

	if true {
		timeBytes, _ := hex.DecodeString(fmt.Sprintf("%x", header.Time))
		splitEls[1] = timeBytes
	} else {
		splitEls[1] = header.Difficulty.Bytes()
	}

	splitted, err := helpers.BytesSplit(rlpHeader, splitEls)
	if err != nil {
		panic(err)
	}

	return BlockInput{
		P1:                    splitted[0],
		PrevHashOrReceiptRoot: splitEls[0],
		P2:                    splitted[1],
		Timestamp:             splitEls[1],
		P3:                    splitted[2],
		Signature:             []byte(header.Signature),
	}

}
