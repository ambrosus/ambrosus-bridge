package amb

import (
	"bytes"
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/helpers"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

func EncodeBlock(header *Header, isEventBlock bool) (*contracts.CheckPoABlockPoA, error) {
	// split rlp encoded header (bytes) by
	// - receiptHash (for event block) / parentHash (for safety block)
	// - Step, Signature (for AURA)

	rlpHeader, err := header.Rlp(false)
	if err != nil {
		return nil, err
	}
	rlpHeaderSeal, err := header.Rlp(true)
	if err != nil {
		return nil, err
	}

	// rlpHeader length about 508 bytes => rlp prefix always 3 bytes length
	p0Bare := rlpHeader[:3]
	p0Seal := rlpHeaderSeal[:3]

	rlpHeader = rlpHeader[3:] // we'll work without prefix further

	// common part
	var prevHashOrReceiptRoot []byte
	if isEventBlock {
		prevHashOrReceiptRoot = header.ReceiptHash.Bytes()
	} else {
		prevHashOrReceiptRoot = header.ParentHash.Bytes()
	}

	rlpParts := bytes.SplitN(rlpHeader, prevHashOrReceiptRoot, 2)
	if len(rlpParts) != 2 {
		return nil, fmt.Errorf("split result length (%v) != 2 ", len(rlpParts))
	}

	// seal part
	step := common.Hex2Bytes(fmt.Sprintf("%x", header.Step)) // int -> bytes
	signature := common.Hex2Bytes(header.Signature)

	stepPrefix, err := rlpPrefix(step)
	if err != nil {
		return nil, err
	}
	signaturePrefix, err := rlpPrefix(signature)
	if err != nil {
		return nil, err
	}

	return &contracts.CheckPoABlockPoA{
		P0Bare: p0Bare,
		P0Seal: p0Seal,

		P1:                    rlpParts[0],
		PrevHashOrReceiptRoot: helpers.BytesToBytes32(prevHashOrReceiptRoot),
		P2:                    rlpParts[1],

		S1:        stepPrefix,
		Step:      step,
		S2:        signaturePrefix,
		Signature: signature,
	}, nil
}

func rlpPrefix(value []byte) ([]byte, error) {
	prefixedValue, err := rlp.EncodeToBytes(value)
	if err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(prefixedValue, value), nil
}
