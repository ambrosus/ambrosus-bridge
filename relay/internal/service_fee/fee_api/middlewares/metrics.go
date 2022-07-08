package middlewares

import (
	"net/http"

	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	m := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	return std.Handler("", m, next)
}
