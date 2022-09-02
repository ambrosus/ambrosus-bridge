package fee_api

import (
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/service_fee/fee_api/middlewares"
	"github.com/ethereum/go-ethereum/common"
	"github.com/kofalt/go-memoize"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

const (
	signatureFeeTimestamp = 30 * 60 // 30 minutes
	cacheExpiration       = time.Minute * 10
)

type BridgeFeeApi interface {
	networks.Bridge
	Sign(message []byte) ([]byte, error)
	GetTransferFee() *big.Int
	GetWrapperAddress() common.Address
	GetMinBridgeFee() decimal.Decimal // GetMinBridgeFee returns the minimal bridge fee that can be used
	GetDefaultTransferFee() *big.Int
}

type FeeAPI struct {
	amb, side BridgeFeeApi

	cache  *memoize.Memoizer
	logger *zerolog.Logger
}

func NewFeeAPI(amb, side BridgeFeeApi, logger zerolog.Logger) *FeeAPI {
	return &FeeAPI{
		amb:    amb,
		side:   side,
		cache:  memoize.NewMemoizer(cacheExpiration, time.Hour),
		logger: &logger,
	}
}

func (p *FeeAPI) Run(endpoint string, ip string, port int) {
	// endpoints
	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, p.feesHandler)

	// init middlewares
	corsMiddleware := p.setupCORS().Handler
	loggingMiddleware := middlewares.LoggingMiddleware(p.logger)

	// apply middlewares
	handler := corsMiddleware(mux)
	handler = loggingMiddleware(handler)
	handler = middlewares.MetricsMiddleware(handler)

	p.logger.Fatal().Err(http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), handler)).Msg("")
}

func (p *FeeAPI) setupCORS() *cors.Cors {
	// setup CORS
	var allowedOrigins = []string{"*"}
	if os.Getenv("STAGE") == "prod" {
		allowedOrigins = []string{"https://*.ambrosus.io"}
	}

	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
	})
}