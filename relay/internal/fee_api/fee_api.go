package fee_api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/rs/cors"
)

// func signData(pk *ecdsa.PrivateKey, tokenPrice float64, tokenAddress string) ([]byte, error) {
// 	var data bytes.Buffer
// 	if err := binary.Write(&data, binary.LittleEndian, tokenPrice); err != nil {
// 		return nil, err
// 	}
// 	data.WriteString(tokenAddress)
// 	return crypto.Sign(crypto.Keccak256(data.Bytes()), pk)
// }
//
// func getTokenPrice(tokenAddress string) (float64, error) {
// 	return 13.37, nil
// }

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
