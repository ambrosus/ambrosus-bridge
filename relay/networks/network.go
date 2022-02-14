package networks

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
)

type Bridge interface {
	// GetLastEventId used by the other side of the bridge for synchronization
	GetLastEventId() (*big.Int, error)

	SubmitTransfer(contracts.TransferProof) error
}
