package amb

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"relay/contracts"
	"relay/helpers"
)

func EncodeBlock(header *Header, isEventBlock bool) *contracts.CheckPoABlockPoA {
	// split rlp encoded header (bytes) by
	// - receiptHash (for event block) / parentHash (for safety block)
	// - Timestamp (for AURA)

	rlpHeader := header.Rlp(false)

	splitEls := make([][]byte, 2)
	if isEventBlock {
		splitEls[0] = header.ReceiptHash.Bytes()
	} else {
		splitEls[0] = header.ParentHash.Bytes()
	}

	splitEls[1] = uint64ToBytes(header.Time)

	splitted, err := helpers.BytesSplit(rlpHeader, splitEls)
	if err != nil {
		panic(err)
	}

	return &contracts.CheckPoABlockPoA{
		P1:                    splitted[0],
		PrevHashOrReceiptRoot: helpers.BytesToBytes32(splitEls[0]),
		P2:                    splitted[1],
		Timestamp:             splitEls[1],
		P3:                    splitted[2],
		// seal
		S1:        nil,
		Signature: []byte(header.Signature),
		S2:        nil,
	}

}

func uint64ToBytes(i *hexutil.Uint64) []byte {
	b, _ := hex.DecodeString(fmt.Sprintf("%x", i))
	return b
}
