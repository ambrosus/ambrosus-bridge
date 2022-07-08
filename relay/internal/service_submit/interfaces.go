package service_submit

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ethereum/go-ethereum/common"
)

type Submitter interface {
	networks.Bridge
	GetEventById(eventId *big.Int) (*bindings.BridgeTransfer, error)
	SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error
}

type Receiver interface {
	// GetLastReceivedEventId used by the other side of the bridge for synchronization
	GetLastReceivedEventId() (*big.Int, error)
	GetMinSafetyBlocksNum() (uint64, error)
	EnsureContractUnpaused()
}

type BridgeReceiveAura interface {
	networks.Bridge
	SubmitTransferAura(*bindings.CheckAuraAuraProof) error
	GetValidatorSet() ([]common.Address, error)
	GetLastProcessedBlockHash() (*common.Hash, error)
	GetMinSafetyBlocksValidators() (uint64, error)
}

type BridgeReceiveEthash interface {
	networks.Bridge
	SubmitTransferPoW(*bindings.CheckPoWPoWProof) error
	SubmitEpochData(*ethash.EpochData) error
	IsEpochSet(epoch uint64) (bool, error)
}

type BridgeReceivePoSA interface {
	networks.Bridge
	SubmitTransferPoSA(proof *bindings.CheckPoSAPoSAProof) error
	GetCurrentEpoch() (uint64, error)
}
