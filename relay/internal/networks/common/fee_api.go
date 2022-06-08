package common

import (
	"math/big"

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
	return b.DefaultTransferFeeWei
}
