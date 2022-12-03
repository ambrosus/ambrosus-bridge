package networks

import (
	"errors"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/bindings/interfaces"
	"github.com/ambrosus/ambrosus-bridge/relay/internal/networks/events"
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
	Events() events.Events

	GetLogger() *zerolog.Logger
	GetName() string
	GetAuth() *bind.TransactOpts
	GetContractAddress() common.Address
	GetRelayAddress() common.Address

	ProcessTx(methodName string, txOpts *bind.TransactOpts, txCallback ContractCallFn) error
}
