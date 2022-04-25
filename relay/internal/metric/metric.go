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
)

func ServeEndpoint(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}
