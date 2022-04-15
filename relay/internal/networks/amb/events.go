package amb

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sort"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// todo name
func (b *Bridge) getBlocksAndEvents(transferEvent *contracts.BridgeTransfer) (*contracts.CheckAuraAuraProof, error) {
	// populated by functions below
	blocksMap := make(map[uint64]contracts.CheckAuraBlockAura)

	transfer, err := b.encodeTransferEvent(blocksMap, transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferEvent: %w", err)
	}

	vsChangeEvents, err := b.getVSChangeEvents(transferEvent)
	if err != nil {
		return nil, fmt.Errorf("getVSChangeEvents: %w", err)
	}
	vsChanges, err := b.encodeVSChangeEvents(blocksMap, vsChangeEvents)
	if err != nil {
		return nil, fmt.Errorf("encodeVSChangeEvents: %w", err)
	}

	// add safety blocks after each event block
	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return nil, fmt.Errorf("getMinSafetyBlocksNum: %w", err)
	}

	blockNums := sortedKeys(blocksMap)
	for _, blockNum := range blockNums {
		for i := uint64(0); i <= safetyBlocks; i++ {
			targetBlockNum := blockNum + i

			// set block type == safety; need to explicitly specify if this is the end of safety chain
			blType := BlTypeSafety
			if i == safetyBlocks {
				blType = BlTypeSafetyEnd
			}

			if bl, ok := blocksMap[targetBlockNum]; ok {
				// if the block existed and was the end of safety chain, then that could change now
				if bl.Type&BlTypeSafetyEnd != 0 {
					bl.Type |= blType
				}
			} else {
				// save block as safety
				encodedBlockWithType, err := b.encodeBlockWithType(targetBlockNum, blType)
				if err != nil {
					return nil, fmt.Errorf("encode block as safety: %w", err)
				}
				blocksMap[targetBlockNum] = *encodedBlockWithType
			}

		}
	}

	blocks := make([]contracts.CheckAuraBlockAura, len(blocksMap))
	indexToBlockNum := sortedKeys(blocksMap)
	var transferEventIndex uint64

	for i, blockNum := range indexToBlockNum {
		blocks[i] = blocksMap[blockNum]

		// change FinalizeVs (if it set) from offset to absolute index in blocks array
		if blocks[i].FinalizedVs != 0 {
			blocks[i].FinalizedVs = uint64(i) - blocks[i].FinalizedVs
		}

		// set transferEventIndex to index in blocks array
		if blockNum == transferEvent.Raw.BlockNumber {
			transferEventIndex = uint64(i)
		}
	}

	return &contracts.CheckAuraAuraProof{
		Blocks:             blocks,
		Transfer:           *transfer,
		VsChanges:          vsChanges,
		TransferEventBlock: transferEventIndex,
	}, nil
}

func (b *Bridge) encodeTransferEvent(blocks map[uint64]contracts.CheckAuraBlockAura, event *contracts.BridgeTransfer) (*contracts.CommonStructsTransferProof, error) {
	proof, err := b.getProof(event)
	if err != nil {
		return nil, err
	}

	encodedBlockWithType, err := b.encodeBlockWithType(event.Raw.BlockNumber, BlTypeTransfer)
	if err != nil {
		return nil, fmt.Errorf("encode block as transfer: %w", err)
	}
	blocks[event.Raw.BlockNumber] = *encodedBlockWithType

	return &contracts.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func (b *Bridge) encodeVSChangeEvents(blocks map[uint64]contracts.CheckAuraBlockAura, events []*contracts.VsInitiateChange) ([]contracts.CheckAuraValidatorSetProof, error) {
	vsChanges := make([]contracts.CheckAuraValidatorSetProof, len(events))

	prevSet, err := b.sideBridge.GetValidatorSet()
	if err != nil {
		return nil, fmt.Errorf("GetValidatorSet: %w", err)
	}

	for i, event := range events {
		encodedEvent, err := b.encodeVSChangeEvent(prevSet, event)
		if err != nil {
			return nil, fmt.Errorf("encodeVSChangeEvent: %w", err)
		}
		vsChanges[i] = *encodedEvent
		prevSet = event.NewSet

		if _, ok := blocks[event.Raw.BlockNumber]; !ok {
			bl, err = b.encodeBlock(event.Raw.BlockNumber)
			if err != nil {
				return nil, err
			}
			blocks[event.Raw.BlockNumber] = *bl
		}

		bl := blocks[event.Raw.BlockNumber]
		// FinalizeVs field must be block **index** in CheckAuraAuraProof.blocks, but we don't have it now
		// so we set **offset** value here and will change it later
		bl.FinalizedVs = 1 // vs_change finalize on next block, offset always 1

	}
	return vsChanges, nil
}

func (b *Bridge) encodeVSChangeEvent(prevSet []common.Address, event *contracts.VsInitiateChange) (*contracts.CheckAuraValidatorSetProof, error) {
	address, index, err := deltaVS(prevSet, event.NewSet)
	if err != nil {
		return nil, fmt.Errorf("deltaVS: %w", err)
	}

	proof, err := b.getProof(event)
	if err != nil {
		return nil, fmt.Errorf("getProof: %w", err)
	}

	return &contracts.CheckAuraValidatorSetProof{
		ReceiptProof: proof,
		DeltaAddress: address,
		DeltaIndex:   index,
	}, nil
}

// add safety blocks after each event block
func (b *Bridge) addSafetyBlocks(blocksMap map[uint64]contracts.CheckAuraBlockAura, minSafetyBlocks uint64) (err error) {

	for blockNum, _ := range blocksMap {
		for i := uint64(0); i <= minSafetyBlocks; i++ {
			targetBlockNum := blockNum + i

			if _, ok := blocksMap[targetBlockNum]; ok {
				continue
			}

			bl, err := b.encodeBlock(targetBlockNum)
			if err != nil {
				return err
			}
			blocksMap[targetBlockNum] = *bl

		}
	}

	return nil
}

func (b *Bridge) fetchVSChangeEvents(event *contracts.BridgeTransfer) ([]*contracts.VsInitiateChange, error) {
	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return nil, fmt.Errorf("getMinSafetyBlocksNum: %w", err)
	}

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

	var res []*contracts.VsInitiateChange
	for logs.Next() {
		res = append(res, logs.Event)
	}

	return res, nil
}

func (b *Bridge) getProof(event receipts_proof.ProofEvent) ([][]byte, error) {
	receipts, err := ethereum.GetReceipts(b.Client, event.Log().BlockHash)
	if err != nil {
		return nil, fmt.Errorf("GetReceipts: %w", err)
	}
	return receipts_proof.CalcProofEvent(receipts, event)
}

func (b *Bridge) encodeBlock(blockNumber uint64) (*contracts.CheckAuraBlockAura, error) {
	block, err := b.HeaderByNumber(big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, fmt.Errorf("HeaderByNumber: %w", err)
	}
	encodedBlock, err := EncodeBlock(block)
	if err != nil {
		return nil, fmt.Errorf("encode: %w", err)
	}
	return encodedBlock, nil
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
	if math.Abs(float64(d)) != 1 {
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

	return common.Address{}, 0, fmt.Errorf("this error shouln't exist")
}

// used for 'ordered' map
func sortedKeys(m map[uint64]contracts.CheckAuraBlockAura) []uint64 {
	keys := make([]uint64, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}
