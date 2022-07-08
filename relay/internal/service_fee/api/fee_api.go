package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/api/middlewares"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee"
	"github.com/rs/cors"
)

type FeeAPI struct {
	Service *fee.Fee
}

func (p *FeeAPI) Run(endpoint string, ip string, port int) {
	// endpoints
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, p.feesHandler)

	// init middlewares
	corsMiddleware := p.setupCORS().Handler
	loggingMiddleware := middlewares.LoggingMiddleware(p.Service.Logger)

	// apply middlewares
	handler := corsMiddleware(mux)
	handler = loggingMiddleware(handler)
	handler = middlewares.MetricsMiddleware(handler)

	p.Service.Logger.Fatal().Err(http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), handler)).Msg("")
}

func (p *FeeAPI) setupCORS() *cors.Cors {
	// setup CORS
	var allowedOrigins = []string{"*"}
	if os.Getenv("STAGE") == "main" {
		allowedOrigins = []string{"https://*.ambrosus.io"}
	}

	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
	})
}
