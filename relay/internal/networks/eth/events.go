package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof"
)

func (b *Bridge) getBlocksAndEvents(transferEvent *contracts.BridgeTransfer, safetyBlocks uint64) (*contracts.CheckPoWPoWProof, error) {
	blocks := make([]contracts.CheckPoWBlockPoW, 0, safetyBlocks+1)

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

	return &contracts.CheckPoWPoWProof{
		Blocks:   blocks,
		Transfer: *transfer,
	}, nil
}

func (b *Bridge) encodeTransferEvent(event *contracts.BridgeTransfer) (*contracts.CommonStructsTransferProof, error) {
	proof, err := b.getProof(event)
	if err != nil {
		return nil, fmt.Errorf("getProof: %w", err)
	}

	return &contracts.CommonStructsTransferProof{
		ReceiptProof: proof,
		EventId:      event.EventId,
		Transfers:    event.Queue,
	}, nil
}

func (b *Bridge) getProof(event receipts_proof.ProofEvent) ([][]byte, error) {
	receipts, err := b.GetReceipts(event.Log().BlockHash)
	if err != nil {
		return nil, fmt.Errorf("GetReceipts: %w", err)
	}
	return receipts_proof.CalcProofEvent(receipts, event)
}
