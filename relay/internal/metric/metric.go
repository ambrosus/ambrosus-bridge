package metric

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RelayBalance = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "relay_balance",
		Help: "Balance of a relay in the bridge (in gwei)",
	}, []string{"bridge_name"})

	UsedGas = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "used_gas",
		Help: "Used gas of a transaction (in gwei)",
	}, []string{"network_name"})

	TxCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tx_count",
		Help: "Sent transactions count",
	}, []string{"bridge_name", "method_name"})
	FailedTxCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "failed_tx_count",
		Help: "Sent transactions count that have been failed",
	}, []string{"bridge_name", "method_name"})

	TransfersCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "transfers_count",
		Help: "Transfers count",
	}, []string{"bridge_name"})
	WithdrawalsCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "withdrawals_count",
		Help: "Withdrawals count",
	}, []string{"bridge_name"})
)

func ServeEndpoint(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}
