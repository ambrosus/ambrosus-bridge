package service_monitoring

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

func (m *Monitoring) balanceLoop() {
	for {
		err := m.sendLowBalances()
		if err != nil {
			m.logger.Error().Err(err).Msg("error sending balances")
		}
		time.Sleep(time.Duration(m.cfg.BalanceCheckTime) * time.Second)
	}
}

func (m *Monitoring) sendLowBalances() error {
	balances, err := m.getBalances()
	if err != nil {
		return fmt.Errorf("getBalances: %w", err)
	}
	lowBalances := m.getLowBalances(balances)
	if len(lowBalances) == 0 {
		return nil
	}

	text := "Low balances:\n"
	for addr, balance := range lowBalances {
		text += fmt.Sprintf(" - %s: %v\n", m.formatAddress(addr), decimal.NewFromBigInt(balance, 0).Div(decimal.NewFromFloat(1e18)))
	}
	m.monitoringLogger.Info().Msg(text)
	return nil
}

func (m *Monitoring) getBalances() (map[common.Address]*big.Int, error) {
	balances := make(map[common.Address]*big.Int)
	relays, err := m.getRelays()
	if err != nil {
		return nil, err
	}
	for _, relay := range relays {
		balance, err := m.bridge.GetClient().BalanceAt(context.Background(), relay, nil)
		if err != nil {
			return nil, fmt.Errorf("BalanceAt(%v): %w", relay, err)
		}
		balances[relay] = balance
	}
	return balances, nil

}

func (m *Monitoring) getLowBalances(balances map[common.Address]*big.Int) map[common.Address]*big.Int {
	lowBalances := make(map[common.Address]*big.Int)
	for relay, balance := range balances {
		if m.isBalanceLow(balance) {
			lowBalances[relay] = balance
		}
	}
	return lowBalances
}

func (m *Monitoring) isBalanceLow(balance *big.Int) bool {
	return balance.Cmp(m.minimumRequiredBalance) <= 0
}
