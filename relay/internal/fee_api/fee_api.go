package fee_api

import (
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

func (p *FeeAPI) Run(endpoint string) {
	http.HandleFunc(endpoint, p.feesHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
