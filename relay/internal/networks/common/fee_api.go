package common

import (
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (b *CommonBridge) Sign(digestHash []byte) ([]byte, error) {
	return crypto.Sign(digestHash, b.Pk)
}

func (b *CommonBridge) GetMinBridgeFee() *big.Float {
	return b.MinBridgeFee
}

func (b *CommonBridge) GetWrapperAddress() (common.Address, error) {
	return b.Contract.WrapperAddress(nil)
}

func (b *CommonBridge) GetDefaultTransferFeeWei() *big.Int {
	// todo: immutable, can be cached
	return b.DefaultTransferFeeWei
}

func (b *CommonBridge) GetTokenData(tokenAddress common.Address) (string, uint8, error) {
	if (tokenAddress == common.Address{}) {
		var err error
		tokenAddress, err = b.GetWrapperAddress()
		if err != nil {
			return "", 0, fmt.Errorf("GetWrapperAddress error %w", err)
		}
	}

	tokenContract, err := contracts.NewToken(tokenAddress, b.GetClient())
	if err != nil {
		return "", 0, fmt.Errorf("get token contract: %w", err)
	}
	tokenSymbol, err := tokenContract.Symbol(nil)
	if err != nil {
		return "", 0, fmt.Errorf("get token symbol: %w", err)
	}
	tokenDecimals, err := tokenContract.Decimals(nil)
	if err != nil {
		return "", 0, fmt.Errorf("get token decimals: %w", err)
	}
	return tokenSymbol, tokenDecimals, nil
}
