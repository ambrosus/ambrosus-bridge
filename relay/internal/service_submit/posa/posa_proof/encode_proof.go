package posa_proof

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	c "github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

const (
	addressLength     = 20
	extraVanityLength = 32
	extraSealLength   = 65
	epochLength       = 200
)

const maxRequestContentLength = 1024*1024*5 - 10240 // -10KB for extra data in request
var ProofTooBig = errors.New("proof is too big")

type PoSAEncoder struct {
	bridge       networks.Bridge
	posaReceiver service_submit.ReceiverPoSA

	chainId *big.Int

	logger *zerolog.Logger

	// cache for fetching blocks. clear every `EncodePoSAProof` call
	fetchBlockCache func(arg uint64) (*types.Header, error)
}

type epochChangeS struct {
	blockWithSet, finalizeBlock uint64
}

func NewPoSAEncoder(bridge networks.Bridge, sideBridge service_submit.ReceiverPoSA, chainId *big.Int) *PoSAEncoder {
	logger := bridge.GetLogger().With().Str("service", "PoSAEncoder").Logger()

	return &PoSAEncoder{
		bridge:       bridge,
		posaReceiver: sideBridge,
		chainId:      chainId,
		logger:       &logger,
	}
}

func (e *PoSAEncoder) EncodePoSAProof(transferEvent *c.BridgeTransfer, safetyBlocks uint64) (*c.CheckPoSAPoSAProof, error) {
	// new cache for every call
	e.fetchBlockCache = helpers.NewCache(e.fetchBlock)

	lastBlock := transferEvent.Raw.BlockNumber + safetyBlocks

	// lastBlock can be decreased if proof is too big
	epochChanges, err := e.getEpochChanges(lastBlock)
	if err != nil {
		return nil, fmt.Errorf("getEpochChanges: %w", err)
	}

	var blocksMap map[uint64]c.CheckPoSABlockPoSA
	var transferProof c.CommonStructsTransferProof

	// this func use so many local variables, so it's better to be closure
	buildProof := func() *c.CheckPoSAPoSAProof {
		blocks, blockNumToIndex := helpers.SortedValuesWithIndices(blocksMap)
		return &c.CheckPoSAPoSAProof{
			Blocks:             blocks,
			Transfer:           transferProof,
			TransferEventBlock: uint64(blockNumToIndex[transferEvent.Raw.BlockNumber]),
		}
	}

	proof := buildProof()

	for _, epochChange := range epochChanges {
		blocksToSave := helpers.Range(epochChange.blockWithSet, epochChange.finalizeBlock)
		if err := e.saveBlocks(blocksMap, blocksToSave...); err != nil {
			return nil, fmt.Errorf("saveBlocks: %w", err)
		}

		newProof := buildProof()
		if err = isProofTooBig(newProof); err != nil {
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
	if err = isProofTooBig(newProof); err != nil {
		proof.TransferEventBlock = ^uint64(0) // max uint64, coz gaps between epoch changes work only BEFORE `TransferEventBlock`
		return proof, err
	}
	return newProof, nil
}

func (e *PoSAEncoder) getEpochChanges(end uint64) ([]epochChangeS, error) {
	var epochChanges []epochChangeS

	currentEpoch, err := e.posaReceiver.GetCurrentEpoch()
	if err != nil {
		return nil, fmt.Errorf("GetCurrentEpoch: %w", err)
	}
	currentSetLen, err := e.fetchVSLength(currentEpoch) // use to determine how many blocks need to save in next epoch
	if err != nil {
		return nil, fmt.Errorf("fetchVSLength: %w", err)
	}

	// save block nums (without current epoch)
	for epoch := currentEpoch + 1; epoch*epochLength < end; epoch++ {
		blockWithSet := epoch * epochLength
		finalizeBlock := blockWithSet + currentSetLen/2 + 1
		epochChanges = append(epochChanges, epochChangeS{blockWithSet, finalizeBlock})

		currentSetLen, err = e.fetchVSLength(epoch) // use to determine how many blocks need to save in next epoch
		if err != nil {
			return nil, fmt.Errorf("fetchVSLength: %w", err)
		}
	}
	return epochChanges, nil
}

func (e *PoSAEncoder) saveBlocks(blocksMap map[uint64]c.CheckPoSABlockPoSA, blockNums ...uint64) error {
	for _, bn := range blockNums {
		if _, ok := blocksMap[bn]; !ok {
			block, err := e.fetchBlockCache(bn)
			if err != nil {
				return fmt.Errorf("fetchBlockCache: %w", err)
			}
			encodedBlock, err := e.EncodeBlock(block)
			if err != nil {
				return fmt.Errorf("EncodeBlock: %w", err)
			}
			blocksMap[bn] = *encodedBlock
		}
	}
	return nil
}

func (e *PoSAEncoder) fetchBlock(blockNum uint64) (*types.Header, error) {
	return e.bridge.GetClient().HeaderByNumber(context.Background(), big.NewInt(int64(blockNum)))
}

func (e *PoSAEncoder) fetchVSLength(epoch uint64) (uint64, error) {
	block, err := e.fetchBlockCache(epoch * epochLength)
	if err != nil {
		return 0, fmt.Errorf("fetchBlockCache: %w", err)
	}
	return getVSLength(block), nil
}

func getVSLength(epochChangeBlock *types.Header) uint64 {
	validatorsLen := len(epochChangeBlock.Extra) - extraSealLength - extraVanityLength
	return uint64(validatorsLen) / addressLength
}

func isProofTooBig(proof *c.CheckPoSAPoSAProof) error {
	size, err := proof.Size()
	if err != nil {
		return fmt.Errorf("proof.Size(): %w", err)
	}
	// todo maxRequestContentLength depends on receiver network
	if size > maxRequestContentLength {
		return ProofTooBig
	}
	return nil
}
