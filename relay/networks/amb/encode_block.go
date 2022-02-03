package amb

import (
	"encoding/hex"
	"fmt"
	"relay/contracts"
	"relay/helpers"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func EncodeBlock(header *Header, isEventBlock bool) *contracts.CheckPoABlockPoA {
	// split rlp encoded header (bytes) by
	// - receiptHash (for event block) / parentHash (for safety block)
	// - Timestamp (for AURA)

	// todo handle errors

	rlpHeaderWithSeal, _ := header.Rlp(true)
	rlpHeaderWithoutSeal, _ := header.Rlp(false)

	p0Bare := rlpHeaderWithoutSeal[:3]
	p0Seal := rlpHeaderWithSeal[:3]
	rlpHeaderWithSeal = rlpHeaderWithSeal[3:] // we'll work without prefix

	splitEls := make([][]byte, 2)

	if isEventBlock {
		splitEls[0] = header.ReceiptHash.Bytes()
	} else {
		splitEls[0] = header.ParentHash.Bytes()
	}

	splitEls[1], _ = hexutil.Decode(header.SealFields[0])

	splitted, err := helpers.BytesSplit(rlpHeaderWithSeal, splitEls)
	if err != nil {
		panic(err)
	}

	stepHex := fmt.Sprintf("%x", header.Step)
	stepHexBytes, _ := hex.DecodeString(stepHex)
	stepPrefix, _ := hexutil.Decode(strings.TrimSuffix(header.SealFields[0], stepHex))
	signaturePrefix, _ := hexutil.Decode(strings.TrimSuffix(header.SealFields[1], header.Signature))
	signatureBytes, _ := hex.DecodeString(header.Signature)

	return &contracts.CheckPoABlockPoA{
		P0Bare:                p0Bare,
		P0Seal:                p0Seal,
		P1:                    splitted[0],
		PrevHashOrReceiptRoot: helpers.BytesToBytes32(splitEls[0]),
		P2:                    splitted[1],
		S1:                    stepPrefix,
		Step:                  stepHexBytes,
		S2:                    signaturePrefix,
		Signature:             signatureBytes,
	}
}

func uint64ToBytes(i *hexutil.Uint64) []byte {
	b, _ := hexutil.Decode(i.String())
	return b
}
