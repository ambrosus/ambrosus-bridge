package price_0x

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
)

const (
	EthUrl = NetworkUrl("https://api.0x.org/swap/v1/price?sellToken=%s&buyToken=%s&sellAmount=%d")
	BscUrl = NetworkUrl("https://bsc.api.0x.org/swap/v1/price?sellToken=%s&buyToken=%s&sellAmount=%d")
)

type NetworkUrl string

type response struct {
	Price float64 `json:"price,string"`

	Reason string `json:"reason"` // when error occurred
}

func CoinToUSDT(networkUrl NetworkUrl, symbol string, decimals uint8) (float64, error) {
	r, err := doRequest(networkUrl, symbol, "USDT", decimals)
	if err != nil {
		return 0, err
	}

	return r.Price, nil
}

func CoinToBUSD(networkUrl NetworkUrl, symbol string, decimals uint8) (float64, error) {
	r, err := doRequest(networkUrl, symbol, "BUSD", decimals)
	if err != nil {
		return 0, err
	}

	return r.Price, nil
}

func doRequest(urlFormat NetworkUrl, sellToken, buyToken string, decimals uint8) (*response, error) {
	amount := int(math.Pow10(int(decimals)))
	url := fmt.Sprintf(string(urlFormat), sellToken, buyToken, amount)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Reason != "" {
		return nil, errors.New(r.Reason)
	}

	return &r, nil
}
