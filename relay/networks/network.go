package networks

import "relay/contracts"

type Bridge interface {
	GetLastEventId() uint64
	GetEventById() common.WithdrawEvent
}

type SubmitPoWF func(eventId uint, blocks contracts.CheckPoWBlockPoW, events contracts.CommonStructsTransfer, proof contracts.ReceiptsProof)
type SubmitPoAF func(eventId uint, blocks contracts.CheckPoABlockPoA, events contracts.CommonStructsTransfer, proof contracts.ReceiptsProof)
