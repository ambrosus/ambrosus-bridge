package bsc

import (
	"context"
	"fmt"
	"math"
	"math/big"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
)

const (
	addressLength     = 20
	extraVanityLength = 32
	extraSealLength   = 65
	epochLength       = 200
)

func (b *Bridge) encodePoSAProof(transferEvent *c.BridgeTransfer, safetyBlocks uint64) (*c.CheckPoSAPoSAProof, error) {
	// populated by functions below
	var blocksMap = make(map[uint64]*c.CheckPoSABlockPoSA)

	// encode transferProof and save event block to blocksMap
	transfer, err := b.encodeTransferEvent(blocksMap, transferEvent, safetyBlocks)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferEvent: %w", err)
	}

	// encode vsChange blocks to blocksMap
	epochChangesNums, err := b.findEpochChangeNums(transferEvent, safetyBlocks)
	if err != nil {
		return nil, fmt.Errorf("findEpochChangeNums: %w", err)
	}
	err = b.encodeEpochChanges(blocksMap, epochChangesNums)
	if err != nil {
		return nil, fmt.Errorf("encodeEpochChanges: %w", err)
	}

	// fill up blocks and get transfer event index
	indexToBlockNum := helpers.SortedKeys(blocksMap)
	var blocks []c.CheckPoSABlockPoSA
	var transferEventIndex uint64

	for i, blockNum := range indexToBlockNum {
		if blockNum == transferEvent.Raw.BlockNumber {
			transferEventIndex = uint64(i) // set transferEventIndex to index in blocks array
		}
		blocks = append(blocks, *blocksMap[blockNum])
	}

	return &c.CheckPoSAPoSAProof{
		Blocks:             blocks,
		Transfer:           *transfer,
		TransferEventBlock: transferEventIndex,
	}, nil
}

func (b *Bridge) encodeTransferEvent(blocks map[uint64]*c.CheckPoSABlockPoSA, event *c.BridgeTransfer, safetyBlocks uint64) (*c.CommonStructsTransferProof, error) {
	proof, err := b.GetProof(event)
	if err != nil {
		return nil, err
	}

	if err := b.saveBlocksRange(blocks, event.Raw.BlockNumber, event.Raw.BlockNumber+safetyBlocks); err != nil {
		return nil, err
	}

	return &c.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func (b *Bridge) encodeEpochChanges(blocksMap map[uint64]*c.CheckPoSABlockPoSA, epochChanges []uint64) error {
	// save blocks into blocksMap
	for _, epochChange := range epochChanges {
		// save epoch change block and get VS length
		epochChangeBlock, err := b.saveBlock(blocksMap, epochChange)
		if err != nil {
			return fmt.Errorf("save epoch change block: %w", err)
		}
		vsLength := getVSLength(epochChangeBlock)

		// start from +1 cuz the epoch change block is already saved
		if err := b.saveBlocksRange(blocksMap, epochChange+1, epochChange+vsLength); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bridge) findEpochChangeNums(event *c.BridgeTransfer, safetyBlocks uint64) ([]uint64, error) {
	start, err := b.getLastProcessedBlockNum(event.EventId)
	if err != nil {
		return nil, fmt.Errorf("getLastProcessedBlockNum: %w", err)
	}
	start = uint64(math.Ceil(float64(start)/epochLength) * epochLength) // first epoch change after last processed block (or this block itself)
	end := event.Raw.BlockNumber + safetyBlocks - 1                     // no need to change epoch for last block (it will be changed in next event processing)

	var epochChanges []uint64
	for blockNum := start; blockNum < end; blockNum += epochLength {
		epochChanges = append(epochChanges, blockNum)
	}
	return epochChanges, nil
}

func getVSLength(epochChangeBlock *c.CheckPoSABlockPoSA) uint64 {
	validatorsLen := len(epochChangeBlock.ExtraData) - extraSealLength - extraVanityLength
	return uint64(validatorsLen) / addressLength
}

// todo move to common
func (b *Bridge) getLastProcessedBlockNum(currEventId *big.Int) (uint64, error) {
	prevEventId := new(big.Int).Sub(currEventId, big.NewInt(1))
	prevEvent, err := b.GetEventById(prevEventId)
	if err != nil {
		return 0, fmt.Errorf("GetEventById: %w", err)
	}

	if prevEventId.Uint64() == 0 {
		return prevEvent.Raw.BlockNumber, nil
	}

	// todo specify block when prevEvent submitted in side network for 100$ correct `minSafetyBlocks` value
	minSafetyBlocks, err := b.SideBridge.GetMinSafetyBlocksNum(nil)
	if err != nil {
		return 0, fmt.Errorf("get block by hash: %w", err)
	}

	return prevEvent.Raw.BlockNumber + minSafetyBlocks, nil
}

// save blocks from `from` to `to` INCLUSIVE
func (b *Bridge) saveBlocksRange(blocksMap map[uint64]*c.CheckPoSABlockPoSA, from, to uint64) error {
	for i := from; i <= to; i++ {
		if _, err := b.saveBlock(blocksMap, i); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bridge) saveBlock(blocksMap map[uint64]*c.CheckPoSABlockPoSA, blockNumber uint64) (*c.CheckPoSABlockPoSA, error) {
	if encodedBlock, ok := blocksMap[blockNumber]; ok {
		return encodedBlock, nil
	}

	block, err := b.Client.HeaderByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, fmt.Errorf("HeaderByNumber: %w", err)
	}
	encodedBlock, err := b.EncodeBlock(block)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}

	blocksMap[blockNumber] = encodedBlock
	return encodedBlock, nil
}
