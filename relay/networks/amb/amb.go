package amb

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/networks"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type Bridge struct {
	Client     *ethclient.Client
	Contract   *contracts.Amb
	VSContract *contracts.Vs
	sideBridge networks.Bridge
	config     *config.Bridge
}

// Creating a new ambrosus bridge.
func New(cfg *config.Bridge) (*Bridge, error) {
	// Creating a new ethereum client.
	client, err := ethclient.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}

	// Creating a new ambrosus bridge contract instance.
	contract, err := contracts.NewAmb(cfg.ContractAddress, client)
	if err != nil {
		return nil, err
	}

	// Creating a new ambrosus VS contract instance.
	vsContract, err := contracts.NewVs(cfg.VSContractAddress, client)
	if err != nil {
		return nil, err
	}

	return &Bridge{Client: client, Contract: contract, VSContract: vsContract, config: cfg}, nil
}

func (b *Bridge) SubmitTransfer(proof contracts.TransferProof) error {
	switch proof.(type) {
	case *contracts.CheckPoWPoWProof:
		// todo
	default:
		// todo error

	}
	return nil
	//auth, err := b.getAuth()
	//if err != nil {
	//	// todo
	//}

	//tx, err := b.Contract.CheckPoW(auth, powProof)
	//if err != nil {
	//	// todo
	//}
	//_ = tx
}

// Getting last contract event id.
func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// Getting contract event by id.
func (b *Bridge) GetEventById(eventId *big.Int) (*contracts.TransferEvent, error) {
	// Filter contract event transfers options.
	opts := &bind.FilterOpts{Context: context.Background()}

	// Filter contact event transfers.
	logs, err := b.Contract.FilterTransfer(opts, []*big.Int{eventId})
	if err != nil {
		return nil, err
	}

	if logs.Next() {
		return &logs.Event.TransferEvent, nil
	}

	return nil, fmt.Errorf("error event: '%d' not found", eventId.Uint64())
}

// todo code below may be common for all networks?

// Running ambrosus bridge.
func (b *Bridge) Run(sideBridge networks.Bridge) {
	b.sideBridge = sideBridge

	for {
		// Listening contract transfer events.
		if err := b.listen(); err != nil {
			log.Error().Err(err).Msg("listen ambrosus error")
		}
	}
}

// Listening contract transfer events.
func (b *Bridge) listen() error {
	// Getting last contract event id.
	lastEventId, err := b.sideBridge.GetLastEventId()
	if err != nil {
		return err
	}

	// Getting contract event by event id.
	lastEvent, err := b.GetEventById(lastEventId)
	if err != nil {
		return err
	}

	startBlock := lastEvent.Raw.BlockNumber + 1

	// Watching transfer options.
	opts := &bind.WatchOpts{Start: &startBlock, Context: context.Background()}

	// Create a new channel for transfers.
	events := make(chan *contracts.AmbTransfer)

	// Subscribe to contract transfer event.
	sub, err := b.Contract.WatchTransfer(opts, events, nil)
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			return err
		case event := <-events:
			// Process ambrosus transfer event.
			if err := b.processEvent(&event.TransferEvent); err != nil {
				return err
			}
		}
	}
}

// Process ambrosus transfer event.
func (b *Bridge) processEvent(event *contracts.TransferEvent) error {
	log.Info().Str("eventID", event.EventId.String()).Msg("proccessing ambrosus event...")

	// TODO: update minSafetyBlocks value from contract

	// Wait for safety blocks.
	if err := b.waitForBlock(event.Raw.BlockNumber + b.config.SafetyBlocks); err != nil {
		return err
	}

	// Check if the event has been removed.
	if err := b.isEventRemoved; err != nil {
		return err(event)
	}

	// Getting blocks and events by transfer event.
	transfer, err := b.getBlocksAndEvents(event)
	if err != nil {
		return err
	}

	// TODO
	_ = transfer
	// b.submitFunc(blocks, transfer, vsChanges)

	return nil
}

// Getting receipts by block hash.
func (b *Bridge) GetReceipts(blockHash common.Hash) ([]*types.Receipt, error) {

	// TODO: we can use goroutines here

	// Getting transaction count.
	txsCount, err := b.Client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		return nil, err
	}

	// Creating a new channel for receipts.
	receipts := make([]*types.Receipt, 0, txsCount)

	// Transaction processing by index.
	for i := uint(0); i < txsCount; i++ {
		// Getting transaction from block by index.
		tx, err := b.Client.TransactionInBlock(context.Background(), blockHash, i)
		if err != nil {
			return nil, err
		}

		// Getting receipts by transactions hash.
		receipt, err := b.Client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			return nil, err
		}

		receipts[i] = receipt
	}

	return receipts, nil
}

// TODO func
func (b *Bridge) getAuth() (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(b.config.PrivateKey, b.config.ChainID)
	if err != nil {
		return nil, err
	}

	// todo check if nonce can set automatically. if so, remove this function
	nonce, err := b.Client.PendingNonceAt(auth.Context, auth.From)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	return auth, nil
}

// Check is event removed in contract.
func (b *Bridge) isEventRemoved(event *contracts.TransferEvent) error {
	// Getting block by number.
	block, err := b.Client.BlockByNumber(context.Background(), big.NewInt(int64(event.Raw.BlockNumber)))
	if err != nil {
		return err
	}

	// Check is block hash != event block hash.
	if block.Hash() != event.Raw.BlockHash {
		return errors.New("block hash != event's block hash")
	}

	return nil
}

// Waiting a new block.
func (b *Bridge) waitForBlock(targetBlockNum uint64) error {

	// TODO maybe timeout (context)

	// Creating a new channel for blocks.
	blocks := make(chan *types.Header)

	// Subscribe to new head.
	sub, err := b.Client.SubscribeNewHead(context.Background(), blocks)
	if err != nil {
		return err
	}

	// Getting current block number.
	currentBlockNum, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	// Waiting a new bloks.
	for currentBlockNum < targetBlockNum {
		select {
		case err := <-sub.Err():
			return err
		case block := <-blocks:
			currentBlockNum = block.Number.Uint64()
		}
	}

	return nil
}
