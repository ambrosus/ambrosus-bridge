package eth

import (
	"context"
	"math/big"
	"relay/config"
	"relay/contracts"
	"relay/networks"
	"relay/receipts_proof"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Bridge struct {
	Client     *ethclient.Client
	Contract   *contracts.Eth
	sideBridge networks.Bridge
	config     *config.Bridge
	submitFunc networks.SubmitPoWF
}

func New(c *config.Bridge) *Bridge {
	client, err := ethclient.Dial(c.Url)
	if err != nil {
		panic(err)
	}
	ethBridge, err := contracts.NewEth(c.ContractAddress, client)
	if err != nil {
		panic(err)
	}
	return &Bridge{
		Client:   client,
		Contract: ethBridge,
		config:   c,
	}
}

func (b *Bridge) SubmitBlockPoA(
	eventId *big.Int,
	blocks []contracts.CheckPoABlockPoA,
	events []contracts.CommonStructsTransfer,
	proof *contracts.ReceiptsProof,
) {
	auth, err := b.getAuth()
	if err != nil {
		// todo
	}

	_, err = b.Contract.CheckPoA(auth, blocks, events, *proof)
	if err != nil {
		// todo
	}

}

func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.Bridge, submit networks.SubmitPoWF) {
	b.sideBridge = sideBridge
	b.submitFunc = submit
	b.CheckOldEvents()
	b.Listen()
}

// не дописано
func (b *Bridge) CheckOldEvents() {
	for {
		// needId := sideBridge.GetLastEventId() + 1
		// // todo get event by id `needId
		//
		// if !event {
		//	return
		// }
		//
		// b.sendEvent()
	}
}

func (b *Bridge) Listen() {
	// Subscribe to events
	watchOpts := &bind.WatchOpts{
		Context: context.Background(),
	}
	eventChannel := make(chan *contracts.EthTransfer) // <-- тут я хз как сделать общий(common) тип для канала
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
			go b.sendEvent(&event.TransferEvent)
		}
	}
}

func (b *Bridge) sendEvent(event *contracts.TransferEvent) {
	// wait for safety blocks
	blockChannel := make(chan *types.Header)
	blockSub, err := b.Client.SubscribeNewHead(context.Background(), blockChannel)
	if err != nil {
		// todo
	}
	var i uint64
	for i <= b.config.SafetyBlocks {
		select {
		case err := <-blockSub.Err():
			// todo
			panic(err)

		case _ = <-blockChannel:
			i++
		}
	}
	// check if the event has been removed
	isEventRemoved, err := b.isEventRemoved(event.EventId)
	if isEventRemoved {
		// todo
	}

	// encode safety blocks
	safetyBlocks := b.getSafetyBlocks(event.Raw.BlockNumber)
	encodedBlocks := b.encodeSafetyBlocks(safetyBlocks)

	// calculate receipt proof
	receipts, err := b.GetReceipts(event.Raw.BlockHash)
	if err != nil {
		// todo
	}
	proof_, err := receipts_proof.CalcProof(receipts, &event.Raw)
	if err != nil {
		// todo
	}
	proof := contracts.ReceiptsProof(proof_)

	b.submitFunc(event.EventId, encodedBlocks, event.Queue, &proof)
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

func (b Bridge) getSafetyBlocks(offset uint64) types.Blocks {
	blocks := make(types.Blocks, b.config.SafetyBlocks)
	for i := uint64(0); i < b.config.SafetyBlocks; i++ {
		block, err := b.Client.BlockByNumber(context.Background(), big.NewInt(int64(offset+i)))
		if err != nil {
			panic(err)
		}
		blocks = append(blocks, block)
	}

	return blocks
}

func (b Bridge) encodeSafetyBlocks(safetyBlocks types.Blocks) []contracts.CheckPoWBlockPoW {
	encodedBlocks := make([]contracts.CheckPoWBlockPoW, b.config.SafetyBlocks)
	encodedBlocks = append(encodedBlocks, *EncodeBlock(safetyBlocks[0].Header(), true))

	for i := uint64(1); i < b.config.SafetyBlocks; i++ {
		encodedBlocks = append(encodedBlocks, *EncodeBlock(safetyBlocks[i].Header(), false))
	}

	return encodedBlocks
}

func (b Bridge) getAuth() (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(b.config.PrivateKey, b.config.ChainID)
	if err != nil {
		return nil, err
	}

	nonce, err := b.Client.PendingNonceAt(auth.Context, auth.From)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	return auth, nil
}

func (b *Bridge) isEventRemoved(eventId *big.Int) (bool, error) {
	event, err := b.GetEventById(eventId)
	if err != nil {
		return false, err
	}
	return event.Raw.Removed, nil
}

func (b *Bridge) GetEventById(eventId *big.Int) (*contracts.TransferEvent, error) {
	opts := &bind.FilterOpts{
		Context: context.Background(),
	}
	logs, err := b.Contract.FilterTransfer(opts, []*big.Int{eventId})
	if err != nil {
		return nil, err
	}

	if logs.Next() {
		return &logs.Event.TransferEvent, nil
	}
	return nil, nil
}
