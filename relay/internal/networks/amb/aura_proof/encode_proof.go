package aura_proof

import (
	"context"
	"fmt"
	"math/big"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/parity"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/rs/zerolog"
)

type AuraEncoder struct {
	bridge       networks.Bridge
	auraReceiver networks.BridgeReceiveAura

	vsContract   *c.Vs
	parityClient *parity.Client

	logger *zerolog.Logger

	// cache for fetching blocks. clear (almost) every `EncodeAuraProof` call
	fetchBlockCache func(arg uint64) (*parity.Header, error)
}

func NewAuraEncoder(bridge networks.Bridge, sideBridge networks.BridgeReceiveAura, vSContract *c.Vs, parityClient *parity.Client) *AuraEncoder {
	return &AuraEncoder{
		bridge:       bridge,
		auraReceiver: sideBridge,
		vsContract:   vSContract,
		parityClient: parityClient,
		logger:       bridge.GetLogger(), // todo maybe sublogger?
	}
}

func (b *AuraEncoder) EncodeAuraProof(transferEvent *c.BridgeTransfer, safetyBlocks uint64) (*c.CheckAuraAuraProof, error) {
	// todo arg for split (reduced size proof)

	// new cache for every call
	// todo don't clear cache for reduced size proof
	b.fetchBlockCache = helpers.NewCache(b.fetchBlock)

	var blocksToSave []uint64

	lastBlock := transferEvent.Raw.BlockNumber + safetyBlocks

	// todo don't encode and save blocks for reduced size proof
	transferProof, err := b.encodeTransferProof(transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferProof: %w", err)
	}

	// save blocks for transfer
	blocksToSave = append(blocksToSave, helpers.Range(transferEvent.Raw.BlockNumber, lastBlock+1)...)

	// lastBlock can be decreased if proof is too big
	// todo better name
	vsChangesExt, err := b.getVsChanges(lastBlock)
	if err != nil {
		return nil, fmt.Errorf("getVsChanges: %w", err)
	}

	// save blocks for vs change events
	safetyBlocksValidators, err := b.auraReceiver.GetMinSafetyBlocksValidators()
	if err != nil {
		return nil, fmt.Errorf("GetMinSafetyBlocksValidators: %w", err)
	}
	for _, vsChange := range vsChangesExt {
		blocksToSave = append(blocksToSave, helpers.Range(vsChange.eventBlock, vsChange.eventBlock+safetyBlocksValidators+1)...)
		// gap
		blocksToSave = append(blocksToSave, vsChange.finalizedBlock)
	}

	// fetch and encode blocksToSave
	blocks, blockNumToIndex, err := b.saveEncodedBlocks(blocksToSave)
	if err != nil {
		return nil, fmt.Errorf("saveEncodedBlocks: %w", err)
	}

	// for each finalized event save it into `vsChanges` and set it `index+1` to block `FinalizedVs` field
	var vsChanges []c.CheckAuraValidatorSetProof
	for _, blockWithEvent := range helpers.SortedKeys(vsChangesExt) {
		vsChangeEvent := vsChangesExt[blockWithEvent]
		finalizedBlockIndex := blockNumToIndex[vsChangeEvent.finalizedBlock]

		proof, err := cb.GetProof(b.bridge.GetClient(), vsChangeEvent.lastEvent)
		if err != nil {
			return nil, fmt.Errorf("GetProof: %w", err)
		}
		vsChanges = append(vsChanges, c.CheckAuraValidatorSetProof{
			ReceiptProof: proof,
			Changes:      vsChangeEvent.changes,
			EventBlock:   big.NewInt(int64(blockNumToIndex[vsChangeEvent.eventBlock])), // todo uint64
		})

		// in this Block contract should finalize all events in vsChanges array up to `FinalizedVs` index
		blocks[finalizedBlockIndex].FinalizedVs = uint64(len(vsChanges))
	}

	return &c.CheckAuraAuraProof{
		Blocks:             blocks,
		Transfer:           *transferProof,
		VsChanges:          vsChanges,
		TransferEventBlock: uint64(blockNumToIndex[transferEvent.Raw.BlockNumber]),
	}, nil

}

func (b *AuraEncoder) encodeTransferProof(event *c.BridgeTransfer) (*c.CommonStructsTransferProof, error) {
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

func (b *AuraEncoder) saveEncodedBlocks(blockNums []uint64) (blocks []c.CheckAuraBlockAura, blockNumToIndex map[uint64]int, err error) {
	blocks = make([]c.CheckAuraBlockAura, len(blockNums))

	for i, bn := range helpers.Sorted(blockNums) {
		block, err := b.fetchBlockCache(bn)
		if err != nil {
			return nil, nil, fmt.Errorf("fetchBlockCache: %w", err)
		}
		encodedBlock, err := EncodeBlock(block)
		if err != nil {
			return nil, nil, fmt.Errorf("EncodeBlock: %w", err)
		}

		blocks[i] = *encodedBlock
		blockNumToIndex[bn] = i
	}

	return blocks, blockNumToIndex, nil
}

func (b *AuraEncoder) fetchBlock(blockNum uint64) (*parity.Header, error) {
	return b.parityClient.ParityHeaderByNumber(context.Background(), big.NewInt(int64(blockNum)))
}
