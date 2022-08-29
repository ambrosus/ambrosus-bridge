package aura_proof

import (
	"context"
	"fmt"
	"math/big"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients/parity"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/rs/zerolog"
)

const gapForExtraData = 1024 * 10

type finalizeService interface {
	GetBlockWhenFinalize(emitBlockNum uint64) (uint64, error)
}

type AuraEncoder struct {
	bridge       networks.Bridge
	auraReceiver service_submit.ReceiverAura

	vsContract              *c.Vs
	parityClient            *parity.Client
	vsContract              *c.Vs
	parityClient            *parity.Client
	finalizeService         finalizeService
	receiverBridgeMaxTxSize uint64

	logger *zerolog.Logger

	// cache for fetching blocks. clear every `EncodeAuraProof` call
	fetchBlockCache func(arg uint64) (*parity.Header, error)
}

func NewAuraEncoder(bridge networks.Bridge, sideBridge service_submit.ReceiverAura,
	vSContract *c.Vs, parityClient *parity.Client, finalizeService finalizeService, receiverBridgeMaxTxSizeKB uint64) *AuraEncoder {
	logger := bridge.GetLogger().With().Str("service", "AuraEncoder").Logger()

	return &AuraEncoder{
		bridge:                  bridge,
		auraReceiver:            sideBridge,
		vsContract:              vSContract,
		parityClient:            parityClient,
		finalizeService:         finalizeService,
		receiverBridgeMaxTxSize: receiverBridgeMaxTxSizeKB * 1024,
		logger:                  &logger,
	}
}

func (e *AuraEncoder) EncodeAuraProof(transferEvent *c.BridgeTransfer, safetyBlocks uint64) (*c.CheckAuraAuraProof, error) {
	// new cache for every call
	e.fetchBlockCache = helpers.NewCache(e.fetchBlock)

	lastBlock := transferEvent.Raw.BlockNumber + safetyBlocks

	safetyBlocksValidators, err := e.auraReceiver.GetMinSafetyBlocksValidators()
	if err != nil {
		return nil, fmt.Errorf("GetMinSafetyBlocksValidators: %w", err)
	}

	vsProofs, err := e.getVsChanges(lastBlock)
	if err != nil {
		return nil, fmt.Errorf("getVsChanges: %w", err)
	}

	var blocksMap = make(map[uint64]c.CheckAuraBlockAura)
	var vsChanges []c.CheckAuraValidatorSetProof
	var transferProof = c.CommonStructsTransferProof{
		ReceiptProof: [][]byte{},
		EventId:      big.NewInt(0),
		Transfers:    []c.CommonStructsTransfer{},
	}

	// this func use so many local variables, so it's better to be closure
	buildProof := func() *c.CheckAuraAuraProof {
		blocks, blockNumToIndex := helpers.SortedValuesWithIndices(blocksMap)
		// set indexes
		for i := 0; i < len(vsChanges); i++ {
			// in this block contract should finalize vsChanges[FinalizedVs-1] event
			finalizedBlockIndex := blockNumToIndex[vsProofs[i].finalizedBlock]
			blocks[finalizedBlockIndex].FinalizedVs = uint64(i + 1)

			vsChanges[i].EventBlock = uint64(blockNumToIndex[vsProofs[i].EventBlock])
		}
		return &c.CheckAuraAuraProof{
			Blocks:             blocks,
			Transfer:           transferProof,
			VsChanges:          vsChanges,
			TransferEventBlock: uint64(blockNumToIndex[transferEvent.Raw.BlockNumber]),
		}
	}

	proof := buildProof()

	// add vsChange events one by one to proof
	// return acceptable size proof if new proof is too big
	for _, vsChange := range vsProofs {
		receiptProof, err := cb.GetProof(e.bridge.GetClient(), vsChange.lastEvent)
		if err != nil {
			return nil, fmt.Errorf("GetProof: %w", err)
		}
		vsChange.CheckAuraValidatorSetProof.ReceiptProof = receiptProof
		vsChanges = append(vsChanges, vsChange.CheckAuraValidatorSetProof)

		blocksToSave := append(helpers.Range(vsChange.EventBlock, vsChange.EventBlock+safetyBlocksValidators+1), vsChange.finalizedBlock)
		if err := e.saveBlocks(blocksMap, blocksToSave...); err != nil {
			return nil, fmt.Errorf("saveBlocks: %w", err)
		}

		newProof := buildProof()
		if err = c.IsProofTooBig(newProof, e.getMaxAllowedProofSize()); err != nil {
			proof.TransferEventBlock = ^uint64(0) // max uint64, coz gaps between vsChanges work only BEFORE `TransferEventBlock`
			return proof, err
		}
		proof = newProof
	}

	// add transfer event to proof
	// return proof without transferEvent if new proof is too big
	transferProof, err = cb.EncodeTransferProof(e.bridge.GetClient(), transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferProof: %w", err)
	}
	blocksToSave := helpers.Range(transferEvent.Raw.BlockNumber, transferEvent.Raw.BlockNumber+safetyBlocks+1)
	if err := e.saveBlocks(blocksMap, blocksToSave...); err != nil {
		return nil, fmt.Errorf("saveBlocks: %w", err)
	}

	newProof := buildProof()
	if err = c.IsProofTooBig(newProof, e.getMaxAllowedProofSize()); err != nil {
		proof.TransferEventBlock = ^uint64(0) // max uint64, coz gaps between vsChanges work only BEFORE `TransferEventBlock`
		return proof, err
	}
	return newProof, nil
}

func (e *AuraEncoder) saveBlocks(blocksMap map[uint64]c.CheckAuraBlockAura, blockNums ...uint64) error {
	for _, bn := range blockNums {
		if _, ok := blocksMap[bn]; !ok {
			block, err := e.fetchBlockCache(bn)
			if err != nil {
				return fmt.Errorf("fetchBlockCache: %w", err)
			}
			encodedBlock, err := EncodeBlock(block)
			if err != nil {
				return fmt.Errorf("EncodeBlock: %w", err)
			}
			blocksMap[bn] = *encodedBlock
		}
	}
	return nil
}

func (e *AuraEncoder) fetchBlock(blockNum uint64) (*parity.Header, error) {
	return e.parityClient.ParityHeaderByNumber(context.Background(), big.NewInt(int64(blockNum)))
}

func (e *AuraEncoder) getMaxAllowedProofSize() uint64 {
	return e.receiverBridgeMaxTxSize - gapForExtraData
}
