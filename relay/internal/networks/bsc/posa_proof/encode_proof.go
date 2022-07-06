package posa_proof

import (
	"context"
	"fmt"
	"math/big"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
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
	posaReceiver networks.BridgeReceivePoSA

	chainId *big.Int

	logger *zerolog.Logger

	// cache for fetching blocks. clear (almost) every `EncodePoSAProof` call
	fetchBlockCache func(arg uint64) (*types.Header, error)
}

func NewPoSAEncoder(bridge networks.Bridge, sideBridge networks.BridgeReceivePoSA, chainId *big.Int) *PoSAEncoder {
	return &PoSAEncoder{
		bridge:       bridge,
		posaReceiver: sideBridge,
		chainId:      chainId,
		logger:       bridge.GetLogger(), // todo maybe sublogger?
	}
}

func (b *PoSAEncoder) EncodePoSAProof(transferEvent *c.BridgeTransfer, safetyBlocks uint64) (*c.CheckPoSAPoSAProof, error) {
	// todo arg for split (reduced size proof)

	// new cache for every call
	// todo don't clear cache for reduced size proof
	b.fetchBlockCache = helpers.NewCache(b.fetchBlock)

	var blocksToSave []uint64

	lastBlock := transferEvent.Raw.BlockNumber + safetyBlocks

	// todo don't encode and save blocks for reduced size proof
	transferProof, err := b.encodeTransferEvent(transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferProof: %w", err)
	}

	// save blocks for transfer
	blocksToSave = append(blocksToSave, helpers.Range(transferEvent.Raw.BlockNumber, lastBlock+1)...)

	// lastBlock can be decreased if proof is too big
	epochBlocks, err := b.encodeEpochChanges(lastBlock)
	if err != nil {
		return nil, fmt.Errorf("encodeEpochChanges: %w", err)
	}
	// save blocks for epoch changes
	blocksToSave = append(blocksToSave, epochBlocks...)

	// fetch and encode blocksToSave
	blocks, blockNumToIndex, err := b.saveEncodedBlocks(blocksToSave)
	if err != nil {
		return nil, fmt.Errorf("saveEncodedBlocks: %w", err)
	}

	return &c.CheckPoSAPoSAProof{
		Blocks:             blocks,
		Transfer:           *transferProof,
		TransferEventBlock: uint64(blockNumToIndex[transferEvent.Raw.BlockNumber]),
	}, nil
}

func (b *PoSAEncoder) encodeTransferEvent(event *c.BridgeTransfer) (*c.CommonStructsTransferProof, error) {
	proof, err := cb.GetProof(b.bridge.GetClient(), event)
	if err != nil {
		return nil, err
	}

	return &c.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func (b *PoSAEncoder) encodeEpochChanges(end uint64) ([]uint64, error) {
	var blocksToSave []uint64

	currentEpoch, err := b.posaReceiver.GetCurrentEpoch()
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpoch: %w", err)
	}
	currentSetLen, err := b.fetchVSLength(currentEpoch) // use to determine how many blocks need to save in next epoch
	if err != nil {
		return nil, fmt.Errorf("fetchVSLength: %w", err)
	}

	// save blocks (without current epoch)
	for epoch := currentEpoch + 1; epoch*epochLength < end; epoch++ {
		blockWithSet := epoch * epochLength
		finalizeBlock := blockWithSet + currentSetLen/2 + 1
		blocksToSave = append(blocksToSave, helpers.Range(blockWithSet, finalizeBlock+1)...)

		currentSetLen, err = b.fetchVSLength(epoch) // use to determine how many blocks need to save in next epoch
		if err != nil {
			return nil, fmt.Errorf("fetchVSLength: %w", err)
		}
	}
	return blocksToSave, nil
}

func (b *PoSAEncoder) saveEncodedBlocks(blockNums []uint64) (blocks []c.CheckPoSABlockPoSA, blockNumToIndex map[uint64]int, err error) {
	blocks = make([]c.CheckPoSABlockPoSA, len(blockNums))

	for i, bn := range helpers.Sorted(blockNums) {
		block, err := b.fetchBlockCache(bn)
		if err != nil {
			return nil, nil, fmt.Errorf("fetchBlockCache: %w", err)
		}
		encodedBlock, err := b.EncodeBlock(block)
		if err != nil {
			return nil, nil, fmt.Errorf("EncodeBlock: %w", err)
		}

		blocks[i] = *encodedBlock
		blockNumToIndex[bn] = i
	}

	return blocks, blockNumToIndex, nil
}

func (b *PoSAEncoder) fetchBlock(blockNum uint64) (*types.Header, error) {
	return b.bridge.GetClient().HeaderByNumber(context.Background(), big.NewInt(int64(blockNum)))
}

func (b *PoSAEncoder) fetchVSLength(epoch uint64) (uint64, error) {
	block, err := b.fetchBlockCache(epoch * epochLength)
	if err != nil {
		return 0, fmt.Errorf("fetchBlockCache: %w", err)
	}
	return getVSLength(block), nil
}

func getVSLength(epochChangeBlock *types.Header) uint64 {
	validatorsLen := len(epochChangeBlock.Extra) - extraSealLength - extraVanityLength
	return uint64(validatorsLen) / addressLength
}
