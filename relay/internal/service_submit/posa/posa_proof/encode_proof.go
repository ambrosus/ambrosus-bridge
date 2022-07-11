package posa_proof

import (
	"context"
	"fmt"
	"math/big"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

const (
	addressLength     = 20
	extraVanityLength = 32
	extraSealLength   = 65
	epochLength       = 200
)

type PoSAEncoder struct {
	bridge       networks.Bridge
	posaReceiver service_submit.ReceiverPoSA

	chainId *big.Int

	logger *zerolog.Logger

	// cache for fetching blocks. clear (almost) every `EncodePoSAProof` call
	fetchBlockCache func(arg uint64) (*types.Header, error)
}

func NewPoSAEncoder(bridge networks.Bridge, sideBridge service_submit.ReceiverPoSA, chainId *big.Int) *PoSAEncoder {
	logger := bridge.GetLogger().With().Str("service", "PoSAEncoder").Logger()

	return &PoSAEncoder{
		bridge:       bridge,
		posaReceiver: sideBridge,
		chainId:      chainId,
		logger:       &logger,
	}
}

func (e *PoSAEncoder) EncodePoSAProof(transferEvent *c.BridgeTransfer, safetyBlocks uint64) (*c.CheckPoSAPoSAProof, error) {
	// todo arg for split (reduced size proof)

	// new cache for every call
	// todo don't clear cache for reduced size proof
	e.fetchBlockCache = helpers.NewCache(e.fetchBlock)

	var blocksToSave []uint64

	lastBlock := transferEvent.Raw.BlockNumber + safetyBlocks

	// todo don't encode and save blocks for reduced size proof
	transferProof, err := cb.EncodeTransferProof(e.bridge.GetClient(), transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferProof: %w", err)
	}

	// save blocks for transfer
	blocksToSave = append(blocksToSave, helpers.Range(transferEvent.Raw.BlockNumber, lastBlock+1)...)

	// lastBlock can be decreased if proof is too big
	epochBlocks, err := e.encodeEpochChanges(lastBlock)
	if err != nil {
		return nil, fmt.Errorf("encodeEpochChanges: %w", err)
	}
	// save blocks for epoch changes
	blocksToSave = append(blocksToSave, epochBlocks...)

	// fetch and encode blocksToSave
	blocks, blockNumToIndex, err := e.saveEncodedBlocks(blocksToSave)
	if err != nil {
		return nil, fmt.Errorf("saveEncodedBlocks: %w", err)
	}

	return &c.CheckPoSAPoSAProof{
		Blocks:             blocks,
		Transfer:           *transferProof,
		TransferEventBlock: uint64(blockNumToIndex[transferEvent.Raw.BlockNumber]),
	}, nil
}

func (e *PoSAEncoder) encodeEpochChanges(end uint64) ([]uint64, error) {
	var blocksToSave []uint64

	currentEpoch, err := e.posaReceiver.GetCurrentEpoch()
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpoch: %w", err)
	}
	currentSetLen, err := e.fetchVSLength(currentEpoch) // use to determine how many blocks need to save in next epoch
	if err != nil {
		return nil, fmt.Errorf("fetchVSLength: %w", err)
	}

	// save blocks (without current epoch)
	for epoch := currentEpoch + 1; epoch*epochLength < end; epoch++ {
		blockWithSet := epoch * epochLength
		finalizeBlock := blockWithSet + currentSetLen/2 + 1
		blocksToSave = append(blocksToSave, helpers.Range(blockWithSet, finalizeBlock+1)...)

		currentSetLen, err = e.fetchVSLength(epoch) // use to determine how many blocks need to save in next epoch
		if err != nil {
			return nil, fmt.Errorf("fetchVSLength: %w", err)
		}
	}
	return blocksToSave, nil
}

func (e *PoSAEncoder) saveEncodedBlocks(blockNums []uint64) (blocks []c.CheckPoSABlockPoSA, blockNumToIndex map[uint64]int, err error) {
	blocks = make([]c.CheckPoSABlockPoSA, len(blockNums))

	for i, bn := range helpers.Sorted(blockNums) {
		block, err := e.fetchBlockCache(bn)
		if err != nil {
			return nil, nil, fmt.Errorf("fetchBlockCache: %w", err)
		}
		encodedBlock, err := e.EncodeBlock(block)
		if err != nil {
			return nil, nil, fmt.Errorf("EncodeBlock: %w", err)
		}

		blocks[i] = *encodedBlock
		blockNumToIndex[bn] = i
	}

	return blocks, blockNumToIndex, nil
}

func (e *PoSAEncoder) fetchBlock(blockNum uint64) (*types.Header, error) {
	return e.bridge.GetClient().HeaderByNumber(context.Background(), big.NewInt(int64(blockNum)))
}

func (e *PoSAEncoder) fetchVSLength(epoch uint64) (uint64, error) {
	block, err := e.fetchBlockCache(epoch * epochLength)
	if err != nil {
		return 0, fmt.Errorf("fetchBlockCache: %w", err)
	}
	return getVSLength(block), nil
}

func getVSLength(epochChangeBlock *types.Header) uint64 {
	validatorsLen := len(epochChangeBlock.Extra) - extraSealLength - extraVanityLength
	return uint64(validatorsLen) / addressLength
}
