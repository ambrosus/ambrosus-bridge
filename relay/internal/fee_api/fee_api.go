package fee_api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/rs/cors"
)

type FeeAPI struct {
	amb, side networks.BridgeFeeApi
}

func NewFeeAPI(amb, side networks.BridgeFeeApi) *FeeAPI {
	return &FeeAPI{amb: amb, side: side}
}

func (p *FeeAPI) Run(endpoint string, ip string, port int) {
	// setup CORS
	var allowedOrigins = []string{"*"}
	if os.Getenv("STAGE") == "main" {
		allowedOrigins = []string{"https://*.ambrosus.io"}
	}
	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
	})

	mux := http.NewServeMux()
	mux.HandleFunc(endpoint, p.feesHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), c.Handler(mux)))
}
