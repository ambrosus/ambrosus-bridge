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

	// todo use
	GasPrice = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gas_price",
		Help: "Gas price of a transaction (in gwei)",
	}, []string{"network_name"})

	UsedGas = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "used_gas",
		Help: "Used gas of a transaction (in gwei)",
	}, []string{"network_name"})
)

func ServeEndpoint(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}
