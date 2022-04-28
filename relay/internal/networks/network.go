package networks

import (
	"errors"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var (
	ErrEventNotFound = errors.New("error event not found")
)

type GetTxErrParams struct {
	Tx    *types.Transaction
	TxErr error
	// MethodName amd TxParams used for getting the error in parity/openethereum and for logging
	MethodName string
	TxParams   []interface{}
}

type Bridge interface {
	// GetLastEventId used by the other side of the bridge for synchronization
	GetLastEventId() (*big.Int, error)
	GetMinSafetyBlocksNum(opts *bind.CallOpts) (uint64, error)
	GetEventById(eventId *big.Int) (*contracts.BridgeTransfer, error)

	SendEvent(event *contracts.BridgeTransfer) error

	// GetTxErr returns error of the transaction
	GetTxErr(params GetTxErrParams) error
}

type BridgeReceiveAura interface {
	Bridge
	SubmitTransferAura(*contracts.CheckAuraAuraProof) error
	GetValidatorSet() ([]common.Address, error)
}

type BridgeReceiveEthash interface {
	Bridge
	SubmitTransferPoW(*contracts.CheckPoWPoWProof) error
	SubmitEpochData(*ethash.EpochData) error
	IsEpochSet(epoch uint64) (bool, error)
}
