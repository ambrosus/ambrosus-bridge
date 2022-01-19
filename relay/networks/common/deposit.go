package common

import (
	"math/big"
	"relay/contracts"
)

type DepositEvent contracts.EthBridgeWithdraw

type PoWBlock struct {
	P1                    []byte
	PrevHashOrReceiptRoot []byte
	P2                    []byte
	Nonce                 []byte
	P3                    []byte
}

type Deposit struct {
	Network string
	EventId *big.Int

	Blocks        []PoWBlock
	Events        []DepositEvent
	ReceiptsProof [][]byte
}
