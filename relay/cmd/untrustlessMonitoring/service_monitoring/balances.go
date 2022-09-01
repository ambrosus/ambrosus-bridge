package service_monitoring

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func (m *Monitoring) balanceLoop() {
	for {
		err := m.sendBalances()
		if err != nil {
			m.logger.Error().Err(err).Msg("error sending balances")
		}
		time.Sleep(time.Duration(m.cfg.BalanceCheckTime) * time.Second)
	}
}

func (m *Monitoring) sendBalances() error {
	balances, err := m.getBalances()
	if err != nil {
		return fmt.Errorf("getBalances: %w", err)
	}
	text := "Balances:\n"
	for addr, balance := range balances {
		text += fmt.Sprintf(" - %s: %v\n", m.formatAddress(addr), balance)
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
