package pow_proof

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_submit"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/rs/zerolog"
)

type PoWEncoder struct {
	bridge      networks.Bridge
	powReceiver service_submit.ReceiverPoW

	ethash *ethash.Ethash
	logger *zerolog.Logger
}

func NewPoWEncoder(bridge networks.Bridge, sideBridge service_submit.ReceiverPoW, ethash *ethash.Ethash) *PoWEncoder {
	return &PoWEncoder{
		bridge:      bridge,
		powReceiver: sideBridge,
		ethash:      ethash,
		logger:      bridge.GetLogger(), // todo maybe sublogger?
	}
}

func (e *PoWEncoder) EncodePoWProof(transferEvent *bindings.BridgeTransfer, safetyBlocks uint64) (*bindings.CheckPoWPoWProof, error) {
	blocks := make([]bindings.CheckPoWBlockPoW, 0, safetyBlocks+1)

	transferProof, err := cb.EncodeTransferProof(e.bridge.GetClient(), transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferEvent: %w", err)
	}

	for i := uint64(0); i <= safetyBlocks; i++ {
		targetBlockNum := big.NewInt(int64(transferEvent.Raw.BlockNumber + i))
		targetBlock, err := e.bridge.GetClient().BlockByNumber(context.Background(), targetBlockNum)
		if err != nil {
			return nil, fmt.Errorf("BlockByNumber: %w", err)
		}

		e.logger.Debug().Msgf("Encoding block %d... (%d/%d)", targetBlock.NumberU64(), i, safetyBlocks)
		encodedBlock, err := e.EncodeBlock(targetBlock.Header(), i == 0)
		if err != nil {
			return nil, fmt.Errorf("EncodeBlock: %w", err)
		}
		e.logger.Debug().Msgf("Encoded block %d", targetBlock.NumberU64())
		blocks = append(blocks, *encodedBlock)
	}

	return &bindings.CheckPoWPoWProof{
		Blocks:   blocks,
		Transfer: *transferProof,
	}, nil
}
