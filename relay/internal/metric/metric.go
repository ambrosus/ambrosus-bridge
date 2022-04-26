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
		Buckets: []float64{0, 21_000, 50_000, 75_000, 100_000, 150_000,
			200_000, 250_000, 300_000, 350_000, 400_000, 450_000,
			500_000, 550_000, 600_000, 650_000, 700_000, 750_000,
			800_000, 850_000, 900_000, 950_000, 1_000_000, 1_100_000,
			1_200_000, 1_300_000, 1_400_000, 1_500_000, 1_600_000,
			1_700_000, 1_800_000, 1_900_000, 2_000_000, 2_200_000,
			2_400_000, 2_600_000, 2_800_000, 3_000_000, 3_500_000,
			4_000_000, 4_500_000, 5_000_000},
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
