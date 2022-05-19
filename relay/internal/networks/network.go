package networks

import (
	"errors"
	"math/big"

	"github.com/ambrosus/ambrosus-bridge/relay/internal/contracts"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethash"
	"github.com/ambrosus/ambrosus-bridge/relay/pkg/ethclients"
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
	Run()
	ValidityWatchdog()

	// GetLastReceivedEventId used by the other side of the bridge for synchronization
	GetLastReceivedEventId() (*big.Int, error)
	GetMinSafetyBlocksNum() (uint64, error)
	GetEventById(eventId *big.Int) (*contracts.BridgeTransfer, error)

	SendEvent(event *contracts.BridgeTransfer, safetyBlocks uint64) error

	// GetTxErr returns error of the transaction
	GetTxErr(params GetTxErrParams) error

	EnsureContractUnpaused()
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

type BridgeReceivePoSA interface {
	Bridge
	SubmitTransferPoSA(proof *contracts.CheckPoSAPoSAProof) error
}

type BridgeFeeApi interface {
	Bridge
	GetClient() ethclients.ClientInterface
	Sign(digestHash []byte) ([]byte, error)
	GetTransferFee() (*big.Int, error)
	CoinPrice() (float64, error) // CoinPrice return that net native coin price in USDT

	// WithdrawCount return count of `Withdraw` events emitted after given event
	WithdrawCount(afterBlockNumber, endBlockNumber uint64) (int, error)
	GetLatestBlockNumber() (uint64, error)
}
