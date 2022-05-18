package price_0x

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type response struct {
	Price float64 `json:"price,string"`
}

func CoinToUSDT(symbol string, decimals uint8) (float64, error) {
	amount := int(math.Pow10(int(decimals)))
	url := fmt.Sprintf("https://api.0x.org/swap/v1/price?sellToken=%s&buyToken=USDT&sellAmount=%d", symbol, amount)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}

	return r.Price, nil
}
