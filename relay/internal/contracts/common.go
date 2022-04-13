// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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
	Blocks    []CheckAuraBlockAura
	Transfer  CommonStructsTransferProof
	VsChanges []CheckAuraValidatorSetProof
}

// CheckAuraBlockAura is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraBlockAura struct {
	P0Seal      []byte
	P0Bare      []byte
	P1          []byte
	ParentHash  [32]byte
	P2          []byte
	ReceiptHash [32]byte
	P3          []byte
	S1          []byte
	Step        []byte
	S2          []byte
	Signature   []byte
	Type        uint8
	DeltaIndex  int64
}

// CheckAuraValidatorSetProof is an auto generated low-level Go binding around an user-defined struct.
type CheckAuraValidatorSetProof struct {
	ReceiptProof [][]byte
	DeltaAddress common.Address
	DeltaIndex   int64
}

// CheckPoWBlockPoW is an auto generated low-level Go binding around an user-defined struct.
type CheckPoWBlockPoW struct {
	P0WithNonce         []byte
	P0WithoutNonce      []byte
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

// CommonMetaData contains all meta data concerning the Common contract.
var CommonMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"queue\",\"type\":\"tuple[]\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"TransferFinish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"TransferSubmit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"proof\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"el\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"proofStart\",\"type\":\"uint256\"}],\"name\":\"CalcReceiptsHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"p\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"eventContractAddress\",\"type\":\"address\"}],\"name\":\"CalcTransferReceiptsHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p0WithNonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p0WithoutNonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parentOrReceiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"difficulty\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"nonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p6\",\"type\":\"bytes\"},{\"internalType\":\"uint256[]\",\"name\":\"dataSetLookup\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"witnessForLookup\",\"type\":\"uint256[]\"}],\"internalType\":\"structCheckPoW.BlockPoW[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"}],\"internalType\":\"structCheckPoW.PoWProof\",\"name\":\"powProof\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"sideBridgeAddress\",\"type\":\"address\"}],\"name\":\"CheckPoW_\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RELAY_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"}],\"name\":\"changeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"feeRecipient_\",\"type\":\"address\"}],\"name\":\"changeFeeRecipient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lockTime_\",\"type\":\"uint256\"}],\"name\":\"changeLockTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minSafetyBlocks_\",\"type\":\"uint256\"}],\"name\":\"changeMinSafetyBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timeframeSeconds_\",\"type\":\"uint256\"}],\"name\":\"changeTimeframeSeconds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inputEventId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epochIndex\",\"type\":\"uint256\"}],\"name\":\"isEpochDataSet\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lockedTransfers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minSafetyBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oldestLockedEventId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"removeLockedTransfers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"wrapper\",\"type\":\"address\"}],\"name\":\"setAmbWrapper\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epochNum\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fullSizeIn128Resultion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"branchDepth\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"merkleNodes\",\"type\":\"uint256[]\"}],\"name\":\"setEpochData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sideBridgeAddress\",\"type\":\"address\"}],\"name\":\"setSideBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sideBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p0WithNonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p0WithoutNonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parentOrReceiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"difficulty\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"nonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p6\",\"type\":\"bytes\"},{\"internalType\":\"uint256[]\",\"name\":\"dataSetLookup\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"witnessForLookup\",\"type\":\"uint256[]\"}],\"internalType\":\"structCheckPoW.BlockPoW[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"}],\"internalType\":\"structCheckPoW.PoWProof\",\"name\":\"powProof\",\"type\":\"tuple\"}],\"name\":\"submitTransferPoW\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeframeSeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenSideAddress\",\"type\":\"address\"}],\"name\":\"tokensAdd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenSideAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensAddBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"}],\"name\":\"tokensRemove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensRemoveBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"unlockTransfers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockTransfersBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p0WithNonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p0WithoutNonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parentOrReceiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"difficulty\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"nonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p6\",\"type\":\"bytes\"},{\"internalType\":\"uint256[]\",\"name\":\"dataSetLookup\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"witnessForLookup\",\"type\":\"uint256[]\"}],\"internalType\":\"structCheckPoW.BlockPoW\",\"name\":\"block_\",\"type\":\"tuple\"}],\"name\":\"verifyEthash\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAmbAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"}],\"name\":\"wrap_withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p0_seal\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p0_bare\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parent_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"receipt_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s1\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"step\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"type_\",\"type\":\"uint8\"},{\"internalType\":\"int64\",\"name\":\"delta_index\",\"type\":\"int64\"}],\"internalType\":\"structCheckAura.BlockAura[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"address\",\"name\":\"delta_address\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"delta_index\",\"type\":\"int64\"}],\"internalType\":\"structCheckAura.ValidatorSetProof[]\",\"name\":\"vs_changes\",\"type\":\"tuple[]\"}],\"internalType\":\"structCheckAura.AuraProof\",\"name\":\"auraProof\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"minSafetyBlocks\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sideBridgeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"validatorSetAddress\",\"type\":\"address\"}],\"name\":\"CheckAura_\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetValidatorSet\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastProcessedBlock\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p0_seal\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p0_bare\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parent_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"receipt_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s1\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"step\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"type_\",\"type\":\"uint8\"},{\"internalType\":\"int64\",\"name\":\"delta_index\",\"type\":\"int64\"}],\"internalType\":\"structCheckAura.BlockAura[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"address\",\"name\":\"delta_address\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"delta_index\",\"type\":\"int64\"}],\"internalType\":\"structCheckAura.ValidatorSetProof[]\",\"name\":\"vs_changes\",\"type\":\"tuple[]\"}],\"internalType\":\"structCheckAura.AuraProof\",\"name\":\"auraProof\",\"type\":\"tuple\"}],\"name\":\"submitTransferAura\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatorSet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// CommonABI is the input ABI used to generate the binding from.
// Deprecated: Use CommonMetaData.ABI instead.
var CommonABI = CommonMetaData.ABI

// Common is an auto generated Go binding around an Ethereum contract.
type Common struct {
	CommonCaller     // Read-only binding to the contract
	CommonTransactor // Write-only binding to the contract
	CommonFilterer   // Log filterer for contract events
}

// CommonCaller is an auto generated read-only Go binding around an Ethereum contract.
type CommonCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommonTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CommonTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommonFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CommonFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommonSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CommonSession struct {
	Contract     *Common           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CommonCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CommonCallerSession struct {
	Contract *CommonCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// CommonTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CommonTransactorSession struct {
	Contract     *CommonTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CommonRaw is an auto generated low-level Go binding around an Ethereum contract.
type CommonRaw struct {
	Contract *Common // Generic contract binding to access the raw methods on
}

// CommonCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CommonCallerRaw struct {
	Contract *CommonCaller // Generic read-only contract binding to access the raw methods on
}

// CommonTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CommonTransactorRaw struct {
	Contract *CommonTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCommon creates a new instance of Common, bound to a specific deployed contract.
func NewCommon(address common.Address, backend bind.ContractBackend) (*Common, error) {
	contract, err := bindCommon(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Common{CommonCaller: CommonCaller{contract: contract}, CommonTransactor: CommonTransactor{contract: contract}, CommonFilterer: CommonFilterer{contract: contract}}, nil
}

// NewCommonCaller creates a new read-only instance of Common, bound to a specific deployed contract.
func NewCommonCaller(address common.Address, caller bind.ContractCaller) (*CommonCaller, error) {
	contract, err := bindCommon(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CommonCaller{contract: contract}, nil
}

// NewCommonTransactor creates a new write-only instance of Common, bound to a specific deployed contract.
func NewCommonTransactor(address common.Address, transactor bind.ContractTransactor) (*CommonTransactor, error) {
	contract, err := bindCommon(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CommonTransactor{contract: contract}, nil
}

// NewCommonFilterer creates a new log filterer instance of Common, bound to a specific deployed contract.
func NewCommonFilterer(address common.Address, filterer bind.ContractFilterer) (*CommonFilterer, error) {
	contract, err := bindCommon(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CommonFilterer{contract: contract}, nil
}

// bindCommon binds a generic wrapper to an already deployed contract.
func bindCommon(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CommonABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Common *CommonRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Common.Contract.CommonCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Common *CommonRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Common.Contract.CommonTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Common *CommonRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Common.Contract.CommonTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Common *CommonCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Common.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Common *CommonTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Common.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Common *CommonTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Common.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Common *CommonCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Common *CommonSession) ADMINROLE() ([32]byte, error) {
	return _Common.Contract.ADMINROLE(&_Common.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Common *CommonCallerSession) ADMINROLE() ([32]byte, error) {
	return _Common.Contract.ADMINROLE(&_Common.CallOpts)
}

// CalcReceiptsHash is a free data retrieval call binding the contract method 0xe7899536.
//
// Solidity: function CalcReceiptsHash(bytes[] proof, bytes32 el, uint256 proofStart) pure returns(bytes32)
func (_Common *CommonCaller) CalcReceiptsHash(opts *bind.CallOpts, proof [][]byte, el [32]byte, proofStart *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "CalcReceiptsHash", proof, el, proofStart)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalcReceiptsHash is a free data retrieval call binding the contract method 0xe7899536.
//
// Solidity: function CalcReceiptsHash(bytes[] proof, bytes32 el, uint256 proofStart) pure returns(bytes32)
func (_Common *CommonSession) CalcReceiptsHash(proof [][]byte, el [32]byte, proofStart *big.Int) ([32]byte, error) {
	return _Common.Contract.CalcReceiptsHash(&_Common.CallOpts, proof, el, proofStart)
}

// CalcReceiptsHash is a free data retrieval call binding the contract method 0xe7899536.
//
// Solidity: function CalcReceiptsHash(bytes[] proof, bytes32 el, uint256 proofStart) pure returns(bytes32)
func (_Common *CommonCallerSession) CalcReceiptsHash(proof [][]byte, el [32]byte, proofStart *big.Int) ([32]byte, error) {
	return _Common.Contract.CalcReceiptsHash(&_Common.CallOpts, proof, el, proofStart)
}

// CalcTransferReceiptsHash is a free data retrieval call binding the contract method 0x90d0308f.
//
// Solidity: function CalcTransferReceiptsHash((bytes[],uint256,(address,address,uint256)[]) p, address eventContractAddress) pure returns(bytes32)
func (_Common *CommonCaller) CalcTransferReceiptsHash(opts *bind.CallOpts, p CommonStructsTransferProof, eventContractAddress common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "CalcTransferReceiptsHash", p, eventContractAddress)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalcTransferReceiptsHash is a free data retrieval call binding the contract method 0x90d0308f.
//
// Solidity: function CalcTransferReceiptsHash((bytes[],uint256,(address,address,uint256)[]) p, address eventContractAddress) pure returns(bytes32)
func (_Common *CommonSession) CalcTransferReceiptsHash(p CommonStructsTransferProof, eventContractAddress common.Address) ([32]byte, error) {
	return _Common.Contract.CalcTransferReceiptsHash(&_Common.CallOpts, p, eventContractAddress)
}

// CalcTransferReceiptsHash is a free data retrieval call binding the contract method 0x90d0308f.
//
// Solidity: function CalcTransferReceiptsHash((bytes[],uint256,(address,address,uint256)[]) p, address eventContractAddress) pure returns(bytes32)
func (_Common *CommonCallerSession) CalcTransferReceiptsHash(p CommonStructsTransferProof, eventContractAddress common.Address) ([32]byte, error) {
	return _Common.Contract.CalcTransferReceiptsHash(&_Common.CallOpts, p, eventContractAddress)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Common *CommonCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Common *CommonSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Common.Contract.DEFAULTADMINROLE(&_Common.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Common *CommonCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Common.Contract.DEFAULTADMINROLE(&_Common.CallOpts)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xc7456a69.
//
// Solidity: function GetValidatorSet() view returns(address[])
func (_Common *CommonCaller) GetValidatorSet(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "GetValidatorSet")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetValidatorSet is a free data retrieval call binding the contract method 0xc7456a69.
//
// Solidity: function GetValidatorSet() view returns(address[])
func (_Common *CommonSession) GetValidatorSet() ([]common.Address, error) {
	return _Common.Contract.GetValidatorSet(&_Common.CallOpts)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xc7456a69.
//
// Solidity: function GetValidatorSet() view returns(address[])
func (_Common *CommonCallerSession) GetValidatorSet() ([]common.Address, error) {
	return _Common.Contract.GetValidatorSet(&_Common.CallOpts)
}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Common *CommonCaller) RELAYROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "RELAY_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Common *CommonSession) RELAYROLE() ([32]byte, error) {
	return _Common.Contract.RELAYROLE(&_Common.CallOpts)
}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Common *CommonCallerSession) RELAYROLE() ([32]byte, error) {
	return _Common.Contract.RELAYROLE(&_Common.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Common *CommonCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Common *CommonSession) Fee() (*big.Int, error) {
	return _Common.Contract.Fee(&_Common.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Common *CommonCallerSession) Fee() (*big.Int, error) {
	return _Common.Contract.Fee(&_Common.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Common *CommonCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Common *CommonSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Common.Contract.GetRoleAdmin(&_Common.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Common *CommonCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Common.Contract.GetRoleAdmin(&_Common.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Common *CommonCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Common *CommonSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Common.Contract.HasRole(&_Common.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Common *CommonCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Common.Contract.HasRole(&_Common.CallOpts, role, account)
}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Common *CommonCaller) InputEventId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "inputEventId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Common *CommonSession) InputEventId() (*big.Int, error) {
	return _Common.Contract.InputEventId(&_Common.CallOpts)
}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Common *CommonCallerSession) InputEventId() (*big.Int, error) {
	return _Common.Contract.InputEventId(&_Common.CallOpts)
}

// IsEpochDataSet is a free data retrieval call binding the contract method 0xc7b81f4f.
//
// Solidity: function isEpochDataSet(uint256 epochIndex) view returns(bool)
func (_Common *CommonCaller) IsEpochDataSet(opts *bind.CallOpts, epochIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "isEpochDataSet", epochIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEpochDataSet is a free data retrieval call binding the contract method 0xc7b81f4f.
//
// Solidity: function isEpochDataSet(uint256 epochIndex) view returns(bool)
func (_Common *CommonSession) IsEpochDataSet(epochIndex *big.Int) (bool, error) {
	return _Common.Contract.IsEpochDataSet(&_Common.CallOpts, epochIndex)
}

// IsEpochDataSet is a free data retrieval call binding the contract method 0xc7b81f4f.
//
// Solidity: function isEpochDataSet(uint256 epochIndex) view returns(bool)
func (_Common *CommonCallerSession) IsEpochDataSet(epochIndex *big.Int) (bool, error) {
	return _Common.Contract.IsEpochDataSet(&_Common.CallOpts, epochIndex)
}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(bytes32)
func (_Common *CommonCaller) LastProcessedBlock(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "lastProcessedBlock")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(bytes32)
func (_Common *CommonSession) LastProcessedBlock() ([32]byte, error) {
	return _Common.Contract.LastProcessedBlock(&_Common.CallOpts)
}

// LastProcessedBlock is a free data retrieval call binding the contract method 0x33de61d2.
//
// Solidity: function lastProcessedBlock() view returns(bytes32)
func (_Common *CommonCallerSession) LastProcessedBlock() ([32]byte, error) {
	return _Common.Contract.LastProcessedBlock(&_Common.CallOpts)
}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Common *CommonCaller) LockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "lockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Common *CommonSession) LockTime() (*big.Int, error) {
	return _Common.Contract.LockTime(&_Common.CallOpts)
}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Common *CommonCallerSession) LockTime() (*big.Int, error) {
	return _Common.Contract.LockTime(&_Common.CallOpts)
}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Common *CommonCaller) LockedTransfers(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "lockedTransfers", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Common *CommonSession) LockedTransfers(arg0 *big.Int) (*big.Int, error) {
	return _Common.Contract.LockedTransfers(&_Common.CallOpts, arg0)
}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Common *CommonCallerSession) LockedTransfers(arg0 *big.Int) (*big.Int, error) {
	return _Common.Contract.LockedTransfers(&_Common.CallOpts, arg0)
}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Common *CommonCaller) MinSafetyBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "minSafetyBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Common *CommonSession) MinSafetyBlocks() (*big.Int, error) {
	return _Common.Contract.MinSafetyBlocks(&_Common.CallOpts)
}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Common *CommonCallerSession) MinSafetyBlocks() (*big.Int, error) {
	return _Common.Contract.MinSafetyBlocks(&_Common.CallOpts)
}

// OldestLockedEventId is a free data retrieval call binding the contract method 0xba8bbbe0.
//
// Solidity: function oldestLockedEventId() view returns(uint256)
func (_Common *CommonCaller) OldestLockedEventId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "oldestLockedEventId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OldestLockedEventId is a free data retrieval call binding the contract method 0xba8bbbe0.
//
// Solidity: function oldestLockedEventId() view returns(uint256)
func (_Common *CommonSession) OldestLockedEventId() (*big.Int, error) {
	return _Common.Contract.OldestLockedEventId(&_Common.CallOpts)
}

// OldestLockedEventId is a free data retrieval call binding the contract method 0xba8bbbe0.
//
// Solidity: function oldestLockedEventId() view returns(uint256)
func (_Common *CommonCallerSession) OldestLockedEventId() (*big.Int, error) {
	return _Common.Contract.OldestLockedEventId(&_Common.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Common *CommonCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Common *CommonSession) Paused() (bool, error) {
	return _Common.Contract.Paused(&_Common.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Common *CommonCallerSession) Paused() (bool, error) {
	return _Common.Contract.Paused(&_Common.CallOpts)
}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Common *CommonCaller) SideBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "sideBridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Common *CommonSession) SideBridgeAddress() (common.Address, error) {
	return _Common.Contract.SideBridgeAddress(&_Common.CallOpts)
}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Common *CommonCallerSession) SideBridgeAddress() (common.Address, error) {
	return _Common.Contract.SideBridgeAddress(&_Common.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Common *CommonCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Common *CommonSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Common.Contract.SupportsInterface(&_Common.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Common *CommonCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Common.Contract.SupportsInterface(&_Common.CallOpts, interfaceId)
}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Common *CommonCaller) TimeframeSeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "timeframeSeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Common *CommonSession) TimeframeSeconds() (*big.Int, error) {
	return _Common.Contract.TimeframeSeconds(&_Common.CallOpts)
}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Common *CommonCallerSession) TimeframeSeconds() (*big.Int, error) {
	return _Common.Contract.TimeframeSeconds(&_Common.CallOpts)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Common *CommonCaller) TokenAddresses(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "tokenAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Common *CommonSession) TokenAddresses(arg0 common.Address) (common.Address, error) {
	return _Common.Contract.TokenAddresses(&_Common.CallOpts, arg0)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Common *CommonCallerSession) TokenAddresses(arg0 common.Address) (common.Address, error) {
	return _Common.Contract.TokenAddresses(&_Common.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) view returns(address)
func (_Common *CommonCaller) ValidatorSet(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "validatorSet", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) view returns(address)
func (_Common *CommonSession) ValidatorSet(arg0 *big.Int) (common.Address, error) {
	return _Common.Contract.ValidatorSet(&_Common.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) view returns(address)
func (_Common *CommonCallerSession) ValidatorSet(arg0 *big.Int) (common.Address, error) {
	return _Common.Contract.ValidatorSet(&_Common.CallOpts, arg0)
}

// VerifyEthash is a free data retrieval call binding the contract method 0x0c31a003.
//
// Solidity: function verifyEthash((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[]) block_) view returns()
func (_Common *CommonCaller) VerifyEthash(opts *bind.CallOpts, block_ CheckPoWBlockPoW) error {
	var out []interface{}
	err := _Common.contract.Call(opts, &out, "verifyEthash", block_)

	if err != nil {
		return err
	}

	return err

}

// VerifyEthash is a free data retrieval call binding the contract method 0x0c31a003.
//
// Solidity: function verifyEthash((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[]) block_) view returns()
func (_Common *CommonSession) VerifyEthash(block_ CheckPoWBlockPoW) error {
	return _Common.Contract.VerifyEthash(&_Common.CallOpts, block_)
}

// VerifyEthash is a free data retrieval call binding the contract method 0x0c31a003.
//
// Solidity: function verifyEthash((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[]) block_) view returns()
func (_Common *CommonCallerSession) VerifyEthash(block_ CheckPoWBlockPoW) error {
	return _Common.Contract.VerifyEthash(&_Common.CallOpts, block_)
}

// CheckAura is a paid mutator transaction binding the contract method 0x88e09c86.
//
// Solidity: function CheckAura_(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,uint8,int64)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof, uint256 minSafetyBlocks, address sideBridgeAddress, address validatorSetAddress) returns()
func (_Common *CommonTransactor) CheckAura(opts *bind.TransactOpts, auraProof CheckAuraAuraProof, minSafetyBlocks *big.Int, sideBridgeAddress common.Address, validatorSetAddress common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "CheckAura_", auraProof, minSafetyBlocks, sideBridgeAddress, validatorSetAddress)
}

// CheckAura is a paid mutator transaction binding the contract method 0x88e09c86.
//
// Solidity: function CheckAura_(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,uint8,int64)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof, uint256 minSafetyBlocks, address sideBridgeAddress, address validatorSetAddress) returns()
func (_Common *CommonSession) CheckAura(auraProof CheckAuraAuraProof, minSafetyBlocks *big.Int, sideBridgeAddress common.Address, validatorSetAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.CheckAura(&_Common.TransactOpts, auraProof, minSafetyBlocks, sideBridgeAddress, validatorSetAddress)
}

// CheckAura is a paid mutator transaction binding the contract method 0x88e09c86.
//
// Solidity: function CheckAura_(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,uint8,int64)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof, uint256 minSafetyBlocks, address sideBridgeAddress, address validatorSetAddress) returns()
func (_Common *CommonTransactorSession) CheckAura(auraProof CheckAuraAuraProof, minSafetyBlocks *big.Int, sideBridgeAddress common.Address, validatorSetAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.CheckAura(&_Common.TransactOpts, auraProof, minSafetyBlocks, sideBridgeAddress, validatorSetAddress)
}

// CheckPoW is a paid mutator transaction binding the contract method 0x973365e0.
//
// Solidity: function CheckPoW_(((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof, address sideBridgeAddress) returns()
func (_Common *CommonTransactor) CheckPoW(opts *bind.TransactOpts, powProof CheckPoWPoWProof, sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "CheckPoW_", powProof, sideBridgeAddress)
}

// CheckPoW is a paid mutator transaction binding the contract method 0x973365e0.
//
// Solidity: function CheckPoW_(((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof, address sideBridgeAddress) returns()
func (_Common *CommonSession) CheckPoW(powProof CheckPoWPoWProof, sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.CheckPoW(&_Common.TransactOpts, powProof, sideBridgeAddress)
}

// CheckPoW is a paid mutator transaction binding the contract method 0x973365e0.
//
// Solidity: function CheckPoW_(((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof, address sideBridgeAddress) returns()
func (_Common *CommonTransactorSession) CheckPoW(powProof CheckPoWPoWProof, sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.CheckPoW(&_Common.TransactOpts, powProof, sideBridgeAddress)
}

// ChangeFee is a paid mutator transaction binding the contract method 0x6a1db1bf.
//
// Solidity: function changeFee(uint256 fee_) returns()
func (_Common *CommonTransactor) ChangeFee(opts *bind.TransactOpts, fee_ *big.Int) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "changeFee", fee_)
}

// ChangeFee is a paid mutator transaction binding the contract method 0x6a1db1bf.
//
// Solidity: function changeFee(uint256 fee_) returns()
func (_Common *CommonSession) ChangeFee(fee_ *big.Int) (*types.Transaction, error) {
	return _Common.Contract.ChangeFee(&_Common.TransactOpts, fee_)
}

// ChangeFee is a paid mutator transaction binding the contract method 0x6a1db1bf.
//
// Solidity: function changeFee(uint256 fee_) returns()
func (_Common *CommonTransactorSession) ChangeFee(fee_ *big.Int) (*types.Transaction, error) {
	return _Common.Contract.ChangeFee(&_Common.TransactOpts, fee_)
}

// ChangeFeeRecipient is a paid mutator transaction binding the contract method 0x23604071.
//
// Solidity: function changeFeeRecipient(address feeRecipient_) returns()
func (_Common *CommonTransactor) ChangeFeeRecipient(opts *bind.TransactOpts, feeRecipient_ common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "changeFeeRecipient", feeRecipient_)
}

// ChangeFeeRecipient is a paid mutator transaction binding the contract method 0x23604071.
//
// Solidity: function changeFeeRecipient(address feeRecipient_) returns()
func (_Common *CommonSession) ChangeFeeRecipient(feeRecipient_ common.Address) (*types.Transaction, error) {
	return _Common.Contract.ChangeFeeRecipient(&_Common.TransactOpts, feeRecipient_)
}

// ChangeFeeRecipient is a paid mutator transaction binding the contract method 0x23604071.
//
// Solidity: function changeFeeRecipient(address feeRecipient_) returns()
func (_Common *CommonTransactorSession) ChangeFeeRecipient(feeRecipient_ common.Address) (*types.Transaction, error) {
	return _Common.Contract.ChangeFeeRecipient(&_Common.TransactOpts, feeRecipient_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Common *CommonTransactor) ChangeLockTime(opts *bind.TransactOpts, lockTime_ *big.Int) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "changeLockTime", lockTime_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Common *CommonSession) ChangeLockTime(lockTime_ *big.Int) (*types.Transaction, error) {
	return _Common.Contract.ChangeLockTime(&_Common.TransactOpts, lockTime_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Common *CommonTransactorSession) ChangeLockTime(lockTime_ *big.Int) (*types.Transaction, error) {
	return _Common.Contract.ChangeLockTime(&_Common.TransactOpts, lockTime_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Common *CommonTransactor) ChangeMinSafetyBlocks(opts *bind.TransactOpts, minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "changeMinSafetyBlocks", minSafetyBlocks_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Common *CommonSession) ChangeMinSafetyBlocks(minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Common.Contract.ChangeMinSafetyBlocks(&_Common.TransactOpts, minSafetyBlocks_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Common *CommonTransactorSession) ChangeMinSafetyBlocks(minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Common.Contract.ChangeMinSafetyBlocks(&_Common.TransactOpts, minSafetyBlocks_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Common *CommonTransactor) ChangeTimeframeSeconds(opts *bind.TransactOpts, timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "changeTimeframeSeconds", timeframeSeconds_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Common *CommonSession) ChangeTimeframeSeconds(timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Common.Contract.ChangeTimeframeSeconds(&_Common.TransactOpts, timeframeSeconds_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Common *CommonTransactorSession) ChangeTimeframeSeconds(timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Common.Contract.ChangeTimeframeSeconds(&_Common.TransactOpts, timeframeSeconds_)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Common *CommonTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Common *CommonSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Common.Contract.GrantRole(&_Common.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Common *CommonTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Common.Contract.GrantRole(&_Common.TransactOpts, role, account)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Common *CommonTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Common *CommonSession) Pause() (*types.Transaction, error) {
	return _Common.Contract.Pause(&_Common.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Common *CommonTransactorSession) Pause() (*types.Transaction, error) {
	return _Common.Contract.Pause(&_Common.TransactOpts)
}

// RemoveLockedTransfers is a paid mutator transaction binding the contract method 0x331a891a.
//
// Solidity: function removeLockedTransfers(uint256 event_id) returns()
func (_Common *CommonTransactor) RemoveLockedTransfers(opts *bind.TransactOpts, event_id *big.Int) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "removeLockedTransfers", event_id)
}

// RemoveLockedTransfers is a paid mutator transaction binding the contract method 0x331a891a.
//
// Solidity: function removeLockedTransfers(uint256 event_id) returns()
func (_Common *CommonSession) RemoveLockedTransfers(event_id *big.Int) (*types.Transaction, error) {
	return _Common.Contract.RemoveLockedTransfers(&_Common.TransactOpts, event_id)
}

// RemoveLockedTransfers is a paid mutator transaction binding the contract method 0x331a891a.
//
// Solidity: function removeLockedTransfers(uint256 event_id) returns()
func (_Common *CommonTransactorSession) RemoveLockedTransfers(event_id *big.Int) (*types.Transaction, error) {
	return _Common.Contract.RemoveLockedTransfers(&_Common.TransactOpts, event_id)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Common *CommonTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Common *CommonSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Common.Contract.RenounceRole(&_Common.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Common *CommonTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Common.Contract.RenounceRole(&_Common.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Common *CommonTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Common *CommonSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Common.Contract.RevokeRole(&_Common.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Common *CommonTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Common.Contract.RevokeRole(&_Common.TransactOpts, role, account)
}

// SetAmbWrapper is a paid mutator transaction binding the contract method 0xe1c6da90.
//
// Solidity: function setAmbWrapper(address wrapper) returns()
func (_Common *CommonTransactor) SetAmbWrapper(opts *bind.TransactOpts, wrapper common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "setAmbWrapper", wrapper)
}

// SetAmbWrapper is a paid mutator transaction binding the contract method 0xe1c6da90.
//
// Solidity: function setAmbWrapper(address wrapper) returns()
func (_Common *CommonSession) SetAmbWrapper(wrapper common.Address) (*types.Transaction, error) {
	return _Common.Contract.SetAmbWrapper(&_Common.TransactOpts, wrapper)
}

// SetAmbWrapper is a paid mutator transaction binding the contract method 0xe1c6da90.
//
// Solidity: function setAmbWrapper(address wrapper) returns()
func (_Common *CommonTransactorSession) SetAmbWrapper(wrapper common.Address) (*types.Transaction, error) {
	return _Common.Contract.SetAmbWrapper(&_Common.TransactOpts, wrapper)
}

// SetEpochData is a paid mutator transaction binding the contract method 0xe88b6626.
//
// Solidity: function setEpochData(uint256 epochNum, uint256 fullSizeIn128Resultion, uint256 branchDepth, uint256[] merkleNodes) returns()
func (_Common *CommonTransactor) SetEpochData(opts *bind.TransactOpts, epochNum *big.Int, fullSizeIn128Resultion *big.Int, branchDepth *big.Int, merkleNodes []*big.Int) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "setEpochData", epochNum, fullSizeIn128Resultion, branchDepth, merkleNodes)
}

// SetEpochData is a paid mutator transaction binding the contract method 0xe88b6626.
//
// Solidity: function setEpochData(uint256 epochNum, uint256 fullSizeIn128Resultion, uint256 branchDepth, uint256[] merkleNodes) returns()
func (_Common *CommonSession) SetEpochData(epochNum *big.Int, fullSizeIn128Resultion *big.Int, branchDepth *big.Int, merkleNodes []*big.Int) (*types.Transaction, error) {
	return _Common.Contract.SetEpochData(&_Common.TransactOpts, epochNum, fullSizeIn128Resultion, branchDepth, merkleNodes)
}

// SetEpochData is a paid mutator transaction binding the contract method 0xe88b6626.
//
// Solidity: function setEpochData(uint256 epochNum, uint256 fullSizeIn128Resultion, uint256 branchDepth, uint256[] merkleNodes) returns()
func (_Common *CommonTransactorSession) SetEpochData(epochNum *big.Int, fullSizeIn128Resultion *big.Int, branchDepth *big.Int, merkleNodes []*big.Int) (*types.Transaction, error) {
	return _Common.Contract.SetEpochData(&_Common.TransactOpts, epochNum, fullSizeIn128Resultion, branchDepth, merkleNodes)
}

// SetSideBridge is a paid mutator transaction binding the contract method 0x21d3d536.
//
// Solidity: function setSideBridge(address _sideBridgeAddress) returns()
func (_Common *CommonTransactor) SetSideBridge(opts *bind.TransactOpts, _sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "setSideBridge", _sideBridgeAddress)
}

// SetSideBridge is a paid mutator transaction binding the contract method 0x21d3d536.
//
// Solidity: function setSideBridge(address _sideBridgeAddress) returns()
func (_Common *CommonSession) SetSideBridge(_sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.SetSideBridge(&_Common.TransactOpts, _sideBridgeAddress)
}

// SetSideBridge is a paid mutator transaction binding the contract method 0x21d3d536.
//
// Solidity: function setSideBridge(address _sideBridgeAddress) returns()
func (_Common *CommonTransactorSession) SetSideBridge(_sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.SetSideBridge(&_Common.TransactOpts, _sideBridgeAddress)
}

// SubmitTransferAura is a paid mutator transaction binding the contract method 0x0e7ced0c.
//
// Solidity: function submitTransferAura(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,uint8,int64)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof) returns()
func (_Common *CommonTransactor) SubmitTransferAura(opts *bind.TransactOpts, auraProof CheckAuraAuraProof) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "submitTransferAura", auraProof)
}

// SubmitTransferAura is a paid mutator transaction binding the contract method 0x0e7ced0c.
//
// Solidity: function submitTransferAura(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,uint8,int64)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof) returns()
func (_Common *CommonSession) SubmitTransferAura(auraProof CheckAuraAuraProof) (*types.Transaction, error) {
	return _Common.Contract.SubmitTransferAura(&_Common.TransactOpts, auraProof)
}

// SubmitTransferAura is a paid mutator transaction binding the contract method 0x0e7ced0c.
//
// Solidity: function submitTransferAura(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,uint8,int64)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof) returns()
func (_Common *CommonTransactorSession) SubmitTransferAura(auraProof CheckAuraAuraProof) (*types.Transaction, error) {
	return _Common.Contract.SubmitTransferAura(&_Common.TransactOpts, auraProof)
}

// SubmitTransferPoW is a paid mutator transaction binding the contract method 0x1ef0bd25.
//
// Solidity: function submitTransferPoW(((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Common *CommonTransactor) SubmitTransferPoW(opts *bind.TransactOpts, powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "submitTransferPoW", powProof)
}

// SubmitTransferPoW is a paid mutator transaction binding the contract method 0x1ef0bd25.
//
// Solidity: function submitTransferPoW(((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Common *CommonSession) SubmitTransferPoW(powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Common.Contract.SubmitTransferPoW(&_Common.TransactOpts, powProof)
}

// SubmitTransferPoW is a paid mutator transaction binding the contract method 0x1ef0bd25.
//
// Solidity: function submitTransferPoW(((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Common *CommonTransactorSession) SubmitTransferPoW(powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Common.Contract.SubmitTransferPoW(&_Common.TransactOpts, powProof)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Common *CommonTransactor) TokensAdd(opts *bind.TransactOpts, tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "tokensAdd", tokenThisAddress, tokenSideAddress)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Common *CommonSession) TokensAdd(tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.TokensAdd(&_Common.TransactOpts, tokenThisAddress, tokenSideAddress)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Common *CommonTransactorSession) TokensAdd(tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.TokensAdd(&_Common.TransactOpts, tokenThisAddress, tokenSideAddress)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Common *CommonTransactor) TokensAddBatch(opts *bind.TransactOpts, tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "tokensAddBatch", tokenThisAddresses, tokenSideAddresses)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Common *CommonSession) TokensAddBatch(tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Common.Contract.TokensAddBatch(&_Common.TransactOpts, tokenThisAddresses, tokenSideAddresses)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Common *CommonTransactorSession) TokensAddBatch(tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Common.Contract.TokensAddBatch(&_Common.TransactOpts, tokenThisAddresses, tokenSideAddresses)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Common *CommonTransactor) TokensRemove(opts *bind.TransactOpts, tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "tokensRemove", tokenThisAddress)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Common *CommonSession) TokensRemove(tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.TokensRemove(&_Common.TransactOpts, tokenThisAddress)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Common *CommonTransactorSession) TokensRemove(tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.TokensRemove(&_Common.TransactOpts, tokenThisAddress)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Common *CommonTransactor) TokensRemoveBatch(opts *bind.TransactOpts, tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "tokensRemoveBatch", tokenThisAddresses)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Common *CommonSession) TokensRemoveBatch(tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Common.Contract.TokensRemoveBatch(&_Common.TransactOpts, tokenThisAddresses)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Common *CommonTransactorSession) TokensRemoveBatch(tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Common.Contract.TokensRemoveBatch(&_Common.TransactOpts, tokenThisAddresses)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 event_id) returns()
func (_Common *CommonTransactor) UnlockTransfers(opts *bind.TransactOpts, event_id *big.Int) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "unlockTransfers", event_id)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 event_id) returns()
func (_Common *CommonSession) UnlockTransfers(event_id *big.Int) (*types.Transaction, error) {
	return _Common.Contract.UnlockTransfers(&_Common.TransactOpts, event_id)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 event_id) returns()
func (_Common *CommonTransactorSession) UnlockTransfers(event_id *big.Int) (*types.Transaction, error) {
	return _Common.Contract.UnlockTransfers(&_Common.TransactOpts, event_id)
}

// UnlockTransfersBatch is a paid mutator transaction binding the contract method 0x8ac1f86f.
//
// Solidity: function unlockTransfersBatch() returns()
func (_Common *CommonTransactor) UnlockTransfersBatch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "unlockTransfersBatch")
}

// UnlockTransfersBatch is a paid mutator transaction binding the contract method 0x8ac1f86f.
//
// Solidity: function unlockTransfersBatch() returns()
func (_Common *CommonSession) UnlockTransfersBatch() (*types.Transaction, error) {
	return _Common.Contract.UnlockTransfersBatch(&_Common.TransactOpts)
}

// UnlockTransfersBatch is a paid mutator transaction binding the contract method 0x8ac1f86f.
//
// Solidity: function unlockTransfersBatch() returns()
func (_Common *CommonTransactorSession) UnlockTransfersBatch() (*types.Transaction, error) {
	return _Common.Contract.UnlockTransfersBatch(&_Common.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Common *CommonTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Common *CommonSession) Unpause() (*types.Transaction, error) {
	return _Common.Contract.Unpause(&_Common.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Common *CommonTransactorSession) Unpause() (*types.Transaction, error) {
	return _Common.Contract.Unpause(&_Common.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAmbAddress, address toAddress, uint256 amount) payable returns()
func (_Common *CommonTransactor) Withdraw(opts *bind.TransactOpts, tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "withdraw", tokenAmbAddress, toAddress, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAmbAddress, address toAddress, uint256 amount) payable returns()
func (_Common *CommonSession) Withdraw(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Common.Contract.Withdraw(&_Common.TransactOpts, tokenAmbAddress, toAddress, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAmbAddress, address toAddress, uint256 amount) payable returns()
func (_Common *CommonTransactorSession) Withdraw(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Common.Contract.Withdraw(&_Common.TransactOpts, tokenAmbAddress, toAddress, amount)
}

// WrapWithdraw is a paid mutator transaction binding the contract method 0x23291c20.
//
// Solidity: function wrap_withdraw(address toAddress) payable returns()
func (_Common *CommonTransactor) WrapWithdraw(opts *bind.TransactOpts, toAddress common.Address) (*types.Transaction, error) {
	return _Common.contract.Transact(opts, "wrap_withdraw", toAddress)
}

// WrapWithdraw is a paid mutator transaction binding the contract method 0x23291c20.
//
// Solidity: function wrap_withdraw(address toAddress) payable returns()
func (_Common *CommonSession) WrapWithdraw(toAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.WrapWithdraw(&_Common.TransactOpts, toAddress)
}

// WrapWithdraw is a paid mutator transaction binding the contract method 0x23291c20.
//
// Solidity: function wrap_withdraw(address toAddress) payable returns()
func (_Common *CommonTransactorSession) WrapWithdraw(toAddress common.Address) (*types.Transaction, error) {
	return _Common.Contract.WrapWithdraw(&_Common.TransactOpts, toAddress)
}

// CommonPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Common contract.
type CommonPausedIterator struct {
	Event *CommonPaused // Event containing the contract specifics and raw log

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
func (it *CommonPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommonPaused)
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
		it.Event = new(CommonPaused)
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
func (it *CommonPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommonPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommonPaused represents a Paused event raised by the Common contract.
type CommonPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Common *CommonFilterer) FilterPaused(opts *bind.FilterOpts) (*CommonPausedIterator, error) {

	logs, sub, err := _Common.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &CommonPausedIterator{contract: _Common.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Common *CommonFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *CommonPaused) (event.Subscription, error) {

	logs, sub, err := _Common.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommonPaused)
				if err := _Common.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Common *CommonFilterer) ParsePaused(log types.Log) (*CommonPaused, error) {
	event := new(CommonPaused)
	if err := _Common.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommonRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Common contract.
type CommonRoleAdminChangedIterator struct {
	Event *CommonRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *CommonRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommonRoleAdminChanged)
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
		it.Event = new(CommonRoleAdminChanged)
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
func (it *CommonRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommonRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommonRoleAdminChanged represents a RoleAdminChanged event raised by the Common contract.
type CommonRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Common *CommonFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*CommonRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Common.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &CommonRoleAdminChangedIterator{contract: _Common.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Common *CommonFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *CommonRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Common.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommonRoleAdminChanged)
				if err := _Common.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Common *CommonFilterer) ParseRoleAdminChanged(log types.Log) (*CommonRoleAdminChanged, error) {
	event := new(CommonRoleAdminChanged)
	if err := _Common.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommonRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Common contract.
type CommonRoleGrantedIterator struct {
	Event *CommonRoleGranted // Event containing the contract specifics and raw log

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
func (it *CommonRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommonRoleGranted)
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
		it.Event = new(CommonRoleGranted)
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
func (it *CommonRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommonRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommonRoleGranted represents a RoleGranted event raised by the Common contract.
type CommonRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Common *CommonFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*CommonRoleGrantedIterator, error) {

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

	logs, sub, err := _Common.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &CommonRoleGrantedIterator{contract: _Common.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Common *CommonFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *CommonRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Common.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommonRoleGranted)
				if err := _Common.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Common *CommonFilterer) ParseRoleGranted(log types.Log) (*CommonRoleGranted, error) {
	event := new(CommonRoleGranted)
	if err := _Common.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommonRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Common contract.
type CommonRoleRevokedIterator struct {
	Event *CommonRoleRevoked // Event containing the contract specifics and raw log

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
func (it *CommonRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommonRoleRevoked)
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
		it.Event = new(CommonRoleRevoked)
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
func (it *CommonRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommonRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommonRoleRevoked represents a RoleRevoked event raised by the Common contract.
type CommonRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Common *CommonFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*CommonRoleRevokedIterator, error) {

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

	logs, sub, err := _Common.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &CommonRoleRevokedIterator{contract: _Common.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Common *CommonFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *CommonRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Common.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommonRoleRevoked)
				if err := _Common.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Common *CommonFilterer) ParseRoleRevoked(log types.Log) (*CommonRoleRevoked, error) {
	event := new(CommonRoleRevoked)
	if err := _Common.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommonTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Common contract.
type CommonTransferIterator struct {
	Event *CommonTransfer // Event containing the contract specifics and raw log

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
func (it *CommonTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommonTransfer)
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
		it.Event = new(CommonTransfer)
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
func (it *CommonTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommonTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommonTransfer represents a Transfer event raised by the Common contract.
type CommonTransfer struct {
	EventId *big.Int
	Queue   []CommonStructsTransfer
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2.
//
// Solidity: event Transfer(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Common *CommonFilterer) FilterTransfer(opts *bind.FilterOpts, event_id []*big.Int) (*CommonTransferIterator, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Common.contract.FilterLogs(opts, "Transfer", event_idRule)
	if err != nil {
		return nil, err
	}
	return &CommonTransferIterator{contract: _Common.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2.
//
// Solidity: event Transfer(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Common *CommonFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CommonTransfer, event_id []*big.Int) (event.Subscription, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Common.contract.WatchLogs(opts, "Transfer", event_idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommonTransfer)
				if err := _Common.contract.UnpackLog(event, "Transfer", log); err != nil {
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
// Solidity: event Transfer(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Common *CommonFilterer) ParseTransfer(log types.Log) (*CommonTransfer, error) {
	event := new(CommonTransfer)
	if err := _Common.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommonTransferFinishIterator is returned from FilterTransferFinish and is used to iterate over the raw logs and unpacked data for TransferFinish events raised by the Common contract.
type CommonTransferFinishIterator struct {
	Event *CommonTransferFinish // Event containing the contract specifics and raw log

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
func (it *CommonTransferFinishIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommonTransferFinish)
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
		it.Event = new(CommonTransferFinish)
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
func (it *CommonTransferFinishIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommonTransferFinishIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommonTransferFinish represents a TransferFinish event raised by the Common contract.
type CommonTransferFinish struct {
	EventId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransferFinish is a free log retrieval operation binding the contract event 0x78ff9b3176bb0d6421590f3816e75cb15a9ffa2b7a1a028f40a3f4e029eb39f2.
//
// Solidity: event TransferFinish(uint256 indexed event_id)
func (_Common *CommonFilterer) FilterTransferFinish(opts *bind.FilterOpts, event_id []*big.Int) (*CommonTransferFinishIterator, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Common.contract.FilterLogs(opts, "TransferFinish", event_idRule)
	if err != nil {
		return nil, err
	}
	return &CommonTransferFinishIterator{contract: _Common.contract, event: "TransferFinish", logs: logs, sub: sub}, nil
}

// WatchTransferFinish is a free log subscription operation binding the contract event 0x78ff9b3176bb0d6421590f3816e75cb15a9ffa2b7a1a028f40a3f4e029eb39f2.
//
// Solidity: event TransferFinish(uint256 indexed event_id)
func (_Common *CommonFilterer) WatchTransferFinish(opts *bind.WatchOpts, sink chan<- *CommonTransferFinish, event_id []*big.Int) (event.Subscription, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Common.contract.WatchLogs(opts, "TransferFinish", event_idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommonTransferFinish)
				if err := _Common.contract.UnpackLog(event, "TransferFinish", log); err != nil {
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
// Solidity: event TransferFinish(uint256 indexed event_id)
func (_Common *CommonFilterer) ParseTransferFinish(log types.Log) (*CommonTransferFinish, error) {
	event := new(CommonTransferFinish)
	if err := _Common.contract.UnpackLog(event, "TransferFinish", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommonTransferSubmitIterator is returned from FilterTransferSubmit and is used to iterate over the raw logs and unpacked data for TransferSubmit events raised by the Common contract.
type CommonTransferSubmitIterator struct {
	Event *CommonTransferSubmit // Event containing the contract specifics and raw log

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
func (it *CommonTransferSubmitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommonTransferSubmit)
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
		it.Event = new(CommonTransferSubmit)
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
func (it *CommonTransferSubmitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommonTransferSubmitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommonTransferSubmit represents a TransferSubmit event raised by the Common contract.
type CommonTransferSubmit struct {
	EventId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransferSubmit is a free log retrieval operation binding the contract event 0x196c47048e38df7a4fe6e581c8f4f2e2ba67ac0dd45b90da756e97bd61d9dd3b.
//
// Solidity: event TransferSubmit(uint256 indexed event_id)
func (_Common *CommonFilterer) FilterTransferSubmit(opts *bind.FilterOpts, event_id []*big.Int) (*CommonTransferSubmitIterator, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Common.contract.FilterLogs(opts, "TransferSubmit", event_idRule)
	if err != nil {
		return nil, err
	}
	return &CommonTransferSubmitIterator{contract: _Common.contract, event: "TransferSubmit", logs: logs, sub: sub}, nil
}

// WatchTransferSubmit is a free log subscription operation binding the contract event 0x196c47048e38df7a4fe6e581c8f4f2e2ba67ac0dd45b90da756e97bd61d9dd3b.
//
// Solidity: event TransferSubmit(uint256 indexed event_id)
func (_Common *CommonFilterer) WatchTransferSubmit(opts *bind.WatchOpts, sink chan<- *CommonTransferSubmit, event_id []*big.Int) (event.Subscription, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Common.contract.WatchLogs(opts, "TransferSubmit", event_idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommonTransferSubmit)
				if err := _Common.contract.UnpackLog(event, "TransferSubmit", log); err != nil {
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
// Solidity: event TransferSubmit(uint256 indexed event_id)
func (_Common *CommonFilterer) ParseTransferSubmit(log types.Log) (*CommonTransferSubmit, error) {
	event := new(CommonTransferSubmit)
	if err := _Common.contract.UnpackLog(event, "TransferSubmit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommonUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Common contract.
type CommonUnpausedIterator struct {
	Event *CommonUnpaused // Event containing the contract specifics and raw log

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
func (it *CommonUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommonUnpaused)
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
		it.Event = new(CommonUnpaused)
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
func (it *CommonUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommonUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommonUnpaused represents a Unpaused event raised by the Common contract.
type CommonUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Common *CommonFilterer) FilterUnpaused(opts *bind.FilterOpts) (*CommonUnpausedIterator, error) {

	logs, sub, err := _Common.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &CommonUnpausedIterator{contract: _Common.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Common *CommonFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *CommonUnpaused) (event.Subscription, error) {

	logs, sub, err := _Common.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommonUnpaused)
				if err := _Common.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Common *CommonFilterer) ParseUnpaused(log types.Log) (*CommonUnpaused, error) {
	event := new(CommonUnpaused)
	if err := _Common.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CommonWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Common contract.
type CommonWithdrawIterator struct {
	Event *CommonWithdraw // Event containing the contract specifics and raw log

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
func (it *CommonWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CommonWithdraw)
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
		it.Event = new(CommonWithdraw)
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
func (it *CommonWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CommonWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CommonWithdraw represents a Withdraw event raised by the Common contract.
type CommonWithdraw struct {
	From      common.Address
	EventId   *big.Int
	FeeAmount *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0xf279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b568.
//
// Solidity: event Withdraw(address indexed from, uint256 event_id, uint256 feeAmount)
func (_Common *CommonFilterer) FilterWithdraw(opts *bind.FilterOpts, from []common.Address) (*CommonWithdrawIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Common.contract.FilterLogs(opts, "Withdraw", fromRule)
	if err != nil {
		return nil, err
	}
	return &CommonWithdrawIterator{contract: _Common.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0xf279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b568.
//
// Solidity: event Withdraw(address indexed from, uint256 event_id, uint256 feeAmount)
func (_Common *CommonFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *CommonWithdraw, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Common.contract.WatchLogs(opts, "Withdraw", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CommonWithdraw)
				if err := _Common.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0xf279e6a1f5e320cca91135676d9cb6e44ca8a08c0b88342bcdb1144f6511b568.
//
// Solidity: event Withdraw(address indexed from, uint256 event_id, uint256 feeAmount)
func (_Common *CommonFilterer) ParseWithdraw(log types.Log) (*CommonWithdraw, error) {
	event := new(CommonWithdraw)
	if err := _Common.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
