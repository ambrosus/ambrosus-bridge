package fee_helper

import (
	"crypto/ecdsa"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/config"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/helpers"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

type FeeHelper struct {
	networks.Bridge

	minBridgeFee   decimal.Decimal // in usd
	minTransferFee decimal.Decimal // in usd

	wrapperAddress common.Address
	privateKey     *ecdsa.PrivateKey

	transferFeeTracker *transferFeeTracker
}

func NewFeeHelper(bridge, sideBridge networks.Bridge, explorer, sideExplorer explorerClient, cfg config.FeeApiNetwork, sideCfg config.FeeApiNetwork) (*FeeHelper, error) {
	wrapperAddress, err := bridge.GetContract().WrapperAddress(nil)
	if err != nil {
		return nil, err
	}

	transferFee, err := newTransferFeeTracker(
		bridge, sideBridge,
		explorer, sideExplorer,
		cfg.TransferFeeIncludedTxsFromAddresses, sideCfg.TransferFeeIncludedTxsFromAddresses,
		cfg.TransferFeeTxsFromBlock, sideCfg.TransferFeeTxsFromBlock,
	)
	if err != nil {
		return nil, err
	}

	privateKey, err := helpers.ParsePK(cfg.PrivateKey)
	if err != nil {
		return nil, err
	}

	return &FeeHelper{
		Bridge:             bridge,
		minBridgeFee:       decimal.NewFromFloat(cfg.MinBridgeFee),
		minTransferFee:     decimal.NewFromFloat(cfg.MinTransferFee),
		privateKey:         privateKey,
		wrapperAddress:     wrapperAddress,
		transferFeeTracker: transferFee,
	}, nil
}

func (b *FeeHelper) Sign(digestHash []byte) ([]byte, error) {
	return crypto.Sign(digestHash, b.privateKey)
}

func (b *FeeHelper) GetTransferFee() (thisGas, sideGas decimal.Decimal) {
	this, side := b.transferFeeTracker.GasPerWithdraw()
	return decimal.NewFromBigInt(this, 0), decimal.NewFromBigInt(side, 0)
}

func (b *FeeHelper) GetWrapperAddress() common.Address {
	return b.wrapperAddress
}

func (b *FeeHelper) GetMinBridgeFee() decimal.Decimal {
	return b.minBridgeFee
}

func (b *FeeHelper) GetMinTransferFee() decimal.Decimal {
	return b.minTransferFee
}
