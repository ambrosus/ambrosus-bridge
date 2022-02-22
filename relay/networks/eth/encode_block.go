package eth

import (
	"bytes"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/helpers"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func EncodeBlock(header *types.Header) (*contracts.CheckPoWBlockPoW, error) {
	// split rlp encoded header (bytes) by
	// - Number
	// - Nonce
	// - Difficulty (for PoW)

	rlpHeader, err := rlp.EncodeToBytes(header)
	if err != nil {
		return nil, err
	}

	splitEls := [][]byte{
		header.Difficulty.Bytes(),
		header.Number.Bytes(),
		header.Nonce[:],
	}

	splitted, err := helpers.BytesSplit(rlpHeader, splitEls)
	if err != nil {
		return nil, err
	}

	dataSetLookUp, witnessForLookup, err := DisputeBlock(header)
	if err != nil {
		return nil, err
	}

	return &contracts.CheckPoWBlockPoW{
		P1:         splitted[0],
		Difficulty: splitEls[0],
		P2:         splitted[1],
		Number:     splitEls[1],
		P3:         splitted[2],
		Nonce:      splitEls[2],
		P4:         splitted[3],

		DataSetLookUp:    dataSetLookUp,
		WitnessForLookUp: witnessForLookup,
	}, nil

}

func DisputeBlock(header *types.Header) ([]*big.Int, []*big.Int, error) {
	blockHeaderWithoutNonce, err := encodeHeaderWithoutNonceToRLP(header)
	if err != nil {
		log.Error().Err(err).Msg("block header not encode")
		return nil, nil, err
	}

	blockHeaderHashWithoutNonce := crypto.Keccak256(blockHeaderWithoutNonce)

	var blockHeaderHashWithoutNonceLength32 [32]byte
	copy(blockHeaderHashWithoutNonceLength32[:], blockHeaderHashWithoutNonce)

	blockMetaData := ethash.NewBlockMetaData(
		header.Number.Uint64(), header.Nonce.Uint64(),
		blockHeaderHashWithoutNonceLength32,
	)

	dataSetLookUp := blockMetaData.DAGElementArray()
	witnessForLookup := blockMetaData.DAGProofArray()

	return dataSetLookUp, witnessForLookup, nil
}

func encodeHeaderWithoutNonceToRLP(header *types.Header) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := rlp.Encode(buffer, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
		header.BaseFee,
	})

	return buffer.Bytes(), err
}
