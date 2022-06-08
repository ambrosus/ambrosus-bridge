package price_0x

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
)

const (
	EthUrl = NetworkUrl("https://api.0x.org/swap/v1/price?sellToken=%s&buyToken=USDT&sellAmount=%d")
	BscUrl = NetworkUrl("https://bsc.api.0x.org/swap/v1/price?sellToken=%s&buyToken=BUSD&sellAmount=%d")
)

type NetworkUrl string

type response struct {
	Price float64 `json:"price,string"`

	Reason string `json:"reason"` // when error occurred
}

// CoinToUSD return usd price for smallest token part (wei 1e-18 / satoshi 1e-9)
func CoinToUSD(networkUrl NetworkUrl, symbol string, decimals uint8) (float64, error) {
	amount := math.Pow10(int(decimals))
	price, err := doRequest(networkUrl, symbol, uint(amount))
	if err != nil {
		return 0, err
	}
	return price / amount, nil
}

func doRequest(urlFormat NetworkUrl, sellToken string, amount uint) (float64, error) {
	url := fmt.Sprintf(string(urlFormat), sellToken, amount)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}

	if r.Reason != "" {
		return 0, errors.New(r.Reason)
	}

	return r.Price, nil
}
