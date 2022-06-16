package bsc

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func (b *Bridge) EncodeBlock(header *types.Header) (*bindings.CheckPoSABlockPoSA, error) {
	signedHeader, err := rlp.EncodeToBytes(header)
	if err != nil {
		return nil, fmt.Errorf("encode signed header: %w", err)
	}

	unsignedHeader, err := b.encodeUnsignedHeader(header)
	if err != nil {
		return nil, fmt.Errorf("encode unsigned header: %w", err)
	}

	// rlpHeader length about 508 bytes => rlp prefix always 3 bytes length
	p0Signed := signedHeader[:3]
	p0Unsigned := unsignedHeader[:3]

	p4Signed := helpers.RlpPrefix(len(header.Extra))

	splitEls := [][]byte{
		header.ParentHash.Bytes(),
		header.ReceiptHash.Bytes(),
		header.Number.Bytes(),
		p4Signed,
		header.Extra,
	}
	rlpParts, err := helpers.BytesSplit(signedHeader, splitEls)
	if err != nil {
		return nil, fmt.Errorf("split rlp header: %w", err)
	}

	return &bindings.CheckPoSABlockPoSA{
		P0Signed:   helpers.BytesToBytes3(p0Signed),
		P0Unsigned: helpers.BytesToBytes3(p0Unsigned),

		ParentHash:  helpers.BytesToBytes32(header.ParentHash.Bytes()),
		P1:          rlpParts[1],
		ReceiptHash: helpers.BytesToBytes32(header.ReceiptHash.Bytes()),
		P2:          rlpParts[2],
		Number:      header.Number.Bytes(),
		P3:          rlpParts[3],

		P4Signed:   p4Signed,
		P4Unsigned: helpers.RlpPrefix(len(header.Extra) - 65),
		ExtraData:  header.Extra,

		P5: rlpParts[5],
	}, nil
}

func (b *Bridge) encodeUnsignedHeader(header *types.Header) ([]byte, error) {
	return rlp.EncodeToBytes([]interface{}{
		b.chainId,
		header.ParentHash, header.UncleHash, header.Coinbase,
		header.Root, header.TxHash, header.ReceiptHash,
		header.Bloom, header.Difficulty, header.Number,
		header.GasLimit, header.GasUsed, header.Time,
		header.Extra[:len(header.Extra)-65],
		header.MixDigest, header.Nonce,
	})
}
