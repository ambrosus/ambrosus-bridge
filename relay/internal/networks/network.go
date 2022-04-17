package networks

import (
	"errors"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrEventNotFound = errors.New("error event not found")
)

type GetTransactionErrorParams struct {
	Tx         *types.Transaction
	TxErr      error
	MethodName string
}

type Bridge interface {
	// GetLastEventId used by the other side of the bridge for synchronization
	GetLastEventId() (*big.Int, error)
	GetMinSafetyBlocksNum() (uint64, error)

	// GetTransactionError returns error of the transaction
	// "methodName" amd "params" are optional and used only
	// for getting the error in parity/openethereum
	GetTransactionError(params GetTransactionErrorParams, txParams ...interface{}) error
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
