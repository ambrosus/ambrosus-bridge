package amb

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/receipts_proof"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
	blocksMap := make(map[uint64]*contracts.CheckAuraBlockAura)

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
	for blockNum := range blocksMap {
		for i := uint64(0); i < b.config.SafetyBlocks; i++ {
			targetBlockNum := blockNum + i

			// set block type == safety; need to explicitly specify if this is the end of safety chain
			blType := int64(BlTypeSafety)
			if i == b.config.SafetyBlocks {
				blType = BlTypeSafetyEnd
			}

			if bl, ok := blocksMap[targetBlockNum]; ok {
				// if the block existed and was the end of safety chain, then that could change now
				if bl.Type == BlTypeSafetyEnd {
					bl.Type = blType
				}
			} else {
				// save block as safety
				blocksMap[targetBlockNum], err = b.encodeBlockWithType(targetBlockNum, blType)
				if err != nil {
					return nil, err
				}
			}

		}
	}

	blocks := make([]*contracts.CheckAuraBlockAura, 0, len(blocksMap))
	for _, block := range blocksMap {
		blocks = append(blocks, block)
	}

	return &contracts.CheckAuraAuraProof{
		Blocks:    blocks,
		Transfer:  transfer,
		VsChanges: vsChanges,
	}, nil
}

func (b *Bridge) encodeTransferEvent(blocks map[uint64]*contracts.CheckAuraBlockAura, event *contracts.TransferEvent) (*contracts.CommonStructsTransferProof, error) {
	proof, err := b.getProof(&event.Raw)
	if err != nil {
		return nil, err
	}

	blocks[event.Raw.BlockNumber], err = b.encodeBlockWithType(event.Raw.BlockNumber, BlTypeTransfer)
	if err != nil {
		return nil, err
	}

	return &contracts.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func (b *Bridge) encodeVSChangeEvents(blocks map[uint64]*contracts.CheckAuraBlockAura, events []*contracts.VsInitiateChange) ([]*contracts.CheckAuraValidatorSetProof, error) {
	vsChanges := make([]*contracts.CheckAuraValidatorSetProof, 0, len(events))

	var prev_event *contracts.VsInitiateChange // todo VS_0 state

	for i, event := range events {
		encodedEvent, err := b.encodeVSChangeEvent(prev_event, event)
		if err != nil {
			return nil, err
		}
		vsChanges[i] = encodedEvent
		prev_event = event

		blocks[event.Raw.BlockNumber], err = b.encodeBlockWithType(event.Raw.BlockNumber, int64(i))
		if err != nil {
			return nil, err
		}
	}
	return vsChanges, nil
}

func (b *Bridge) encodeVSChangeEvent(prevEvent, event *contracts.VsInitiateChange) (*contracts.CheckAuraValidatorSetProof, error) {
	address, index, err := deltaVS(prevEvent.NewSet, event.NewSet)
	if err != nil {
		return nil, err
	}

	proof, err := b.getProof(&event.Raw)
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

	start := prevEvent.Raw.BlockNumber
	end := event.Raw.BlockNumber + b.config.SafetyBlocks

	opts := &bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}

	// todo validator set is separate contract
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

func (b *Bridge) getProof(log *types.Log) ([][]byte, error) {
	receipts, err := b.GetReceipts(log.BlockHash)
	if err != nil {
		return nil, err
	}
	return receipts_proof.CalcProof(receipts, log)
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
	encodedBlock.Type = type_
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
