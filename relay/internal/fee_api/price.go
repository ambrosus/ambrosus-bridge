package fee_api

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/price"
	"github.com/ethereum/go-ethereum/common"
)

func (p *FeeAPI) getTokenPrice(bridge BridgeFeeApi, tokenAddress common.Address) (float64, error) {
	tokenPriceI, err, _ := p.cache.Memoize(bridge.GetName()+tokenAddress.Hex(), func() (interface{}, error) {
		return tokenPrice(bridge, tokenAddress)
	})
	if err != nil {
		return 0, err
	}
	return tokenPriceI.(float64), nil
}

func tokenPrice(bridge BridgeFeeApi, tokenAddress common.Address) (float64, error) {
	if (tokenAddress == common.Address{}) {
		tokenAddress = bridge.GetWrapperAddress()
	}

	tokenContract, err := bindings.NewToken(tokenAddress, bridge.GetClient())
	if err != nil {
		return 0, fmt.Errorf("get token contract: %w", err)
	}
	tokenSymbol, err := tokenContract.Symbol(nil)
	if err != nil {
		return 0, fmt.Errorf("get token symbol: %w", err)
	}
	tokenDecimals, err := tokenContract.Decimals(nil)
	if err != nil {
		return 0, fmt.Errorf("get token decimals: %w", err)
	}
	tokenInfo := price.TokenInfo{Symbol: tokenSymbol, Decimals: tokenDecimals, Address: tokenAddress}
	return price.TokenToUSD(&tokenInfo)
}
