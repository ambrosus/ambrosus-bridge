package metric

import (
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
)

func SetContractBalance(gauge prometheus.Gauge, client *ethclient.Client, contractAddress common.Address) error {
	balance, err := ethereum.GetBalanceGWei(client, contractAddress)
	if err != nil {
		return err
	}

	gauge.Set(balance)
	return nil
}
