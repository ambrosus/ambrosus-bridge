package posa_proof

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
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
	// populated by functions below
	var blocksMap = make(map[uint64]*c.CheckPoSABlockPoSA)

	// encode transferProof and save event block to blocksMap
	transfer, err := b.encodeTransferEvent(blocksMap, transferEvent, safetyBlocks)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferEvent: %w", err)
	}

	// encode vsChange blocks to blocksMap
	epochChangesLastBlock := transferEvent.Raw.BlockNumber + safetyBlocks - 1 // no need to change epoch for last block (it will be changed in next event processing)
	err = b.encodeEpochChanges(blocksMap, epochChangesLastBlock)
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

// if proof too long we need to split it into smaller parts
func (b *PoSAEncoder) splitVsChanges(proof *c.CheckPoSAPoSAProof) *c.CheckPoSAPoSAProof {
	b.logger.Warn().Int("blocks", len(proof.Blocks)).Msgf("PoSA proof too long")
	blocks := proof.Blocks[:len(proof.Blocks)/2] // drop half of blocks
	// todo maybe keep ~3000 blocks instead of half

	// last block should be vsFinalize (just before new epoch block)
	for i := len(blocks); i > 0; i-- {
		if binary.BigEndian.Uint64(blocks[i].Number)%200 == 0 {
			blocks = blocks[:i]
		}
	}
	return &c.CheckPoSAPoSAProof{
		Blocks: blocks,
	}
}

func (b *PoSAEncoder) encodeTransferEvent(blocks map[uint64]*c.CheckPoSABlockPoSA, event *c.BridgeTransfer, safetyBlocks uint64) (*c.CommonStructsTransferProof, error) {
	proof, err := cb.GetProof(b.bridge.GetClient(), event)
	if err != nil {
		return nil, err
	}

	// save `safetyBlocks` blocks after event block
	if err := b.saveBlocksRange(blocks, event.Raw.BlockNumber, event.Raw.BlockNumber+safetyBlocks); err != nil {
		return nil, err
	}

	return &c.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func (b *PoSAEncoder) encodeEpochChanges(blocks map[uint64]*c.CheckPoSABlockPoSA, end uint64) error {
	currentEpoch, err := b.posaReceiver.GetCurrentEpoch()
	if err != nil {
		return fmt.Errorf("GetCurrentEpoch: %w", err)
	}
	epochChangeBlock, err := b.saveBlock(map[uint64]*c.CheckPoSABlockPoSA{}, currentEpoch*epochLength)
	if err != nil {
		return fmt.Errorf("saveBlock: %w", err)
	}
	prevVsLen := getVSLength(epochChangeBlock) // use to determine how many blocks need to save in next epoch

	// save blocks into blocks (without current epoch)
	for epochBlock := (currentEpoch + 1) * epochLength; epochBlock < end; epochBlock += epochLength {
		// save epoch change block and get VS length
		epochChangeBlock, err = b.saveBlock(blocks, epochBlock)
		if err != nil {
			return fmt.Errorf("save epoch change block: %w", err)
		}

		// start from +1 cuz the epoch change block is already saved
		if err = b.saveBlocksRange(blocks, epochBlock+1, epochBlock+prevVsLen/2); err != nil {
			return err
		}

		prevVsLen = getVSLength(epochChangeBlock)
	}
	return nil
}

func getVSLength(epochChangeBlock *c.CheckPoSABlockPoSA) uint64 {
	validatorsLen := len(epochChangeBlock.ExtraData) - extraSealLength - extraVanityLength
	return uint64(validatorsLen) / addressLength
}

// save blocks from `from` to `to` INCLUSIVE
func (b *PoSAEncoder) saveBlocksRange(blocksMap map[uint64]*c.CheckPoSABlockPoSA, from, to uint64) error {
	for i := from; i <= to; i++ {
		if _, err := b.saveBlock(blocksMap, i); err != nil {
			return err
		}
	}
	return nil
}

func (b *PoSAEncoder) saveBlock(blocksMap map[uint64]*c.CheckPoSABlockPoSA, blockNumber uint64) (*c.CheckPoSABlockPoSA, error) {
	if encodedBlock, ok := blocksMap[blockNumber]; ok {
		return encodedBlock, nil
	}

	block, err := b.bridge.GetClient().HeaderByNumber(context.Background(), big.NewInt(int64(blockNumber)))
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
