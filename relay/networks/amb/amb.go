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
	"golang.org/x/sync/errgroup"
)

type Bridge struct {
	Client      *ethclient.Client
	Contract    *contracts.Amb
	ContractRaw *contracts.AmbRaw
	VSContract  *contracts.Vs
	HttpUrl     string // TODO: delete this field
	sideBridge  networks.BridgeReceiveAura
	auth        *bind.TransactOpts
}

func (b *Bridge) SubmitEpochData() error {
	//TODO implement me
	panic("implement me")
}

// Creating a new ambrosus bridge.
func New(cfg *config.AMBConfig) (*Bridge, error) {
	// Creating a new ethereum client.
	client, err := ethclient.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}

	// Creating a new ambrosus bridge contract instance.
	contract, err := contracts.NewAmb(common.HexToAddress(cfg.ContractAddr), client)
	if err != nil {
		return nil, err
	}

	// Creating a new ambrosus VS contract instance.
	vsContract, err := contracts.NewVs(common.HexToAddress(cfg.VSContractAddr), client)
	if err != nil {
		return nil, err
	}

	var auth *bind.TransactOpts

	if cfg.PrivateKey != nil {
		chainId, err := client.ChainID(context.Background())
		if err != nil {
			return nil, err
		}

		auth, err = bind.NewKeyedTransactorWithChainID(cfg.PrivateKey, chainId)
		if err != nil {
			return nil, err
		}
	}

	return &Bridge{
		Client:      client,
		Contract:    contract,
		ContractRaw: &contracts.AmbRaw{Contract: contract},
		VSContract:  vsContract,
		HttpUrl:     "https://network.ambrosus.io",
		auth:        auth,
	}, nil
}

func (b *Bridge) SubmitTransferPoW(proof *contracts.CheckPoWPoWProof) error {
	tx, txErr := b.Contract.SubmitTransfer(b.auth, *proof)

	// todo find way to make this part common for different contract methods
	if txErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message
		// so, make low-level call method to get the full error message

		err := b.ContractRaw.Call(&bind.CallOpts{
			From: b.auth.From,
		}, nil, "submitTransfer", *proof)

		if err != nil {
			return fmt.Errorf("%s", parseError(err))
		}
		return fmt.Errorf("%s", parseError(txErr))
	}

	return b.waitForTxMined(tx)
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

	return nil, networks.ErrEventNotFound
}

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.BridgeReceiveAura) {
	b.sideBridge = sideBridge

	for {
		if err := b.listen(); err != nil {
			log.Error().Err(err).Msg("listen ambrosus error")
		}
	}
}

func (b *Bridge) checkOldEvents() error {
	lastEventId, err := b.sideBridge.GetLastEventId()
	if err != nil {
		return err
	}

	i := big.NewInt(1)
	for {
		nextEventId := big.NewInt(0).Add(lastEventId, i)
		nextEvent, err := b.GetEventById(nextEventId)
		if err != nil {
			if errors.Is(err, networks.ErrEventNotFound) {
				// no more old events
				return nil
			}
			return err
		}

		err = b.sendEvent(nextEvent)
		if err != nil {
			return err
		}

		i = big.NewInt(0).Add(i, big.NewInt(1))
	}
}

func (b *Bridge) listen() error {
	err := b.checkOldEvents()
	if err != nil {
		return err
	}

	// Subscribe to events
	watchOpts := &bind.WatchOpts{Context: context.Background()}
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
	// Wait for safety blocks.
	safetyBlocks, err := b.getSafetyBlocksNum()
	if err != nil {
		return err
	}

	if err := b.waitForBlock(event.Raw.BlockNumber + safetyBlocks); err != nil {
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

	return b.sideBridge.SubmitTransferAura(ambTransfer)
}

func (b *Bridge) GetReceipts(blockHash common.Hash) ([]*types.Receipt, error) {
	txsCount, err := b.Client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		return nil, err
	}

	receipts := make([]*types.Receipt, txsCount)

	errGroup := new(errgroup.Group)
	for i := uint(0); i < txsCount; i++ {
		i := i // https://golang.org/doc/faq#closures_and_goroutines ¯\_(ツ)_/¯
		errGroup.Go(func() error {
			tx, err := b.Client.TransactionInBlock(context.Background(), blockHash, i)
			if err != nil {
				return err
			}
			receipt, err := b.Client.TransactionReceipt(context.Background(), tx.Hash())
			if err != nil {
				return err
			}

			receipts[i] = receipt
			return nil
		})
	}

	err = errGroup.Wait()
	if err != nil {
		return nil, err
	}
	return receipts, nil
}

func (b *Bridge) isEventRemoved(event *contracts.TransferEvent) error {
	block, err := b.HeaderByNumber(big.NewInt(int64(event.Raw.BlockNumber)))
	if err != nil {
		return err
	}

	if block.Hash(true) != event.Raw.BlockHash {
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

func (b *Bridge) getSafetyBlocksNum() (uint64, error) {
	safetyBlocks, err := b.Contract.MinSafetyBlocks(nil)
	if err != nil {
		return 0, err
	}
	return safetyBlocks.Uint64(), nil
}
