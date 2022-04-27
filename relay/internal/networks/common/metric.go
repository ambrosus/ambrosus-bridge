package common

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
)

func (b *CommonBridge) SetUsedGasMetric(methodName string, usedGas uint64, gasPrice *big.Int) {
	metric.UsedGas.WithLabelValues(b.Name, methodName).Observe(float64(usedGas))

	gasCost := new(big.Int).Mul(big.NewInt(int64(usedGas)), gasPrice)
	metric.GasCost.WithLabelValues(b.Name, methodName).Observe(weiToGwei(gasCost))
}

func (b *CommonBridge) IncTxCountMetric(methodName string) {
	metric.TxCount.WithLabelValues(b.Name, methodName).Inc()
}

func (b *CommonBridge) IncFailedTxCountMetric(methodName string) {
	metric.FailedTxCount.WithLabelValues(b.Name, methodName).Inc()
}

func (b *CommonBridge) AddWithdrawalsCountMetric(count int) {
	metric.TransfersCount.WithLabelValues(b.Name).Inc()
	metric.WithdrawalsCount.WithLabelValues(b.Name).Add(float64(count))
}

func (b *CommonBridge) SetRelayBalanceMetric() {
	balance, err := b.getBalanceGWei(b.Auth.From)
	if err != nil {
		b.Logger.Error().Err(err).Msg("get balance error")
		return
	}

	metric.RelayBalance.WithLabelValues(b.Name).Set(balance)
}
