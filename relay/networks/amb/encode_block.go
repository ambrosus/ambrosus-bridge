package amb

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"relay/contracts"
	"relay/helpers"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func EncodeBlock(header *Header, isEventBlock bool) (*contracts.CheckPoABlockPoA, error) {
	// split rlp encoded header (bytes) by
	// - receiptHash (for event block) / parentHash (for safety block)
	// - Step, signature (for AURA)

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
	stepHexBytes, signatureBytes, err := getStepNSignBytes(header, stepHex)
	if err != nil {
		return nil, err
	}

	stepPrefix, signaturePrefix, err := getStepNSignPrefixes(header, stepHex)
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

func getStepNSignBytes(header *Header, stepHex string) ([]byte, []byte, error) {
	stepHexBytes, err := hex.DecodeString(stepHex)
	if err != nil {
		return nil, nil, err
	}
	signatureBytes, err := hex.DecodeString(header.Signature)
	if err != nil {
		return nil, nil, err
	}
	return stepHexBytes, signatureBytes, nil
}

func getStepNSignPrefixes(header *Header, stepHex string) ([]byte, []byte, error) {
	stepPrefix, err := rlpPrefixFromStrings(header.SealFields[0], stepHex)
	if err != nil {
		return nil, nil, err
	}
	signaturePrefix, err := rlpPrefixFromStrings(header.SealFields[1], header.Signature)
	if err != nil {
		return nil, nil, err
	}
	return stepPrefix, signaturePrefix, nil
}

func rlpPrefixFromStrings(encodedHexValue string, hexValue string) ([]byte, error) {
	encodedValue, err := hexutil.Decode(encodedHexValue)
	if err != nil {
		return nil, err
	}
	value, err := hexutil.Decode(hexValue)
	if err != nil {
		return nil, err
	}

	return rlpPrefix(encodedValue, value), nil
}

func rlpPrefix(encodedValue []byte, value []byte) []byte {
	res := bytes.TrimSuffix(encodedValue, value)
	return res
}

func uint64ToBytes(i *hexutil.Uint64) []byte {
	b, _ := hexutil.Decode(i.String())
	return b
}
