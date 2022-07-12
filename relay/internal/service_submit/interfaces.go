package service_submit

import (
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ethereum/go-ethereum/common"
)

type Submitter interface {
	networks.Bridge
	SendEvent(event *bindings.BridgeTransfer, safetyBlocks uint64) error
}

type Receiver interface {
	networks.Bridge
}

type ReceiverAura interface {
	Receiver
	SubmitTransferAura(*bindings.CheckAuraAuraProof) error
	SubmitValidatorSetChangesAura(*bindings.CheckAuraAuraProof) error
	GetValidatorSet() ([]common.Address, error)
	GetLastProcessedBlockHash() (*common.Hash, error)
	GetMinSafetyBlocksValidators() (uint64, error)
}

type ReceiverPoW interface {
	Receiver
	SubmitTransferPoW(*bindings.CheckPoWPoWProof) error
	SubmitEpochData(*ethash.EpochData) error
	IsEpochSet(epoch uint64) (bool, error)
}

type ReceiverPoSA interface {
	Receiver
	SubmitTransferPoSA(proof *bindings.CheckPoSAPoSAProof) error
	GetCurrentEpoch() (uint64, error)
}
