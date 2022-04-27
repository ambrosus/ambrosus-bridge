package common

import (
	"context"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/metric"
	"github.com/ethereum/go-ethereum/params"
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
	balance, err := b.Client.BalanceAt(context.Background(), b.Auth.From, nil)
	if err != nil {
		b.Logger.Error().Err(err).Msg("get balance error")
		return
	}
	metric.RelayBalance.WithLabelValues(b.Name).Set(weiToGwei(balance))
}

func weiToGwei(wei *big.Int) float64 {
	gWei := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.GWei))
	gWeiF, _ := gWei.Float64()
	return gWeiF
}
