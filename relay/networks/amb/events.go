package amb

import (
	"context"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/receipts_proof"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// todo эти структуры сгенерируются абигеном и будут в пакете contracts, но пока так

type ValidatorSet_Event struct {
	Receipt_proof contracts.ReceiptsProof
	Delta_address common.Address
	Delta_index   uint64    // 12байт шоб спаковать с 20байт адресом
	Raw           types.Log // Blockchain specific contextual infos
}

type Transfer_Event struct {
	Receipt_proof contracts.ReceiptsProof
	Event_id      *big.Int
	Transfers     []contracts.CommonStructsTransfer
	Raw           types.Log // Blockchain specific contextual infos
}

func (b *Bridge) listenTransferEvent() {
	// Subscribe to events
	watchOpts := &bind.WatchOpts{Context: context.Background()}
	eventChannel := make(chan *contracts.AmbTransfer) // <-- тут я хз как сделать общий(common) тип для канала
	eventSub, err := b.Contract.WatchTransfer(watchOpts, eventChannel, nil)
	if err != nil {
		panic(err)
	}

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			panic(err)

		case event := <-eventChannel:
			b.getBlocksAndEvents(&event.TransferEvent)
		}
	}
}

// todo name
func (b *Bridge) getBlocksAndEvents(event *contracts.TransferEvent) {
	prevEventId := big.NewInt(0).Sub(event.EventId, big.NewInt(1))
	prevEvent, err := b.GetEventById(prevEventId)
	if err != nil {
		// todo
	}

	// start = get blockNum of prev event
	startSearch := prevEvent.Raw.BlockNumber
	// stop = blockNum of current event + minSafetyBlocks
	stopSearch := event.Raw.BlockNumber + b.config.SafetyBlocks

	// get vs change events from start to stop blocks
	vsChangeEvents, err := b.getVSChangeEvents(startSearch, stopSearch)
	if err != nil {
		// todo
	}

	// encode events

	// get list of block numbers that we need
	blocksNum := make(map[uint64]*Header)
	for _, vsChangeEvent := range vsChangeEvents {
		blockNum := vsChangeEvent.Raw.BlockNumber
		block, err := b.HeaderByNumber(big.NewInt(int64(blockNum)))
		if err != nil {
			// todo
		}
		blocksNum[blockNum] = block
	}

	// get this blocks, encode, set block type

	// submit(blocks, transfer, vsChanges)
}

func (b *Bridge) getVSChangeEvents(start, stop uint64) ([]*ValidatorSet_Event, error) {
	opts := &bind.FilterOpts{
		Start:   start,
		End:     &stop,
		Context: context.Background(),
	}
	logs, err := b.Contract.FilterValidatorSet_Event(opts)
	if err != nil {
		return nil, err
	}

	var res []*ValidatorSet_Event
	for logs.Next() {
		res = append(res, logs.Event)
	}

	return res, nil
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
