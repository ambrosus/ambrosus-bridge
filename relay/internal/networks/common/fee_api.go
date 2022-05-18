package common

import (
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
	// res, err := b.GasPerWithdraw(1)
	// if err != nil {
	// 	return nil, err
	// }
	// return big.NewInt(res), nil

	return big.NewInt(228), nil
}
