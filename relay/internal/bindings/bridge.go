// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// CheckAuraAuraProof is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraAuraProof struct {
	Blocks             []CheckAuraBlockAura
	Transfer           CommonStructsTransferProof
	VsChanges          []CheckAuraValidatorSetProof
	TransferEventBlock uint64
}

// CheckAuraBlockAura is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraBlockAura struct {
	P0Seal      [3]byte
	P0Bare      [3]byte
	ParentHash  [32]byte
	P2          []byte
	ReceiptHash [32]byte
	P3          []byte
	Step        [4]byte
	Signature   []byte
	FinalizedVs uint64
}

// CheckAuraValidatorSetChange is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraValidatorSetChange struct {
	DeltaAddress common.Address
	DeltaIndex   int64
}

// CheckAuraValidatorSetProof is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraValidatorSetProof struct {
	ReceiptProof [][]byte
	Changes      []CheckAuraValidatorSetChange
}

// CheckPoSABlockPoSA is an auto generated low-level Go binding around an user-defined struct.
type CheckPoSABlockPoSA struct {
	P0Signed    [3]byte
	P0Unsigned  [3]byte
	ParentHash  [32]byte
	P1          []byte
	ReceiptHash [32]byte
	P2          []byte
	Number      []byte
	P3          []byte
	P4Signed    []byte
	P4Unsigned  []byte
	ExtraData   []byte
	P5          []byte
}

// CheckPoSAPoSAProof is an auto generated low-level Go binding around an user-defined struct.
type CheckPoSAPoSAProof struct {
	Blocks             []CheckPoSABlockPoSA
	Transfer           CommonStructsTransferProof
	TransferEventBlock uint64
}

// CheckPoWBlockPoW is an auto generated low-level Go binding around an user-defined struct.
type CheckPoWBlockPoW struct {
	P0WithNonce         [3]byte
	P0WithoutNonce      [3]byte
	P1                  []byte
	ParentOrReceiptHash [32]byte
	P2                  []byte
	Difficulty          []byte
	P3                  []byte
	Number              []byte
	P4                  []byte
	P5                  []byte
	Nonce               []byte
	P6                  []byte
	DataSetLookup       []*big.Int
	WitnessForLookup    []*big.Int
}

// CheckPoWPoWProof is an auto generated low-level Go binding around an user-defined struct.
type CheckPoWPoWProof struct {
	Blocks   []CheckPoWBlockPoW
	Transfer CommonStructsTransferProof
}

// CommonStructsConstructorArgs is an auto generated low-level Go binding around an user-defined struct.
type CommonStructsConstructorArgs struct {
	SideBridgeAddress    common.Address
	AdminAddress         common.Address
	RelayAddress         common.Address
	WrappingTokenAddress common.Address
	TokenThisAddresses   []common.Address
	TokenSideAddresses   []common.Address
	TransferFeeRecipient common.Address
	BridgeFeeRecipient   common.Address
	TimeframeSeconds     *big.Int
	LockTime             *big.Int
	MinSafetyBlocks      *big.Int
}

// CommonStructsLockedTransfers is an auto generated low-level Go binding around an user-defined struct.
type CommonStructsLockedTransfers struct {
	Transfers    []CommonStructsTransfer
	EndTimestamp *big.Int
}

// CommonStructsTransfer is an auto generated low-level Go binding around an user-defined struct.
type CommonStructsTransfer struct {
	TokenAddress common.Address
	ToAddress    common.Address
	Amount       *big.Int
}

// CommonStructsTransferProof is an auto generated low-level Go binding around an user-defined struct.
type CommonStructsTransferProof struct {
	ReceiptProof [][]byte
	EventId      *big.Int
	Transfers    []CommonStructsTransfer
}

// BridgeMetaData contains all meta data concerning the Bridge contract.
var BridgeMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"queue\",\"type\":\"tuple[]\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"TransferFinish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"TransferSubmit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenFrom\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenTo\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"transferFeeAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bridgeFeeAmount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RELAY_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"feeRecipient_\",\"type\":\"address\"}],\"name\":\"changeBridgeFeeRecipient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lockTime_\",\"type\":\"uint256\"}],\"name\":\"changeLockTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minSafetyBlocks_\",\"type\":\"uint256\"}],\"name\":\"changeMinSafetyBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"signatureFeeCheckNumber_\",\"type\":\"uint256\"}],\"name\":\"changeSignatureFeeCheckNumber\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timeframeSeconds_\",\"type\":\"uint256\"}],\"name\":\"changeTimeframeSeconds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"feeRecipient_\",\"type\":\"address\"}],\"name\":\"changeTransferFeeRecipient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"getLockedTransfers\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.LockedTransfers\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sideBridgeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"adminAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"relayAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"wrappingTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenSideAddresses\",\"type\":\"address[]\"},{\"internalType\":\"addresspayable\",\"name\":\"transferFeeRecipient\",\"type\":\"address\"},{\"internalType\":\"addresspayable\",\"name\":\"bridgeFeeRecipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"timeframeSeconds\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSafetyBlocks\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.ConstructorArgs\",\"name\":\"args\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"minimumDifficulty\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inputEventId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epochIndex\",\"type\":\"uint256\"}],\"name\":\"isEpochDataSet\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isQueueEmpty\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastTimeframe\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lockedTransfers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minSafetyBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oldestLockedEventId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"removeLockedTransfers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epochNum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fullSizeIn128Resultion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"branchDepth\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"merkleNodes\",\"type\":\"uint256[]\"}],\"name\":\"setEpochData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sideBridgeAddress\",\"type\":\"address\"}],\"name\":\"setSideBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sideBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes3\",\"name\":\"p0WithNonce\",\"type\":\"bytes3\"},{\"internalType\":\"bytes3\",\"name\":\"p0WithoutNonce\",\"type\":\"bytes3\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parentOrReceiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"difficulty\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"nonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p6\",\"type\":\"bytes\"},{\"internalType\":\"uint256[]\",\"name\":\"dataSetLookup\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"witnessForLookup\",\"type\":\"uint256[]\"}],\"internalType\":\"structCheckPoW.BlockPoW[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receiptProof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"}],\"internalType\":\"structCheckPoW.PoWProof\",\"name\":\"powProof\",\"type\":\"tuple\"}],\"name\":\"submitTransferPoW\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeframeSeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenSideAddress\",\"type\":\"address\"}],\"name\":\"tokensAdd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenSideAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensAddBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"}],\"name\":\"tokensRemove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensRemoveBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"triggerTransfers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"}],\"name\":\"unlockTransfers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockTransfersBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"unwrapSide\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"transferFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bridgeFee\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"transferFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bridgeFee\",\"type\":\"uint256\"}],\"name\":\"wrapWithdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"wrapperAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidatorSet\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastProcessedBlock\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes3\",\"name\":\"p0Seal\",\"type\":\"bytes3\"},{\"internalType\":\"bytes3\",\"name\":\"p0Bare\",\"type\":\"bytes3\"},{\"internalType\":\"bytes32\",\"name\":\"parentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"receiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes4\",\"name\":\"step\",\"type\":\"bytes4\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint64\",\"name\":\"finalizedVs\",\"type\":\"uint64\"}],\"internalType\":\"structCheckAura.BlockAura[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receiptProof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receiptProof\",\"type\":\"bytes[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"deltaAddress\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"deltaIndex\",\"type\":\"int64\"}],\"internalType\":\"structCheckAura.ValidatorSetChange[]\",\"name\":\"changes\",\"type\":\"tuple[]\"}],\"internalType\":\"structCheckAura.ValidatorSetProof[]\",\"name\":\"vsChanges\",\"type\":\"tuple[]\"},{\"internalType\":\"uint64\",\"name\":\"transferEventBlock\",\"type\":\"uint64\"}],\"internalType\":\"structCheckAura.AuraProof\",\"name\":\"auraProof\",\"type\":\"tuple\"}],\"name\":\"submitTransferAura\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatorSet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes3\",\"name\":\"p0Signed\",\"type\":\"bytes3\"},{\"internalType\":\"bytes3\",\"name\":\"p0Unsigned\",\"type\":\"bytes3\"},{\"internalType\":\"bytes32\",\"name\":\"parentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"receiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4Signed\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4Unsigned\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"}],\"internalType\":\"structCheckPoSA.BlockPoSA[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receiptProof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"transferEventBlock\",\"type\":\"uint64\"}],\"internalType\":\"structCheckPoSA.PoSAProof\",\"name\":\"posaProof\",\"type\":\"tuple\"}],\"name\":\"submitTransferPoSA\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes3\",\"name\":\"p0Signed\",\"type\":\"bytes3\"},{\"internalType\":\"bytes3\",\"name\":\"p0Unsigned\",\"type\":\"bytes3\"},{\"internalType\":\"bytes32\",\"name\":\"parentHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"receiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4Signed\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4Unsigned\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"}],\"internalType\":\"structCheckPoSA.BlockPoSA[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receiptProof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"eventId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"},{\"internalType\":\"uint64\",\"name\":\"transferEventBlock\",\"type\":\"uint64\"}],\"internalType\":\"structCheckPoSA.PoSAProof\",\"name\":\"posaProof\",\"type\":\"tuple\"}],\"name\":\"submitValidatorSetChangesPoSA\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BridgeABI is the input ABI used to generate the binding from.
// Deprecated: Use BridgeMetaData.ABI instead.
var BridgeABI = BridgeMetaData.ABI

// Bridge is an auto generated Go binding around an Ethereum contract.
type Bridge struct {
	BridgeCaller     // Read-only binding to the contract
	BridgeTransactor // Write-only binding to the contract
	BridgeFilterer   // Log filterer for contract events
}

// BridgeCaller is an auto generated read-only Go binding around an Ethereum contract.
type BridgeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BridgeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BridgeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BridgeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BridgeSession struct {
	Contract     *Bridge           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BridgeCallerSession struct {
	Contract *BridgeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BridgeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BridgeTransactorSession struct {
	Contract     *BridgeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BridgeRaw is an auto generated low-level Go binding around an Ethereum contract.
type BridgeRaw struct {
	Contract *Bridge // Generic contract binding to access the raw methods on
}

// BridgeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BridgeCallerRaw struct {
	Contract *BridgeCaller // Generic read-only contract binding to access the raw methods on
}

// BridgeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BridgeTransactorRaw struct {
	Contract *BridgeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBridge creates a new instance of Bridge, bound to a specific deployed contract.
func NewBridge(address common.Address, backend bind.ContractBackend) (*Bridge, error) {
	contract, err := bindBridge(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bridge{BridgeCaller: BridgeCaller{contract: contract}, BridgeTransactor: BridgeTransactor{contract: contract}, BridgeFilterer: BridgeFilterer{contract: contract}}, nil
}

// NewBridgeCaller creates a new read-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeCaller(address common.Address, caller bind.ContractCaller) (*BridgeCaller, error) {
	contract, err := bindBridge(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeCaller{contract: contract}, nil
}

// NewBridgeTransactor creates a new write-only instance of Bridge, bound to a specific deployed contract.
func NewBridgeTransactor(address common.Address, transactor bind.ContractTransactor) (*BridgeTransactor, error) {
	contract, err := bindBridge(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BridgeTransactor{contract: contract}, nil
}

// NewBridgeFilterer creates a new log filterer instance of Bridge, bound to a specific deployed contract.
func NewBridgeFilterer(address common.Address, filterer bind.ContractFilterer) (*BridgeFilterer, error) {
	contract, err := bindBridge(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BridgeFilterer{contract: contract}, nil
}

// bindBridge binds a generic wrapper to an already deployed contract.
func bindBridge(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BridgeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.BridgeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.BridgeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bridge *BridgeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bridge.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bridge *BridgeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bridge *BridgeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bridge.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Bridge *BridgeCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Bridge *BridgeSession) ADMINROLE() ([32]byte, error) {
	return _Bridge.Contract.ADMINROLE(&_Bridge.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Bridge *BridgeCallerSession) ADMINROLE() ([32]byte, error) {
	return _Bridge.Contract.ADMINROLE(&_Bridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Bridge *BridgeCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Bridge *BridgeSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Bridge.Contract.DEFAULTADMINROLE(&_Bridge.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Bridge *BridgeCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Bridge.Contract.DEFAULTADMINROLE(&_Bridge.CallOpts)
}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Bridge *BridgeCaller) RELAYROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "RELAY_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Bridge *BridgeSession) RELAYROLE() ([32]byte, error) {
	return _Bridge.Contract.RELAYROLE(&_Bridge.CallOpts)
}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Bridge *BridgeCallerSession) RELAYROLE() ([32]byte, error) {
	return _Bridge.Contract.RELAYROLE(&_Bridge.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Bridge *BridgeCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "currentEpoch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Bridge *BridgeSession) CurrentEpoch() (*big.Int, error) {
	return _Bridge.Contract.CurrentEpoch(&_Bridge.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() view returns(uint256)
func (_Bridge *BridgeCallerSession) CurrentEpoch() (*big.Int, error) {
	return _Bridge.Contract.CurrentEpoch(&_Bridge.CallOpts)
}

// GetLockedTransfers is a free data retrieval call binding the contract method 0xe31be087.
//
// Solidity: function getLockedTransfers(uint256 eventId) view returns(((address,address,uint256)[],uint256))
func (_Bridge *BridgeCaller) GetLockedTransfers(opts *bind.CallOpts, eventId *big.Int) (CommonStructsLockedTransfers, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getLockedTransfers", eventId)

	if err != nil {
		return *new(CommonStructsLockedTransfers), err
	}

	out0 := *abi.ConvertType(out[0], new(CommonStructsLockedTransfers)).(*CommonStructsLockedTransfers)

	return out0, err

}

// GetLockedTransfers is a free data retrieval call binding the contract method 0xe31be087.
//
// Solidity: function getLockedTransfers(uint256 eventId) view returns(((address,address,uint256)[],uint256))
func (_Bridge *BridgeSession) GetLockedTransfers(eventId *big.Int) (CommonStructsLockedTransfers, error) {
	return _Bridge.Contract.GetLockedTransfers(&_Bridge.CallOpts, eventId)
}

// GetLockedTransfers is a free data retrieval call binding the contract method 0xe31be087.
//
// Solidity: function getLockedTransfers(uint256 eventId) view returns(((address,address,uint256)[],uint256))
func (_Bridge *BridgeCallerSession) GetLockedTransfers(eventId *big.Int) (CommonStructsLockedTransfers, error) {
	return _Bridge.Contract.GetLockedTransfers(&_Bridge.CallOpts, eventId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Bridge *BridgeCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Bridge *BridgeSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Bridge.Contract.GetRoleAdmin(&_Bridge.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Bridge *BridgeCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Bridge.Contract.GetRoleAdmin(&_Bridge.CallOpts, role)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xcf331250.
//
// Solidity: function getValidatorSet() view returns(address[])
func (_Bridge *BridgeCaller) GetValidatorSet(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "getValidatorSet")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetValidatorSet is a free data retrieval call binding the contract method 0xcf331250.
//
// Solidity: function getValidatorSet() view returns(address[])
func (_Bridge *BridgeSession) GetValidatorSet() ([]common.Address, error) {
	return _Bridge.Contract.GetValidatorSet(&_Bridge.CallOpts)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xcf331250.
//
// Solidity: function getValidatorSet() view returns(address[])
func (_Bridge *BridgeCallerSession) GetValidatorSet() ([]common.Address, error) {
	return _Bridge.Contract.GetValidatorSet(&_Bridge.CallOpts)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Bridge *BridgeCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Bridge *BridgeSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Bridge.Contract.HasRole(&_Bridge.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Bridge *BridgeCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Bridge.Contract.HasRole(&_Bridge.CallOpts, role, account)
}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Bridge *BridgeCaller) InputEventId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "inputEventId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Bridge *BridgeSession) InputEventId() (*big.Int, error) {
	return _Bridge.Contract.InputEventId(&_Bridge.CallOpts)
}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Bridge *BridgeCallerSession) InputEventId() (*big.Int, error) {
	return _Bridge.Contract.InputEventId(&_Bridge.CallOpts)
}

// IsEpochDataSet is a free data retrieval call binding the contract method 0xc7b81f4f.
//
// Solidity: function isEpochDataSet(uint256 epochIndex) view returns(bool)
func (_Bridge *BridgeCaller) IsEpochDataSet(opts *bind.CallOpts, epochIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "isEpochDataSet", epochIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEpochDataSet is a free data retrieval call binding the contract method 0xc7b81f4f.
//
// Solidity: function isEpochDataSet(uint256 epochIndex) view returns(bool)
func (_Bridge *BridgeSession) IsEpochDataSet(epochIndex *big.Int) (bool, error) {
	return _Bridge.Contract.IsEpochDataSet(&_Bridge.CallOpts, epochIndex)
}

// IsEpochDataSet is a free data retrieval call binding the contract method 0xc7b81f4f.
//
// Solidity: function isEpochDataSet(uint256 epochIndex) view returns(bool)
func (_Bridge *BridgeCallerSession) IsEpochDataSet(epochIndex *big.Int) (bool, error) {
	return _Bridge.Contract.IsEpochDataSet(&_Bridge.CallOpts, epochIndex)
}

// IsQueueEmpty is a free data retrieval call binding the contract method 0x2bbba847.
//
// Solidity: function isQueueEmpty() view returns(bool)
func (_Bridge *BridgeCaller) IsQueueEmpty(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "isQueueEmpty")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsQueueEmpty is a free data retrieval call binding the contract method 0x2bbba847.
//
// Solidity: function isQueueEmpty() view returns(bool)
func (_Bridge *BridgeSession) IsQueueEmpty() (bool, error) {
	return _Bridge.Contract.IsQueueEmpty(&_Bridge.CallOpts)
}

// IsQueueEmpty is a free data retrieval call binding the contract method 0x2bbba847.
//
// Solidity: function isQueueEmpty() view returns(bool)
func (_Bridge *BridgeCallerSession) IsQueueEmpty() (bool, error) {
	return _Bridge.Contract.IsQueueEmpty(&_Bridge.CallOpts)
}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(bytes32)
func (_Bridge *BridgeCaller) LastProcessedBlock(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "lastProcessedBlock")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(bytes32)
func (_Bridge *BridgeSession) LastProcessedBlock() ([32]byte, error) {
	return _Bridge.Contract.LastProcessedBlock(&_Bridge.CallOpts)
}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(bytes32)
func (_Bridge *BridgeCallerSession) LastProcessedBlock() ([32]byte, error) {
	return _Bridge.Contract.LastProcessedBlock(&_Bridge.CallOpts)
}

// LastTimeframe is a free data retrieval call binding the contract method 0x86d7679e.
//
// Solidity: function lastTimeframe() view returns(uint256)
func (_Bridge *BridgeCaller) LastTimeframe(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "lastTimeframe")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastTimeframe is a free data retrieval call binding the contract method 0x86d7679e.
//
// Solidity: function lastTimeframe() view returns(uint256)
func (_Bridge *BridgeSession) LastTimeframe() (*big.Int, error) {
	return _Bridge.Contract.LastTimeframe(&_Bridge.CallOpts)
}

// LastTimeframe is a free data retrieval call binding the contract method 0x86d7679e.
//
// Solidity: function lastTimeframe() view returns(uint256)
func (_Bridge *BridgeCallerSession) LastTimeframe() (*big.Int, error) {
	return _Bridge.Contract.LastTimeframe(&_Bridge.CallOpts)
}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Bridge *BridgeCaller) LockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "lockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Bridge *BridgeSession) LockTime() (*big.Int, error) {
	return _Bridge.Contract.LockTime(&_Bridge.CallOpts)
}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Bridge *BridgeCallerSession) LockTime() (*big.Int, error) {
	return _Bridge.Contract.LockTime(&_Bridge.CallOpts)
}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Bridge *BridgeCaller) LockedTransfers(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "lockedTransfers", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Bridge *BridgeSession) LockedTransfers(arg0 *big.Int) (*big.Int, error) {
	return _Bridge.Contract.LockedTransfers(&_Bridge.CallOpts, arg0)
}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Bridge *BridgeCallerSession) LockedTransfers(arg0 *big.Int) (*big.Int, error) {
	return _Bridge.Contract.LockedTransfers(&_Bridge.CallOpts, arg0)
}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Bridge *BridgeCaller) MinSafetyBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "minSafetyBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Bridge *BridgeSession) MinSafetyBlocks() (*big.Int, error) {
	return _Bridge.Contract.MinSafetyBlocks(&_Bridge.CallOpts)
}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Bridge *BridgeCallerSession) MinSafetyBlocks() (*big.Int, error) {
	return _Bridge.Contract.MinSafetyBlocks(&_Bridge.CallOpts)
}

// OldestLockedEventId is a free data retrieval call binding the contract method 0xba8bbbe0.
//
// Solidity: function oldestLockedEventId() view returns(uint256)
func (_Bridge *BridgeCaller) OldestLockedEventId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "oldestLockedEventId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OldestLockedEventId is a free data retrieval call binding the contract method 0xba8bbbe0.
//
// Solidity: function oldestLockedEventId() view returns(uint256)
func (_Bridge *BridgeSession) OldestLockedEventId() (*big.Int, error) {
	return _Bridge.Contract.OldestLockedEventId(&_Bridge.CallOpts)
}

// OldestLockedEventId is a free data retrieval call binding the contract method 0xba8bbbe0.
//
// Solidity: function oldestLockedEventId() view returns(uint256)
func (_Bridge *BridgeCallerSession) OldestLockedEventId() (*big.Int, error) {
	return _Bridge.Contract.OldestLockedEventId(&_Bridge.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Bridge *BridgeCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Bridge *BridgeSession) Paused() (bool, error) {
	return _Bridge.Contract.Paused(&_Bridge.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Bridge *BridgeCallerSession) Paused() (bool, error) {
	return _Bridge.Contract.Paused(&_Bridge.CallOpts)
}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Bridge *BridgeCaller) SideBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "sideBridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Bridge *BridgeSession) SideBridgeAddress() (common.Address, error) {
	return _Bridge.Contract.SideBridgeAddress(&_Bridge.CallOpts)
}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Bridge *BridgeCallerSession) SideBridgeAddress() (common.Address, error) {
	return _Bridge.Contract.SideBridgeAddress(&_Bridge.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Bridge *BridgeCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Bridge *BridgeSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Bridge.Contract.SupportsInterface(&_Bridge.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Bridge *BridgeCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Bridge.Contract.SupportsInterface(&_Bridge.CallOpts, interfaceId)
}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Bridge *BridgeCaller) TimeframeSeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "timeframeSeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Bridge *BridgeSession) TimeframeSeconds() (*big.Int, error) {
	return _Bridge.Contract.TimeframeSeconds(&_Bridge.CallOpts)
}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Bridge *BridgeCallerSession) TimeframeSeconds() (*big.Int, error) {
	return _Bridge.Contract.TimeframeSeconds(&_Bridge.CallOpts)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Bridge *BridgeCaller) TokenAddresses(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "tokenAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Bridge *BridgeSession) TokenAddresses(arg0 common.Address) (common.Address, error) {
	return _Bridge.Contract.TokenAddresses(&_Bridge.CallOpts, arg0)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Bridge *BridgeCallerSession) TokenAddresses(arg0 common.Address) (common.Address, error) {
	return _Bridge.Contract.TokenAddresses(&_Bridge.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) view returns(address)
func (_Bridge *BridgeCaller) ValidatorSet(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "validatorSet", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) view returns(address)
func (_Bridge *BridgeSession) ValidatorSet(arg0 *big.Int) (common.Address, error) {
	return _Bridge.Contract.ValidatorSet(&_Bridge.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) view returns(address)
func (_Bridge *BridgeCallerSession) ValidatorSet(arg0 *big.Int) (common.Address, error) {
	return _Bridge.Contract.ValidatorSet(&_Bridge.CallOpts, arg0)
}

// WrapperAddress is a free data retrieval call binding the contract method 0x7df8b802.
//
// Solidity: function wrapperAddress() view returns(address)
func (_Bridge *BridgeCaller) WrapperAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Bridge.contract.Call(opts, &out, "wrapperAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WrapperAddress is a free data retrieval call binding the contract method 0x7df8b802.
//
// Solidity: function wrapperAddress() view returns(address)
func (_Bridge *BridgeSession) WrapperAddress() (common.Address, error) {
	return _Bridge.Contract.WrapperAddress(&_Bridge.CallOpts)
}

// WrapperAddress is a free data retrieval call binding the contract method 0x7df8b802.
//
// Solidity: function wrapperAddress() view returns(address)
func (_Bridge *BridgeCallerSession) WrapperAddress() (common.Address, error) {
	return _Bridge.Contract.WrapperAddress(&_Bridge.CallOpts)
}

// ChangeBridgeFeeRecipient is a paid mutator transaction binding the contract method 0x3928b296.
//
// Solidity: function changeBridgeFeeRecipient(address feeRecipient_) returns()
func (_Bridge *BridgeTransactor) ChangeBridgeFeeRecipient(opts *bind.TransactOpts, feeRecipient_ common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "changeBridgeFeeRecipient", feeRecipient_)
}

// ChangeBridgeFeeRecipient is a paid mutator transaction binding the contract method 0x3928b296.
//
// Solidity: function changeBridgeFeeRecipient(address feeRecipient_) returns()
func (_Bridge *BridgeSession) ChangeBridgeFeeRecipient(feeRecipient_ common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeBridgeFeeRecipient(&_Bridge.TransactOpts, feeRecipient_)
}

// ChangeBridgeFeeRecipient is a paid mutator transaction binding the contract method 0x3928b296.
//
// Solidity: function changeBridgeFeeRecipient(address feeRecipient_) returns()
func (_Bridge *BridgeTransactorSession) ChangeBridgeFeeRecipient(feeRecipient_ common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeBridgeFeeRecipient(&_Bridge.TransactOpts, feeRecipient_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Bridge *BridgeTransactor) ChangeLockTime(opts *bind.TransactOpts, lockTime_ *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "changeLockTime", lockTime_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Bridge *BridgeSession) ChangeLockTime(lockTime_ *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeLockTime(&_Bridge.TransactOpts, lockTime_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Bridge *BridgeTransactorSession) ChangeLockTime(lockTime_ *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeLockTime(&_Bridge.TransactOpts, lockTime_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Bridge *BridgeTransactor) ChangeMinSafetyBlocks(opts *bind.TransactOpts, minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "changeMinSafetyBlocks", minSafetyBlocks_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Bridge *BridgeSession) ChangeMinSafetyBlocks(minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeMinSafetyBlocks(&_Bridge.TransactOpts, minSafetyBlocks_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Bridge *BridgeTransactorSession) ChangeMinSafetyBlocks(minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeMinSafetyBlocks(&_Bridge.TransactOpts, minSafetyBlocks_)
}

// ChangeSignatureFeeCheckNumber is a paid mutator transaction binding the contract method 0xf7005587.
//
// Solidity: function changeSignatureFeeCheckNumber(uint256 signatureFeeCheckNumber_) returns()
func (_Bridge *BridgeTransactor) ChangeSignatureFeeCheckNumber(opts *bind.TransactOpts, signatureFeeCheckNumber_ *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "changeSignatureFeeCheckNumber", signatureFeeCheckNumber_)
}

// ChangeSignatureFeeCheckNumber is a paid mutator transaction binding the contract method 0xf7005587.
//
// Solidity: function changeSignatureFeeCheckNumber(uint256 signatureFeeCheckNumber_) returns()
func (_Bridge *BridgeSession) ChangeSignatureFeeCheckNumber(signatureFeeCheckNumber_ *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeSignatureFeeCheckNumber(&_Bridge.TransactOpts, signatureFeeCheckNumber_)
}

// ChangeSignatureFeeCheckNumber is a paid mutator transaction binding the contract method 0xf7005587.
//
// Solidity: function changeSignatureFeeCheckNumber(uint256 signatureFeeCheckNumber_) returns()
func (_Bridge *BridgeTransactorSession) ChangeSignatureFeeCheckNumber(signatureFeeCheckNumber_ *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeSignatureFeeCheckNumber(&_Bridge.TransactOpts, signatureFeeCheckNumber_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Bridge *BridgeTransactor) ChangeTimeframeSeconds(opts *bind.TransactOpts, timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "changeTimeframeSeconds", timeframeSeconds_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Bridge *BridgeSession) ChangeTimeframeSeconds(timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeTimeframeSeconds(&_Bridge.TransactOpts, timeframeSeconds_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Bridge *BridgeTransactorSession) ChangeTimeframeSeconds(timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeTimeframeSeconds(&_Bridge.TransactOpts, timeframeSeconds_)
}

// ChangeTransferFeeRecipient is a paid mutator transaction binding the contract method 0xb87330e9.
//
// Solidity: function changeTransferFeeRecipient(address feeRecipient_) returns()
func (_Bridge *BridgeTransactor) ChangeTransferFeeRecipient(opts *bind.TransactOpts, feeRecipient_ common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "changeTransferFeeRecipient", feeRecipient_)
}

// ChangeTransferFeeRecipient is a paid mutator transaction binding the contract method 0xb87330e9.
//
// Solidity: function changeTransferFeeRecipient(address feeRecipient_) returns()
func (_Bridge *BridgeSession) ChangeTransferFeeRecipient(feeRecipient_ common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeTransferFeeRecipient(&_Bridge.TransactOpts, feeRecipient_)
}

// ChangeTransferFeeRecipient is a paid mutator transaction binding the contract method 0xb87330e9.
//
// Solidity: function changeTransferFeeRecipient(address feeRecipient_) returns()
func (_Bridge *BridgeTransactorSession) ChangeTransferFeeRecipient(feeRecipient_ common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.ChangeTransferFeeRecipient(&_Bridge.TransactOpts, feeRecipient_)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Bridge *BridgeSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.GrantRole(&_Bridge.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.GrantRole(&_Bridge.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x3f45c799.
//
// Solidity: function initialize((address,address,address,address,address[],address[],address,address,uint256,uint256,uint256) args, uint256 minimumDifficulty) returns()
func (_Bridge *BridgeTransactor) Initialize(opts *bind.TransactOpts, args CommonStructsConstructorArgs, minimumDifficulty *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "initialize", args, minimumDifficulty)
}

// Initialize is a paid mutator transaction binding the contract method 0x3f45c799.
//
// Solidity: function initialize((address,address,address,address,address[],address[],address,address,uint256,uint256,uint256) args, uint256 minimumDifficulty) returns()
func (_Bridge *BridgeSession) Initialize(args CommonStructsConstructorArgs, minimumDifficulty *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Initialize(&_Bridge.TransactOpts, args, minimumDifficulty)
}

// Initialize is a paid mutator transaction binding the contract method 0x3f45c799.
//
// Solidity: function initialize((address,address,address,address,address[],address[],address,address,uint256,uint256,uint256) args, uint256 minimumDifficulty) returns()
func (_Bridge *BridgeTransactorSession) Initialize(args CommonStructsConstructorArgs, minimumDifficulty *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Initialize(&_Bridge.TransactOpts, args, minimumDifficulty)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Bridge *BridgeTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Bridge *BridgeSession) Pause() (*types.Transaction, error) {
	return _Bridge.Contract.Pause(&_Bridge.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Bridge *BridgeTransactorSession) Pause() (*types.Transaction, error) {
	return _Bridge.Contract.Pause(&_Bridge.TransactOpts)
}

// RemoveLockedTransfers is a paid mutator transaction binding the contract method 0x331a891a.
//
// Solidity: function removeLockedTransfers(uint256 eventId) returns()
func (_Bridge *BridgeTransactor) RemoveLockedTransfers(opts *bind.TransactOpts, eventId *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "removeLockedTransfers", eventId)
}

// RemoveLockedTransfers is a paid mutator transaction binding the contract method 0x331a891a.
//
// Solidity: function removeLockedTransfers(uint256 eventId) returns()
func (_Bridge *BridgeSession) RemoveLockedTransfers(eventId *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.RemoveLockedTransfers(&_Bridge.TransactOpts, eventId)
}

// RemoveLockedTransfers is a paid mutator transaction binding the contract method 0x331a891a.
//
// Solidity: function removeLockedTransfers(uint256 eventId) returns()
func (_Bridge *BridgeTransactorSession) RemoveLockedTransfers(eventId *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.RemoveLockedTransfers(&_Bridge.TransactOpts, eventId)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Bridge *BridgeSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RenounceRole(&_Bridge.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RenounceRole(&_Bridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Bridge *BridgeSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RevokeRole(&_Bridge.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Bridge *BridgeTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.RevokeRole(&_Bridge.TransactOpts, role, account)
}

// SetEpochData is a paid mutator transaction binding the contract method 0xe88b6626.
//
// Solidity: function setEpochData(uint256 epochNum, uint256 fullSizeIn128Resultion, uint256 branchDepth, uint256[] merkleNodes) returns()
func (_Bridge *BridgeTransactor) SetEpochData(opts *bind.TransactOpts, epochNum *big.Int, fullSizeIn128Resultion *big.Int, branchDepth *big.Int, merkleNodes []*big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setEpochData", epochNum, fullSizeIn128Resultion, branchDepth, merkleNodes)
}

// SetEpochData is a paid mutator transaction binding the contract method 0xe88b6626.
//
// Solidity: function setEpochData(uint256 epochNum, uint256 fullSizeIn128Resultion, uint256 branchDepth, uint256[] merkleNodes) returns()
func (_Bridge *BridgeSession) SetEpochData(epochNum *big.Int, fullSizeIn128Resultion *big.Int, branchDepth *big.Int, merkleNodes []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.SetEpochData(&_Bridge.TransactOpts, epochNum, fullSizeIn128Resultion, branchDepth, merkleNodes)
}

// SetEpochData is a paid mutator transaction binding the contract method 0xe88b6626.
//
// Solidity: function setEpochData(uint256 epochNum, uint256 fullSizeIn128Resultion, uint256 branchDepth, uint256[] merkleNodes) returns()
func (_Bridge *BridgeTransactorSession) SetEpochData(epochNum *big.Int, fullSizeIn128Resultion *big.Int, branchDepth *big.Int, merkleNodes []*big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.SetEpochData(&_Bridge.TransactOpts, epochNum, fullSizeIn128Resultion, branchDepth, merkleNodes)
}

// SetSideBridge is a paid mutator transaction binding the contract method 0x21d3d536.
//
// Solidity: function setSideBridge(address _sideBridgeAddress) returns()
func (_Bridge *BridgeTransactor) SetSideBridge(opts *bind.TransactOpts, _sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "setSideBridge", _sideBridgeAddress)
}

// SetSideBridge is a paid mutator transaction binding the contract method 0x21d3d536.
//
// Solidity: function setSideBridge(address _sideBridgeAddress) returns()
func (_Bridge *BridgeSession) SetSideBridge(_sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetSideBridge(&_Bridge.TransactOpts, _sideBridgeAddress)
}

// SetSideBridge is a paid mutator transaction binding the contract method 0x21d3d536.
//
// Solidity: function setSideBridge(address _sideBridgeAddress) returns()
func (_Bridge *BridgeTransactorSession) SetSideBridge(_sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.SetSideBridge(&_Bridge.TransactOpts, _sideBridgeAddress)
}

// SubmitTransferAura is a paid mutator transaction binding the contract method 0x119b980d.
//
// Solidity: function submitTransferAura(((bytes3,bytes3,bytes32,bytes,bytes32,bytes,bytes4,bytes,uint64)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],(address,int64)[])[],uint64) auraProof) returns()
func (_Bridge *BridgeTransactor) SubmitTransferAura(opts *bind.TransactOpts, auraProof CheckAuraAuraProof) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "submitTransferAura", auraProof)
}

// SubmitTransferAura is a paid mutator transaction binding the contract method 0x119b980d.
//
// Solidity: function submitTransferAura(((bytes3,bytes3,bytes32,bytes,bytes32,bytes,bytes4,bytes,uint64)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],(address,int64)[])[],uint64) auraProof) returns()
func (_Bridge *BridgeSession) SubmitTransferAura(auraProof CheckAuraAuraProof) (*types.Transaction, error) {
	return _Bridge.Contract.SubmitTransferAura(&_Bridge.TransactOpts, auraProof)
}

// SubmitTransferAura is a paid mutator transaction binding the contract method 0x119b980d.
//
// Solidity: function submitTransferAura(((bytes3,bytes3,bytes32,bytes,bytes32,bytes,bytes4,bytes,uint64)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],(address,int64)[])[],uint64) auraProof) returns()
func (_Bridge *BridgeTransactorSession) SubmitTransferAura(auraProof CheckAuraAuraProof) (*types.Transaction, error) {
	return _Bridge.Contract.SubmitTransferAura(&_Bridge.TransactOpts, auraProof)
}

// SubmitTransferPoSA is a paid mutator transaction binding the contract method 0xaa9cea38.
//
// Solidity: function submitTransferPoSA(((bytes3,bytes3,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes)[],(bytes[],uint256,(address,address,uint256)[]),uint64) posaProof) returns()
func (_Bridge *BridgeTransactor) SubmitTransferPoSA(opts *bind.TransactOpts, posaProof CheckPoSAPoSAProof) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "submitTransferPoSA", posaProof)
}

// SubmitTransferPoSA is a paid mutator transaction binding the contract method 0xaa9cea38.
//
// Solidity: function submitTransferPoSA(((bytes3,bytes3,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes)[],(bytes[],uint256,(address,address,uint256)[]),uint64) posaProof) returns()
func (_Bridge *BridgeSession) SubmitTransferPoSA(posaProof CheckPoSAPoSAProof) (*types.Transaction, error) {
	return _Bridge.Contract.SubmitTransferPoSA(&_Bridge.TransactOpts, posaProof)
}

// SubmitTransferPoSA is a paid mutator transaction binding the contract method 0xaa9cea38.
//
// Solidity: function submitTransferPoSA(((bytes3,bytes3,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes)[],(bytes[],uint256,(address,address,uint256)[]),uint64) posaProof) returns()
func (_Bridge *BridgeTransactorSession) SubmitTransferPoSA(posaProof CheckPoSAPoSAProof) (*types.Transaction, error) {
	return _Bridge.Contract.SubmitTransferPoSA(&_Bridge.TransactOpts, posaProof)
}

// SubmitTransferPoW is a paid mutator transaction binding the contract method 0xe1d862be.
//
// Solidity: function submitTransferPoW(((bytes3,bytes3,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Bridge *BridgeTransactor) SubmitTransferPoW(opts *bind.TransactOpts, powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "submitTransferPoW", powProof)
}

// SubmitTransferPoW is a paid mutator transaction binding the contract method 0xe1d862be.
//
// Solidity: function submitTransferPoW(((bytes3,bytes3,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Bridge *BridgeSession) SubmitTransferPoW(powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Bridge.Contract.SubmitTransferPoW(&_Bridge.TransactOpts, powProof)
}

// SubmitTransferPoW is a paid mutator transaction binding the contract method 0xe1d862be.
//
// Solidity: function submitTransferPoW(((bytes3,bytes3,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Bridge *BridgeTransactorSession) SubmitTransferPoW(powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Bridge.Contract.SubmitTransferPoW(&_Bridge.TransactOpts, powProof)
}

// SubmitValidatorSetChangesPoSA is a paid mutator transaction binding the contract method 0xd83a358f.
//
// Solidity: function submitValidatorSetChangesPoSA(((bytes3,bytes3,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes)[],(bytes[],uint256,(address,address,uint256)[]),uint64) posaProof) returns()
func (_Bridge *BridgeTransactor) SubmitValidatorSetChangesPoSA(opts *bind.TransactOpts, posaProof CheckPoSAPoSAProof) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "submitValidatorSetChangesPoSA", posaProof)
}

// SubmitValidatorSetChangesPoSA is a paid mutator transaction binding the contract method 0xd83a358f.
//
// Solidity: function submitValidatorSetChangesPoSA(((bytes3,bytes3,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes)[],(bytes[],uint256,(address,address,uint256)[]),uint64) posaProof) returns()
func (_Bridge *BridgeSession) SubmitValidatorSetChangesPoSA(posaProof CheckPoSAPoSAProof) (*types.Transaction, error) {
	return _Bridge.Contract.SubmitValidatorSetChangesPoSA(&_Bridge.TransactOpts, posaProof)
}

// SubmitValidatorSetChangesPoSA is a paid mutator transaction binding the contract method 0xd83a358f.
//
// Solidity: function submitValidatorSetChangesPoSA(((bytes3,bytes3,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes)[],(bytes[],uint256,(address,address,uint256)[]),uint64) posaProof) returns()
func (_Bridge *BridgeTransactorSession) SubmitValidatorSetChangesPoSA(posaProof CheckPoSAPoSAProof) (*types.Transaction, error) {
	return _Bridge.Contract.SubmitValidatorSetChangesPoSA(&_Bridge.TransactOpts, posaProof)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Bridge *BridgeTransactor) TokensAdd(opts *bind.TransactOpts, tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "tokensAdd", tokenThisAddress, tokenSideAddress)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Bridge *BridgeSession) TokensAdd(tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TokensAdd(&_Bridge.TransactOpts, tokenThisAddress, tokenSideAddress)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Bridge *BridgeTransactorSession) TokensAdd(tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TokensAdd(&_Bridge.TransactOpts, tokenThisAddress, tokenSideAddress)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Bridge *BridgeTransactor) TokensAddBatch(opts *bind.TransactOpts, tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "tokensAddBatch", tokenThisAddresses, tokenSideAddresses)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Bridge *BridgeSession) TokensAddBatch(tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TokensAddBatch(&_Bridge.TransactOpts, tokenThisAddresses, tokenSideAddresses)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Bridge *BridgeTransactorSession) TokensAddBatch(tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TokensAddBatch(&_Bridge.TransactOpts, tokenThisAddresses, tokenSideAddresses)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Bridge *BridgeTransactor) TokensRemove(opts *bind.TransactOpts, tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "tokensRemove", tokenThisAddress)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Bridge *BridgeSession) TokensRemove(tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TokensRemove(&_Bridge.TransactOpts, tokenThisAddress)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Bridge *BridgeTransactorSession) TokensRemove(tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TokensRemove(&_Bridge.TransactOpts, tokenThisAddress)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Bridge *BridgeTransactor) TokensRemoveBatch(opts *bind.TransactOpts, tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "tokensRemoveBatch", tokenThisAddresses)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Bridge *BridgeSession) TokensRemoveBatch(tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TokensRemoveBatch(&_Bridge.TransactOpts, tokenThisAddresses)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Bridge *BridgeTransactorSession) TokensRemoveBatch(tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Bridge.Contract.TokensRemoveBatch(&_Bridge.TransactOpts, tokenThisAddresses)
}

// TriggerTransfers is a paid mutator transaction binding the contract method 0x746b5c42.
//
// Solidity: function triggerTransfers() returns()
func (_Bridge *BridgeTransactor) TriggerTransfers(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "triggerTransfers")
}

// TriggerTransfers is a paid mutator transaction binding the contract method 0x746b5c42.
//
// Solidity: function triggerTransfers() returns()
func (_Bridge *BridgeSession) TriggerTransfers() (*types.Transaction, error) {
	return _Bridge.Contract.TriggerTransfers(&_Bridge.TransactOpts)
}

// TriggerTransfers is a paid mutator transaction binding the contract method 0x746b5c42.
//
// Solidity: function triggerTransfers() returns()
func (_Bridge *BridgeTransactorSession) TriggerTransfers() (*types.Transaction, error) {
	return _Bridge.Contract.TriggerTransfers(&_Bridge.TransactOpts)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 eventId) returns()
func (_Bridge *BridgeTransactor) UnlockTransfers(opts *bind.TransactOpts, eventId *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "unlockTransfers", eventId)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 eventId) returns()
func (_Bridge *BridgeSession) UnlockTransfers(eventId *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.UnlockTransfers(&_Bridge.TransactOpts, eventId)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 eventId) returns()
func (_Bridge *BridgeTransactorSession) UnlockTransfers(eventId *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.UnlockTransfers(&_Bridge.TransactOpts, eventId)
}

// UnlockTransfersBatch is a paid mutator transaction binding the contract method 0x8ac1f86f.
//
// Solidity: function unlockTransfersBatch() returns()
func (_Bridge *BridgeTransactor) UnlockTransfersBatch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "unlockTransfersBatch")
}

// UnlockTransfersBatch is a paid mutator transaction binding the contract method 0x8ac1f86f.
//
// Solidity: function unlockTransfersBatch() returns()
func (_Bridge *BridgeSession) UnlockTransfersBatch() (*types.Transaction, error) {
	return _Bridge.Contract.UnlockTransfersBatch(&_Bridge.TransactOpts)
}

// UnlockTransfersBatch is a paid mutator transaction binding the contract method 0x8ac1f86f.
//
// Solidity: function unlockTransfersBatch() returns()
func (_Bridge *BridgeTransactorSession) UnlockTransfersBatch() (*types.Transaction, error) {
	return _Bridge.Contract.UnlockTransfersBatch(&_Bridge.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Bridge *BridgeTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Bridge *BridgeSession) Unpause() (*types.Transaction, error) {
	return _Bridge.Contract.Unpause(&_Bridge.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Bridge *BridgeTransactorSession) Unpause() (*types.Transaction, error) {
	return _Bridge.Contract.Unpause(&_Bridge.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf8b74433.
//
// Solidity: function withdraw(address tokenThisAddress, address toAddress, uint256 amount, bool unwrapSide, bytes feeSignature, uint256 transferFee, uint256 bridgeFee) payable returns()
func (_Bridge *BridgeTransactor) Withdraw(opts *bind.TransactOpts, tokenThisAddress common.Address, toAddress common.Address, amount *big.Int, unwrapSide bool, feeSignature []byte, transferFee *big.Int, bridgeFee *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "withdraw", tokenThisAddress, toAddress, amount, unwrapSide, feeSignature, transferFee, bridgeFee)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf8b74433.
//
// Solidity: function withdraw(address tokenThisAddress, address toAddress, uint256 amount, bool unwrapSide, bytes feeSignature, uint256 transferFee, uint256 bridgeFee) payable returns()
func (_Bridge *BridgeSession) Withdraw(tokenThisAddress common.Address, toAddress common.Address, amount *big.Int, unwrapSide bool, feeSignature []byte, transferFee *big.Int, bridgeFee *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Withdraw(&_Bridge.TransactOpts, tokenThisAddress, toAddress, amount, unwrapSide, feeSignature, transferFee, bridgeFee)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf8b74433.
//
// Solidity: function withdraw(address tokenThisAddress, address toAddress, uint256 amount, bool unwrapSide, bytes feeSignature, uint256 transferFee, uint256 bridgeFee) payable returns()
func (_Bridge *BridgeTransactorSession) Withdraw(tokenThisAddress common.Address, toAddress common.Address, amount *big.Int, unwrapSide bool, feeSignature []byte, transferFee *big.Int, bridgeFee *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.Withdraw(&_Bridge.TransactOpts, tokenThisAddress, toAddress, amount, unwrapSide, feeSignature, transferFee, bridgeFee)
}

// WrapWithdraw is a paid mutator transaction binding the contract method 0x91c5bc0a.
//
// Solidity: function wrapWithdraw(address toAddress, bytes feeSignature, uint256 transferFee, uint256 bridgeFee) payable returns()
func (_Bridge *BridgeTransactor) WrapWithdraw(opts *bind.TransactOpts, toAddress common.Address, feeSignature []byte, transferFee *big.Int, bridgeFee *big.Int) (*types.Transaction, error) {
	return _Bridge.contract.Transact(opts, "wrapWithdraw", toAddress, feeSignature, transferFee, bridgeFee)
}

// WrapWithdraw is a paid mutator transaction binding the contract method 0x91c5bc0a.
//
// Solidity: function wrapWithdraw(address toAddress, bytes feeSignature, uint256 transferFee, uint256 bridgeFee) payable returns()
func (_Bridge *BridgeSession) WrapWithdraw(toAddress common.Address, feeSignature []byte, transferFee *big.Int, bridgeFee *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.WrapWithdraw(&_Bridge.TransactOpts, toAddress, feeSignature, transferFee, bridgeFee)
}

// WrapWithdraw is a paid mutator transaction binding the contract method 0x91c5bc0a.
//
// Solidity: function wrapWithdraw(address toAddress, bytes feeSignature, uint256 transferFee, uint256 bridgeFee) payable returns()
func (_Bridge *BridgeTransactorSession) WrapWithdraw(toAddress common.Address, feeSignature []byte, transferFee *big.Int, bridgeFee *big.Int) (*types.Transaction, error) {
	return _Bridge.Contract.WrapWithdraw(&_Bridge.TransactOpts, toAddress, feeSignature, transferFee, bridgeFee)
}

// BridgeRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Bridge contract.
type BridgeRoleAdminChangedIterator struct {
	Event *BridgeRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRoleAdminChanged represents a RoleAdminChanged event raised by the Bridge contract.
type BridgeRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Bridge *BridgeFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*BridgeRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &BridgeRoleAdminChangedIterator{contract: _Bridge.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Bridge *BridgeFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *BridgeRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRoleAdminChanged)
				if err := _Bridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Bridge *BridgeFilterer) ParseRoleAdminChanged(log types.Log) (*BridgeRoleAdminChanged, error) {
	event := new(BridgeRoleAdminChanged)
	if err := _Bridge.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Bridge contract.
type BridgeRoleGrantedIterator struct {
	Event *BridgeRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRoleGranted represents a RoleGranted event raised by the Bridge contract.
type BridgeRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BridgeRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BridgeRoleGrantedIterator{contract: _Bridge.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *BridgeRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRoleGranted)
				if err := _Bridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) ParseRoleGranted(log types.Log) (*BridgeRoleGranted, error) {
	event := new(BridgeRoleGranted)
	if err := _Bridge.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Bridge contract.
type BridgeRoleRevokedIterator struct {
	Event *BridgeRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeRoleRevoked represents a RoleRevoked event raised by the Bridge contract.
type BridgeRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*BridgeRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &BridgeRoleRevokedIterator{contract: _Bridge.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *BridgeRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeRoleRevoked)
				if err := _Bridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Bridge *BridgeFilterer) ParseRoleRevoked(log types.Log) (*BridgeRoleRevoked, error) {
	event := new(BridgeRoleRevoked)
	if err := _Bridge.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Bridge contract.
type BridgeTransferIterator struct {
	Event *BridgeTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTransfer represents a Transfer event raised by the Bridge contract.
type BridgeTransfer struct {
	EventId *big.Int
	Queue   []CommonStructsTransfer
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2.
//
// Solidity: event Transfer(uint256 indexed eventId, (address,address,uint256)[] queue)
func (_Bridge *BridgeFilterer) FilterTransfer(opts *bind.FilterOpts, eventId []*big.Int) (*BridgeTransferIterator, error) {

	var eventIdRule []interface{}
	for _, eventIdItem := range eventId {
		eventIdRule = append(eventIdRule, eventIdItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Transfer", eventIdRule)
	if err != nil {
		return nil, err
	}
	return &BridgeTransferIterator{contract: _Bridge.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2.
//
// Solidity: event Transfer(uint256 indexed eventId, (address,address,uint256)[] queue)
func (_Bridge *BridgeFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *BridgeTransfer, eventId []*big.Int) (event.Subscription, error) {

	var eventIdRule []interface{}
	for _, eventIdItem := range eventId {
		eventIdRule = append(eventIdRule, eventIdItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Transfer", eventIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTransfer)
				if err := _Bridge.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2.
//
// Solidity: event Transfer(uint256 indexed eventId, (address,address,uint256)[] queue)
func (_Bridge *BridgeFilterer) ParseTransfer(log types.Log) (*BridgeTransfer, error) {
	event := new(BridgeTransfer)
	if err := _Bridge.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeTransferFinishIterator is returned from FilterTransferFinish and is used to iterate over the raw logs and unpacked data for TransferFinish events raised by the Bridge contract.
type BridgeTransferFinishIterator struct {
	Event *BridgeTransferFinish // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeTransferFinishIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTransferFinish)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeTransferFinish)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeTransferFinishIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTransferFinishIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTransferFinish represents a TransferFinish event raised by the Bridge contract.
type BridgeTransferFinish struct {
	EventId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransferFinish is a free log retrieval operation binding the contract event 0x78ff9b3176bb0d6421590f3816e75cb15a9ffa2b7a1a028f40a3f4e029eb39f2.
//
// Solidity: event TransferFinish(uint256 indexed eventId)
func (_Bridge *BridgeFilterer) FilterTransferFinish(opts *bind.FilterOpts, eventId []*big.Int) (*BridgeTransferFinishIterator, error) {

	var eventIdRule []interface{}
	for _, eventIdItem := range eventId {
		eventIdRule = append(eventIdRule, eventIdItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "TransferFinish", eventIdRule)
	if err != nil {
		return nil, err
	}
	return &BridgeTransferFinishIterator{contract: _Bridge.contract, event: "TransferFinish", logs: logs, sub: sub}, nil
}

// WatchTransferFinish is a free log subscription operation binding the contract event 0x78ff9b3176bb0d6421590f3816e75cb15a9ffa2b7a1a028f40a3f4e029eb39f2.
//
// Solidity: event TransferFinish(uint256 indexed eventId)
func (_Bridge *BridgeFilterer) WatchTransferFinish(opts *bind.WatchOpts, sink chan<- *BridgeTransferFinish, eventId []*big.Int) (event.Subscription, error) {

	var eventIdRule []interface{}
	for _, eventIdItem := range eventId {
		eventIdRule = append(eventIdRule, eventIdItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "TransferFinish", eventIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTransferFinish)
				if err := _Bridge.contract.UnpackLog(event, "TransferFinish", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferFinish is a log parse operation binding the contract event 0x78ff9b3176bb0d6421590f3816e75cb15a9ffa2b7a1a028f40a3f4e029eb39f2.
//
// Solidity: event TransferFinish(uint256 indexed eventId)
func (_Bridge *BridgeFilterer) ParseTransferFinish(log types.Log) (*BridgeTransferFinish, error) {
	event := new(BridgeTransferFinish)
	if err := _Bridge.contract.UnpackLog(event, "TransferFinish", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeTransferSubmitIterator is returned from FilterTransferSubmit and is used to iterate over the raw logs and unpacked data for TransferSubmit events raised by the Bridge contract.
type BridgeTransferSubmitIterator struct {
	Event *BridgeTransferSubmit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeTransferSubmitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeTransferSubmit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeTransferSubmit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeTransferSubmitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeTransferSubmitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeTransferSubmit represents a TransferSubmit event raised by the Bridge contract.
type BridgeTransferSubmit struct {
	EventId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransferSubmit is a free log retrieval operation binding the contract event 0x196c47048e38df7a4fe6e581c8f4f2e2ba67ac0dd45b90da756e97bd61d9dd3b.
//
// Solidity: event TransferSubmit(uint256 indexed eventId)
func (_Bridge *BridgeFilterer) FilterTransferSubmit(opts *bind.FilterOpts, eventId []*big.Int) (*BridgeTransferSubmitIterator, error) {

	var eventIdRule []interface{}
	for _, eventIdItem := range eventId {
		eventIdRule = append(eventIdRule, eventIdItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "TransferSubmit", eventIdRule)
	if err != nil {
		return nil, err
	}
	return &BridgeTransferSubmitIterator{contract: _Bridge.contract, event: "TransferSubmit", logs: logs, sub: sub}, nil
}

// WatchTransferSubmit is a free log subscription operation binding the contract event 0x196c47048e38df7a4fe6e581c8f4f2e2ba67ac0dd45b90da756e97bd61d9dd3b.
//
// Solidity: event TransferSubmit(uint256 indexed eventId)
func (_Bridge *BridgeFilterer) WatchTransferSubmit(opts *bind.WatchOpts, sink chan<- *BridgeTransferSubmit, eventId []*big.Int) (event.Subscription, error) {

	var eventIdRule []interface{}
	for _, eventIdItem := range eventId {
		eventIdRule = append(eventIdRule, eventIdItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "TransferSubmit", eventIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeTransferSubmit)
				if err := _Bridge.contract.UnpackLog(event, "TransferSubmit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransferSubmit is a log parse operation binding the contract event 0x196c47048e38df7a4fe6e581c8f4f2e2ba67ac0dd45b90da756e97bd61d9dd3b.
//
// Solidity: event TransferSubmit(uint256 indexed eventId)
func (_Bridge *BridgeFilterer) ParseTransferSubmit(log types.Log) (*BridgeTransferSubmit, error) {
	event := new(BridgeTransferSubmit)
	if err := _Bridge.contract.UnpackLog(event, "TransferSubmit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Bridge contract.
type BridgeUnpausedIterator struct {
	Event *BridgeUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeUnpaused represents a Unpaused event raised by the Bridge contract.
type BridgeUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Bridge *BridgeFilterer) FilterUnpaused(opts *bind.FilterOpts) (*BridgeUnpausedIterator, error) {

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &BridgeUnpausedIterator{contract: _Bridge.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Bridge *BridgeFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *BridgeUnpaused) (event.Subscription, error) {

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeUnpaused)
				if err := _Bridge.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Bridge *BridgeFilterer) ParseUnpaused(log types.Log) (*BridgeUnpaused, error) {
	event := new(BridgeUnpaused)
	if err := _Bridge.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BridgeWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Bridge contract.
type BridgeWithdrawIterator struct {
	Event *BridgeWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BridgeWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BridgeWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BridgeWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BridgeWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BridgeWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BridgeWithdraw represents a Withdraw event raised by the Bridge contract.
type BridgeWithdraw struct {
	From              common.Address
	EventId           *big.Int
	TokenFrom         common.Address
	TokenTo           common.Address
	Amount            *big.Int
	TransferFeeAmount *big.Int
	BridgeFeeAmount   *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x3b823d3b837159352faddd1e0b9e9f584ca2bfe976238f6a9e6da4615cf2f444.
//
// Solidity: event Withdraw(address indexed from, uint256 eventId, address tokenFrom, address tokenTo, uint256 amount, uint256 transferFeeAmount, uint256 bridgeFeeAmount)
func (_Bridge *BridgeFilterer) FilterWithdraw(opts *bind.FilterOpts, from []common.Address) (*BridgeWithdrawIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Bridge.contract.FilterLogs(opts, "Withdraw", fromRule)
	if err != nil {
		return nil, err
	}
	return &BridgeWithdrawIterator{contract: _Bridge.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x3b823d3b837159352faddd1e0b9e9f584ca2bfe976238f6a9e6da4615cf2f444.
//
// Solidity: event Withdraw(address indexed from, uint256 eventId, address tokenFrom, address tokenTo, uint256 amount, uint256 transferFeeAmount, uint256 bridgeFeeAmount)
func (_Bridge *BridgeFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *BridgeWithdraw, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Bridge.contract.WatchLogs(opts, "Withdraw", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BridgeWithdraw)
				if err := _Bridge.contract.UnpackLog(event, "Withdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdraw is a log parse operation binding the contract event 0x3b823d3b837159352faddd1e0b9e9f584ca2bfe976238f6a9e6da4615cf2f444.
//
// Solidity: event Withdraw(address indexed from, uint256 eventId, address tokenFrom, address tokenTo, uint256 amount, uint256 transferFeeAmount, uint256 bridgeFeeAmount)
func (_Bridge *BridgeFilterer) ParseWithdraw(log types.Log) (*BridgeWithdraw, error) {
	event := new(BridgeWithdraw)
	if err := _Bridge.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
