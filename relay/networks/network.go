package networks

import (
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/contracts"
	"github.com/ethereum/go-ethereum/common"
)

type Bridge interface {
	// GetLastEventId used by the other side of the bridge for synchronization
	GetLastEventId() (*big.Int, error)
}

type BridgeReceiveAura interface {
	Bridge
	SubmitTransferAura(*contracts.CheckAuraAuraProof) error
	GetValidatorSet() ([]common.Address, error)
}

type BridgeReceiveEthash interface {
	Bridge
	SubmitTransferPoW(*contracts.CheckPoWPoWProof) error
	SubmitEpochData() error // todo
}
