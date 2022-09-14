package fee

import (
	"fmt"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/price"
	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
)

type priceGetter interface {
	tokenPrice(bridge BridgeFeeApi, tokenAddress common.Address) (decimal.Decimal, error)
}

type priceGetterS struct{}

func (p *Fee) getPrices(bridge, sideBridge BridgeFeeApi, tokenAddress common.Address) (thisCoinPrice, sideCoinPrice, tokenUsdPrice decimal.Decimal, err error) {
	thisCoinPrice, err = p.getTokenPrice(bridge, common.Address{})
	if err != nil {
		err = fmt.Errorf("getTokenPrice native bridge: %w", err)
		return
	}
	sideCoinPrice, err = p.getTokenPrice(sideBridge, common.Address{})
	if err != nil {
		err = fmt.Errorf("getTokenPrice native sideBridge: %w", err)
		return
	}

	if bridge == p.amb { // if bridge is amb, then we must use token address from sideBridge
		tokenAddress, err = bridge.GetContract().TokenAddresses(nil, tokenAddress)
		if err != nil {
			err = fmt.Errorf("TokenAddresses: %w", err)
			return
		}
		bridge = p.side
	}

	tokenUsdPrice, err = p.getTokenPrice(bridge, tokenAddress)
	if err != nil {
		err = fmt.Errorf("getTokenPrice %v: %w", tokenAddress, err)
		return
	}

	return
}

func (p *Fee) getTokenPrice(bridge BridgeFeeApi, tokenAddress common.Address) (decimal.Decimal, error) {
	tokenPriceI, err, _ := p.cache.Memoize(bridge.GetName()+tokenAddress.Hex(), func() (interface{}, error) {
		return p.priceGetter.tokenPrice(bridge, tokenAddress)
	})
	if err != nil {
		return decimal.Decimal{}, err
	}
	return tokenPriceI.(decimal.Decimal), nil
}

func (*priceGetterS) tokenPrice(bridge BridgeFeeApi, tokenAddress common.Address) (decimal.Decimal, error) {
	if (tokenAddress == common.Address{}) {
		tokenAddress = bridge.GetWrapperAddress()
	}

	tokenContract, err := bindings.NewToken(tokenAddress, bridge.GetClient())
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("get token contract: %w", err)
	}
	tokenSymbol, err := tokenContract.Symbol(nil)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("get token symbol: %w", err)
	}
	tokenDecimals, err := tokenContract.Decimals(nil)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("get token decimals: %w", err)
	}

	tokenInfo := price.TokenInfo{Symbol: tokenSymbol, Decimals: tokenDecimals, Address: tokenAddress}
	tokenPrice, err := price.TokenToUSD(&tokenInfo)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("get token price: %w", err)
	}
	return decimal.NewFromFloat(tokenPrice), nil
}
