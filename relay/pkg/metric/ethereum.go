package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RelayBalanceGWeiGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "bridge_relay_balance_gwei",
	Help: "Balance of a relay in the bridge",
}, []string{"bridge_name"})
