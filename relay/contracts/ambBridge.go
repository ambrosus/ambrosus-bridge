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

// AmbMetaData contains all meta data concerning the Amb contract.
var AmbMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sideBridgeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"relayAddress\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenSideAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeframeSeconds_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockTime_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSafetyBlocks_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"error\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"errorInfo\",\"type\":\"uint256\"}],\"name\":\"SetEpochData\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"b\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"c\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"d\",\"type\":\"uint256\"}],\"name\":\"Test\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"queue\",\"type\":\"tuple[]\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"proof\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes32\",\"name\":\"el\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"proofStart\",\"type\":\"uint256\"}],\"name\":\"CalcReceiptsHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"p\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"eventContractAddress\",\"type\":\"address\"}],\"name\":\"CalcTransferReceiptsHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parentOrReceiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"difficulty\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"nonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p6\",\"type\":\"bytes\"},{\"internalType\":\"uint256[]\",\"name\":\"dataSetLookup\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"witnessForLookup\",\"type\":\"uint256[]\"}],\"internalType\":\"structCheckPoW.BlockPoW[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"}],\"internalType\":\"structCheckPoW.PoWProof\",\"name\":\"powProof\",\"type\":\"tuple\"}],\"name\":\"CheckPoW_\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RELAY_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parentOrReceiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"difficulty\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"nonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p6\",\"type\":\"bytes\"},{\"internalType\":\"uint256[]\",\"name\":\"dataSetLookup\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"witnessForLookup\",\"type\":\"uint256[]\"}],\"internalType\":\"structCheckPoW.BlockPoW\",\"name\":\"block\",\"type\":\"tuple\"}],\"name\":\"blockHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"}],\"name\":\"bytesToUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"}],\"name\":\"changeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lockTime_\",\"type\":\"uint256\"}],\"name\":\"changeLockTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minSafetyBlocks_\",\"type\":\"uint256\"}],\"name\":\"changeMinSafetyBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timeframeSeconds_\",\"type\":\"uint256\"}],\"name\":\"changeTimeframeSeconds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAmbAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"transferEvent\",\"type\":\"bool\"}],\"name\":\"emitTestEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inputEventId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epochIndex\",\"type\":\"uint256\"}],\"name\":\"isEpochDataSet\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lockedTransfers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minSafetyBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"epoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fullSizeIn128Resultion\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"branchDepth\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"merkleNodes\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numElems\",\"type\":\"uint256\"}],\"name\":\"setEpochData\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sideBridgeAddress\",\"type\":\"address\"}],\"name\":\"setSideBridge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sideBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parentOrReceiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"difficulty\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"nonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p6\",\"type\":\"bytes\"},{\"internalType\":\"uint256[]\",\"name\":\"dataSetLookup\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"witnessForLookup\",\"type\":\"uint256[]\"}],\"internalType\":\"structCheckPoW.BlockPoW[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"}],\"internalType\":\"structCheckPoW.PoWProof\",\"name\":\"powProof\",\"type\":\"tuple\"}],\"name\":\"submitTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeframeSeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenSideAddress\",\"type\":\"address\"}],\"name\":\"tokensAdd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenSideAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensAddBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"}],\"name\":\"tokensRemove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensRemoveBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"unlockTransfers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parentOrReceiptHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"difficulty\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"number\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p4\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p5\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"nonce\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p6\",\"type\":\"bytes\"},{\"internalType\":\"uint256[]\",\"name\":\"dataSetLookup\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"witnessForLookup\",\"type\":\"uint256[]\"}],\"internalType\":\"structCheckPoW.BlockPoW\",\"name\":\"block\",\"type\":\"tuple\"}],\"name\":\"verifyEthash\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAmbAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// AmbABI is the input ABI used to generate the binding from.
// Deprecated: Use AmbMetaData.ABI instead.
var AmbABI = AmbMetaData.ABI

// Amb is an auto generated Go binding around an Ethereum contract.
type Amb struct {
	AmbCaller     // Read-only binding to the contract
	AmbTransactor // Write-only binding to the contract
	AmbFilterer   // Log filterer for contract events
}

// AmbCaller is an auto generated read-only Go binding around an Ethereum contract.
type AmbCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AmbTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AmbTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AmbFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AmbFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AmbSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AmbSession struct {
	Contract     *Amb              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AmbCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AmbCallerSession struct {
	Contract *AmbCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AmbTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AmbTransactorSession struct {
	Contract     *AmbTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AmbRaw is an auto generated low-level Go binding around an Ethereum contract.
type AmbRaw struct {
	Contract *Amb // Generic contract binding to access the raw methods on
}

// AmbCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AmbCallerRaw struct {
	Contract *AmbCaller // Generic read-only contract binding to access the raw methods on
}

// AmbTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AmbTransactorRaw struct {
	Contract *AmbTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAmb creates a new instance of Amb, bound to a specific deployed contract.
func NewAmb(address common.Address, backend bind.ContractBackend) (*Amb, error) {
	contract, err := bindAmb(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Amb{AmbCaller: AmbCaller{contract: contract}, AmbTransactor: AmbTransactor{contract: contract}, AmbFilterer: AmbFilterer{contract: contract}}, nil
}

// NewAmbCaller creates a new read-only instance of Amb, bound to a specific deployed contract.
func NewAmbCaller(address common.Address, caller bind.ContractCaller) (*AmbCaller, error) {
	contract, err := bindAmb(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AmbCaller{contract: contract}, nil
}

// NewAmbTransactor creates a new write-only instance of Amb, bound to a specific deployed contract.
func NewAmbTransactor(address common.Address, transactor bind.ContractTransactor) (*AmbTransactor, error) {
	contract, err := bindAmb(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AmbTransactor{contract: contract}, nil
}

// NewAmbFilterer creates a new log filterer instance of Amb, bound to a specific deployed contract.
func NewAmbFilterer(address common.Address, filterer bind.ContractFilterer) (*AmbFilterer, error) {
	contract, err := bindAmb(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AmbFilterer{contract: contract}, nil
}

// bindAmb binds a generic wrapper to an already deployed contract.
func bindAmb(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AmbABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Amb *AmbRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Amb.Contract.AmbCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Amb *AmbRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Amb.Contract.AmbTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Amb *AmbRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Amb.Contract.AmbTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Amb *AmbCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Amb.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Amb *AmbTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Amb.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Amb *AmbTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Amb.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Amb *AmbCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Amb *AmbSession) ADMINROLE() ([32]byte, error) {
	return _Amb.Contract.ADMINROLE(&_Amb.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Amb *AmbCallerSession) ADMINROLE() ([32]byte, error) {
	return _Amb.Contract.ADMINROLE(&_Amb.CallOpts)
}

// CalcReceiptsHash is a free data retrieval call binding the contract method 0xe7899536.
//
// Solidity: function CalcReceiptsHash(bytes[] proof, bytes32 el, uint256 proofStart) view returns(bytes32)
func (_Amb *AmbCaller) CalcReceiptsHash(opts *bind.CallOpts, proof [][]byte, el [32]byte, proofStart *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "CalcReceiptsHash", proof, el, proofStart)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalcReceiptsHash is a free data retrieval call binding the contract method 0xe7899536.
//
// Solidity: function CalcReceiptsHash(bytes[] proof, bytes32 el, uint256 proofStart) view returns(bytes32)
func (_Amb *AmbSession) CalcReceiptsHash(proof [][]byte, el [32]byte, proofStart *big.Int) ([32]byte, error) {
	return _Amb.Contract.CalcReceiptsHash(&_Amb.CallOpts, proof, el, proofStart)
}

// CalcReceiptsHash is a free data retrieval call binding the contract method 0xe7899536.
//
// Solidity: function CalcReceiptsHash(bytes[] proof, bytes32 el, uint256 proofStart) view returns(bytes32)
func (_Amb *AmbCallerSession) CalcReceiptsHash(proof [][]byte, el [32]byte, proofStart *big.Int) ([32]byte, error) {
	return _Amb.Contract.CalcReceiptsHash(&_Amb.CallOpts, proof, el, proofStart)
}

// CalcTransferReceiptsHash is a free data retrieval call binding the contract method 0x90d0308f.
//
// Solidity: function CalcTransferReceiptsHash((bytes[],uint256,(address,address,uint256)[]) p, address eventContractAddress) view returns(bytes32)
func (_Amb *AmbCaller) CalcTransferReceiptsHash(opts *bind.CallOpts, p CommonStructsTransferProof, eventContractAddress common.Address) ([32]byte, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "CalcTransferReceiptsHash", p, eventContractAddress)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalcTransferReceiptsHash is a free data retrieval call binding the contract method 0x90d0308f.
//
// Solidity: function CalcTransferReceiptsHash((bytes[],uint256,(address,address,uint256)[]) p, address eventContractAddress) view returns(bytes32)
func (_Amb *AmbSession) CalcTransferReceiptsHash(p CommonStructsTransferProof, eventContractAddress common.Address) ([32]byte, error) {
	return _Amb.Contract.CalcTransferReceiptsHash(&_Amb.CallOpts, p, eventContractAddress)
}

// CalcTransferReceiptsHash is a free data retrieval call binding the contract method 0x90d0308f.
//
// Solidity: function CalcTransferReceiptsHash((bytes[],uint256,(address,address,uint256)[]) p, address eventContractAddress) view returns(bytes32)
func (_Amb *AmbCallerSession) CalcTransferReceiptsHash(p CommonStructsTransferProof, eventContractAddress common.Address) ([32]byte, error) {
	return _Amb.Contract.CalcTransferReceiptsHash(&_Amb.CallOpts, p, eventContractAddress)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Amb *AmbCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Amb *AmbSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Amb.Contract.DEFAULTADMINROLE(&_Amb.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Amb *AmbCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Amb.Contract.DEFAULTADMINROLE(&_Amb.CallOpts)
}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Amb *AmbCaller) RELAYROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "RELAY_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Amb *AmbSession) RELAYROLE() ([32]byte, error) {
	return _Amb.Contract.RELAYROLE(&_Amb.CallOpts)
}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Amb *AmbCallerSession) RELAYROLE() ([32]byte, error) {
	return _Amb.Contract.RELAYROLE(&_Amb.CallOpts)
}

// BlockHash is a free data retrieval call binding the contract method 0x7d55a274.
//
// Solidity: function blockHash((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[]) block) view returns(bytes32)
func (_Amb *AmbCaller) BlockHash(opts *bind.CallOpts, block CheckPoWBlockPoW) ([32]byte, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "blockHash", block)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BlockHash is a free data retrieval call binding the contract method 0x7d55a274.
//
// Solidity: function blockHash((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[]) block) view returns(bytes32)
func (_Amb *AmbSession) BlockHash(block CheckPoWBlockPoW) ([32]byte, error) {
	return _Amb.Contract.BlockHash(&_Amb.CallOpts, block)
}

// BlockHash is a free data retrieval call binding the contract method 0x7d55a274.
//
// Solidity: function blockHash((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[]) block) view returns(bytes32)
func (_Amb *AmbCallerSession) BlockHash(block CheckPoWBlockPoW) ([32]byte, error) {
	return _Amb.Contract.BlockHash(&_Amb.CallOpts, block)
}

// BytesToUint is a free data retrieval call binding the contract method 0x02d06d05.
//
// Solidity: function bytesToUint(bytes b) view returns(uint256)
func (_Amb *AmbCaller) BytesToUint(opts *bind.CallOpts, b []byte) (*big.Int, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "bytesToUint", b)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BytesToUint is a free data retrieval call binding the contract method 0x02d06d05.
//
// Solidity: function bytesToUint(bytes b) view returns(uint256)
func (_Amb *AmbSession) BytesToUint(b []byte) (*big.Int, error) {
	return _Amb.Contract.BytesToUint(&_Amb.CallOpts, b)
}

// BytesToUint is a free data retrieval call binding the contract method 0x02d06d05.
//
// Solidity: function bytesToUint(bytes b) view returns(uint256)
func (_Amb *AmbCallerSession) BytesToUint(b []byte) (*big.Int, error) {
	return _Amb.Contract.BytesToUint(&_Amb.CallOpts, b)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Amb *AmbCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Amb *AmbSession) Fee() (*big.Int, error) {
	return _Amb.Contract.Fee(&_Amb.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Amb *AmbCallerSession) Fee() (*big.Int, error) {
	return _Amb.Contract.Fee(&_Amb.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Amb *AmbCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Amb *AmbSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Amb.Contract.GetRoleAdmin(&_Amb.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Amb *AmbCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Amb.Contract.GetRoleAdmin(&_Amb.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Amb *AmbCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Amb *AmbSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Amb.Contract.HasRole(&_Amb.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Amb *AmbCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Amb.Contract.HasRole(&_Amb.CallOpts, role, account)
}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Amb *AmbCaller) InputEventId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "inputEventId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Amb *AmbSession) InputEventId() (*big.Int, error) {
	return _Amb.Contract.InputEventId(&_Amb.CallOpts)
}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Amb *AmbCallerSession) InputEventId() (*big.Int, error) {
	return _Amb.Contract.InputEventId(&_Amb.CallOpts)
}

// IsEpochDataSet is a free data retrieval call binding the contract method 0xc7b81f4f.
//
// Solidity: function isEpochDataSet(uint256 epochIndex) view returns(bool)
func (_Amb *AmbCaller) IsEpochDataSet(opts *bind.CallOpts, epochIndex *big.Int) (bool, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "isEpochDataSet", epochIndex)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEpochDataSet is a free data retrieval call binding the contract method 0xc7b81f4f.
//
// Solidity: function isEpochDataSet(uint256 epochIndex) view returns(bool)
func (_Amb *AmbSession) IsEpochDataSet(epochIndex *big.Int) (bool, error) {
	return _Amb.Contract.IsEpochDataSet(&_Amb.CallOpts, epochIndex)
}

// IsEpochDataSet is a free data retrieval call binding the contract method 0xc7b81f4f.
//
// Solidity: function isEpochDataSet(uint256 epochIndex) view returns(bool)
func (_Amb *AmbCallerSession) IsEpochDataSet(epochIndex *big.Int) (bool, error) {
	return _Amb.Contract.IsEpochDataSet(&_Amb.CallOpts, epochIndex)
}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Amb *AmbCaller) LockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "lockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Amb *AmbSession) LockTime() (*big.Int, error) {
	return _Amb.Contract.LockTime(&_Amb.CallOpts)
}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Amb *AmbCallerSession) LockTime() (*big.Int, error) {
	return _Amb.Contract.LockTime(&_Amb.CallOpts)
}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Amb *AmbCaller) LockedTransfers(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "lockedTransfers", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Amb *AmbSession) LockedTransfers(arg0 *big.Int) (*big.Int, error) {
	return _Amb.Contract.LockedTransfers(&_Amb.CallOpts, arg0)
}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Amb *AmbCallerSession) LockedTransfers(arg0 *big.Int) (*big.Int, error) {
	return _Amb.Contract.LockedTransfers(&_Amb.CallOpts, arg0)
}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Amb *AmbCaller) MinSafetyBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "minSafetyBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Amb *AmbSession) MinSafetyBlocks() (*big.Int, error) {
	return _Amb.Contract.MinSafetyBlocks(&_Amb.CallOpts)
}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Amb *AmbCallerSession) MinSafetyBlocks() (*big.Int, error) {
	return _Amb.Contract.MinSafetyBlocks(&_Amb.CallOpts)
}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Amb *AmbCaller) SideBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "sideBridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Amb *AmbSession) SideBridgeAddress() (common.Address, error) {
	return _Amb.Contract.SideBridgeAddress(&_Amb.CallOpts)
}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Amb *AmbCallerSession) SideBridgeAddress() (common.Address, error) {
	return _Amb.Contract.SideBridgeAddress(&_Amb.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Amb *AmbCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Amb *AmbSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Amb.Contract.SupportsInterface(&_Amb.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Amb *AmbCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Amb.Contract.SupportsInterface(&_Amb.CallOpts, interfaceId)
}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Amb *AmbCaller) TimeframeSeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "timeframeSeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Amb *AmbSession) TimeframeSeconds() (*big.Int, error) {
	return _Amb.Contract.TimeframeSeconds(&_Amb.CallOpts)
}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Amb *AmbCallerSession) TimeframeSeconds() (*big.Int, error) {
	return _Amb.Contract.TimeframeSeconds(&_Amb.CallOpts)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Amb *AmbCaller) TokenAddresses(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "tokenAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Amb *AmbSession) TokenAddresses(arg0 common.Address) (common.Address, error) {
	return _Amb.Contract.TokenAddresses(&_Amb.CallOpts, arg0)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Amb *AmbCallerSession) TokenAddresses(arg0 common.Address) (common.Address, error) {
	return _Amb.Contract.TokenAddresses(&_Amb.CallOpts, arg0)
}

// VerifyEthash is a free data retrieval call binding the contract method 0xf9e14a52.
//
// Solidity: function verifyEthash((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[]) block) view returns()
func (_Amb *AmbCaller) VerifyEthash(opts *bind.CallOpts, block CheckPoWBlockPoW) error {
	var out []interface{}
	err := _Amb.contract.Call(opts, &out, "verifyEthash", block)

	if err != nil {
		return err
	}

	return err

}

// VerifyEthash is a free data retrieval call binding the contract method 0xf9e14a52.
//
// Solidity: function verifyEthash((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[]) block) view returns()
func (_Amb *AmbSession) VerifyEthash(block CheckPoWBlockPoW) error {
	return _Amb.Contract.VerifyEthash(&_Amb.CallOpts, block)
}

// VerifyEthash is a free data retrieval call binding the contract method 0xf9e14a52.
//
// Solidity: function verifyEthash((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[]) block) view returns()
func (_Amb *AmbCallerSession) VerifyEthash(block CheckPoWBlockPoW) error {
	return _Amb.Contract.VerifyEthash(&_Amb.CallOpts, block)
}

// CheckPoW is a paid mutator transaction binding the contract method 0xc2ddeca9.
//
// Solidity: function CheckPoW_(((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Amb *AmbTransactor) CheckPoW(opts *bind.TransactOpts, powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "CheckPoW_", powProof)
}

// CheckPoW is a paid mutator transaction binding the contract method 0xc2ddeca9.
//
// Solidity: function CheckPoW_(((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Amb *AmbSession) CheckPoW(powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Amb.Contract.CheckPoW(&_Amb.TransactOpts, powProof)
}

// CheckPoW is a paid mutator transaction binding the contract method 0xc2ddeca9.
//
// Solidity: function CheckPoW_(((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Amb *AmbTransactorSession) CheckPoW(powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Amb.Contract.CheckPoW(&_Amb.TransactOpts, powProof)
}

// ChangeFee is a paid mutator transaction binding the contract method 0x6a1db1bf.
//
// Solidity: function changeFee(uint256 fee_) returns()
func (_Amb *AmbTransactor) ChangeFee(opts *bind.TransactOpts, fee_ *big.Int) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "changeFee", fee_)
}

// ChangeFee is a paid mutator transaction binding the contract method 0x6a1db1bf.
//
// Solidity: function changeFee(uint256 fee_) returns()
func (_Amb *AmbSession) ChangeFee(fee_ *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.ChangeFee(&_Amb.TransactOpts, fee_)
}

// ChangeFee is a paid mutator transaction binding the contract method 0x6a1db1bf.
//
// Solidity: function changeFee(uint256 fee_) returns()
func (_Amb *AmbTransactorSession) ChangeFee(fee_ *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.ChangeFee(&_Amb.TransactOpts, fee_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Amb *AmbTransactor) ChangeLockTime(opts *bind.TransactOpts, lockTime_ *big.Int) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "changeLockTime", lockTime_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Amb *AmbSession) ChangeLockTime(lockTime_ *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.ChangeLockTime(&_Amb.TransactOpts, lockTime_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Amb *AmbTransactorSession) ChangeLockTime(lockTime_ *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.ChangeLockTime(&_Amb.TransactOpts, lockTime_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Amb *AmbTransactor) ChangeMinSafetyBlocks(opts *bind.TransactOpts, minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "changeMinSafetyBlocks", minSafetyBlocks_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Amb *AmbSession) ChangeMinSafetyBlocks(minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.ChangeMinSafetyBlocks(&_Amb.TransactOpts, minSafetyBlocks_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Amb *AmbTransactorSession) ChangeMinSafetyBlocks(minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.ChangeMinSafetyBlocks(&_Amb.TransactOpts, minSafetyBlocks_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Amb *AmbTransactor) ChangeTimeframeSeconds(opts *bind.TransactOpts, timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "changeTimeframeSeconds", timeframeSeconds_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Amb *AmbSession) ChangeTimeframeSeconds(timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.ChangeTimeframeSeconds(&_Amb.TransactOpts, timeframeSeconds_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Amb *AmbTransactorSession) ChangeTimeframeSeconds(timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.ChangeTimeframeSeconds(&_Amb.TransactOpts, timeframeSeconds_)
}

// EmitTestEvent is a paid mutator transaction binding the contract method 0xcbe31e75.
//
// Solidity: function emitTestEvent(address tokenAmbAddress, address toAddress, uint256 amount, bool transferEvent) returns()
func (_Amb *AmbTransactor) EmitTestEvent(opts *bind.TransactOpts, tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int, transferEvent bool) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "emitTestEvent", tokenAmbAddress, toAddress, amount, transferEvent)
}

// EmitTestEvent is a paid mutator transaction binding the contract method 0xcbe31e75.
//
// Solidity: function emitTestEvent(address tokenAmbAddress, address toAddress, uint256 amount, bool transferEvent) returns()
func (_Amb *AmbSession) EmitTestEvent(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int, transferEvent bool) (*types.Transaction, error) {
	return _Amb.Contract.EmitTestEvent(&_Amb.TransactOpts, tokenAmbAddress, toAddress, amount, transferEvent)
}

// EmitTestEvent is a paid mutator transaction binding the contract method 0xcbe31e75.
//
// Solidity: function emitTestEvent(address tokenAmbAddress, address toAddress, uint256 amount, bool transferEvent) returns()
func (_Amb *AmbTransactorSession) EmitTestEvent(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int, transferEvent bool) (*types.Transaction, error) {
	return _Amb.Contract.EmitTestEvent(&_Amb.TransactOpts, tokenAmbAddress, toAddress, amount, transferEvent)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Amb *AmbTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Amb *AmbSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Amb.Contract.GrantRole(&_Amb.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Amb *AmbTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Amb.Contract.GrantRole(&_Amb.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Amb *AmbTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Amb *AmbSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Amb.Contract.RenounceRole(&_Amb.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Amb *AmbTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Amb.Contract.RenounceRole(&_Amb.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Amb *AmbTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Amb *AmbSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Amb.Contract.RevokeRole(&_Amb.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Amb *AmbTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Amb.Contract.RevokeRole(&_Amb.TransactOpts, role, account)
}

// SetEpochData is a paid mutator transaction binding the contract method 0xc891a29d.
//
// Solidity: function setEpochData(uint256 epoch, uint256 fullSizeIn128Resultion, uint256 branchDepth, uint256[] merkleNodes, uint256 start, uint256 numElems) returns()
func (_Amb *AmbTransactor) SetEpochData(opts *bind.TransactOpts, epoch *big.Int, fullSizeIn128Resultion *big.Int, branchDepth *big.Int, merkleNodes []*big.Int, start *big.Int, numElems *big.Int) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "setEpochData", epoch, fullSizeIn128Resultion, branchDepth, merkleNodes, start, numElems)
}

// SetEpochData is a paid mutator transaction binding the contract method 0xc891a29d.
//
// Solidity: function setEpochData(uint256 epoch, uint256 fullSizeIn128Resultion, uint256 branchDepth, uint256[] merkleNodes, uint256 start, uint256 numElems) returns()
func (_Amb *AmbSession) SetEpochData(epoch *big.Int, fullSizeIn128Resultion *big.Int, branchDepth *big.Int, merkleNodes []*big.Int, start *big.Int, numElems *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.SetEpochData(&_Amb.TransactOpts, epoch, fullSizeIn128Resultion, branchDepth, merkleNodes, start, numElems)
}

// SetEpochData is a paid mutator transaction binding the contract method 0xc891a29d.
//
// Solidity: function setEpochData(uint256 epoch, uint256 fullSizeIn128Resultion, uint256 branchDepth, uint256[] merkleNodes, uint256 start, uint256 numElems) returns()
func (_Amb *AmbTransactorSession) SetEpochData(epoch *big.Int, fullSizeIn128Resultion *big.Int, branchDepth *big.Int, merkleNodes []*big.Int, start *big.Int, numElems *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.SetEpochData(&_Amb.TransactOpts, epoch, fullSizeIn128Resultion, branchDepth, merkleNodes, start, numElems)
}

// SetSideBridge is a paid mutator transaction binding the contract method 0x21d3d536.
//
// Solidity: function setSideBridge(address _sideBridgeAddress) returns()
func (_Amb *AmbTransactor) SetSideBridge(opts *bind.TransactOpts, _sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "setSideBridge", _sideBridgeAddress)
}

// SetSideBridge is a paid mutator transaction binding the contract method 0x21d3d536.
//
// Solidity: function setSideBridge(address _sideBridgeAddress) returns()
func (_Amb *AmbSession) SetSideBridge(_sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Amb.Contract.SetSideBridge(&_Amb.TransactOpts, _sideBridgeAddress)
}

// SetSideBridge is a paid mutator transaction binding the contract method 0x21d3d536.
//
// Solidity: function setSideBridge(address _sideBridgeAddress) returns()
func (_Amb *AmbTransactorSession) SetSideBridge(_sideBridgeAddress common.Address) (*types.Transaction, error) {
	return _Amb.Contract.SetSideBridge(&_Amb.TransactOpts, _sideBridgeAddress)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0xc248a913.
//
// Solidity: function submitTransfer(((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Amb *AmbTransactor) SubmitTransfer(opts *bind.TransactOpts, powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "submitTransfer", powProof)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0xc248a913.
//
// Solidity: function submitTransfer(((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Amb *AmbSession) SubmitTransfer(powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Amb.Contract.SubmitTransfer(&_Amb.TransactOpts, powProof)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0xc248a913.
//
// Solidity: function submitTransfer(((bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes,bytes,bytes,uint256[],uint256[])[],(bytes[],uint256,(address,address,uint256)[])) powProof) returns()
func (_Amb *AmbTransactorSession) SubmitTransfer(powProof CheckPoWPoWProof) (*types.Transaction, error) {
	return _Amb.Contract.SubmitTransfer(&_Amb.TransactOpts, powProof)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Amb *AmbTransactor) TokensAdd(opts *bind.TransactOpts, tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "tokensAdd", tokenThisAddress, tokenSideAddress)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Amb *AmbSession) TokensAdd(tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Amb.Contract.TokensAdd(&_Amb.TransactOpts, tokenThisAddress, tokenSideAddress)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Amb *AmbTransactorSession) TokensAdd(tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Amb.Contract.TokensAdd(&_Amb.TransactOpts, tokenThisAddress, tokenSideAddress)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Amb *AmbTransactor) TokensAddBatch(opts *bind.TransactOpts, tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "tokensAddBatch", tokenThisAddresses, tokenSideAddresses)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Amb *AmbSession) TokensAddBatch(tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Amb.Contract.TokensAddBatch(&_Amb.TransactOpts, tokenThisAddresses, tokenSideAddresses)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Amb *AmbTransactorSession) TokensAddBatch(tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Amb.Contract.TokensAddBatch(&_Amb.TransactOpts, tokenThisAddresses, tokenSideAddresses)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Amb *AmbTransactor) TokensRemove(opts *bind.TransactOpts, tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "tokensRemove", tokenThisAddress)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Amb *AmbSession) TokensRemove(tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Amb.Contract.TokensRemove(&_Amb.TransactOpts, tokenThisAddress)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Amb *AmbTransactorSession) TokensRemove(tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Amb.Contract.TokensRemove(&_Amb.TransactOpts, tokenThisAddress)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Amb *AmbTransactor) TokensRemoveBatch(opts *bind.TransactOpts, tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "tokensRemoveBatch", tokenThisAddresses)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Amb *AmbSession) TokensRemoveBatch(tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Amb.Contract.TokensRemoveBatch(&_Amb.TransactOpts, tokenThisAddresses)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Amb *AmbTransactorSession) TokensRemoveBatch(tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Amb.Contract.TokensRemoveBatch(&_Amb.TransactOpts, tokenThisAddresses)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 event_id) returns()
func (_Amb *AmbTransactor) UnlockTransfers(opts *bind.TransactOpts, event_id *big.Int) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "unlockTransfers", event_id)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 event_id) returns()
func (_Amb *AmbSession) UnlockTransfers(event_id *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.UnlockTransfers(&_Amb.TransactOpts, event_id)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 event_id) returns()
func (_Amb *AmbTransactorSession) UnlockTransfers(event_id *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.UnlockTransfers(&_Amb.TransactOpts, event_id)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAmbAddress, address toAddress, uint256 amount) payable returns()
func (_Amb *AmbTransactor) Withdraw(opts *bind.TransactOpts, tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "withdraw", tokenAmbAddress, toAddress, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAmbAddress, address toAddress, uint256 amount) payable returns()
func (_Amb *AmbSession) Withdraw(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.Withdraw(&_Amb.TransactOpts, tokenAmbAddress, toAddress, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAmbAddress, address toAddress, uint256 amount) payable returns()
func (_Amb *AmbTransactorSession) Withdraw(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.Withdraw(&_Amb.TransactOpts, tokenAmbAddress, toAddress, amount)
}

// AmbRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Amb contract.
type AmbRoleAdminChangedIterator struct {
	Event *AmbRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *AmbRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AmbRoleAdminChanged)
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
		it.Event = new(AmbRoleAdminChanged)
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
func (it *AmbRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AmbRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AmbRoleAdminChanged represents a RoleAdminChanged event raised by the Amb contract.
type AmbRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Amb *AmbFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*AmbRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Amb.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &AmbRoleAdminChangedIterator{contract: _Amb.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Amb *AmbFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *AmbRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Amb.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AmbRoleAdminChanged)
				if err := _Amb.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Amb *AmbFilterer) ParseRoleAdminChanged(log types.Log) (*AmbRoleAdminChanged, error) {
	event := new(AmbRoleAdminChanged)
	if err := _Amb.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AmbRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Amb contract.
type AmbRoleGrantedIterator struct {
	Event *AmbRoleGranted // Event containing the contract specifics and raw log

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
func (it *AmbRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AmbRoleGranted)
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
		it.Event = new(AmbRoleGranted)
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
func (it *AmbRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AmbRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AmbRoleGranted represents a RoleGranted event raised by the Amb contract.
type AmbRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Amb *AmbFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AmbRoleGrantedIterator, error) {

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

	logs, sub, err := _Amb.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AmbRoleGrantedIterator{contract: _Amb.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Amb *AmbFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *AmbRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Amb.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AmbRoleGranted)
				if err := _Amb.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Amb *AmbFilterer) ParseRoleGranted(log types.Log) (*AmbRoleGranted, error) {
	event := new(AmbRoleGranted)
	if err := _Amb.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AmbRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Amb contract.
type AmbRoleRevokedIterator struct {
	Event *AmbRoleRevoked // Event containing the contract specifics and raw log

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
func (it *AmbRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AmbRoleRevoked)
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
		it.Event = new(AmbRoleRevoked)
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
func (it *AmbRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AmbRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AmbRoleRevoked represents a RoleRevoked event raised by the Amb contract.
type AmbRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Amb *AmbFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*AmbRoleRevokedIterator, error) {

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

	logs, sub, err := _Amb.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &AmbRoleRevokedIterator{contract: _Amb.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Amb *AmbFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *AmbRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Amb.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AmbRoleRevoked)
				if err := _Amb.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Amb *AmbFilterer) ParseRoleRevoked(log types.Log) (*AmbRoleRevoked, error) {
	event := new(AmbRoleRevoked)
	if err := _Amb.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AmbSetEpochDataIterator is returned from FilterSetEpochData and is used to iterate over the raw logs and unpacked data for SetEpochData events raised by the Amb contract.
type AmbSetEpochDataIterator struct {
	Event *AmbSetEpochData // Event containing the contract specifics and raw log

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
func (it *AmbSetEpochDataIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AmbSetEpochData)
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
		it.Event = new(AmbSetEpochData)
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
func (it *AmbSetEpochDataIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AmbSetEpochDataIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AmbSetEpochData represents a SetEpochData event raised by the Amb contract.
type AmbSetEpochData struct {
	Sender    common.Address
	Error     *big.Int
	ErrorInfo *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetEpochData is a free log retrieval operation binding the contract event 0x5cd723400be8430351b9cbaa5ea421b3fb2528c6a7650c493f895e7d97750da1.
//
// Solidity: event SetEpochData(address indexed sender, uint256 error, uint256 errorInfo)
func (_Amb *AmbFilterer) FilterSetEpochData(opts *bind.FilterOpts, sender []common.Address) (*AmbSetEpochDataIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Amb.contract.FilterLogs(opts, "SetEpochData", senderRule)
	if err != nil {
		return nil, err
	}
	return &AmbSetEpochDataIterator{contract: _Amb.contract, event: "SetEpochData", logs: logs, sub: sub}, nil
}

// WatchSetEpochData is a free log subscription operation binding the contract event 0x5cd723400be8430351b9cbaa5ea421b3fb2528c6a7650c493f895e7d97750da1.
//
// Solidity: event SetEpochData(address indexed sender, uint256 error, uint256 errorInfo)
func (_Amb *AmbFilterer) WatchSetEpochData(opts *bind.WatchOpts, sink chan<- *AmbSetEpochData, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Amb.contract.WatchLogs(opts, "SetEpochData", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AmbSetEpochData)
				if err := _Amb.contract.UnpackLog(event, "SetEpochData", log); err != nil {
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

// ParseSetEpochData is a log parse operation binding the contract event 0x5cd723400be8430351b9cbaa5ea421b3fb2528c6a7650c493f895e7d97750da1.
//
// Solidity: event SetEpochData(address indexed sender, uint256 error, uint256 errorInfo)
func (_Amb *AmbFilterer) ParseSetEpochData(log types.Log) (*AmbSetEpochData, error) {
	event := new(AmbSetEpochData)
	if err := _Amb.contract.UnpackLog(event, "SetEpochData", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AmbTestIterator is returned from FilterTest and is used to iterate over the raw logs and unpacked data for Test events raised by the Amb contract.
type AmbTestIterator struct {
	Event *AmbTest // Event containing the contract specifics and raw log

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
func (it *AmbTestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AmbTest)
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
		it.Event = new(AmbTest)
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
func (it *AmbTestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AmbTestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AmbTest represents a Test event raised by the Amb contract.
type AmbTest struct {
	A   *big.Int
	B   common.Address
	C   string
	D   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTest is a free log retrieval operation binding the contract event 0xf690cc04b54c126882a8927cb940bececfeabf9b0b4b0069db2427fbf9e86dd6.
//
// Solidity: event Test(uint256 indexed a, address indexed b, string c, uint256 d)
func (_Amb *AmbFilterer) FilterTest(opts *bind.FilterOpts, a []*big.Int, b []common.Address) (*AmbTestIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _Amb.contract.FilterLogs(opts, "Test", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return &AmbTestIterator{contract: _Amb.contract, event: "Test", logs: logs, sub: sub}, nil
}

// WatchTest is a free log subscription operation binding the contract event 0xf690cc04b54c126882a8927cb940bececfeabf9b0b4b0069db2427fbf9e86dd6.
//
// Solidity: event Test(uint256 indexed a, address indexed b, string c, uint256 d)
func (_Amb *AmbFilterer) WatchTest(opts *bind.WatchOpts, sink chan<- *AmbTest, a []*big.Int, b []common.Address) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _Amb.contract.WatchLogs(opts, "Test", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AmbTest)
				if err := _Amb.contract.UnpackLog(event, "Test", log); err != nil {
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

// ParseTest is a log parse operation binding the contract event 0xf690cc04b54c126882a8927cb940bececfeabf9b0b4b0069db2427fbf9e86dd6.
//
// Solidity: event Test(uint256 indexed a, address indexed b, string c, uint256 d)
func (_Amb *AmbFilterer) ParseTest(log types.Log) (*AmbTest, error) {
	event := new(AmbTest)
	if err := _Amb.contract.UnpackLog(event, "Test", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AmbTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Amb contract.
type AmbTransferIterator struct {
	Event *AmbTransfer // Event containing the contract specifics and raw log

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
func (it *AmbTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AmbTransfer)
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
		it.Event = new(AmbTransfer)
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
func (it *AmbTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AmbTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2.
//
// Solidity: event Transfer(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Amb *AmbFilterer) FilterTransfer(opts *bind.FilterOpts, event_id []*big.Int) (*AmbTransferIterator, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Amb.contract.FilterLogs(opts, "Transfer", event_idRule)
	if err != nil {
		return nil, err
	}
	return &AmbTransferIterator{contract: _Amb.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2.
//
// Solidity: event Transfer(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Amb *AmbFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *AmbTransfer, event_id []*big.Int) (event.Subscription, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Amb.contract.WatchLogs(opts, "Transfer", event_idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AmbTransfer)
				if err := _Amb.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_Amb *AmbFilterer) ParseTransfer(log types.Log) (*AmbTransfer, error) {
	event := new(AmbTransfer)
	if err := _Amb.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AmbWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Amb contract.
type AmbWithdrawIterator struct {
	Event *AmbWithdraw // Event containing the contract specifics and raw log

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
func (it *AmbWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AmbWithdraw)
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
		it.Event = new(AmbWithdraw)
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
func (it *AmbWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AmbWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AmbWithdraw represents a Withdraw event raised by the Amb contract.
type AmbWithdraw struct {
	From    common.Address
	EventId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed from, uint256 event_id)
func (_Amb *AmbFilterer) FilterWithdraw(opts *bind.FilterOpts, from []common.Address) (*AmbWithdrawIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Amb.contract.FilterLogs(opts, "Withdraw", fromRule)
	if err != nil {
		return nil, err
	}
	return &AmbWithdrawIterator{contract: _Amb.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed from, uint256 event_id)
func (_Amb *AmbFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *AmbWithdraw, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Amb.contract.WatchLogs(opts, "Withdraw", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AmbWithdraw)
				if err := _Amb.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed from, uint256 event_id)
func (_Amb *AmbFilterer) ParseWithdraw(log types.Log) (*AmbWithdraw, error) {
	event := new(AmbWithdraw)
	if err := _Amb.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
