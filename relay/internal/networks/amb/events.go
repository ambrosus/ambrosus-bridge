package amb

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (b *Bridge) getBlocksAndEvents(transferEvent *c.BridgeTransfer, safetyBlocks uint64) (*c.CheckAuraAuraProof, error) {
	// populated by functions below
	blocksMap := make(map[uint64]*c.CheckAuraBlockAura)

	// encode transferProof and save event block to blocksMap
	transfer, err := b.encodeTransferEvent(blocksMap, transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferEvent: %w", err)
	}

	// encode vsChangeProofs and save event blocks to blocksMap
	vsChangeEvents, err := b.fetchVSChangeEvents(transferEvent, safetyBlocks)
	if err != nil {
		return nil, fmt.Errorf("fetchVSChangeEvents: %w", err)
	}
	vsChanges, err := b.encodeVSChangeEvents(blocksMap, vsChangeEvents)
	if err != nil {
		return nil, fmt.Errorf("encodeVSChangeEvents: %w", err)
	}

	// save safety blocks to blocksMap
	err = b.addSafetyBlocks(blocksMap, safetyBlocks)
	if err != nil {
		return nil, fmt.Errorf("encodeSafetyBlocks: %w", err)
	}

	// sort blocks in blocksMap and use resulting indexes
	blocks := make([]c.CheckAuraBlockAura, len(blocksMap))
	indexToBlockNum := sortedKeys(blocksMap)
	var transferEventIndex uint64

	for i, blockNum := range indexToBlockNum {
		if blockNum == transferEvent.Raw.BlockNumber {
			transferEventIndex = uint64(i) // set transferEventIndex to index in blocks array
		} else if blockNum > transferEvent.Raw.BlockNumber+safetyBlocks {
			blocks = blocks[:i] // in some cases we can fetch more blocks that we need
			break
		}

		blocks[i] = *blocksMap[blockNum]
	}

	return &c.CheckAuraAuraProof{
		Blocks:             blocks,
		Transfer:           transfer,
		VsChanges:          vsChanges,
		TransferEventBlock: transferEventIndex,
	}, nil
}

func (b *Bridge) encodeTransferEvent(blocks map[uint64]*c.CheckAuraBlockAura, event *c.BridgeTransfer) (c.CommonStructsTransferProof, error) {
	proof, err := b.getProof(event)
	if err != nil {
		return c.CommonStructsTransferProof{}, err
	}

	if err := b.saveBlock(event.Raw.BlockNumber, blocks); err != nil {
		return c.CommonStructsTransferProof{}, err
	}

	return c.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func (b *Bridge) encodeVSChangeEvents(blocks map[uint64]*c.CheckAuraBlockAura, events []*c.VsInitiateChange) ([]c.CheckAuraValidatorSetProof, error) {
	vsChanges := make([]c.CheckAuraValidatorSetProof, len(events))

	prevSet, err := b.sideBridge.GetValidatorSet()
	if err != nil {
		return nil, fmt.Errorf("GetValidatorSet: %w", err)
	}

	for i, event := range events {
		vsChanges[i], err = b.encodeVSChangeEvent(prevSet, event)
		if err != nil {
			return nil, fmt.Errorf("encodeVSChangeEvent: %w", err)
		}
		prevSet = event.NewSet

		if err := b.saveBlock(event.Raw.BlockNumber, blocks); err != nil {
			return nil, err
		}

		// block in which VS will be finalized
		if err := b.saveBlock(event.Raw.BlockNumber+2, blocks); err != nil {
			return nil, err
		}

		// in this block (one after event block) contracts should finalize all events
		// in vsChanges array up to `FinalizedVs` index (this event)
		blocks[event.Raw.BlockNumber+2].FinalizedVs = uint64(i) + 1

	}
	return vsChanges, nil
}

// add safety blocks after each event block
func (b *Bridge) addSafetyBlocks(blocksMap map[uint64]*c.CheckAuraBlockAura, minSafetyBlocks uint64) error {
	// we should iterate over keys because on writing new values to map
	// we'll iterate also over those new values, but we don't need that
	blockNums := sortedKeys(blocksMap)
	for _, blockNum := range blockNums {
		for i := uint64(0); i <= minSafetyBlocks; i++ {
			if err := b.saveBlock(blockNum+i, blocksMap); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Bridge) encodeVSChangeEvent(prevSet []common.Address, event *c.VsInitiateChange) (c.CheckAuraValidatorSetProof, error) {
	address, index, err := deltaVS(prevSet, event.NewSet)
	if err != nil {
		return c.CheckAuraValidatorSetProof{}, fmt.Errorf("deltaVS: %w", err)
	}

	proof, err := b.getProof(event)
	if err != nil {
		return c.CheckAuraValidatorSetProof{}, fmt.Errorf("getProof: %w", err)
	}

	return c.CheckAuraValidatorSetProof{
		ReceiptProof: proof,
		DeltaAddress: address,
		DeltaIndex:   index,
	}, nil
}

func (b *Bridge) fetchVSChangeEvents(event *c.BridgeTransfer, safetyBlocks uint64) ([]*c.VsInitiateChange, error) {
	start, err := b.getLastProcessedBlockNum()
	if err != nil {
		return nil, fmt.Errorf("getLastProcessedBlockNum: %w", err)
	}
	end := event.Raw.BlockNumber + safetyBlocks - 1 // we don't need safetyEnd block with VSChange event

	opts := &bind.FilterOpts{
		Start:   start.Uint64(),
		End:     &end,
		Context: context.Background(),
	}

	logs, err := b.VSContract.FilterInitiateChange(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter initiate changes: %w", err)
	}

	var res []*c.VsInitiateChange
	for logs.Next() {
		res = append(res, logs.Event)
	}

	return res, nil
}

func (b *Bridge) getProof(event receipts_proof.ProofEvent) ([][]byte, error) {
	receipts, err := b.GetReceipts(event.Log().BlockHash)
	if err != nil {
		return nil, fmt.Errorf("GetReceipts: %w", err)
	}
	return receipts_proof.CalcProofEvent(receipts, event)
}

func (b *Bridge) saveBlock(blockNumber uint64, blocksMap map[uint64]*c.CheckAuraBlockAura) error {
	if _, ok := blocksMap[blockNumber]; ok {
		return nil
	}

	block, err := b.HeaderByNumber(big.NewInt(int64(blockNumber)))
	if err != nil {
		return fmt.Errorf("HeaderByNumber: %w", err)
	}
	encodedBlock, err := EncodeBlock(block)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	blocksMap[blockNumber] = encodedBlock
	return nil
}

func (b *Bridge) getLastProcessedBlockNum() (*big.Int, error) {
	blockHash, err := b.sideBridge.GetLastProcessedBlockHash()
	if err != nil {
		return nil, fmt.Errorf("GetLastProcessedBlockHash: %w", err)
	}

	block, err := b.Client.BlockByHash(context.Background(), *blockHash)
	if err != nil {
		return nil, fmt.Errorf("get block by hash: %w", err)
	}

	return block.Number(), nil
}

func deltaVS(prev, curr []common.Address) (common.Address, int64, error) {
	d := len(curr) - len(prev)
	if d != 1 && d != -1 {
		return common.Address{}, 0, fmt.Errorf("delta has more (or less) than 1 change")
	}

	for i, prevEl := range prev {
		if i >= len(curr) { // deleted at the end
			return prev[i], int64(-i - 1), nil
		}

		if curr[i] != prevEl {
			if d == 1 { // added
				return curr[i], int64(i), nil
			} else { // deleted
				return prev[i], int64(-i - 1), nil
			}
		}
	}

	// add at the end
	i := len(curr) - 1
	return curr[i], int64(i), nil

	//return common.Address{}, 0, fmt.Errorf("this error shouln't exist")
}

// used for 'ordered' map
func sortedKeys(m map[uint64]*c.CheckAuraBlockAura) []uint64 {
	keys := make([]uint64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	return keys
}
