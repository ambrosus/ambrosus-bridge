package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
)

func (b *Bridge) encodePoWProof(transferEvent *bindings.BridgeTransfer, safetyBlocks uint64) (*bindings.CheckPoWPoWProof, error) {
	blocks := make([]bindings.CheckPoWBlockPoW, 0, safetyBlocks+1)

	transfer, err := b.encodeTransferEvent(transferEvent)
	if err != nil {
		return nil, fmt.Errorf("encodeTransferEvent: %w", err)
	}

	for i := uint64(0); i <= safetyBlocks; i++ {
		targetBlockNum := big.NewInt(int64(transferEvent.Raw.BlockNumber + i))
		targetBlock, err := b.Client.BlockByNumber(context.Background(), targetBlockNum)
		if err != nil {
			return nil, fmt.Errorf("BlockByNumber: %w", err)
		}

		b.Logger.Debug().Msgf("Encoding block %d... (%d/%d)", targetBlock.NumberU64(), i, safetyBlocks)
		encodedBlock, err := b.EncodeBlock(targetBlock.Header(), i == 0)
		if err != nil {
			return nil, fmt.Errorf("EncodeBlock: %w", err)
		}
		b.Logger.Debug().Msgf("Encoded block %d", targetBlock.NumberU64())
		blocks = append(blocks, *encodedBlock)
	}

	return &bindings.CheckPoWPoWProof{
		Blocks:   blocks,
		Transfer: *transfer,
	}, nil
}

func (b *Bridge) encodeTransferEvent(event *bindings.BridgeTransfer) (*bindings.CommonStructsTransferProof, error) {
	proof, err := b.GetProof(event)
	if err != nil {
		return nil, fmt.Errorf("GetProof: %w", err)
	}

	return &bindings.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}
