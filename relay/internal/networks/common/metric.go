package common

import "github.com/ambrosus/ambrosus-bridge/relay/internal/metric"

func (b *CommonBridge) SetUsedGasMetric(usedGas uint64) {
	metric.UsedGas.WithLabelValues(b.Name).Observe(float64(usedGas))
}

func (b *CommonBridge) IncTxCountMetric() {
	metric.TxCount.WithLabelValues(b.Name).Inc()
}

func (b *CommonBridge) IncFailedTxCountMetric() {
	metric.FailedTxCount.WithLabelValues(b.Name).Inc()
}

func (b *CommonBridge) AddWithdrawalsCountMetric(count int) {
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
