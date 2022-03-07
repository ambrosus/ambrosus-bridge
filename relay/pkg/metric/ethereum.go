package metric

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/prometheus/client_golang/prometheus"
)

func SetContractBalance(gauge prometheus.Gauge, client *ethclient.Client, contractAddress common.Address) error {
	balance, err := client.BalanceAt(context.Background(), contractAddress, nil)
	if err != nil {
		return err
	}
	balanceGWei := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(params.GWei))
	balanceFloat64, _ := balanceGWei.Float64()
	gauge.Set(balanceFloat64)
	return nil
}
