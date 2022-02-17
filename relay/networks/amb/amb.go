package amb

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ambrosus/ambrosus-bridge/relay/config"
	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/networks"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

// maybe move to helpers pkg
type JsonError interface {
	Error() string
	ErrorCode() int
	ErrorData() interface{}
}

type Bridge struct {
	Client      *ethclient.Client
	Contract    *contracts.Amb
	ContractRaw *contracts.AmbRaw
	VSContract  *contracts.Vs
	sideBridge  networks.Bridge
	config      *config.Bridge
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

	return &Bridge{
		Client:      client,
		Contract:    contract,
		ContractRaw: &contracts.AmbRaw{Contract: contract},
		VSContract:  vsContract,
		config:      cfg,
	}, nil
}

func (b *Bridge) SubmitTransfer(proof contracts.TransferProof) error {
	var castProof *contracts.CheckPoWPoWProof
	switch proof.(type) {
	case *contracts.CheckPoWPoWProof:
		// todo
		castProof = proof.(*contracts.CheckPoWPoWProof)
	default:
		// todo error
		return fmt.Errorf("")
	}

	auth, err := b.getAuth()
	if err != nil {
		return err
	}

	tx, txErr := b.Contract.SubmitTransfer(auth, *castProof)

	if txErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message
		// so, make low-level call method to get the full error message
		err = b.ContractRaw.Call(&bind.CallOpts{
			From: auth.From,
		}, nil, "submitTransfer", *castProof)

		if err != nil {
			errStr := getJsonErrData(err)

			if strings.HasPrefix(errStr, "Reverted") {
				errBytes, err := hex.DecodeString(errStr[11:])
				if err != nil {
					return err
				}
				errStr = string(errBytes)
			}

			return fmt.Errorf("%s", errStr)
		} else {
			// врятли такой кейс будет.
			// но если и будет, то шото странно: при вызове контракта норм способом ошибка есть
			// а при вызове через ненорм способ - нету.
			// ну тогда вернём ту ошибку, которая при норм способе
			errStr := getJsonErrData(txErr)
			return fmt.Errorf("%s", errStr)
		}
	}

	receipt, err := bind.WaitMined(context.Background(), b.Client, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful {
		// we've got here probably due to low gas limit,
		// and revert() that hasn't been caught at eth_estimateGas
		err := getFailureReason(b.Client, auth.From, tx)
		if err != nil {
			errStr := getJsonErrData(err)
			return fmt.Errorf("%s", errStr)
		}
	}

	return nil
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

// todo delete
func (b *Bridge) GetValidatorSet() ([]common.Address, error) {
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

// maybe move the functions below to helpers pkg
func getFailureReason(client *ethclient.Client, from common.Address, tx *types.Transaction) error {
	_, err := client.CallContract(context.Background(), ethereum.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}, nil)

	return err
}

func getJsonErrData(err error) string {
	var jsonErr = err.(JsonError)
	return jsonErr.ErrorData().(string)
}
