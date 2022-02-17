package amb

import (
	"context"
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
	// auth, err := b.getAuth()
	// if err != nil {
	//	// todo
	// }

	// tx, err := b.Contract.CheckPoW(auth, powProof)
	// if err != nil {
	//	// todo
	// }
	// _ = tx
}

// Getting last contract event id.
func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
}

// Getting contract event by id.
func (b *Bridge) GetEventById(eventId *big.Int) (*contracts.TransferEvent, error) {
	opts := &bind.FilterOpts{Context: context.Background()}

	logs, err := b.Contract.FilterTransfer(opts, []*big.Int{eventId})
	if err != nil {
		return nil, err
	}

	if logs.Next() {
		return &logs.Event.TransferEvent, nil
	}
	// todo err not found?
	return nil, nil
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.Bridge) {
	b.sideBridge = sideBridge

	for {
		if err := b.listen(); err != nil {
			log.Error().Err(err).Msg("listen ambrosus error")
		}
	}
}

func (b *Bridge) listen() error {
	lastEventId, err := b.sideBridge.GetLastEventId()
	if err != nil {
		return err
	}

	lastEvent, err := b.GetEventById(lastEventId)
	if err != nil {
		return err
	}

	startBlock := lastEvent.Raw.BlockNumber + 1

	// Subscribe to events
	watchOpts := &bind.WatchOpts{Start: &startBlock, Context: context.Background()}
	eventChannel := make(chan *contracts.AmbTransfer) // <-- тут я хз как сделать общий(common) тип для канала
	eventSub, err := b.Contract.WatchTransfer(watchOpts, eventChannel, nil)
	if err != nil {
		return err
	}

	defer eventSub.Unsubscribe()

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return err
		case event := <-eventChannel:
			if err := b.sendEvent(&event.TransferEvent); err != nil {
				return err
			}
		}
	}
}

func (b *Bridge) sendEvent(event *contracts.TransferEvent) error {
	// todo update minSafetyBlocks value from contract

	// Wait for safety blocks.
	if err := b.waitForBlock(event.Raw.BlockNumber + b.config.SafetyBlocks); err != nil {
		return err
	}

	// Check if the event has been removed.
	if err := b.isEventRemoved(event); err != nil {
		return err
	}

	ambTransfer, err := b.getBlocksAndEvents(event)
	if err != nil {
		return err
	}

	// todo
	_ = ambTransfer
	// b.submitFunc(blocks, transfer, vsChanges)

	return nil
}

func (b *Bridge) GetReceipts(blockHash common.Hash) ([]*types.Receipt, error) {
	// todo we can use goroutines here
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

func (b *Bridge) isEventRemoved(event *contracts.TransferEvent) error {
	block, err := b.Client.BlockByNumber(context.Background(), big.NewInt(int64(event.Raw.BlockNumber)))
	if err != nil {
		return err
	}

	if block.Hash() != event.Raw.BlockHash {
		return fmt.Errorf("block hash != event's block hash")
	}
	return nil
}

func (b *Bridge) waitForBlock(targetBlockNum uint64) error {
	// todo maybe timeout (context)
	blockChannel := make(chan *types.Header)
	blockSub, err := b.Client.SubscribeNewHead(context.Background(), blockChannel)
	if err != nil {
		return err
	}

	currentBlockNum, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	for currentBlockNum < targetBlockNum {
		select {
		case err := <-blockSub.Err():
			return err

		case block := <-blockChannel:
			currentBlockNum = block.Number.Uint64()
		}
	}

	return nil
}
