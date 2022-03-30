package networks

import (
	"errors"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrEventNotFound = errors.New("error event not found")
	ErrEpochData     = errors.New("Required epoch data-pow not set")
)

type Bridge interface {
	// GetLastEventId used by the other side of the bridge for synchronization
	GetLastEventId() (*big.Int, error)
	GetMinSafetyBlocksNum() (uint64, error)
}

type BridgeReceiveAura interface {
	Bridge
	SubmitTransferAura(*contracts.CheckAuraAuraProof) error
	GetValidatorSet() ([]common.Address, error)
	GetLastProcessedBlockNum() (*big.Int, error)
}

type BridgeReceiveEthash interface {
	Bridge
	SubmitTransferPoW(*contracts.CheckPoWPoWProof) error
	SubmitEpochData(*big.Int, *big.Int, *big.Int, []*big.Int, *big.Int, *big.Int) error
}
