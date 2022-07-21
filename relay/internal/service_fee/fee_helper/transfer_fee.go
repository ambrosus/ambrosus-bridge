package fee_helper

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	cb "github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/errgroup"
)

const (
	eventsForGasCalc = 20
)

type transferFeeTracker struct {
	bridge     networks.Bridge
	sideBridge networks.Bridge

	latestProcessedEvent uint64

	totalWithdrawCount *big.Int
	totalGas           *big.Int
}

func newTransferFeeTracker(bridge, sideBridge networks.Bridge) (*transferFeeTracker, error) {
	p := &transferFeeTracker{
		bridge:             bridge,
		sideBridge:         sideBridge,
		totalWithdrawCount: big.NewInt(0),
		totalGas:           big.NewInt(0),
	}

	if err := p.init(); err != nil {
		return nil, err
	}
	go p.WatchUnlocksLoop()

	return p, nil
}

func (p *transferFeeTracker) GasPerWithdraw() *big.Int {
	p.bridge.GetLogger().Debug().Msgf("GasPerWithdraw: totalGas: %d, totalWithdrawCount: %d", p.totalGas, p.totalWithdrawCount)
	// if there's no transfers then return nil
	if p.totalWithdrawCount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	return new(big.Int).Div(p.totalGas, p.totalWithdrawCount)
}

func (p *transferFeeTracker) init() error {
	latestEventId, err := getOldestLockedEventId(p.sideBridge.GetContract())
	if err != nil {
		return err
	}

	// there's no unlocked events, so we can't get gas cost per withdraw
	if latestEventId.Cmp(big.NewInt(1)) == 0 {
		return nil
	}

	latestProcessedEvent := latestEventId.Int64() - eventsForGasCalc - 1 // -1 cuz in `processEvents` we'll +1 to it
	if latestProcessedEvent < 0 {
		latestProcessedEvent = 0
	}

	p.latestProcessedEvent = uint64(latestProcessedEvent)
	return p.processEvents(latestEventId.Uint64() - 1) // -1 cuz we need latest unlocked instead of locked event id
}

func (p *transferFeeTracker) processEvents(newEventId uint64) error {
	// get events ids for getting submits, unlocks and transfers for "batch" requests
	var eventIds []*big.Int
	for i := p.latestProcessedEvent + 1; i <= newEventId; i++ { // +1 cuz we've already calculated the gas cost for that event
		eventIds = append(eventIds, big.NewInt(int64(i)))
	}

	// get events batch requests
	transfers, err := getTransfersByIds(p.bridge.GetContract(), eventIds)
	if err != nil {
		return fmt.Errorf("get transfers by ids: %v", err)
	}
	submits, err := getTransferSubmitsByIds(p.sideBridge.GetContract(), eventIds)
	if err != nil {
		return fmt.Errorf("get transfer submits by ids: %v", err)
	}
	unlocks, err := getTransferUnlocksByIds(p.sideBridge.GetContract(), eventIds)
	if err != nil {
		return fmt.Errorf("get transfer unlocks by ids: %v", err)
	}

	// save tx hashes made by (side) relay

	var relayTxHashes []common.Hash
	for _, event := range submits {
		relayTxHashes = append(relayTxHashes, event.Raw.TxHash)
	}
	for _, event := range unlocks {
		relayTxHashes = append(relayTxHashes, event.Raw.TxHash)
	}

	// calc how much gas used in this txs
	gas, _, err := usedGas(p.sideBridge.GetClient(), helpers.Unique(relayTxHashes))
	if err != nil {
		return err
	}

	// calc withdraws count in transfer events
	var withdrawsCount int
	for _, event := range transfers {
		withdrawsCount += len(event.Queue)
	}

	if withdrawsCount == 0 {
		return nil
	}

	p.totalGas = p.totalGas.Add(p.totalGas, gas)
	p.totalWithdrawCount = p.totalWithdrawCount.Add(p.totalWithdrawCount, big.NewInt(int64(withdrawsCount)))
	p.latestProcessedEvent = newEventId
	p.bridge.GetLogger().Debug().Msgf("from new event we got gas: %d, withdrawsCount: %d. totalGas: %d, totalWithdrawsCount: %d", gas, withdrawsCount, p.totalGas, p.totalWithdrawCount)

	return nil
}

func usedGas(client ethclients.ClientInterface, txs []common.Hash) (*big.Int, *big.Int, error) {
	// fetch transactions and sum gas

	totalGas := new(big.Int)
	totalGasCost := new(big.Int)

	lock := sync.Mutex{}
	sem := make(chan interface{}, 10) // max 20 simultaneous requests
	errGroup := new(errgroup.Group)

	for _, txHash := range txs {
		txHash := txHash
		errGroup.Go(func() error {
			sem <- 0
			defer func() { <-sem }()

			tx, _, err := client.TransactionByHash(context.Background(), txHash)
			if err != nil {
				return fmt.Errorf("get transaction by hash: %w", err)
			}

			txGasPrice, err := cb.GetTxGasPrice(client, tx)
			if err != nil {
				return fmt.Errorf("get tx gas price: %w", err)
			}
			receipt, err := client.TransactionReceipt(context.Background(), txHash)
			if err != nil {
				return fmt.Errorf("get transaction receipt: %w", err)
			}

			txCost := new(big.Int).Mul(txGasPrice, big.NewInt(int64(receipt.GasUsed)))

			lock.Lock()
			totalGas.Add(totalGas, big.NewInt(int64(receipt.GasUsed)))
			totalGasCost.Add(totalGasCost, txCost)
			lock.Unlock()
			return nil
		})
	}
	if err := errGroup.Wait(); err != nil {
		return nil, nil, fmt.Errorf("calc used gas: %w", err)
	}

	return totalGasCost, totalGas, nil
}

// -------------------- watcher --------------------------

func (p *transferFeeTracker) WatchUnlocksLoop() {
	for {
		if err := p.watchUnlocks(); err != nil {
			p.sideBridge.GetLogger().Error().Err(err).Msg("price tracker watchUnlocks error")
		}
		time.Sleep(time.Minute)
	}
}

func (p *transferFeeTracker) watchUnlocks() error {

	eventCh := make(chan *bindings.BridgeTransferFinish)
	eventSub, err := p.sideBridge.GetWsContract().WatchTransferFinish(nil, eventCh, nil)
	if err != nil {
		return fmt.Errorf("watchTransferFinish: %w", err)
	}
	defer eventSub.Unsubscribe()

	// main loop
	for {
		select {
		case err := <-eventSub.Err():
			return fmt.Errorf("watching unlock transfers: %w", err)
		case event := <-eventCh:
			if event.Raw.Removed {
				continue
			}

			p.sideBridge.GetLogger().Info().Str("event_id", event.EventId.String()).Msg("Found new TransferFinish event")
			if err := p.processEvents(event.EventId.Uint64()); err != nil {
				return err
			}
		}
	}
}

// --------------------- side bridge getters --------------------

func getOldestLockedEventId(contract interfaces.BridgeContract) (*big.Int, error) {
	return contract.OldestLockedEventId(nil)
}

func getTransfersByIds(contract interfaces.BridgeContract, eventIds []*big.Int) (transfers []*bindings.BridgeTransfer, err error) {
	logTransfer, err := contract.FilterTransfer(nil, eventIds)
	if err != nil {
		return nil, fmt.Errorf("filter transfer: %w", err)
	}

	for logTransfer.Next() {
		if !logTransfer.Event.Raw.Removed {
			transfers = append(transfers, logTransfer.Event)
		}
	}
	return transfers, nil
}

func getTransferSubmitsByIds(contract interfaces.BridgeContract, eventIds []*big.Int) (submits []*bindings.BridgeTransferSubmit, err error) {
	logSubmit, err := contract.FilterTransferSubmit(nil, eventIds)
	if err != nil {
		return nil, fmt.Errorf("filter transfer submit: %w", err)
	}

	for logSubmit.Next() {
		if !logSubmit.Event.Raw.Removed {
			submits = append(submits, logSubmit.Event)
		}
	}
	return submits, nil
}

func getTransferUnlocksByIds(contract interfaces.BridgeContract, eventIds []*big.Int) (unlocks []*bindings.BridgeTransferFinish, err error) {
	logUnlock, err := contract.FilterTransferFinish(nil, eventIds)
	if err != nil {
		return nil, fmt.Errorf("filter transfer finish: %w", err)
	}

	for logUnlock.Next() {
		if !logUnlock.Event.Raw.Removed {
			unlocks = append(unlocks, logUnlock.Event)
		}
	}
	return unlocks, nil
}
