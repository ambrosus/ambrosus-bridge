package amb

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/receipts_proof"
	"github.com/ethereum/go-ethereum/common"
)

// todo эти структуры сгенерируются абигеном и будут в пакете contracts, но пока так

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
func (b *Bridge) getBlocksAndEvents() {
	// todo get new transfer event
	// start = get blockNum of prev event
	// stop = blockNum of current event + minSafetyBlocks
	// get vs change events from start to stop blocks
	// encode events
	// get list of block numbers that we need
	// get this blocks, encode, set block type
	// submit(blocks, transfer, vsChanges)
}

func (b *Bridge) encodeTransferEvent(event *contracts.TransferEvent) (*Transfer_Event, error) {
	receipts, err := b.GetReceipts(event.Raw.BlockHash)
	if err != nil {
		return nil, err
	}
	proof, err := receipts_proof.CalcProof(receipts, &event.Raw)
	if err != nil {
		return nil, err
	}

	return &Transfer_Event{
		Receipt_proof: proof,
		Event_id:      event.EventId,
		Transfers:     event.Queue,
	}, nil
}

func (b *Bridge) encodeVSChangeEvent(event *contracts.TransferEvent) (*ValidatorSet_Event, error) {
	// todo

	receipts, err := b.GetReceipts(event.Raw.BlockHash)
	if err != nil {
		// todo
	}
	proof, err := receipts_proof.CalcProof(receipts, &event.Raw)
	if err != nil {
		// todo
	}

	return &ValidatorSet_Event{
		Receipt_proof: proof,
	}, nil
}
