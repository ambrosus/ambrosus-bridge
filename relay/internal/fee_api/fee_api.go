package fee_api

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kofalt/go-memoize"
	"github.com/rs/cors"
)

const (
	signatureFeeTimestamp = 30 * 60 // 30 minutes
	cacheExpiration       = time.Minute * 10
)

type BridgeFeeApi interface {
	Sign(message []byte) ([]byte, error)
	TokenPrice(tokenAddress common.Address) (float64, error) // TokenPrice returns token`s price in USD
	GetTransferFee() *big.Int
	GetWrapperAddress() common.Address
	GetMinBridgeFee() *big.Float // GetMinBridgeFee returns the minimal bridge fee that can be used
}

type FeeAPI struct {
	amb, side BridgeFeeApi

	cache *memoize.Memoizer
}

func NewFeeAPI(amb, side BridgeFeeApi) *FeeAPI {
	return &FeeAPI{
		amb:   amb,
		side:  side,
		cache: memoize.NewMemoizer(cacheExpiration, time.Hour),
	}
}

func (p *FeeAPI) Run(endpoint string, ip string, port int) {
	c := p.setupCORS()

	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, p.feesHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), c.Handler(mux)))
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
