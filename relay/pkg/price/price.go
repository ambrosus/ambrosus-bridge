package price

import (
	"github.com/ethereum/go-ethereum/common"
)

type TokenInfo struct {
	Symbol   string
	Decimals uint8
	Address  common.Address
}

// TokenToUSD return usd price for smallest token part (wei 1e-18 / satoshi 1e-9)
func TokenToUSD(token *TokenInfo) (price float64, err error) {
	if token.Symbol == "SAMB" {
		price, err = GetAmb()
	} else if token.Symbol == "USDT" || token.Symbol == "BUSD" {
		price, err = GetKucoin(token)
	} else {
		price, err = Get0x(token)
	}

	return price, err
}
