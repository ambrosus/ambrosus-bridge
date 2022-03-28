package amb

import (
	"bytes"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"

	"github.com/ethereum/go-ethereum/rlp"
)

func EncodeBlock(header *Header) (*contracts.CheckAuraBlockAura, error) {
	// split rlp encoded header (bytes) by
	// - receiptHash, parentHash
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
	splitEls := [][]byte{
		header.ParentHash.Bytes(),
		header.ReceiptHash.Bytes(),
	}
	rlpParts, err := helpers.BytesSplit(rlpHeader, splitEls)
	if err != nil {
		return nil, err
	}

	// seal part
	step := header.Step()
	signature := header.Signature()

	stepPrefix, err := rlpPrefix(step)
	if err != nil {
		return nil, err
	}
	signaturePrefix, err := rlpPrefix(signature)
	if err != nil {
		return nil, err
	}

	return &contracts.CheckAuraBlockAura{
		P0Bare: p0Bare,
		P0Seal: p0Seal,

		P1:          rlpParts[0],
		ParentHash:  helpers.BytesToBytes32(splitEls[0]),
		P2:          rlpParts[1],
		ReceiptHash: helpers.BytesToBytes32(splitEls[1]),
		P3:          rlpParts[2],

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
