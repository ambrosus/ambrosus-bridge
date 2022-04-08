package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/receipts_proof"
)

// todo name
func (b *Bridge) getBlocksAndEvents(transferEvent *contracts.TransferEvent) (*contracts.CheckPoWPoWProof, error) {
	safetyBlocks, err := b.sideBridge.GetMinSafetyBlocksNum()
	if err != nil {
		return nil, fmt.Errorf("GetMinSafetyBlocksNum in %v: %w", b.sideBridge.Name(), err)
	}
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

		encodedBlock, err := EncodeBlock(targetBlock.Header(), i == 0)
		if err != nil {
			return nil, fmt.Errorf("EncodeBlock: %w", err)
		}
		blocks = append(blocks, *encodedBlock)
	}

	return &contracts.CheckPoWPoWProof{
		Blocks:   blocks,
		Transfer: *transfer,
	}, nil
}

func (b *Bridge) encodeTransferEvent(event *contracts.TransferEvent) (*contracts.CommonStructsTransferProof, error) {
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
	receipts, err := ethereum.GetReceipts(b.Client, event.Log().BlockHash)
	if err != nil {
		return nil, err
	}
	return receipts_proof.CalcProofEvent(receipts, event)
}
