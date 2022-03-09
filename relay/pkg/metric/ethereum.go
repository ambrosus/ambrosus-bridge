package metric

import (
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rs/zerolog/log"
)

var ContractBalanceGWeiGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "bridge_contract_balance_gwei",
	Help: "Balance of a contract in the bridge",
}, []string{"bridge_name"})

func SetContractBalance(bridgeName string, client *ethclient.Client, contractAddress common.Address) {
	balance, err := ethereum.GetBalanceGWei(client, contractAddress)
	if err != nil {
		log.Error().Err(err).Msg("error when getting contract balance in GWei")
		return
	}

	ContractBalanceGWeiGauge.WithLabelValues(bridgeName).Set(balance)
}
