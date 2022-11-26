package metric

import (
	"context"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

func SetUsedGasMetric(b networks.Bridge, methodName string, usedGas uint64, gasPrice *big.Int) {
	UsedGas.WithLabelValues(b.GetName(), methodName).Observe(float64(usedGas))

	gasCost := new(big.Int).Mul(big.NewInt(int64(usedGas)), gasPrice)
	GasCost.WithLabelValues(b.GetName(), methodName).Observe(weiToGwei(gasCost))
}

func IncTxCountMetric(b networks.Bridge, methodName string) {
	TxCount.WithLabelValues(b.GetName(), methodName).Inc()
}

func IncFailedTxCountMetric(b networks.Bridge, methodName string) {
	FailedTxCount.WithLabelValues(b.GetName(), methodName).Inc()
}

func AddWithdrawalsCountMetric(b networks.Bridge, count int) {
	TransfersCount.WithLabelValues(b.GetName()).Inc()
	WithdrawalsCount.WithLabelValues(b.GetName()).Add(float64(count))
}

func SetRelayBalanceMetric(b networks.Bridge) {
	balance, err := b.GetClient().BalanceAt(context.Background(), b.GetAuth().From, nil)
	if err != nil {
		b.GetLogger().Error().Err(err).Msg("get balance error")
		return
	}
	RelayBalance.WithLabelValues(b.GetName()).Set(weiToGwei(balance))
}

func SetAmbFaucetBalanceMetric(b networks.Bridge, moneyAccount common.Address) {
	balance, err := b.GetClient().BalanceAt(context.Background(), moneyAccount, nil)
	if err != nil {
		b.GetLogger().Error().Err(err).Msg("get amb faucet balance error")
		return
	}
	AmbFaucetBalance.WithLabelValues(moneyAccount.String()).Set(weiToGwei(balance))
}

func weiToGwei(wei *big.Int) float64 {
	gWei := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.GWei))
	gWeiF, _ := gWei.Float64()
	return gWeiF
}
