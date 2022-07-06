package pow_proof

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/rs/zerolog"
)

type PoWEncoder struct {
	bridge      networks.Bridge
	powReceiver networks.BridgeReceiveEthash

	ethash *ethash.Ethash
	logger *zerolog.Logger
}

func NewPoWEncoder(bridge networks.Bridge, sideBridge networks.BridgeReceiveEthash, ethash *ethash.Ethash) *PoWEncoder {
	return &PoWEncoder{
		bridge:      bridge,
		powReceiver: sideBridge,
		ethash:      ethash,
		logger:      bridge.GetLogger(), // todo maybe sublogger?
	}
}

func (b *PoWEncoder) EncodePoWProof(transferEvent *bindings.BridgeTransfer, safetyBlocks uint64) (*bindings.CheckPoWPoWProof, error) {
	blocks := make([]bindings.CheckPoWBlockPoW, 0, safetyBlocks+1)

	transfer, err := b.encodeTransferEvent(transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferEvent: %w", err)
	}

	for i := uint64(0); i <= safetyBlocks; i++ {
		targetBlockNum := big.NewInt(int64(transferEvent.Raw.BlockNumber + i))
		targetBlock, err := b.bridge.GetClient().BlockByNumber(context.Background(), targetBlockNum)
		if err != nil {
			return nil, fmt.Errorf("BlockByNumber: %w", err)
		}

		b.logger.Debug().Msgf("Encoding block %d... (%d/%d)", targetBlock.NumberU64(), i, safetyBlocks)
		encodedBlock, err := b.EncodeBlock(targetBlock.Header(), i == 0)
		if err != nil {
			return nil, fmt.Errorf("EncodeBlock: %w", err)
		}
		b.logger.Debug().Msgf("Encoded block %d", targetBlock.NumberU64())
		blocks = append(blocks, *encodedBlock)
	}

	return &bindings.CheckPoWPoWProof{
		Blocks:   blocks,
		Transfer: *transfer,
	}, nil
}

func (b *PoWEncoder) encodeTransferEvent(event *bindings.BridgeTransfer) (*bindings.CommonStructsTransferProof, error) {
	proof, err := cb.GetProof(b.bridge.GetClient(), event)
	if err != nil {
		return nil, fmt.Errorf("GetProof: %w", err)
	}

	return &bindings.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}
