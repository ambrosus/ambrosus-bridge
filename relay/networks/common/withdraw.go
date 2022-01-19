package common

import (
	"math/big"
	"relay/contracts"
)

//type BlockInput contracts.EthBridgeBlock
type WithdrawEvent contracts.AmbBridgeWithdraw

type AmbBlock struct {
	P1                    []byte
	PrevHashOrReceiptRoot []byte
	P2                    []byte
	Timestamp             []byte
	P3                    []byte
	Seal                  []byte
	Signature             []byte
}

type Withdraw struct {
	Network string
	EventId *big.Int

	Blocks        []*AmbBlock
	Events        []*WithdrawEvent
	ReceiptsProof [][]byte
}
