package fee

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type BridgeFee struct {
	bridge networks.Bridge

	minBridgeFee       *big.Float
	defaultTransferFee *big.Int

	wrapperAddress common.Address
	privateKey     *ecdsa.PrivateKey

	transferFeeTracker *transferFeeTracker
}

func NewBridgeFee(bridge, sideBridge networks.Bridge, cfg config.FeeApiNetwork) (*BridgeFee, error) {
	wrapperAddress, err := bridge.GetContract().WrapperAddress(nil)
	if err != nil {
		return nil, err
	}

	transferFee, err := newTransferFeeTracker(bridge, sideBridge)
	if err != nil {
		return nil, err
	}

	privateKey, err := helpers.ParsePK(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &BridgeFee{
		bridge:             bridge,
		minBridgeFee:       big.NewFloat(cfg.MinBridgeFee),
		defaultTransferFee: big.NewInt(cfg.DefaultTransferFee),
		privateKey:         privateKey,
		wrapperAddress:     wrapperAddress,
		transferFeeTracker: transferFee,
	}, nil
}

func (b *BridgeFee) Sign(digestHash []byte) ([]byte, error) {
	return crypto.Sign(digestHash, b.privateKey)
}

func (b *BridgeFee) TokenPrice(tokenAddress common.Address) (float64, error) {
	panic("abstract method")
}

func (b *BridgeFee) GetTransferFee() *big.Int {
	feeSideNative := b.transferFeeTracker.GasPerWithdraw()
	if feeSideNative == nil {
		feeSideNative = b.defaultTransferFee
	}
	return feeSideNative
}

func (b *BridgeFee) GetWrapperAddress() common.Address {
	return b.wrapperAddress
}

func (b *BridgeFee) GetMinBridgeFee() *big.Float {
	return b.minBridgeFee
}

// GetTokenData used in TokenPrice method
func (b *BridgeFee) GetTokenData(tokenAddress common.Address) (string, uint8, error) {
	if (tokenAddress == common.Address{}) {
		tokenAddress = b.wrapperAddress
	}

	tokenContract, err := contracts.NewToken(tokenAddress, b.bridge.GetClient())
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
