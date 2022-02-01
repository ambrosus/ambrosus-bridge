package amb

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"relay/config"
	"relay/contracts"
	"relay/networks"
	"relay/receipts_proof"
)

type Bridge struct {
	Client     *ethclient.Client
	Contract   *contracts.Amb
	sideBridge networks.Bridge
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
		Client:   client,
		Contract: ambBridge,
		config:   c,
	}
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
	return b.Contract.InputEventId(nil)
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.Bridge, submit networks.SubmitPoAF) {
	// todo save args to struct?
	b.sideBridge = sideBridge
	b.submitFunc = submit
	b.CheckOldEvents()
	b.Listen()
}

// не дописано
func (b *Bridge) CheckOldEvents() {
	for {
		//needId := sideBridge.GetLastEventId() + 1
		//// todo get event by id `needId
		//
		//if !event {
		//	return
		//}
		//
		//b.sendEvent()
	}
}

func (b *Bridge) Listen() {
	// Subscribe to events
	watchOpts := &bind.WatchOpts{
		Context: context.Background(),
	}
	eventChannel := make(chan *contracts.AmbTransferEvent) // <-- тут я хз как сделать общий(common) тип для канала
	eventSub, err := b.Contract.WatchTransferEvent(watchOpts, eventChannel, nil)
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
	//sleepTime := b.config.SafetyBlocks * b.config.BlockTime
	//time.Sleep(time.Duration(sleepTime) * time.Second)

	// todo encode blocks
	safetyBlocks := b.getSafetyBlocks(event.Raw.BlockNumber)
	encodedBlocks := b.encodeSafetyBlocks(safetyBlocks)

	receipts, err := b.GetReceipts(event.Raw.BlockHash)
	if err != nil {
		// todo
	}
	proof_, err := receipts_proof.CalcProof(receipts, &event.Raw)
	if err != nil {
		// todo
	}
	proof := contracts.ReceiptsProof(proof_)
	b.submitFunc(event.EventId, encodedBlocks, &event.Queue, &proof)
}

func (b *Bridge) GetReceipts(blockHash common.Hash) ([]*types.Receipt, error) {
	txsCount, err := b.Client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		return nil, err
	}

	receipts := make([]*types.Receipt, 0, txsCount)

	for i := uint(0); i < txsCount; i++ {
		tx, err := b.Client.TransactionInBlock(context.Background(), blockHash, i)
		if err != nil {
			return nil, err
		}
		receipt, err := b.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return nil, err
		}
		receipts = append(receipts, receipt)
	}
	return receipts, nil
}

func (b Bridge) getSafetyBlocks(offset uint64) []*Header {
	blocks := make([]*Header, b.config.SafetyBlocks)
	for i := uint64(0); i < b.config.SafetyBlocks; i++ {
		block, err := b.HeaderByNumber(big.NewInt(int64(offset + i)))
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

	for i := uint64(1); i < b.config.SafetyBlocks; i++ {
		encodedBlocks = append(encodedBlocks, EncodeBlock(safetyBlocks[i], false))
	}

	return encodedBlocks
}
