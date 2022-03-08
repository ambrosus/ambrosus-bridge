package metric

import (
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var ContractBalanceGWeiGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "bridge_contract_balance_gwei",
	Help: "Balance of a contract in the bridge",
}, []string{"bridge_name"})

func SetContractBalance(bridgeName string, client *ethclient.Client, contractAddress common.Address) error {
	balance, err := ethereum.GetBalanceGWei(client, contractAddress)
	if err != nil {
		return err
	}

	ContractBalanceGWeiGauge.WithLabelValues(bridgeName).Set(balance)
	return nil
}
