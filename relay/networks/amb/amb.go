package amb

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"relay/config"
	"relay/contracts"
	"relay/networks"
	"relay/receipts_proof"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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
	sideBridge  networks.Bridge
	config      *config.Bridge
	submitFunc  networks.SubmitPoAF
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
		Client:      client,
		Contract:    ambBridge,
		ContractRaw: &contracts.AmbRaw{Contract: ambBridge},
		config:      c,
	}
}

func (b *Bridge) SubmitBlockPoW(
	eventId *big.Int,
	blocks []contracts.CheckPoWBlockPoW,
	events []contracts.CommonStructsTransfer,
	proof *contracts.ReceiptsProof,
) error {
	auth, err := b.getAuth()
	if err != nil {
		return err
	}

	tx, txErr := b.Contract.SubmitTransfer(auth, eventId, blocks, events, *proof, big.NewInt(int64(len(blocks))))

	if txErr != nil {
		// we've got here probably due to error at eth_estimateGas (e.g. revert(), require())
		// openethereum doesn't give us a full error message
		// so, make low-level call method to get the full error message
		err = b.ContractRaw.Call(&bind.CallOpts{
			From: auth.From,
		}, nil, "submitTransfer", eventId, blocks, events, *proof, big.NewInt(int64(len(blocks))))

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

func (b *Bridge) GetLastEventId() (*big.Int, error) {
	return b.Contract.InputEventId(nil)
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

// todo code below may be common for all networks?

func (b *Bridge) Run(sideBridge networks.Bridge, submit networks.SubmitPoAF) {
	// todo save args to struct?
	b.sideBridge = sideBridge
	b.submitFunc = submit
	b.CheckOldEvents()
	b.Listen()
}

func (b *Bridge) CheckOldEvents() {
	lastEventId, err := b.sideBridge.GetLastEventId()
	if err != nil {
		// todo
		panic(err)
	}

	i := big.NewInt(1)
	for {
		nextEventId := big.NewInt(0).Add(lastEventId, i)
		nextEvent, err := b.GetEventById(nextEventId)
		if err != nil {
			// todo
			panic(err)
		}

		if nextEvent == nil {
			return
		}

		go b.sendEvent(nextEvent)
		i = big.NewInt(0).Add(i, big.NewInt(1))
	}
}

func (b *Bridge) Listen() {
	// Subscribe to events
	watchOpts := &bind.WatchOpts{
		Context: context.Background(),
	}
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
			go b.sendEvent(&event.TransferEvent)
		}
	}
}

func (b *Bridge) sendEvent(event *contracts.TransferEvent) {
	// wait for safety blocks
	if err := b.waitForSafetyBlocks(event); err != nil {
		// todo
	}

	// check if the event has been removed
	if err := b.isEventRemoved; err != nil {
		// todo
	}

	blocks := b.encodeBlocks(event.Raw.BlockNumber)

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

	err = b.submitFunc(event.EventId, blocks, event.Queue, &proof)
	if err != nil {
		// todo
	}
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

func (b *Bridge) encodeBlocks(offset uint64) []contracts.CheckPoABlockPoA {
	encodedBlocks := make([]contracts.CheckPoABlockPoA, b.config.SafetyBlocks)

	for i := uint64(0); i < b.config.SafetyBlocks; i++ {
		block, err := b.HeaderByNumber(big.NewInt(int64(offset + i)))
		if err != nil {
			// todo
			panic(err)
		}
		encodedBlock, err := EncodeBlock(block, i == 0) // first block main, others safety
		if err != nil {
			// todo
		}
		encodedBlocks = append(encodedBlocks, *encodedBlock)
	}

	return encodedBlocks
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

func (b *Bridge) waitForSafetyBlocks(event *contracts.TransferEvent) error {
	blockChannel := make(chan *types.Header)
	blockSub, err := b.Client.SubscribeNewHead(context.Background(), blockChannel)
	if err != nil {
		return err
	}

	// init block num
	currentBlockNum, err := b.Client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	// if current block's number - event's block's number less or equal than number of safety blocks...
	// 0 <= 10 -> 1 <= 10 -> 2 <= 10 -> ... -> 11 !<= 10 -> exit from loop
	for currentBlockNum-event.Raw.BlockNumber <= b.config.SafetyBlocks {
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
