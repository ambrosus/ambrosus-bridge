package amb

import (
	"context"
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

// todo эти структуры сгенерируются абигеном и будут в пакете contracts, но пока так

type AmbTransfer struct {
	Blocks    []*contracts.CheckPoABlockPoA
	Transfer  *Transfer_Event
	VsChanges []*ValidatorSet_Event
}

type ValidatorSet_Event struct {
	Receipt_proof contracts.ReceiptsProof
	Delta_address common.Address
	Delta_index   uint64 // 12байт шоб спаковать с 20байт адресом
}

type Transfer_Event struct {
	Receipt_proof contracts.ReceiptsProof
	Event_id      *big.Int
	Transfers     []contracts.CommonStructsTransfer
}

// todo name
func (b *Bridge) getBlocksAndEvents(transferEvent *contracts.TransferEvent) (*AmbTransfer, error) {
	// populated by functions below
	blocksMap := make(map[uint64]*contracts.CheckPoABlockPoA)

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
			blType := BlTypeSafety
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

	blocks := make([]*contracts.CheckPoABlockPoA, 0, len(blocksMap))
	for _, block := range blocksMap {
		blocks = append(blocks, block)
	}

	return &AmbTransfer{
		Blocks:    blocks,
		Transfer:  transfer,
		VsChanges: vsChanges,
	}, nil
}

func (b *Bridge) encodeTransferEvent(blocks map[uint64]*contracts.CheckPoABlockPoA, event *contracts.TransferEvent) (*Transfer_Event, error) {
	proof, err := b.getProof(&event.Raw)
	if err != nil {
		return nil, err
	}

	blocks[event.Raw.BlockNumber], err = b.encodeBlockWithType(event.Raw.BlockNumber, BlTypeTransfer)
	if err != nil {
		return nil, err
	}

	return &Transfer_Event{
		Receipt_proof: proof,
		Event_id:      event.EventId,
		Transfers:     event.Queue,
	}, nil
}

func (b *Bridge) encodeVSChangeEvents(blocks map[uint64]*contracts.CheckPoABlockPoA, events []*contracts.InitiateChange) ([]*ValidatorSet_Event, error) {
	vsChanges := make([]*ValidatorSet_Event, 0, len(events))

	var prev_event *contracts.CheckPoABlockPoA // todo VS_0 state

	for i, event := range events {
		encodedEvent, err := b.encodeVSChangeEvent(prev_event, event)
		if err != nil {
			return nil, err
		}
		vsChanges[i] = encodedEvent
		prev_event = event

		blocks[event.Raw.BlockNumber], err = b.encodeBlockWithType(event.Raw.BlockNumber, i)
		if err != nil {
			return nil, err
		}
	}
	return vsChanges, nil
}

func (b *Bridge) encodeVSChangeEvent(prev_event, event *contracts.InitiateChange) (*ValidatorSet_Event, error) {
	// todo delta

	proof, err := b.getProof(&event.Raw)
	if err != nil {
		return nil, err
	}

	return &ValidatorSet_Event{
		Receipt_proof: proof,
	}, nil
}

func (b *Bridge) getVSChangeEvents(event *contracts.TransferEvent) ([]*contracts.InitiateChange, error) {
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
	logs, err := b.Contract.FilterInitiateChange(opts)
	if err != nil {
		return nil, err
	}

	var res []*ValidatorSet_Event
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

func (b *Bridge) encodeBlockWithType(blockNumber uint64, type_ int) (*contracts.CheckPoABlockPoA, error) {
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
