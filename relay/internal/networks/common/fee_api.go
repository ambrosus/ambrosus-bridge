package common

import (
	"context"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ethereum/go-ethereum/crypto"
)

func (b *CommonBridge) GetClient() ethclients.ClientInterface {
	return b.Client
}

func (b *CommonBridge) Sign(digestHash []byte) ([]byte, error) {
	return crypto.Sign(digestHash, b.Pk)
}

func (b *CommonBridge) GetTransferFee() (*big.Int, error) {
	res, err := b.GasPerWithdraw(&b.PriceTrackerData)
	if err != nil {
		return nil, err
	}
	return big.NewInt(int64(res)), nil
}

func (b *CommonBridge) GetLatestBlockNumber() (uint64, error) {
	return b.Client.BlockNumber(context.Background())
}

func (b *CommonBridge) GetMinBridgeFee() *big.Float {
	return b.MinBridgeFee
}
