package common

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const SignatureFeeTimestamp = 30 * 60

func (b *CommonBridge) Sign(tokenAddress string, transferFee *big.Int, bridgeFee *big.Int) ([]byte, error) {
	var data bytes.Buffer
	var b32 [32]byte // solidity fills uint256 with zero

	data.Write(common.HexToAddress(tokenAddress).Bytes())

	timestamp := time.Now().Unix() / SignatureFeeTimestamp
	data.Write(big.NewInt(timestamp).FillBytes(b32[:]))

	data.Write(transferFee.FillBytes(b32[:]))
	data.Write(bridgeFee.FillBytes(b32[:]))

	return crypto.Sign(accounts.TextHash(crypto.Keccak256(data.Bytes())), b.Pk)
}

func (b *CommonBridge) GetBridgeFee(tokenAddress string) (*big.Int, error) {
	return big.NewInt(1337), nil
}

func (b *CommonBridge) GetTransferFee() (*big.Int, error) {
	// res, err := b.GasPerWithdraw(1)
	// if err != nil {
	// 	return nil, err
	// }
	// return big.NewInt(res), nil

	return big.NewInt(228), nil
}
