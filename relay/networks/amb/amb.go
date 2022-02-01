package amb

import (
	"context"
	"math/big"
	"relay/config"
	"relay/contracts"
	"relay/helpers"
	"relay/networks"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Bridge struct {
	client     *ethclient.Client
	contract   *contracts.Amb
	sideBridge *networks.Bridge
	config     *config.Bridge
	submitFunc networks.SubmitPoAF
}

func New(c *config.Bridge) *Bridge {
	client, err := ethclient.Dial(c.Url)
	if err != nil {
		panic(err)
	}
	ambBridge, err := contracts.NewAmb(c.ContractAddress, client)
	if err != nil {
		panic(err)
	}
	return &Bridge{
		client:   client,
		contract: ambBridge,
		config:   c,
	}
}

func (b Bridge) getSafetyBlocks(offset uint64) []*Header {
	blocks := make([]*Header, b.config.SafetyBlocks)
	var i uint64
	for i = 0; i < b.config.SafetyBlocks; i++ {
		block, err := HeaderByNumber(big.NewInt(int64(i + offset)))
		if err != nil {
			panic(err)
		}
		blocks = append(blocks, block)
	}

	return blocks
}

func (b Bridge) encodeSafetyBlocks(safetyBlocks []*Header) []*contracts.CheckPoABlockPoA {
	encodedBlocks := make([]*contracts.CheckPoABlockPoA, b.config.SafetyBlocks)
	encodedBlocks = append(encodedBlocks, EncodeBlock(safetyBlocks[0], true))

	var i uint64
	for i = 1; i < b.config.SafetyBlocks; i++ {
		encodedBlocks = append(encodedBlocks, EncodeBlock(safetyBlocks[i], false))
	}

	return encodedBlocks
}

func (b *Bridge) SubmitBlockPoW(
	eventId *big.Int,
	blocks []*contracts.CheckPoWBlockPoW,
	events *[]contracts.CommonStructsTransfer,
	proof *contracts.ReceiptsProof,
) {
	// todo estimate gas
	// todo send
	// todo wait status ok
}

func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.contract.InputEventId(nil)
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge *networks.Bridge, submit networks.SubmitPoAF) {
	// todo save args to struct?
	b.sideBridge = sideBridge
	b.submitFunc = submit
	b.CheckOldEvents()
	b.Listen()
}

// не дописано
func (b *Bridge) CheckOldEvents() {
	for {
		needId := sideBridge.GetLastEventId() + 1
		// todo get event by id `needId

		if !event {
			return
		}

		b.sendEvent()
	}
}

func (b *Bridge) Listen() {
	// Subscribe to events
	watchOpts := &bind.WatchOpts{
		Context: context.Background(),
	}
	eventChannel := make(chan *contracts.AmbTransferEvent) // <-- тут я хз как сделать общий(common) тип для канала
	eventSub, err := b.contract.WatchTransferEvent(watchOpts, eventChannel, nil)
	if err != nil {
		panic(err)
	}

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			panic(err)

		case event := <-eventChannel:
			go b.sendEvent(&event.TransferEvent)
		}
	}
}

func (b *Bridge) sendEvent(event *contracts.TransferEvent) {
	// todo wait for safety blocks
	sleepTime := b.config.SafetyBlocks * b.config.BlockTime
	time.Sleep(time.Duration(sleepTime) * time.Second)

	// todo encode blocks
	safetyBlocks := b.getSafetyBlocks(event.Raw.BlockNumber)
	encodedBlocks := b.encodeSafetyBlocks(safetyBlocks)

	receipts := types.Receipts(helpers.FindReceipts(b.client, event.Raw.BlockHash))
	proof := contracts.ReceiptsProof(helpers.CalcProof(&receipts, event.Raw.Data))

	b.submitFunc(event.EventId, encodedBlocks, &event.Queue, &proof)
}
