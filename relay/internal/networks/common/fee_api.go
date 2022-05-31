package common

import (
	"context"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/fee_api"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func (b *CommonBridge) GetClient() ethclients.ClientInterface {
	return b.Client
}

func (b *CommonBridge) Sign(digestHash []byte) ([]byte, error) {
	return crypto.Sign(digestHash, b.Pk)
}

func (b *CommonBridge) GetTransferFee(thisCoinPrice, sideCoinPrice float64) (*big.Int, error) {
	// get gas cost per withdraw in side bridge currency
	gasCostInSide, err := b.GasPerWithdraw(&b.PriceTrackerData)
	if err != nil {
		return nil, err
	}
	if gasCostInSide == nil {
		return b.DefaultTransferFeeWei, nil
	}

	// convert it to native bridge currency
	gasCostInUsd := fee_api.Coin2Usd(gasCostInSide, sideCoinPrice, 18)
	gasCostInNative := fee_api.Usd2Coin(gasCostInUsd, thisCoinPrice, 18)
	return gasCostInNative, nil
}

func (b *CommonBridge) GetLatestBlockNumber() (uint64, error) {
	return b.Client.BlockNumber(context.Background())
}

func (b *CommonBridge) GetMinBridgeFee() *big.Float {
	return b.MinBridgeFee
}

func (b *CommonBridge) GetWrapperAddress() (common.Address, error) {
	return b.Contract.WrapperAddress(nil)
}
