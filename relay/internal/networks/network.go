package networks

import (
	"errors"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrEventNotFound = errors.New("error event not found")
)

type Bridge interface {
	// GetLastEventId used by the other side of the bridge for synchronization
	GetLastEventId() (*big.Int, error)
	GetMinSafetyBlocksNum() (uint64, error)
	SendEvent(event *contracts.BridgeTransfer) error
}

type BridgeReceiveAura interface {
	Bridge
	SubmitTransferAura(*contracts.CheckAuraAuraProof) error
	GetValidatorSet() ([]common.Address, error)
	GetLastProcessedBlockHash() (*common.Hash, error)
}

type BridgeReceiveEthash interface {
	Bridge
	SubmitTransferPoW(*contracts.CheckPoWPoWProof) error
	SubmitEpochData(*ethash.EpochData) error
	IsEpochSet(epoch uint64) (bool, error)
}
