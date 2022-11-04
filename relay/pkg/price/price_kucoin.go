package price

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

const kucoinPriceUrlFormat = "https://api.kucoin.com/api/v1/prices?base=USD&currencies=%s"

type kucoinResponse struct {
	Data map[string]string `json:"data"`
}

func GetKucoin(token *TokenInfo) (price float64, err error) {
	resp, err := http.Get(fmt.Sprintf(kucoinPriceUrlFormat, token.Symbol))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r kucoinResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return 0, err
	}

	priceStr, ok := r.Data[token.Symbol]
	if !ok {
		return 0, fmt.Errorf("response doesn't have key with token symbol")
	}
	price, err = strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0, err
	}

	amount := math.Pow10(int(token.Decimals))
	return price / amount, nil
}
