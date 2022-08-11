package fee_helper

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type FeeHelper struct {
	networks.Bridge

	minBridgeFee           decimal.Decimal
	sideDefaultTransferFee *big.Int

	wrapperAddress common.Address
	privateKey     *ecdsa.PrivateKey

	transferFeeTracker *transferFeeTracker
}

func NewFeeHelper(bridge, sideBridge networks.Bridge, cfg config.FeeApiNetwork, sideCfg config.FeeApiNetwork) (*FeeHelper, error) {
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

	sideDefaultTransferFee, ok := new(big.Int).SetString(sideCfg.DefaultTransferFee, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse sideDefaultTransferFee (%s)", cfg.DefaultTransferFee)
	}

	return &FeeHelper{
		Bridge:                 bridge,
		minBridgeFee:           decimal.NewFromFloat(cfg.MinBridgeFee),
		sideDefaultTransferFee: sideDefaultTransferFee,
		privateKey:             privateKey,
		wrapperAddress:         wrapperAddress,
		transferFeeTracker:     transferFee,
	}, nil
}

func (b *FeeHelper) Sign(digestHash []byte) ([]byte, error) {
	return crypto.Sign(digestHash, b.privateKey)
}

func (b *FeeHelper) GetTransferFee(thisCoinPrice, sideCoinPrice decimal.Decimal) *big.Int {
	return b.transferFeeTracker.GasPerWithdraw(thisCoinPrice, sideCoinPrice)
}

func (b *FeeHelper) GetWrapperAddress() common.Address {
	return b.wrapperAddress
}

func (b *FeeHelper) GetMinBridgeFee() decimal.Decimal {
	return b.minBridgeFee
}

func (b *FeeHelper) GetDefaultTransferFee() *big.Int {
	return b.sideDefaultTransferFee
}
