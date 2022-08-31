package fee_helper

import (
	"crypto/ecdsa"
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
	sideDefaultTransferFee decimal.Decimal

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

	return &FeeHelper{
		Bridge:                 bridge,
		minBridgeFee:           decimal.NewFromFloat(cfg.MinBridgeFee),
		sideDefaultTransferFee: decimal.NewFromFloat(sideCfg.DefaultTransferFee),
		privateKey:             privateKey,
		wrapperAddress:         wrapperAddress,
		transferFeeTracker:     transferFee,
	}, nil
}

func (b *FeeHelper) Sign(digestHash []byte) ([]byte, error) {
	return crypto.Sign(digestHash, b.privateKey)
}

func (b *FeeHelper) GetTransferFee() *big.Int {
	return b.transferFeeTracker.GasPerWithdraw()
}

func (b *FeeHelper) GetWrapperAddress() common.Address {
	return b.wrapperAddress
}

func (b *FeeHelper) GetMinBridgeFee() decimal.Decimal {
	return b.minBridgeFee
}

func (b *FeeHelper) GetDefaultTransferFee() decimal.Decimal {
	return b.sideDefaultTransferFee
}
