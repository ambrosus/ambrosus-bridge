package aura_proof

import (
	"context"
	"fmt"
	"math/big"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/amb"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/parity"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/rs/zerolog"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type blockExt struct {
	block             *c.CheckAuraBlockAura
	finalizedVsEvents []c.CheckAuraValidatorSetChange
	lastEvent         *c.VsInitiateChange
}

type AuraEncoder struct {
	bridge       *amb.Bridge
	auraReceiver networks.BridgeReceiveAura

	vsContract   *c.Vs
	parityClient *parity.Client

	logger *zerolog.Logger
}

func NewAuraEncoder(bridge *amb.Bridge, sideBridge networks.BridgeReceiveAura, vSContract *c.Vs, parityClient *parity.Client) *AuraEncoder {
	return &AuraEncoder{
		bridge:       bridge,
		auraReceiver: sideBridge,
		vsContract:   vSContract,
		parityClient: parityClient,
		logger:       bridge.GetLogger(), // todo maybe sublogger?
	}
}

func (b *AuraEncoder) EncodeAuraProof(transferEvent *c.BridgeTransfer, safetyBlocks uint64) (*c.CheckAuraAuraProof, error) {
	// populated by functions below
	blocksMap := make(map[uint64]*blockExt)

	// encode transferProof and save event block to blocksMap
	transfer, err := b.encodeTransferEvent(blocksMap, transferEvent, safetyBlocks)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferEvent: %w", err)
	}

	// encode vsChangeProofs and save event blocks to blocksMap
	vsChangeEvents, err := b.fetchVSChangeEvents(transferEvent, safetyBlocks)
	if err != nil {
		return nil, fmt.Errorf("fetchVSChangeEvents: %w", err)
	}
	err = b.encodeVSChangeEvents(blocksMap, vsChangeEvents)
	if err != nil {
		return nil, fmt.Errorf("encodeVSChangeEvents: %w", err)
	}

	// sort blocks in blocksMap and use resulting indexes
	indexToBlockNum := helpers.SortedKeys(blocksMap)
	var blocks []c.CheckAuraBlockAura
	var vsChanges []c.CheckAuraValidatorSetProof
	var transferEventIndex uint64
	blocksToIndex := make(map[uint64]uint64) // need for `EventBlock` in vsChange

	for i, blockNum := range indexToBlockNum {
		blocksToIndex[blockNum] = uint64(i)

		// fill up 'blocks'
		blocks = append(blocks, *blocksMap[blockNum].block)

		// fill up 'vsChanges'
		if blocksMap[blockNum].lastEvent != nil {
			proof, err := b.GetProof(blocksMap[blockNum].lastEvent)
			if err != nil {
				return nil, fmt.Errorf("GetProof: %w", err)
			}
			vsChanges = append(vsChanges, c.CheckAuraValidatorSetProof{
				ReceiptProof: proof,
				Changes:      blocksMap[blockNum].finalizedVsEvents,
				EventBlock:   big.NewInt(int64(blocksToIndex[blocksMap[blockNum].lastEvent.Raw.BlockNumber])),
			})

			// in this block contract should finalize all events in vsChanges array up to `FinalizedVs` index
			blocks[i].FinalizedVs = uint64(len(vsChanges))
		}

		// set transferEventIndex to index in blocks array
		if blockNum == transferEvent.Raw.BlockNumber {
			transferEventIndex = uint64(i)
		}
	}
	return &c.CheckAuraAuraProof{
		Blocks:             blocks,
		Transfer:           *transfer,
		VsChanges:          vsChanges,
		TransferEventBlock: transferEventIndex,
	}, nil
}

// if proof too long we need to split it into smaller parts
func (b *AuraEncoder) splitVsChanges(proof *c.CheckAuraAuraProof) *c.CheckAuraAuraProof {
	b.logger.Warn().Int("blocks", len(proof.Blocks)).Msgf("Aura proof too long")
	blocks := proof.Blocks[:len(proof.Blocks)/2] // drop half of blocks
	// todo maybe keep ~3000 blocks instead of half

	// last block should be vsFinalize
	var vsChanges []c.CheckAuraValidatorSetProof
	for i := uint64(len(blocks)); i > 0; i-- {
		if blocks[i].FinalizedVs == 0 {
			blocks = blocks[:i+1]
			vsChanges = proof.VsChanges[:blocks[i].FinalizedVs+1]
		}
	}

	return &c.CheckAuraAuraProof{
		Blocks:    blocks,
		VsChanges: vsChanges,
	}
}

func (b *AuraEncoder) encodeTransferEvent(blocks map[uint64]*blockExt, event *c.BridgeTransfer, safetyBlocks uint64) (*c.CommonStructsTransferProof, error) {
	proof, err := b.bridge.GetProof(event)
	if err != nil {
		return nil, err
	}

	// save `safetyBlocks` blocks after event block
	if err = b.saveBlocksRange(blocks, event.Raw.BlockNumber, event.Raw.BlockNumber+safetyBlocks); err != nil {
		return nil, err
	}

	return &c.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func (b *AuraEncoder) encodeVSChangeEvents(blocks map[uint64]*blockExt, events []*c.VsInitiateChange) error {
	prevSet, err := b.auraReceiver.GetValidatorSet()
	if err != nil {
		return fmt.Errorf("GetValidatorSet: %w", err)
	}

	minSafetyBlocksValidators, err := b.auraReceiver.GetMinSafetyBlocksValidators()
	if err != nil {
		return fmt.Errorf("GetMinSafetyBlocksValidators: %w", err)
	}

	var lastBlock uint64
	var txsBeforeFinalize uint64
	for _, event := range events {
		address, index, err := deltaVS(prevSet, event.NewSet)
		if err != nil {
			return fmt.Errorf("deltaVS: %w", err)
		}
		vsChange := c.CheckAuraValidatorSetChange{DeltaAddress: address, DeltaIndex: index}

		if lastBlock != event.Raw.BlockNumber {
			txsBeforeFinalize = uint64(len(prevSet))/2 + 1
			lastBlock = event.Raw.BlockNumber
		}
		finalizedBlockNum := event.Raw.BlockNumber + txsBeforeFinalize
		safetyEndBlockNum := event.Raw.BlockNumber + minSafetyBlocksValidators

		// save blocks up to finalized block
		if err = b.saveBlocksRange(blocks, event.Raw.BlockNumber, safetyEndBlockNum); err != nil {
			return err
		}

		// block in which VS will be finalized
		if err := b.saveBlock(blocks, finalizedBlockNum); err != nil {
			return err
		}
		blockWhenFinalize := blocks[finalizedBlockNum]
		blockWhenFinalize.finalizedVsEvents = append(blockWhenFinalize.finalizedVsEvents, vsChange)
		blockWhenFinalize.lastEvent = event

		prevSet = event.NewSet
	}
	return nil
}

func (b *AuraEncoder) fetchVSChangeEvents(event *c.BridgeTransfer, safetyBlocks uint64) ([]*c.VsInitiateChange, error) {
	start, err := b.getLastProcessedBlockNum()
	if err != nil {
		return nil, fmt.Errorf("getLastProcessedBlockNum: %w", err)
	}
	end := event.Raw.BlockNumber + safetyBlocks - 1 // no need to change vs for last block (it will be changed in next event processing)

	opts := &bind.FilterOpts{
		Start:   start.Uint64(),
		End:     &end,
		Context: context.Background(),
	}

	logs, err := b.vsContract.FilterInitiateChange(opts, nil)
	if err != nil {
		return nil, fmt.Errorf("filter initiate changes: %w", err)
	}

	var res []*c.VsInitiateChange
	for logs.Next() {
		res = append(res, logs.Event)
	}

	return res, nil
}

func (b *AuraEncoder) getLastProcessedBlockNum() (*big.Int, error) {
	blockHash, err := b.auraReceiver.GetLastProcessedBlockHash()
	if err != nil {
		return nil, fmt.Errorf("GetLastProcessedBlockHash: %w", err)
	}

	block, err := b.bridge.GetClient().BlockByHash(context.Background(), *blockHash)
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

	// return common.Address{}, 0, fmt.Errorf("this error shouln't exist")
}

// save blocks from `from` to `to` INCLUSIVE
func (b *AuraEncoder) saveBlocksRange(blocksMap map[uint64]*blockExt, from, to uint64) error {
	for i := from; i <= to; i++ {
		if err := b.saveBlock(blocksMap, i); err != nil {
			return err
		}
	}
	return nil
}

func (b *AuraEncoder) saveBlock(blocksMap map[uint64]*blockExt, blockNumber uint64) error {
	if _, ok := blocksMap[blockNumber]; ok {
		return nil
	}

	block, err := b.parityClient.ParityHeaderByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		return fmt.Errorf("HeaderByNumber: %w", err)
	}
	encodedBlock, err := EncodeBlock(block)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}

	blocksMap[blockNumber] = &blockExt{block: encodedBlock}
	return nil
}
