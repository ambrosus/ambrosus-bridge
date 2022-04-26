package binance_price

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	amb = coin("AMBBTC")
	Eth = coin("ETHBTC")
	Bnb = coin("BNBBTC")
)

type coin string

type response struct {
	Price float64 `json:"price,string"`
}

func CoinToAmb(symbol coin) (float64, error) {
	coinBtc, err := getPrice(symbol)
	if err != nil {
		return 0, err
	}
	ambBtc, err := getPrice(amb)
	if err != nil {
		return 0, err
	}

	return coinBtc / ambBtc, nil
}

func getPrice(symbol coin) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r response
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return 0, err
	}

	return r.Price, nil
}
