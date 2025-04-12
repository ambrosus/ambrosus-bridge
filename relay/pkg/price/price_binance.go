package price

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type response struct {
	Price float64 `json:"price,string"`

	Reason string `json:"reason"` // when error occurred
}

func GetBinance(token *TokenInfo, ticker string) (price float64, err error) {
	amount := math.Pow10(int(token.Decimals))

	client := http.Client{}

	url := fmt.Sprint("https://api.binance.com/api/v1/ticker/price?symbol=", ticker, "USDT")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}

	return r.Price / amount, err
}
