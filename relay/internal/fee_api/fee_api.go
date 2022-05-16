package fee_api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
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
	networks.BridgeFeeApi
}

func NewFeeAPI(bridgeFeeApi networks.BridgeFeeApi) *FeeAPI {
	return &FeeAPI{BridgeFeeApi: bridgeFeeApi}
}

func (p *FeeAPI) Run(endpoint string, ip string, port int) {
	http.HandleFunc(endpoint, p.feesHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", ip, port), nil))
}
