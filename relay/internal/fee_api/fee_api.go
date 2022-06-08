package fee_api

import (
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/common"
	"github.com/kofalt/go-memoize"
	"github.com/rs/cors"
)

type BridgePriceTracker interface {
	GasPerWithdraw() *big.Int
}

type FeeAPI struct {
	amb, side     networks.BridgeFeeApi
	ambPT, sidePT BridgePriceTracker

	cache *memoize.Memoizer
}

const (
	cacheExpiration = time.Minute * 10
)

func NewFeeAPI(amb, side networks.BridgeFeeApi) (*FeeAPI, error) {
	ambPT, err := common.NewTransferFeeTracker(amb, side.(networks.TransferFeeCalc))
	if err != nil {
		return nil, err
	}
	sidePT, err := common.NewTransferFeeTracker(side, amb.(networks.TransferFeeCalc))
	if err != nil {
		return nil, err
	}

	return &FeeAPI{
		amb:    amb,
		side:   side,
		ambPT:  ambPT,
		sidePT: sidePT,
		cache:  memoize.NewMemoizer(cacheExpiration, time.Hour),
	}, nil
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
