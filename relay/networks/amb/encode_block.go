package amb

import (
	"encoding/hex"
	"fmt"
	"relay/contracts"
	"relay/helpers"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func EncodeBlock(header *Header, isEventBlock bool) (*contracts.CheckPoABlockPoA, error) {
	// split rlp encoded header (bytes) by
	// - receiptHash (for event block) / parentHash (for safety block)
	// - Timestamp (for AURA)

	// todo handle errors

	rlpHeaderWithSeal, err := header.Rlp(true)
	if err != nil {
		return nil, err
	}
	rlpHeaderWithoutSeal, err := header.Rlp(false)
	if err != nil {
		return nil, err
	}

	p0Bare := rlpHeaderWithoutSeal[:3]
	p0Seal := rlpHeaderWithSeal[:3]
	rlpHeaderWithSeal = rlpHeaderWithSeal[3:] // we'll work without prefix

	splitEls := make([][]byte, 2)

	if isEventBlock {
		splitEls[0] = header.ReceiptHash.Bytes()
	} else {
		splitEls[0] = header.ParentHash.Bytes()
	}

	splitEls[1], err = hexutil.Decode(header.SealFields[0])
	if err != nil {
		return nil, err
	}

	splitted, err := helpers.BytesSplit(rlpHeaderWithSeal, splitEls)
	if err != nil {
		return nil, err
	}

	stepHex := fmt.Sprintf("%x", header.Step)
	stepHexBytes, err := hex.DecodeString(stepHex)
	if err != nil {
		return nil, err
	}
	signatureBytes, err := hex.DecodeString(header.Signature)
	if err != nil {
		return nil, err
	}

	stepPrefix, err := rlpPrefix(header.SealFields[0], stepHex)
	if err != nil {
		return nil, err
	}
	signaturePrefix, err := rlpPrefix(header.SealFields[1], header.Signature)
	if err != nil {
		return nil, err
	}

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
	}, nil
}

func rlpPrefix(withPrefix string, withoutPrefix string) ([]byte, error) {
	res, err := hexutil.Decode(strings.TrimSuffix(withPrefix, withoutPrefix))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func uint64ToBytes(i *hexutil.Uint64) []byte {
	b, _ := hexutil.Decode(i.String())
	return b
}
