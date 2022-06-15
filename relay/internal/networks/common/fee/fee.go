package fee

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

type BridgeFee struct {
	networks.Bridge

	minBridgeFee       decimal.Decimal
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
		Bridge:             bridge,
		minBridgeFee:       decimal.NewFromFloat(cfg.MinBridgeFee),
		defaultTransferFee: big.NewInt(cfg.DefaultTransferFee),
		privateKey:         privateKey,
		wrapperAddress:     wrapperAddress,
		transferFeeTracker: transferFee,
	}, nil
}

func (b *BridgeFee) Sign(digestHash []byte) ([]byte, error) {
	return crypto.Sign(digestHash, b.privateKey)
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

func (b *BridgeFee) GetMinBridgeFee() decimal.Decimal {
	return b.minBridgeFee
}
