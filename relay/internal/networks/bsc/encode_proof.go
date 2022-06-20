package bsc

import (
	"context"
	"fmt"
	"math/big"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
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
func splitVsChanges(proof *c.CheckPoSAPoSAProof, currentEpoch uint64) []*c.CheckPoSAPoSAProof {
	var res []*c.CheckPoSAPoSAProof
	var limit = 500
	var blocks = proof.Blocks
	if len(blocks) < limit {
		return nil
	}

	for i := 0; i < len(blocks); i += limit {
		number := new(big.Int).SetBytes(blocks[i].Number)
		if number.Cmp(big.NewInt(int64(currentEpoch*epochLength))) <= 0 {
			continue
		}

		if len(blocks)-i < limit {
			res = append(res, &c.CheckPoSAPoSAProof{
				Blocks:             blocks[i:],
				Transfer:           proof.Transfer,
				TransferEventBlock: proof.TransferEventBlock,
			})
			break
		}

		end := i + limit
		if end >= len(blocks) {
			end = len(blocks) - 1
		}

		for j := end; j >= i; j-- {
			number := new(big.Int).SetBytes(blocks[j].Number)
			if number.Mod(number, big.NewInt(epochLength)).Cmp(big.NewInt(0)) == 0 {
				end = j
				break
			}
		}

		// if i == end {
		// 	res = append(res, &c.CheckPoSAPoSAProof{})
		// 	break
		// }

		res = append(res, &c.CheckPoSAPoSAProof{
			Blocks: blocks[i:end],
			Transfer: c.CommonStructsTransferProof{
				ReceiptProof: [][]byte{},
				EventId:      big.NewInt(0),
				Transfers:    []c.CommonStructsTransfer{},
			},
			TransferEventBlock: 0,
		})
		i = end - limit
	}

	return res
}

func (b *Bridge) encodeTransferEvent(blocks map[uint64]*c.CheckPoSABlockPoSA, event *c.BridgeTransfer, safetyBlocks uint64) (*c.CommonStructsTransferProof, error) {
	proof, err := b.GetProof(event)
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

func (b *Bridge) encodeEpochChanges(blocks map[uint64]*c.CheckPoSABlockPoSA, end uint64) error {
	currentEpoch, err := b.sideBridge.GetCurrentEpoch()
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
