package networks

import (
	"errors"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

//go:generate mockgen -source=network.go -destination=mocks/mock.go

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

// ContractCallFn is a callback type for calling paid contract's method.
type ContractCallFn func(opts *bind.TransactOpts) (*types.Transaction, error)

type Bridge interface {
	GetClient() ethclients.ClientInterface
	GetWsClient() ethclients.ClientInterface

	GetContract() interfaces.BridgeContract
	GetWsContract() interfaces.BridgeContract
	GetLogger() *zerolog.Logger
	GetName() string

	ProcessTx(methodName string, txCallback ContractCallFn) error
	ShouldHavePk()
	EnsureContractUnpaused()
}

type Submitter interface {
	Bridge
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
	Bridge
	SubmitTransferAura(*bindings.CheckAuraAuraProof) error
	GetValidatorSet() ([]common.Address, error)
	GetLastProcessedBlockHash() (*common.Hash, error)
	GetMinSafetyBlocksValidators() (uint64, error)
}

type BridgeReceiveEthash interface {
	Bridge
	SubmitTransferPoW(*bindings.CheckPoWPoWProof) error
	SubmitEpochData(*ethash.EpochData) error
	IsEpochSet(epoch uint64) (bool, error)
}

type BridgeReceivePoSA interface {
	Bridge
	SubmitTransferPoSA(proof *bindings.CheckPoSAPoSAProof) error
	GetCurrentEpoch() (uint64, error)
}
