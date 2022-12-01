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
	"github.com/rs/zerolog"
)

var (
	bigZero         = big.NewInt(0)
	bridgeAbi, _    = bindings.BridgeMetaData.GetAbi()
	triggerMethodID = bridgeAbi.Methods["triggerTransfers"].ID
)

type explorerClient interface {
	// TxListByFromListToAddresses should return all transactions in desc sort filtering by `fromList` and `to` fields
	// if `txFilters.untilTxHash` is not nil, return tx until the `untilTxHash` is reached NOT INCLUDING
	TxListByFromListToAddresses(from []string, to string, txFilters explorers_clients.TxFilters) ([]*explorers_clients.Transaction, error)
}

type transferFeeTracker struct {
	bridge     networks.Bridge
	sideBridge networks.Bridge
	logger     *zerolog.Logger
	sideLogger *zerolog.Logger

	explorer     explorerClient
	sideExplorer explorerClient

	transferFeeIncludedTxsFromAddresses     []string
	sideTransferFeeIncludedTxsFromAddresses []string

	transferFeeTxsFromBlock     uint64
	sideTransferFeeTxsFromBlock uint64

	latestProcessedEvent      uint64
	latestSideProcessedTxHash *common.Hash
	latestThisProcessedTxHash *common.Hash

	totalWithdrawCount *big.Int
	totalSideGas       *big.Int
	totalThisGas       *big.Int
}

func newTransferFeeTracker(
	bridge, sideBridge networks.Bridge,
	explorer, sideExplorer explorerClient,
	transferFeeIncludedTxsFromAddresses []string, sideTransferFeeIncludedTxsFromAddresses []string,
	transferFeeTxsFromBlock, sideTransferFeeTxsFromBlock uint64,
) (*transferFeeTracker, error) {
	logger := bridge.GetLogger().With().Str("service", "TransferFeeTracker").Logger()
	sideLogger := sideBridge.GetLogger().With().Str("service", "TransferFeeTracker").Logger()

	p := &transferFeeTracker{
		bridge:                                  bridge,
		sideBridge:                              sideBridge,
		logger:                                  &logger,
		sideLogger:                              &sideLogger,
		explorer:                                explorer,
		sideExplorer:                            sideExplorer,
		transferFeeIncludedTxsFromAddresses:     transferFeeIncludedTxsFromAddresses,
		sideTransferFeeIncludedTxsFromAddresses: sideTransferFeeIncludedTxsFromAddresses,
		transferFeeTxsFromBlock:                 transferFeeTxsFromBlock,
		sideTransferFeeTxsFromBlock:             sideTransferFeeTxsFromBlock,
		totalWithdrawCount:                      big.NewInt(0),
		totalSideGas:                            big.NewInt(0),
		totalThisGas:                            big.NewInt(0),
	}

	if err := p.init(); err != nil {
		return nil, err
	}
	go p.WatchUnlocksLoop()

	return p, nil
}

func (p *transferFeeTracker) GasPerWithdraw() (thisGas, sideGas *big.Int) {
	p.logger.Debug().Msgf("GasPerWithdraw: totalSideGas: %d, totalThisGas: %d, totalWithdrawCount: %d", p.totalSideGas, p.totalThisGas, p.totalWithdrawCount)
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

	p.latestProcessedEvent = 0                         // 0 cuz in `processEvents` we'll +1 to it
	return p.processEvents(latestEventId.Uint64() - 1) // -1 cuz we need latest unlocked instead of locked event id
}

func (p *transferFeeTracker) processEvents(newEventId uint64) error {
	withdrawsCount, err := p.bridge.Events().GetWithdrawsCount(p.latestProcessedEvent+1, newEventId)
	if err != nil {
		return fmt.Errorf("getWithdrawsCount: %w", err)
	}
	if withdrawsCount == 0 {
		p.latestProcessedEvent = newEventId
		return nil
	}

	// get side bridge txs from explorer (for submit/unlock methods)
	sideBridgeTxList, err := p.sideExplorer.TxListByFromListToAddresses(
		p.sideTransferFeeIncludedTxsFromAddresses,
		p.sideBridge.GetContractAddress().Hex(),
		explorers_clients.TxFilters{FromBlock: p.sideTransferFeeTxsFromBlock, UntilTxHash: p.latestSideProcessedTxHash},
	)
	if err != nil {
		return fmt.Errorf("get side bridge tx list: %w", err)
	}

	// get this bridge txs from explorer (for triggerTransfers method)
	thisBridgeTxList, err := p.explorer.TxListByFromListToAddresses(
		p.transferFeeIncludedTxsFromAddresses,
		p.bridge.GetContractAddress().Hex(),
		explorers_clients.TxFilters{FromBlock: p.transferFeeTxsFromBlock, UntilTxHash: p.latestThisProcessedTxHash},
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
	p.totalThisGas = p.totalThisGas.Add(p.totalThisGas, totalThisGas)
	p.totalWithdrawCount = p.totalWithdrawCount.Add(p.totalWithdrawCount, big.NewInt(int64(withdrawsCount)))

	p.latestProcessedEvent = newEventId

	p.logger.Debug().Msgf("from new event we got gas: %d, withdrawsCount: %d. totalSideGas: %d, totalWithdrawsCount: %d", totalSideGas, withdrawsCount, p.totalSideGas, p.totalWithdrawCount)

	return nil
}

// -------------------- watcher --------------------------

func (p *transferFeeTracker) WatchUnlocksLoop() {
	for {
		if err := p.watchUnlocks(); err != nil {
			p.sideLogger.Error().Err(fmt.Errorf("watchUnlocks: %w", err)).Msg("")
		}
		time.Sleep(time.Minute)
	}
}

func (p *transferFeeTracker) checkOldUnlocks() error {
	latestLockedEventId, err := getOldestLockedEventId(p.sideBridge.GetContract())
	if err != nil {
		return err
	}
	latestUnlockedEventId := latestLockedEventId.Uint64() - 1

	if p.latestProcessedEvent != latestUnlockedEventId {
		p.sideLogger.Info().Str("event_id", fmt.Sprint(latestUnlockedEventId)).Msg("Process old unlock...")
		p.processEvents(latestUnlockedEventId)
	}
	return nil
}

func (p *transferFeeTracker) watchUnlocks() error {
	if err := p.checkOldUnlocks(); err != nil {
		return fmt.Errorf("checkOldUnlocks: %w", err)
	}

	// main loop
	for {
		event, err := p.sideBridge.Events().WatchTransferFinish()
		if err != nil {
			return fmt.Errorf("watching unlock transfers: %w", err)
		}

		p.sideLogger.Info().Str("event_id", event.EventId.String()).Msg("Found new TransferFinish event")
		if err := p.processEvents(event.EventId.Uint64()); err != nil {
			return err
		}
	}
}

// --------------------- side bridge getters --------------------

func getOldestLockedEventId(contract interfaces.BridgeContract) (*big.Int, error) {
	return contract.OldestLockedEventId(nil)
}

//

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
