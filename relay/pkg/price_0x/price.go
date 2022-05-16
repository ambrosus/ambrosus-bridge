package price_0x

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type response struct {
	Price float64 `json:"price,string"`
}

func CoinToUSDT(symbol string, amount float64) (float64, error) {
	url := fmt.Sprintf(
		"https://api.0x.org/swap/v1/price?sellToken=%s&buyToken=USDT&sellAmount=%s",
		symbol,
		strconv.FormatFloat(amount, 'f', -1, 64),
	)
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
