package amb

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
)

func EncodeBlock(header *Header) (*contracts.CheckAuraBlockAura, error) {
	// split rlp encoded header (bytes) by
	// - receiptHash, parentHash
	// - Step, Signature (for AURA)

	rlpHeader, err := header.Rlp(false)
	if err != nil {
		return nil, fmt.Errorf("rlp header: %w", err)
	}
	rlpHeaderSeal, err := header.Rlp(true)
	if err != nil {
		return nil, fmt.Errorf("rlp header seal: %w", err)
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
		return nil, fmt.Errorf("split rlp header: %w", err)
	}

	return &contracts.CheckAuraBlockAura{
		P0Bare: helpers.BytesToBytes3(p0Bare),
		P0Seal: helpers.BytesToBytes3(p0Seal),

		ParentHash:  helpers.BytesToBytes32(splitEls[0]),
		P2:          rlpParts[1],
		ReceiptHash: helpers.BytesToBytes32(splitEls[1]),
		P3:          rlpParts[2],

		Step:      header.Step(),
		Signature: header.Signature(),
	}, nil
}
