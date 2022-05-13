package price_api

import (
	"log"
	"net/http"
)

type Signer interface {
	Sign(tokenPrice float64, tokenAddress string) ([]byte, error)
}

type PriceGetter interface {
	GetPrice(tokenAddress string) (float64, error)
}

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

type PriceAPI struct {
	signer      Signer
	priceGetter PriceGetter
}

func NewPriceAPI(signer Signer, priceGetter PriceGetter) *PriceAPI {
	return &PriceAPI{
		signer:      signer,
		priceGetter: priceGetter,
	}
}

func (p *PriceAPI) Run(endpoint string) {
	http.HandleFunc(endpoint, p.priceHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
