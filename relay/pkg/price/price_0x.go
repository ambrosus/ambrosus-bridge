package price

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
)

const (
	EthUrl = NetworkUrl("https://api.0x.org/swap/v1/price?sellToken=%s&buyToken=USDT&buyAmount=1000000")                 // USDT has 6 decimals
	BscUrl = NetworkUrl("https://bsc.api.0x.org/swap/v1/price?sellToken=%s&buyToken=BUSD&buyAmount=1000000000000000000") // BUSD has 18 decimals
)

var NetworkUrls = []NetworkUrl{EthUrl, BscUrl}
var ErrValidationFailed = errors.New("Validation Failed")

type NetworkUrl string
type response struct {
	Price float64 `json:"price,string"`

	Reason string `json:"reason"` // when error occurred
}

// Get0x returns usd price for smallest token part (wei 1e-18 / satoshi 1e-9)
func Get0x(token *TokenInfo) (price float64, err error) {
	amount := math.Pow10(int(token.Decimals))

	for _, url := range NetworkUrls {
		price, err = doRequest(url, token.Symbol)

		// when token not found - try next url
		if err != nil && err.Error() == ErrValidationFailed.Error() {
			continue
		}

		if err != nil {
			return 0, err
		}

		break
	}

	return 1 / price / amount, err
}

func doRequest(urlFormat NetworkUrl, sellToken string) (float64, error) {
	apiKey := os.Getenv("0X_API_KEY")
	client := http.Client{}

	url := fmt.Sprintf(string(urlFormat), sellToken)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("0x-api-key", apiKey)

	resp, err := client.Do(req)
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
