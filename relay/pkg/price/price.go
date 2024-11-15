package price

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type TokenInfo struct {
	Symbol   string
	Decimals uint8
	Address  common.Address
}

// TokenToUSD return usd price for smallest token part (wei 1e-18 / satoshi 1e-9)
func TokenToUSD(token *TokenInfo) (price float64, err error) {
	if token.Symbol == "SAMB" || token.Symbol == "AMB" {
		price, err = GetAmb()
	} else if token.Symbol == "USDT" || token.Symbol == "BUSD" {
		price, err = GetKucoin(token)
	} else if token.Symbol == "WBNB" {
		price, err = GetWBNB(token)
	} else {
		price, err = Get0x(token)
	}

	return price, err
}

func GetWBNB(token *TokenInfo) (price float64, err error) {
	amount := math.Pow10(int(token.Decimals))

	client := http.Client{}

	req, err := http.NewRequest("GET", "https://api.binance.com/api/v1/ticker/price?symbol=BNBUSDT", nil)
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
