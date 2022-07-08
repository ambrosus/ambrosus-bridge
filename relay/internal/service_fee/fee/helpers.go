package fee

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shopspring/decimal"
)

func buildMessage(tokenAddress common.Address, transferFee, bridgeFee, amount *big.Int) []byte {
	var data bytes.Buffer
	var b32 [32]byte // solidity fills uint256 with zero

	data.Write(tokenAddress.Bytes())

	timestamp := time.Now().Unix() / signatureFeeTimestamp
	data.Write(big.NewInt(timestamp).FillBytes(b32[:]))

	data.Write(transferFee.FillBytes(b32[:]))
	data.Write(bridgeFee.FillBytes(b32[:]))
	data.Write(amount.FillBytes(b32[:]))

	return accounts.TextHash(crypto.Keccak256(data.Bytes()))
}

// amountUsd / priceUsd
func usd2Coin(amountUsd decimal.Decimal, priceUsd decimal.Decimal) decimal.Decimal {
	return amountUsd.Div(priceUsd)
}

// amountWei * priceUsd
func coin2Usd(amountWei decimal.Decimal, priceUsd decimal.Decimal) decimal.Decimal {
	return amountWei.Mul(priceUsd)
}

// amount * bps / 10_000
func calcBps(amount decimal.Decimal, bps int64) decimal.Decimal {
	return amount.Mul(decimal.NewFromInt(bps)).Div(decimal.NewFromInt(10_000))
}
