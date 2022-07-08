package api

import (
	"fmt"
	"math/big"
	"net/http"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/api/middlewares"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

type FeeAPI struct {
	Service FeeService
}

type FeeService interface {
	GetFees(tokenAddress common.Address, reqAmount *big.Int, isAmb, isAmountWithFees bool) (
		bridgeFeeBigInt, transferFeeBigInt, amountBigInt *big.Int, signature []byte, err error)
}

func (p *FeeAPI) Run(endpoint string, ip string, port int, logger *zerolog.Logger) error {
	// endpoints
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, p.feesHandler)

	// init middlewares
	corsMiddleware := p.setupCORS().Handler
	loggingMiddleware := middlewares.LoggingMiddleware(logger)

	// apply middlewares
	handler := corsMiddleware(mux)
	handler = loggingMiddleware(handler)
	handler = middlewares.MetricsMiddleware(handler)

	return http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), handler)
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
