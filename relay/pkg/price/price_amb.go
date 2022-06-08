package price

import (
	"encoding/json"
	"net/http"
)

type ambResponse struct {
	Data ambResponse_ `json:"data"`
}
type ambResponse_ struct {
	Price float64 `json:"total_price_usd"`
}

func AmbToUSD() (float64, error) {
	resp, err := http.Get("https://token.ambrosus.io/price")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r ambResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return 0, err
	}

	return r.Data.Price, nil
}
