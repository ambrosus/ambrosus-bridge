package networks

import (
	"math/big"
	"relay/contracts"
)

type Bridge interface {
	GetLastEventId() (*big.Int, error)
	// GetEventById() common.WithdrawEvent
}

type SubmitPoWF func(eventId *big.Int, blocks []*contracts.CheckPoABlockPoA, events *[]contracts.CommonStructsTransfer, proof *contracts.ReceiptsProof)
type SubmitPoAF func(eventId *big.Int, blocks []*contracts.CheckPoABlockPoA, events *[]contracts.CommonStructsTransfer, proof *contracts.ReceiptsProof)
