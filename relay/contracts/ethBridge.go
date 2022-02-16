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

// EthMetaData contains all meta data concerning the Eth contract.
var EthMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sideBridgeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"relayAddress\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenSideAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeframeSeconds_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lockTime_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minSafetyBlocks_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"a\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"b\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"c\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"d\",\"type\":\"uint256\"}],\"name\":\"Test\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"queue\",\"type\":\"tuple[]\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p0_seal\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p0_bare\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parent_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"receipts_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s1\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"step\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"type_\",\"type\":\"int256\"}],\"internalType\":\"structCheckAura.BlockAura[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"address\",\"name\":\"delta_address\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"delta_index\",\"type\":\"int64\"}],\"internalType\":\"structCheckAura.ValidatorSetProof[]\",\"name\":\"vs_changes\",\"type\":\"tuple[]\"}],\"internalType\":\"structCheckAura.AuraProof\",\"name\":\"auraProof\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"minSafetyBlocks\",\"type\":\"uint256\"}],\"name\":\"CheckAura_\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RELAY_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"}],\"name\":\"bytesToUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"}],\"name\":\"changeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"lockTime_\",\"type\":\"uint256\"}],\"name\":\"changeLockTime\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"minSafetyBlocks_\",\"type\":\"uint256\"}],\"name\":\"changeMinSafetyBlocks\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"timeframeSeconds_\",\"type\":\"uint256\"}],\"name\":\"changeTimeframeSeconds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAmbAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"emitTestEvent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inputEventId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lockTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"lockedTransfers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"endTimestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minSafetyBlocks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sideBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p0_seal\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p0_bare\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"parent_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"receipts_hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s1\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"step\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"type_\",\"type\":\"int256\"}],\"internalType\":\"structCheckAura.BlockAura[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"transfers\",\"type\":\"tuple[]\"}],\"internalType\":\"structCommonStructs.TransferProof\",\"name\":\"transfer\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes[]\",\"name\":\"receipt_proof\",\"type\":\"bytes[]\"},{\"internalType\":\"address\",\"name\":\"delta_address\",\"type\":\"address\"},{\"internalType\":\"int64\",\"name\":\"delta_index\",\"type\":\"int64\"}],\"internalType\":\"structCheckAura.ValidatorSetProof[]\",\"name\":\"vs_changes\",\"type\":\"tuple[]\"}],\"internalType\":\"structCheckAura.AuraProof\",\"name\":\"auraProof\",\"type\":\"tuple\"}],\"name\":\"submitTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"timeframeSeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenSideAddress\",\"type\":\"address\"}],\"name\":\"tokensAdd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenSideAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensAddBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"}],\"name\":\"tokensRemove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensRemoveBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"unlockTransfers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAmbAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// EthABI is the input ABI used to generate the binding from.
// Deprecated: Use EthMetaData.ABI instead.
var EthABI = EthMetaData.ABI

// Eth is an auto generated Go binding around an Ethereum contract.
type Eth struct {
	EthCaller     // Read-only binding to the contract
	EthTransactor // Write-only binding to the contract
	EthFilterer   // Log filterer for contract events
}

// EthCaller is an auto generated read-only Go binding around an Ethereum contract.
type EthCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EthTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EthFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EthSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EthSession struct {
	Contract     *Eth              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EthCallerSession struct {
	Contract *EthCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// EthTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EthTransactorSession struct {
	Contract     *EthTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EthRaw is an auto generated low-level Go binding around an Ethereum contract.
type EthRaw struct {
	Contract *Eth // Generic contract binding to access the raw methods on
}

// EthCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EthCallerRaw struct {
	Contract *EthCaller // Generic read-only contract binding to access the raw methods on
}

// EthTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EthTransactorRaw struct {
	Contract *EthTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEth creates a new instance of Eth, bound to a specific deployed contract.
func NewEth(address common.Address, backend bind.ContractBackend) (*Eth, error) {
	contract, err := bindEth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Eth{EthCaller: EthCaller{contract: contract}, EthTransactor: EthTransactor{contract: contract}, EthFilterer: EthFilterer{contract: contract}}, nil
}

// NewEthCaller creates a new read-only instance of Eth, bound to a specific deployed contract.
func NewEthCaller(address common.Address, caller bind.ContractCaller) (*EthCaller, error) {
	contract, err := bindEth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EthCaller{contract: contract}, nil
}

// NewEthTransactor creates a new write-only instance of Eth, bound to a specific deployed contract.
func NewEthTransactor(address common.Address, transactor bind.ContractTransactor) (*EthTransactor, error) {
	contract, err := bindEth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EthTransactor{contract: contract}, nil
}

// NewEthFilterer creates a new log filterer instance of Eth, bound to a specific deployed contract.
func NewEthFilterer(address common.Address, filterer bind.ContractFilterer) (*EthFilterer, error) {
	contract, err := bindEth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EthFilterer{contract: contract}, nil
}

// bindEth binds a generic wrapper to an already deployed contract.
func bindEth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EthABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eth *EthRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Eth.Contract.EthCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eth *EthRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eth.Contract.EthTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eth *EthRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eth.Contract.EthTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eth *EthCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Eth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eth *EthTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eth *EthTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eth.Contract.contract.Transact(opts, method, params...)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Eth *EthCaller) ADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Eth *EthSession) ADMINROLE() ([32]byte, error) {
	return _Eth.Contract.ADMINROLE(&_Eth.CallOpts)
}

// ADMINROLE is a free data retrieval call binding the contract method 0x75b238fc.
//
// Solidity: function ADMIN_ROLE() view returns(bytes32)
func (_Eth *EthCallerSession) ADMINROLE() ([32]byte, error) {
	return _Eth.Contract.ADMINROLE(&_Eth.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Eth *EthCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Eth *EthSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Eth.Contract.DEFAULTADMINROLE(&_Eth.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Eth *EthCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Eth.Contract.DEFAULTADMINROLE(&_Eth.CallOpts)
}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Eth *EthCaller) RELAYROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "RELAY_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Eth *EthSession) RELAYROLE() ([32]byte, error) {
	return _Eth.Contract.RELAYROLE(&_Eth.CallOpts)
}

// RELAYROLE is a free data retrieval call binding the contract method 0x04421823.
//
// Solidity: function RELAY_ROLE() view returns(bytes32)
func (_Eth *EthCallerSession) RELAYROLE() ([32]byte, error) {
	return _Eth.Contract.RELAYROLE(&_Eth.CallOpts)
}

// BytesToUint is a free data retrieval call binding the contract method 0x02d06d05.
//
// Solidity: function bytesToUint(bytes b) view returns(uint256)
func (_Eth *EthCaller) BytesToUint(opts *bind.CallOpts, b []byte) (*big.Int, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "bytesToUint", b)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BytesToUint is a free data retrieval call binding the contract method 0x02d06d05.
//
// Solidity: function bytesToUint(bytes b) view returns(uint256)
func (_Eth *EthSession) BytesToUint(b []byte) (*big.Int, error) {
	return _Eth.Contract.BytesToUint(&_Eth.CallOpts, b)
}

// BytesToUint is a free data retrieval call binding the contract method 0x02d06d05.
//
// Solidity: function bytesToUint(bytes b) view returns(uint256)
func (_Eth *EthCallerSession) BytesToUint(b []byte) (*big.Int, error) {
	return _Eth.Contract.BytesToUint(&_Eth.CallOpts, b)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Eth *EthCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Eth *EthSession) Fee() (*big.Int, error) {
	return _Eth.Contract.Fee(&_Eth.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Eth *EthCallerSession) Fee() (*big.Int, error) {
	return _Eth.Contract.Fee(&_Eth.CallOpts)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Eth *EthCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Eth *EthSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Eth.Contract.GetRoleAdmin(&_Eth.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Eth *EthCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Eth.Contract.GetRoleAdmin(&_Eth.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Eth *EthCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Eth *EthSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Eth.Contract.HasRole(&_Eth.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Eth *EthCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Eth.Contract.HasRole(&_Eth.CallOpts, role, account)
}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Eth *EthCaller) InputEventId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "inputEventId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Eth *EthSession) InputEventId() (*big.Int, error) {
	return _Eth.Contract.InputEventId(&_Eth.CallOpts)
}

// InputEventId is a free data retrieval call binding the contract method 0x99b5bb64.
//
// Solidity: function inputEventId() view returns(uint256)
func (_Eth *EthCallerSession) InputEventId() (*big.Int, error) {
	return _Eth.Contract.InputEventId(&_Eth.CallOpts)
}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Eth *EthCaller) LockTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "lockTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Eth *EthSession) LockTime() (*big.Int, error) {
	return _Eth.Contract.LockTime(&_Eth.CallOpts)
}

// LockTime is a free data retrieval call binding the contract method 0x0d668087.
//
// Solidity: function lockTime() view returns(uint256)
func (_Eth *EthCallerSession) LockTime() (*big.Int, error) {
	return _Eth.Contract.LockTime(&_Eth.CallOpts)
}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Eth *EthCaller) LockedTransfers(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "lockedTransfers", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Eth *EthSession) LockedTransfers(arg0 *big.Int) (*big.Int, error) {
	return _Eth.Contract.LockedTransfers(&_Eth.CallOpts, arg0)
}

// LockedTransfers is a free data retrieval call binding the contract method 0x4a1856de.
//
// Solidity: function lockedTransfers(uint256 ) view returns(uint256 endTimestamp)
func (_Eth *EthCallerSession) LockedTransfers(arg0 *big.Int) (*big.Int, error) {
	return _Eth.Contract.LockedTransfers(&_Eth.CallOpts, arg0)
}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Eth *EthCaller) MinSafetyBlocks(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "minSafetyBlocks")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Eth *EthSession) MinSafetyBlocks() (*big.Int, error) {
	return _Eth.Contract.MinSafetyBlocks(&_Eth.CallOpts)
}

// MinSafetyBlocks is a free data retrieval call binding the contract method 0x924cf6e0.
//
// Solidity: function minSafetyBlocks() view returns(uint256)
func (_Eth *EthCallerSession) MinSafetyBlocks() (*big.Int, error) {
	return _Eth.Contract.MinSafetyBlocks(&_Eth.CallOpts)
}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Eth *EthCaller) SideBridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "sideBridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Eth *EthSession) SideBridgeAddress() (common.Address, error) {
	return _Eth.Contract.SideBridgeAddress(&_Eth.CallOpts)
}

// SideBridgeAddress is a free data retrieval call binding the contract method 0xf33fe10f.
//
// Solidity: function sideBridgeAddress() view returns(address)
func (_Eth *EthCallerSession) SideBridgeAddress() (common.Address, error) {
	return _Eth.Contract.SideBridgeAddress(&_Eth.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Eth *EthCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Eth *EthSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Eth.Contract.SupportsInterface(&_Eth.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Eth *EthCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Eth.Contract.SupportsInterface(&_Eth.CallOpts, interfaceId)
}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Eth *EthCaller) TimeframeSeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "timeframeSeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Eth *EthSession) TimeframeSeconds() (*big.Int, error) {
	return _Eth.Contract.TimeframeSeconds(&_Eth.CallOpts)
}

// TimeframeSeconds is a free data retrieval call binding the contract method 0xbaeebe75.
//
// Solidity: function timeframeSeconds() view returns(uint256)
func (_Eth *EthCallerSession) TimeframeSeconds() (*big.Int, error) {
	return _Eth.Contract.TimeframeSeconds(&_Eth.CallOpts)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Eth *EthCaller) TokenAddresses(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "tokenAddresses", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Eth *EthSession) TokenAddresses(arg0 common.Address) (common.Address, error) {
	return _Eth.Contract.TokenAddresses(&_Eth.CallOpts, arg0)
}

// TokenAddresses is a free data retrieval call binding the contract method 0xb6d3385e.
//
// Solidity: function tokenAddresses(address ) view returns(address)
func (_Eth *EthCallerSession) TokenAddresses(arg0 common.Address) (common.Address, error) {
	return _Eth.Contract.TokenAddresses(&_Eth.CallOpts, arg0)
}

// CheckAura is a paid mutator transaction binding the contract method 0x70ed397a.
//
// Solidity: function CheckAura_(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,int256)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof, uint256 minSafetyBlocks) returns()
func (_Eth *EthTransactor) CheckAura(opts *bind.TransactOpts, auraProof CheckAuraAuraProof, minSafetyBlocks *big.Int) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "CheckAura_", auraProof, minSafetyBlocks)
}

// CheckAura is a paid mutator transaction binding the contract method 0x70ed397a.
//
// Solidity: function CheckAura_(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,int256)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof, uint256 minSafetyBlocks) returns()
func (_Eth *EthSession) CheckAura(auraProof CheckAuraAuraProof, minSafetyBlocks *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.CheckAura(&_Eth.TransactOpts, auraProof, minSafetyBlocks)
}

// CheckAura is a paid mutator transaction binding the contract method 0x70ed397a.
//
// Solidity: function CheckAura_(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,int256)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof, uint256 minSafetyBlocks) returns()
func (_Eth *EthTransactorSession) CheckAura(auraProof CheckAuraAuraProof, minSafetyBlocks *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.CheckAura(&_Eth.TransactOpts, auraProof, minSafetyBlocks)
}

// ChangeFee is a paid mutator transaction binding the contract method 0x6a1db1bf.
//
// Solidity: function changeFee(uint256 fee_) returns()
func (_Eth *EthTransactor) ChangeFee(opts *bind.TransactOpts, fee_ *big.Int) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "changeFee", fee_)
}

// ChangeFee is a paid mutator transaction binding the contract method 0x6a1db1bf.
//
// Solidity: function changeFee(uint256 fee_) returns()
func (_Eth *EthSession) ChangeFee(fee_ *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.ChangeFee(&_Eth.TransactOpts, fee_)
}

// ChangeFee is a paid mutator transaction binding the contract method 0x6a1db1bf.
//
// Solidity: function changeFee(uint256 fee_) returns()
func (_Eth *EthTransactorSession) ChangeFee(fee_ *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.ChangeFee(&_Eth.TransactOpts, fee_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Eth *EthTransactor) ChangeLockTime(opts *bind.TransactOpts, lockTime_ *big.Int) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "changeLockTime", lockTime_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Eth *EthSession) ChangeLockTime(lockTime_ *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.ChangeLockTime(&_Eth.TransactOpts, lockTime_)
}

// ChangeLockTime is a paid mutator transaction binding the contract method 0x96cf5227.
//
// Solidity: function changeLockTime(uint256 lockTime_) returns()
func (_Eth *EthTransactorSession) ChangeLockTime(lockTime_ *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.ChangeLockTime(&_Eth.TransactOpts, lockTime_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Eth *EthTransactor) ChangeMinSafetyBlocks(opts *bind.TransactOpts, minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "changeMinSafetyBlocks", minSafetyBlocks_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Eth *EthSession) ChangeMinSafetyBlocks(minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.ChangeMinSafetyBlocks(&_Eth.TransactOpts, minSafetyBlocks_)
}

// ChangeMinSafetyBlocks is a paid mutator transaction binding the contract method 0xfd5d2ef3.
//
// Solidity: function changeMinSafetyBlocks(uint256 minSafetyBlocks_) returns()
func (_Eth *EthTransactorSession) ChangeMinSafetyBlocks(minSafetyBlocks_ *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.ChangeMinSafetyBlocks(&_Eth.TransactOpts, minSafetyBlocks_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Eth *EthTransactor) ChangeTimeframeSeconds(opts *bind.TransactOpts, timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "changeTimeframeSeconds", timeframeSeconds_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Eth *EthSession) ChangeTimeframeSeconds(timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.ChangeTimeframeSeconds(&_Eth.TransactOpts, timeframeSeconds_)
}

// ChangeTimeframeSeconds is a paid mutator transaction binding the contract method 0x42180fb8.
//
// Solidity: function changeTimeframeSeconds(uint256 timeframeSeconds_) returns()
func (_Eth *EthTransactorSession) ChangeTimeframeSeconds(timeframeSeconds_ *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.ChangeTimeframeSeconds(&_Eth.TransactOpts, timeframeSeconds_)
}

// EmitTestEvent is a paid mutator transaction binding the contract method 0xbc5b8754.
//
// Solidity: function emitTestEvent(address tokenAmbAddress, address toAddress, uint256 amount) returns()
func (_Eth *EthTransactor) EmitTestEvent(opts *bind.TransactOpts, tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "emitTestEvent", tokenAmbAddress, toAddress, amount)
}

// EmitTestEvent is a paid mutator transaction binding the contract method 0xbc5b8754.
//
// Solidity: function emitTestEvent(address tokenAmbAddress, address toAddress, uint256 amount) returns()
func (_Eth *EthSession) EmitTestEvent(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.EmitTestEvent(&_Eth.TransactOpts, tokenAmbAddress, toAddress, amount)
}

// EmitTestEvent is a paid mutator transaction binding the contract method 0xbc5b8754.
//
// Solidity: function emitTestEvent(address tokenAmbAddress, address toAddress, uint256 amount) returns()
func (_Eth *EthTransactorSession) EmitTestEvent(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.EmitTestEvent(&_Eth.TransactOpts, tokenAmbAddress, toAddress, amount)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Eth *EthTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Eth *EthSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Eth.Contract.GrantRole(&_Eth.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Eth *EthTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Eth.Contract.GrantRole(&_Eth.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Eth *EthTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Eth *EthSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Eth.Contract.RenounceRole(&_Eth.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Eth *EthTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Eth.Contract.RenounceRole(&_Eth.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Eth *EthTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Eth *EthSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Eth.Contract.RevokeRole(&_Eth.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Eth *EthTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Eth.Contract.RevokeRole(&_Eth.TransactOpts, role, account)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0x1f4844f3.
//
// Solidity: function submitTransfer(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,int256)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof) returns()
func (_Eth *EthTransactor) SubmitTransfer(opts *bind.TransactOpts, auraProof CheckAuraAuraProof) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "submitTransfer", auraProof)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0x1f4844f3.
//
// Solidity: function submitTransfer(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,int256)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof) returns()
func (_Eth *EthSession) SubmitTransfer(auraProof CheckAuraAuraProof) (*types.Transaction, error) {
	return _Eth.Contract.SubmitTransfer(&_Eth.TransactOpts, auraProof)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0x1f4844f3.
//
// Solidity: function submitTransfer(((bytes,bytes,bytes,bytes32,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,int256)[],(bytes[],uint256,(address,address,uint256)[]),(bytes[],address,int64)[]) auraProof) returns()
func (_Eth *EthTransactorSession) SubmitTransfer(auraProof CheckAuraAuraProof) (*types.Transaction, error) {
	return _Eth.Contract.SubmitTransfer(&_Eth.TransactOpts, auraProof)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Eth *EthTransactor) TokensAdd(opts *bind.TransactOpts, tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "tokensAdd", tokenThisAddress, tokenSideAddress)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Eth *EthSession) TokensAdd(tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Eth.Contract.TokensAdd(&_Eth.TransactOpts, tokenThisAddress, tokenSideAddress)
}

// TokensAdd is a paid mutator transaction binding the contract method 0x853890ae.
//
// Solidity: function tokensAdd(address tokenThisAddress, address tokenSideAddress) returns()
func (_Eth *EthTransactorSession) TokensAdd(tokenThisAddress common.Address, tokenSideAddress common.Address) (*types.Transaction, error) {
	return _Eth.Contract.TokensAdd(&_Eth.TransactOpts, tokenThisAddress, tokenSideAddress)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Eth *EthTransactor) TokensAddBatch(opts *bind.TransactOpts, tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "tokensAddBatch", tokenThisAddresses, tokenSideAddresses)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Eth *EthSession) TokensAddBatch(tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Eth.Contract.TokensAddBatch(&_Eth.TransactOpts, tokenThisAddresses, tokenSideAddresses)
}

// TokensAddBatch is a paid mutator transaction binding the contract method 0x09fce356.
//
// Solidity: function tokensAddBatch(address[] tokenThisAddresses, address[] tokenSideAddresses) returns()
func (_Eth *EthTransactorSession) TokensAddBatch(tokenThisAddresses []common.Address, tokenSideAddresses []common.Address) (*types.Transaction, error) {
	return _Eth.Contract.TokensAddBatch(&_Eth.TransactOpts, tokenThisAddresses, tokenSideAddresses)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Eth *EthTransactor) TokensRemove(opts *bind.TransactOpts, tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "tokensRemove", tokenThisAddress)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Eth *EthSession) TokensRemove(tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Eth.Contract.TokensRemove(&_Eth.TransactOpts, tokenThisAddress)
}

// TokensRemove is a paid mutator transaction binding the contract method 0x8e5df9c7.
//
// Solidity: function tokensRemove(address tokenThisAddress) returns()
func (_Eth *EthTransactorSession) TokensRemove(tokenThisAddress common.Address) (*types.Transaction, error) {
	return _Eth.Contract.TokensRemove(&_Eth.TransactOpts, tokenThisAddress)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Eth *EthTransactor) TokensRemoveBatch(opts *bind.TransactOpts, tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "tokensRemoveBatch", tokenThisAddresses)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Eth *EthSession) TokensRemoveBatch(tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Eth.Contract.TokensRemoveBatch(&_Eth.TransactOpts, tokenThisAddresses)
}

// TokensRemoveBatch is a paid mutator transaction binding the contract method 0x5249a705.
//
// Solidity: function tokensRemoveBatch(address[] tokenThisAddresses) returns()
func (_Eth *EthTransactorSession) TokensRemoveBatch(tokenThisAddresses []common.Address) (*types.Transaction, error) {
	return _Eth.Contract.TokensRemoveBatch(&_Eth.TransactOpts, tokenThisAddresses)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 event_id) returns()
func (_Eth *EthTransactor) UnlockTransfers(opts *bind.TransactOpts, event_id *big.Int) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "unlockTransfers", event_id)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 event_id) returns()
func (_Eth *EthSession) UnlockTransfers(event_id *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.UnlockTransfers(&_Eth.TransactOpts, event_id)
}

// UnlockTransfers is a paid mutator transaction binding the contract method 0xf862b7eb.
//
// Solidity: function unlockTransfers(uint256 event_id) returns()
func (_Eth *EthTransactorSession) UnlockTransfers(event_id *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.UnlockTransfers(&_Eth.TransactOpts, event_id)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAmbAddress, address toAddress, uint256 amount) payable returns()
func (_Eth *EthTransactor) Withdraw(opts *bind.TransactOpts, tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "withdraw", tokenAmbAddress, toAddress, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAmbAddress, address toAddress, uint256 amount) payable returns()
func (_Eth *EthSession) Withdraw(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.Withdraw(&_Eth.TransactOpts, tokenAmbAddress, toAddress, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address tokenAmbAddress, address toAddress, uint256 amount) payable returns()
func (_Eth *EthTransactorSession) Withdraw(tokenAmbAddress common.Address, toAddress common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.Withdraw(&_Eth.TransactOpts, tokenAmbAddress, toAddress, amount)
}

// EthRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Eth contract.
type EthRoleAdminChangedIterator struct {
	Event *EthRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *EthRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthRoleAdminChanged)
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
		it.Event = new(EthRoleAdminChanged)
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
func (it *EthRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthRoleAdminChanged represents a RoleAdminChanged event raised by the Eth contract.
type EthRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Eth *EthFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*EthRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Eth.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &EthRoleAdminChangedIterator{contract: _Eth.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Eth *EthFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *EthRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Eth.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthRoleAdminChanged)
				if err := _Eth.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Eth *EthFilterer) ParseRoleAdminChanged(log types.Log) (*EthRoleAdminChanged, error) {
	event := new(EthRoleAdminChanged)
	if err := _Eth.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Eth contract.
type EthRoleGrantedIterator struct {
	Event *EthRoleGranted // Event containing the contract specifics and raw log

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
func (it *EthRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthRoleGranted)
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
		it.Event = new(EthRoleGranted)
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
func (it *EthRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthRoleGranted represents a RoleGranted event raised by the Eth contract.
type EthRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Eth *EthFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*EthRoleGrantedIterator, error) {

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

	logs, sub, err := _Eth.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EthRoleGrantedIterator{contract: _Eth.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Eth *EthFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *EthRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Eth.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthRoleGranted)
				if err := _Eth.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Eth *EthFilterer) ParseRoleGranted(log types.Log) (*EthRoleGranted, error) {
	event := new(EthRoleGranted)
	if err := _Eth.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Eth contract.
type EthRoleRevokedIterator struct {
	Event *EthRoleRevoked // Event containing the contract specifics and raw log

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
func (it *EthRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthRoleRevoked)
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
		it.Event = new(EthRoleRevoked)
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
func (it *EthRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthRoleRevoked represents a RoleRevoked event raised by the Eth contract.
type EthRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Eth *EthFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*EthRoleRevokedIterator, error) {

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

	logs, sub, err := _Eth.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &EthRoleRevokedIterator{contract: _Eth.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Eth *EthFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *EthRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Eth.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthRoleRevoked)
				if err := _Eth.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Eth *EthFilterer) ParseRoleRevoked(log types.Log) (*EthRoleRevoked, error) {
	event := new(EthRoleRevoked)
	if err := _Eth.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthTestIterator is returned from FilterTest and is used to iterate over the raw logs and unpacked data for Test events raised by the Eth contract.
type EthTestIterator struct {
	Event *EthTest // Event containing the contract specifics and raw log

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
func (it *EthTestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthTest)
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
		it.Event = new(EthTest)
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
func (it *EthTestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthTestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthTest represents a Test event raised by the Eth contract.
type EthTest struct {
	A   *big.Int
	B   common.Address
	C   string
	D   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterTest is a free log retrieval operation binding the contract event 0xf690cc04b54c126882a8927cb940bececfeabf9b0b4b0069db2427fbf9e86dd6.
//
// Solidity: event Test(uint256 indexed a, address indexed b, string c, uint256 d)
func (_Eth *EthFilterer) FilterTest(opts *bind.FilterOpts, a []*big.Int, b []common.Address) (*EthTestIterator, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _Eth.contract.FilterLogs(opts, "Test", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return &EthTestIterator{contract: _Eth.contract, event: "Test", logs: logs, sub: sub}, nil
}

// WatchTest is a free log subscription operation binding the contract event 0xf690cc04b54c126882a8927cb940bececfeabf9b0b4b0069db2427fbf9e86dd6.
//
// Solidity: event Test(uint256 indexed a, address indexed b, string c, uint256 d)
func (_Eth *EthFilterer) WatchTest(opts *bind.WatchOpts, sink chan<- *EthTest, a []*big.Int, b []common.Address) (event.Subscription, error) {

	var aRule []interface{}
	for _, aItem := range a {
		aRule = append(aRule, aItem)
	}
	var bRule []interface{}
	for _, bItem := range b {
		bRule = append(bRule, bItem)
	}

	logs, sub, err := _Eth.contract.WatchLogs(opts, "Test", aRule, bRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthTest)
				if err := _Eth.contract.UnpackLog(event, "Test", log); err != nil {
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
func (_Eth *EthFilterer) ParseTest(log types.Log) (*EthTest, error) {
	event := new(EthTest)
	if err := _Eth.contract.UnpackLog(event, "Test", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Eth contract.
type EthTransferIterator struct {
	Event *EthTransfer // Event containing the contract specifics and raw log

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
func (it *EthTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthTransfer)
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
		it.Event = new(EthTransfer)
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
func (it *EthTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2.
//
// Solidity: event Transfer(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Eth *EthFilterer) FilterTransfer(opts *bind.FilterOpts, event_id []*big.Int) (*EthTransferIterator, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Eth.contract.FilterLogs(opts, "Transfer", event_idRule)
	if err != nil {
		return nil, err
	}
	return &EthTransferIterator{contract: _Eth.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xe15729a2f427aa4572dab35eb692c902fcbce57d41642013259c741380809ae2.
//
// Solidity: event Transfer(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Eth *EthFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *EthTransfer, event_id []*big.Int) (event.Subscription, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Eth.contract.WatchLogs(opts, "Transfer", event_idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthTransfer)
				if err := _Eth.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_Eth *EthFilterer) ParseTransfer(log types.Log) (*EthTransfer, error) {
	event := new(EthTransfer)
	if err := _Eth.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
