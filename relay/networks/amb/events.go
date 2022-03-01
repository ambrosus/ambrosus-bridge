package amb

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ambrosus/ambrosus-bridge/relay/receipts_proof"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

const (
	BlTypeSafetyEnd = -3
	BlTypeSafety    = -2
	BlTypeTransfer  = -1
	// type values >= 0 mean the validatorSetEvent index in VsChanges list
)

// todo name
func (b *Bridge) getBlocksAndEvents(transferEvent *contracts.TransferEvent) (*contracts.CheckAuraAuraProof, error) {
	// populated by functions below
	blocksMap := make(map[uint64]contracts.CheckAuraBlockAura)

	transfer, err := b.encodeTransferEvent(blocksMap, transferEvent)
	if err != nil {
		return nil, err
	}

	vsChangeEvents, err := b.getVSChangeEvents(transferEvent)
	if err != nil {
		return nil, err
	}
	vsChanges, err := b.encodeVSChangeEvents(blocksMap, vsChangeEvents)
	if err != nil {
		return nil, err
	}

	// add safety blocks after each event block
	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return nil, err
	}
	for blockNum := range blocksMap {
		for i := uint64(0); i <= safetyBlocks; i++ {
			targetBlockNum := blockNum + i

			// set block type == safety; need to explicitly specify if this is the end of safety chain
			blType := int64(BlTypeSafety)
			if i == safetyBlocks {
				blType = BlTypeSafetyEnd
			}

			if bl, ok := blocksMap[targetBlockNum]; ok {
				// if the block existed and was the end of safety chain, then that could change now
				if bl.Type.Int64() == BlTypeSafetyEnd {
					bl.Type = big.NewInt(blType)
				}
			} else {
				// save block as safety
				encodedBlockWithType, err := b.encodeBlockWithType(targetBlockNum, blType)
				if err != nil {
					return nil, err
				}
				blocksMap[targetBlockNum] = *encodedBlockWithType
			}

		}
	}

	blocks := make([]contracts.CheckAuraBlockAura, 0, len(blocksMap))
	for _, block := range blocksMap {
		blocks = append(blocks, block)
	}

	return &contracts.CheckAuraAuraProof{
		Blocks:    blocks,
		Transfer:  *transfer,
		VsChanges: vsChanges,
	}, nil
}

func (b *Bridge) encodeTransferEvent(blocks map[uint64]contracts.CheckAuraBlockAura, event *contracts.TransferEvent) (*contracts.CommonStructsTransferProof, error) {
	proof, err := b.getProof(event)
	if err != nil {
		return nil, err
	}

	encodedBlockWithType, err := b.encodeBlockWithType(event.Raw.BlockNumber, BlTypeTransfer)
	if err != nil {
		return nil, err
	}
	blocks[event.Raw.BlockNumber] = *encodedBlockWithType

	return &contracts.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func (b *Bridge) encodeVSChangeEvents(blocks map[uint64]contracts.CheckAuraBlockAura, events []*contracts.VsInitiateChange) ([]contracts.CheckAuraValidatorSetProof, error) {
	vsChanges := make([]contracts.CheckAuraValidatorSetProof, 0, len(events))

	prevSet, err := b.sideBridge.GetValidatorSet()
	if err != nil {
		return nil, err
	}

	for i, event := range events {
		encodedEvent, err := b.encodeVSChangeEvent(prevSet, event)
		if err != nil {
			return nil, err
		}
		vsChanges[i] = *encodedEvent
		prevSet = event.NewSet

		encodedBlockWithType, err := b.encodeBlockWithType(event.Raw.BlockNumber, int64(i))
		if err != nil {
			return nil, err
		}
		blocks[event.Raw.BlockNumber] = *encodedBlockWithType
	}
	return vsChanges, nil
}

func (b *Bridge) encodeVSChangeEvent(prevSet []common.Address, event *contracts.VsInitiateChange) (*contracts.CheckAuraValidatorSetProof, error) {
	address, index, err := deltaVS(prevSet, event.NewSet)
	if err != nil {
		return nil, err
	}

	proof, err := b.getProof(event)
	if err != nil {
		return nil, err
	}

	return &contracts.CheckAuraValidatorSetProof{
		ReceiptProof: proof,
		DeltaAddress: address,
		DeltaIndex:   index,
	}, nil
}

func (b *Bridge) getVSChangeEvents(event *contracts.TransferEvent) ([]*contracts.VsInitiateChange, error) {
	prevEventId := big.NewInt(0).Sub(event.EventId, big.NewInt(1))
	prevEvent, err := b.GetEventById(prevEventId)
	if err != nil {
		return nil, err
	}

	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return nil, err
	}

	start := prevEvent.Raw.BlockNumber
	end := event.Raw.BlockNumber + safetyBlocks

	opts := &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}

	logs, err := b.VSContract.FilterInitiateChange(opts, nil)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return receipts_proof.CalcProofEvent(receipts, event)
}

func (b *Bridge) encodeBlockWithType(blockNumber uint64, type_ int64) (*contracts.CheckAuraBlockAura, error) {
	block, err := b.HeaderByNumber(big.NewInt(int64(blockNumber)))
	if err != nil {
		return nil, err
	}
	encodedBlock, err := EncodeBlock(block)
	if err != nil {
		return nil, err
	}
	encodedBlock.Type = big.NewInt(type_)
	return encodedBlock, nil
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
