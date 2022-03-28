package eth

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/rs/zerolog/log"
)

func EncodeBlock(header *types.Header, isEventBlock bool) (*contracts.CheckPoWBlockPoW, error) {
	encodedBlock, err := splitBlock(header, isEventBlock)
	if err != nil {
		return nil, err
	}

	encodedBlock.DataSetLookup, encodedBlock.WitnessForLookup, err = getLookupData(header)
	if err != nil {
		return nil, err
	}

	return encodedBlock, nil
}

func splitBlock(header *types.Header, isEventBlock bool) (*contracts.CheckPoWBlockPoW, error) {
	// split rlp encoded header (bytes) by
	// - receiptHash (for event block) / parentHash (for safety block)
	// - Difficulty
	// - Number
	// - Nonce
	// Also blockHeaderWithoutNonce calculated in solidity as concat all fields but P5, Nonce
	// so P4 and P5 can't be concatenated together

	rlpWithNonce, err := headerRlp(header, true)
	if err != nil {
		return nil, err
	}

	rlpWithoutNonce, err := headerRlp(header, false)
	if err != nil {
		return nil, err
	}

	// rlpHeader length about 508 bytes => rlp prefix always 3 bytes length
	p0WithNonce := rlpWithNonce[:3]
	p0WithoutNonce := rlpWithoutNonce[:3]

	rlpHeader := rlpWithNonce[3:] // we'll work without prefix further

	splitEls := [][]byte{
		nil,
		header.Difficulty.Bytes(),
		header.Number.Bytes(),
		header.Extra,
		header.Nonce[:],
	}
	if isEventBlock {
		splitEls[0] = header.ReceiptHash.Bytes()
	} else {
		splitEls[0] = header.ParentHash.Bytes()
	}

	split, err := helpers.BytesSplit(rlpHeader, splitEls)
	if err != nil {
		return nil, err
	}

	return &contracts.CheckPoWBlockPoW{
		P0WithNonce:    p0WithNonce,
		P0WithoutNonce: p0WithoutNonce,

		P1:                  split[0],
		ParentOrReceiptHash: helpers.BytesToBytes32(splitEls[0]),
		P2:                  split[1],
		Difficulty:          splitEls[1],
		P3:                  split[2],
		Number:              splitEls[2],
		P4:                  helpers.BytesConcat(split[3], header.Extra),
		P5:                  split[4],
		Nonce:               splitEls[4],
		P6:                  split[5],
	}, nil

}

func getLookupData(header *types.Header) ([]*big.Int, []*big.Int, error) {
	blockHeaderWithoutNonce, err := headerRlp(header, false)
	if err != nil {
		log.Error().Err(err).Msg("block header not encode")
		return nil, nil, err
	}
	hashWithoutNonce := helpers.BytesToBytes32(crypto.Keccak256(blockHeaderWithoutNonce))

	blockMetaData := ethash.NewBlockMetaData(header.Number.Uint64(), header.Nonce.Uint64(), hashWithoutNonce)

	dataSetLookUp := blockMetaData.DAGElementArray()
	witnessForLookup := blockMetaData.DAGProofArray()

	return dataSetLookUp, witnessForLookup, nil
}

func headerRlp(header *types.Header, withNonce bool) ([]byte, error) {
	headerAsSlice := []interface{}{
		header.ParentHash, header.UncleHash, header.Coinbase,
		header.Root, header.TxHash, header.ReceiptHash,
		header.Bloom, header.Difficulty, header.Number,
		header.GasLimit, header.GasUsed, header.Time, header.Extra,
	}
	if withNonce {
		headerAsSlice = append(headerAsSlice, header.MixDigest, header.Nonce)
	}
	// Note: BaseFee is +- new field, old blocks without BaseFee should be treated as if this field does not exist
	if header.BaseFee != nil {
		headerAsSlice = append(headerAsSlice, header.BaseFee)
	}
	return rlp.EncodeToBytes(headerAsSlice)
}
