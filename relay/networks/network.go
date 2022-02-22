package networks

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ethereum/go-ethereum/common"
)

type Bridge interface {
	// GetLastEventId used by the other side of the bridge for synchronization
	GetLastEventId() (*big.Int, error)

	GetValidatorSet() ([]common.Address, error)

	SubmitTransfer(contracts.TransferProof) error
}
