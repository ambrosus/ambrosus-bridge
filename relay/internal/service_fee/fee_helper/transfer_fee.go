package fee_helper

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_helper/explorers_clients"
	"github.com/ethereum/go-ethereum/common"
)

const (
	eventsForGasCalc = 20
)

var (
	bigZero         = big.NewInt(0)
	bridgeAbi, _    = bindings.BridgeMetaData.GetAbi()
	triggerMethodID = bridgeAbi.Methods["triggerTransfers"].ID
)

type explorerClient interface {
	// TxListByFromToAddresses should return all transactions in desc sort filtering by `from` and `to` fields
	// if `untilTxHash` is not nil, return tx until the `untilTxHash` is reached NOT INCLUDING
	TxListByFromToAddresses(from, to string, untilTxHash *common.Hash) ([]*explorers_clients.Transaction, error)
}

type transferFeeTracker struct {
	bridge     networks.Bridge
	sideBridge networks.Bridge

	explorer     explorerClient
	sideExplorer explorerClient

	latestProcessedEvent      uint64
	latestSideProcessedTxHash *common.Hash
	latestThisProcessedTxHash *common.Hash

	totalWithdrawCount *big.Int
	totalSideGas       *big.Int
	totalThisGas       *big.Int
}

func newTransferFeeTracker(bridge, sideBridge networks.Bridge, explorer, sideExplorer explorerClient) (*transferFeeTracker, error) {
	p := &transferFeeTracker{
		bridge:             bridge,
		sideBridge:         sideBridge,
		explorer:           explorer,
		sideExplorer:       sideExplorer,
		totalWithdrawCount: big.NewInt(0),
		totalSideGas:       big.NewInt(0),
		totalThisGas:       big.NewInt(0),
	}

	if err := p.init(); err != nil {
		return nil, err
	}
	go p.WatchUnlocksLoop()

	return p, nil
}

func (p *transferFeeTracker) GasPerWithdraw() (thisGas, sideGas *big.Int) {
	p.bridge.GetLogger().Debug().Msgf("GasPerWithdraw: totalSideGas: %d, totalThisGas: %d, totalWithdrawCount: %d", p.totalSideGas, p.totalThisGas, p.totalWithdrawCount)
	// if there's no transfers then return zeroes
	if p.totalWithdrawCount.Cmp(big.NewInt(0)) == 0 {
		return bigZero, bigZero
	}

	return new(big.Int).Div(p.totalThisGas, p.totalWithdrawCount), new(big.Int).Div(p.totalSideGas, p.totalWithdrawCount)
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
	withdrawsCount, err := getWithdrawsCount(p.bridge.GetContract(), p.latestProcessedEvent+1, newEventId)
	if err != nil {
		return fmt.Errorf("getWithdrawsCount: %w", err)
	}
	if withdrawsCount == 0 {
		return nil
	}

	// get side bridge txs from explorer (for submit/unlock methods)
	sideBridgeTxList, err := p.sideExplorer.TxListByFromToAddresses(
		p.sideBridge.GetAuth().From.Hex(),
		p.sideBridge.GetContractAddress().Hex(),
		p.latestSideProcessedTxHash,
	)
	if err != nil {
		return fmt.Errorf("get side bridge tx list: %w", err)
	}

	// get this bridge txs from explorer (for triggerTransfers method)
	thisBridgeTxList, err := p.sideExplorer.TxListByFromToAddresses(
		p.sideBridge.GetAuth().From.Hex(),
		p.sideBridge.GetContractAddress().Hex(),
		p.latestThisProcessedTxHash,
	)
	if err != nil {
		return fmt.Errorf("get this bridge tx list: %w", err)
	}

	if len(sideBridgeTxList) != 0 {
		h := common.HexToHash(sideBridgeTxList[0].Hash)
		p.latestSideProcessedTxHash = &h
	}
	if len(thisBridgeTxList) != 0 {
		h := common.HexToHash(thisBridgeTxList[0].Hash)
		p.latestThisProcessedTxHash = &h
	}

	// calc total gas

	totalSideGas := calcGasCost(filterTxsWithoutTriggers(sideBridgeTxList))
	totalThisGas := calcGasCost(filterTxsWithTriggers(thisBridgeTxList))

	p.totalSideGas = p.totalSideGas.Add(p.totalSideGas, totalSideGas)
	p.totalThisGas = p.totalSideGas.Add(p.totalThisGas, totalThisGas)
	p.totalWithdrawCount = p.totalWithdrawCount.Add(p.totalWithdrawCount, big.NewInt(int64(withdrawsCount)))

	p.latestProcessedEvent = newEventId

	p.bridge.GetLogger().Debug().Msgf("from new event we got gas: %d, withdrawsCount: %d. totalSideGas: %d, totalWithdrawsCount: %d", totalSideGas, withdrawsCount, p.totalSideGas, p.totalWithdrawCount)

	return nil
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

//

func getWithdrawsCount(contract interfaces.BridgeContract, fromEvent, toEvent uint64) (int, error) {
	// get events ids for getting submits, unlocks and transfers for "batch" requests
	var eventIds []*big.Int
	for i := fromEvent; i <= toEvent; i++ { // +1 cuz we've already calculated the gas cost for that event
		eventIds = append(eventIds, big.NewInt(int64(i)))
	}

	// get events batch requests
	transfers, err := getTransfersByIds(contract, eventIds)
	if err != nil {
		return 0, fmt.Errorf("get transfers by ids: %w", err)
	}
	// calc withdraws count in transfer events
	var withdrawsCount int
	for _, event := range transfers {
		withdrawsCount += len(event.Queue)
	}
	return withdrawsCount, nil
}

func calcGasCost(txs []*explorers_clients.Transaction) *big.Int {
	totalGas := new(big.Int)
	for _, tx := range txs {
		gas := new(big.Int).Mul(tx.GasPrice, new(big.Int).SetUint64(tx.GasUsed))
		totalGas = totalGas.Add(totalGas, gas)
	}
	return totalGas
}

func isTrigger(tx *explorers_clients.Transaction) bool {
	return bytes.Equal(common.FromHex(tx.Input), triggerMethodID)
}

func filterTxsWithTriggers(txs []*explorers_clients.Transaction) []*explorers_clients.Transaction {
	return explorers_clients.FilterTxsByCallback(txs, func(tx *explorers_clients.Transaction) bool { return isTrigger(tx) })
}

func filterTxsWithoutTriggers(txs []*explorers_clients.Transaction) []*explorers_clients.Transaction {
	return explorers_clients.FilterTxsByCallback(txs, func(tx *explorers_clients.Transaction) bool { return !isTrigger(tx) })
}
