package price

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	Amb = coin("AMBUSDT")
	Eth = coin("ETHUSDT")
	Bnb = coin("BNBUSDT")
)

type coin string

func CoinToUsdt(symbol coin) (float64, error) {
	if symbol == Amb {
		return getAmbPrice()
	}

	return getBinancePrice(symbol)
}

type binanceResponse struct {
	Price float64 `json:"price,string"`
}
type ambResponse struct {
	Data ambResponse_ `json:"data"`
}
type ambResponse_ struct {
	Price float64 `json:"total_price_usd"`
}

func getAmbPrice() (float64, error) {
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

func getBinancePrice(symbol coin) (float64, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r binanceResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return 0, err
	}

	return r.Price, nil
}
