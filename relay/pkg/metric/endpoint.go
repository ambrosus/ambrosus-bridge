package metric

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func ServeEndpoint(ip string, port int) error {
	addr := fmt.Sprintf("%s:%d", ip, port)
	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(addr, nil); err != nil {
		return err
	}
	return nil
}
