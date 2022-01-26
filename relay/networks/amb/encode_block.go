package amb

import (
	"encoding/hex"
	"relay/contracts"
	"relay/helpers"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func EncodeBlock(header *Header, isEventBlock bool) *contracts.CheckPoABlockPoA {
	// split rlp encoded header (bytes) by
	// - receiptHash (for event block) / parentHash (for safety block)
	// - Timestamp (for AURA)

	rlpHeaderWithSeal := header.Rlp(true)
	rlpHeaderWithoutSeal := header.Rlp(false)

	p0Bare := rlpHeaderWithoutSeal[:3]
	p0Seal := rlpHeaderWithSeal[:3]
	rlpHeaderWithSeal = rlpHeaderWithSeal[3:] // we'll work without prefix

	splitEls := make([][]byte, 4)

	if isEventBlock {
		splitEls[0] = header.ReceiptHash.Bytes()
	} else {
		splitEls[0] = header.ParentHash.Bytes()
	}

	splitEls[1] = uint64ToBytes(header.Time)
	splitEls[2], _ = hexutil.Decode(header.SealFields[0])
	splitEls[3], _ = hex.DecodeString(header.Signature)

	splitted, err := helpers.BytesSplit(rlpHeaderWithSeal, splitEls)
	if err != nil {
		panic(err)
	}

	return &contracts.CheckPoABlockPoA{
		P0Bare:                p0Bare,
		P0Seal:                p0Seal,
		P1:                    splitted[0],
		PrevHashOrReceiptRoot: helpers.BytesToBytes32(splitEls[0]),
		P2:                    splitted[1],
		Timestamp:             splitEls[1],
		P3:                    splitted[2],
		// seal
		S1:        append(splitEls[2], splitted[3]...),
		Signature: splitEls[3],
	}
}

func uint64ToBytes(i *hexutil.Uint64) []byte {
	b, _ := hexutil.Decode(i.String())
	return b
}
