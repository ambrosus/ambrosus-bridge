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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_sideBridgeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"relayAddress\",\"type\":\"address\"},{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenSideAddresses\",\"type\":\"address[]\"},{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"unimportant\",\"type\":\"uint256\"}],\"name\":\"Test\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"queue\",\"type\":\"tuple[]\"}],\"name\":\"TransferEvent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p0_seal\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p0_bare\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"prevHashOrReceiptRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"timestamp\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s1\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"s2\",\"type\":\"bytes\"}],\"internalType\":\"structCheckPoA.BlockPoA[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer[]\",\"name\":\"events\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"proof\",\"type\":\"bytes[]\"}],\"name\":\"CheckPoA_\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"proof\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"eventToSearch\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"receiptsRoot\",\"type\":\"bytes32\"}],\"name\":\"CheckReceiptsProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"RELAY_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structCommonStructs.Transfer\",\"name\":\"transfer\",\"type\":\"tuple\"}],\"name\":\"Transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"}],\"name\":\"bytesToUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"proof\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"eventToSearch\",\"type\":\"bytes\"}],\"name\":\"calcReceiptsRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee_\",\"type\":\"uint256\"}],\"name\":\"changeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"event_id\",\"type\":\"uint256\"}],\"name\":\"eventTest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"inputEventId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sideBridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenSideAddress\",\"type\":\"address\"}],\"name\":\"tokensAdd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tokenSideAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensAddBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenThisAddress\",\"type\":\"address\"}],\"name\":\"tokensRemove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenThisAddresses\",\"type\":\"address[]\"}],\"name\":\"tokensRemoveBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAmbAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
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

// CalcReceiptsRoot is a free data retrieval call binding the contract method 0x131c7397.
//
// Solidity: function calcReceiptsRoot(bytes[] proof, bytes eventToSearch) view returns(bytes32)
func (_Eth *EthCaller) CalcReceiptsRoot(opts *bind.CallOpts, proof [][]byte, eventToSearch []byte) ([32]byte, error) {
	var out []interface{}
	err := _Eth.contract.Call(opts, &out, "calcReceiptsRoot", proof, eventToSearch)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalcReceiptsRoot is a free data retrieval call binding the contract method 0x131c7397.
//
// Solidity: function calcReceiptsRoot(bytes[] proof, bytes eventToSearch) view returns(bytes32)
func (_Eth *EthSession) CalcReceiptsRoot(proof [][]byte, eventToSearch []byte) ([32]byte, error) {
	return _Eth.Contract.CalcReceiptsRoot(&_Eth.CallOpts, proof, eventToSearch)
}

// CalcReceiptsRoot is a free data retrieval call binding the contract method 0x131c7397.
//
// Solidity: function calcReceiptsRoot(bytes[] proof, bytes eventToSearch) view returns(bytes32)
func (_Eth *EthCallerSession) CalcReceiptsRoot(proof [][]byte, eventToSearch []byte) ([32]byte, error) {
	return _Eth.Contract.CalcReceiptsRoot(&_Eth.CallOpts, proof, eventToSearch)
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

// CheckPoA is a paid mutator transaction binding the contract method 0x6580ed60.
//
// Solidity: function CheckPoA_((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes)[] blocks, (address,address,uint256)[] events, bytes[] proof) returns()
func (_Eth *EthTransactor) CheckPoA(opts *bind.TransactOpts, blocks []CheckPoABlockPoA, events []CommonStructsTransfer, proof [][]byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "CheckPoA_", blocks, events, proof)
}

// CheckPoA is a paid mutator transaction binding the contract method 0x6580ed60.
//
// Solidity: function CheckPoA_((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes)[] blocks, (address,address,uint256)[] events, bytes[] proof) returns()
func (_Eth *EthSession) CheckPoA(blocks []CheckPoABlockPoA, events []CommonStructsTransfer, proof [][]byte) (*types.Transaction, error) {
	return _Eth.Contract.CheckPoA(&_Eth.TransactOpts, blocks, events, proof)
}

// CheckPoA is a paid mutator transaction binding the contract method 0x6580ed60.
//
// Solidity: function CheckPoA_((bytes,bytes,bytes,bytes32,bytes,bytes,bytes,bytes,bytes,bytes)[] blocks, (address,address,uint256)[] events, bytes[] proof) returns()
func (_Eth *EthTransactorSession) CheckPoA(blocks []CheckPoABlockPoA, events []CommonStructsTransfer, proof [][]byte) (*types.Transaction, error) {
	return _Eth.Contract.CheckPoA(&_Eth.TransactOpts, blocks, events, proof)
}

// CheckReceiptsProof is a paid mutator transaction binding the contract method 0xd44d273c.
//
// Solidity: function CheckReceiptsProof(bytes[] proof, bytes eventToSearch, bytes32 receiptsRoot) returns()
func (_Eth *EthTransactor) CheckReceiptsProof(opts *bind.TransactOpts, proof [][]byte, eventToSearch []byte, receiptsRoot [32]byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "CheckReceiptsProof", proof, eventToSearch, receiptsRoot)
}

// CheckReceiptsProof is a paid mutator transaction binding the contract method 0xd44d273c.
//
// Solidity: function CheckReceiptsProof(bytes[] proof, bytes eventToSearch, bytes32 receiptsRoot) returns()
func (_Eth *EthSession) CheckReceiptsProof(proof [][]byte, eventToSearch []byte, receiptsRoot [32]byte) (*types.Transaction, error) {
	return _Eth.Contract.CheckReceiptsProof(&_Eth.TransactOpts, proof, eventToSearch, receiptsRoot)
}

// CheckReceiptsProof is a paid mutator transaction binding the contract method 0xd44d273c.
//
// Solidity: function CheckReceiptsProof(bytes[] proof, bytes eventToSearch, bytes32 receiptsRoot) returns()
func (_Eth *EthTransactorSession) CheckReceiptsProof(proof [][]byte, eventToSearch []byte, receiptsRoot [32]byte) (*types.Transaction, error) {
	return _Eth.Contract.CheckReceiptsProof(&_Eth.TransactOpts, proof, eventToSearch, receiptsRoot)
}

// Transfer is a paid mutator transaction binding the contract method 0xc6f666c8.
//
// Solidity: function Transfer((address,address,uint256) transfer) returns()
func (_Eth *EthTransactor) Transfer(opts *bind.TransactOpts, transfer CommonStructsTransfer) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "Transfer", transfer)
}

// Transfer is a paid mutator transaction binding the contract method 0xc6f666c8.
//
// Solidity: function Transfer((address,address,uint256) transfer) returns()
func (_Eth *EthSession) Transfer(transfer CommonStructsTransfer) (*types.Transaction, error) {
	return _Eth.Contract.Transfer(&_Eth.TransactOpts, transfer)
}

// Transfer is a paid mutator transaction binding the contract method 0xc6f666c8.
//
// Solidity: function Transfer((address,address,uint256) transfer) returns()
func (_Eth *EthTransactorSession) Transfer(transfer CommonStructsTransfer) (*types.Transaction, error) {
	return _Eth.Contract.Transfer(&_Eth.TransactOpts, transfer)
}

// AcceptBlock is a paid mutator transaction binding the contract method 0x54583364.
//
// Solidity: function acceptBlock() returns()
func (_Eth *EthTransactor) AcceptBlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "acceptBlock")
}

// AcceptBlock is a paid mutator transaction binding the contract method 0x54583364.
//
// Solidity: function acceptBlock() returns()
func (_Eth *EthSession) AcceptBlock() (*types.Transaction, error) {
	return _Eth.Contract.AcceptBlock(&_Eth.TransactOpts)
}

// AcceptBlock is a paid mutator transaction binding the contract method 0x54583364.
//
// Solidity: function acceptBlock() returns()
func (_Eth *EthTransactorSession) AcceptBlock() (*types.Transaction, error) {
	return _Eth.Contract.AcceptBlock(&_Eth.TransactOpts)
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

// EventTest is a paid mutator transaction binding the contract method 0x94fb3d29.
//
// Solidity: function eventTest(uint256 event_id) returns()
func (_Eth *EthTransactor) EventTest(opts *bind.TransactOpts, event_id *big.Int) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "eventTest", event_id)
}

// EventTest is a paid mutator transaction binding the contract method 0x94fb3d29.
//
// Solidity: function eventTest(uint256 event_id) returns()
func (_Eth *EthSession) EventTest(event_id *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.EventTest(&_Eth.TransactOpts, event_id)
}

// EventTest is a paid mutator transaction binding the contract method 0x94fb3d29.
//
// Solidity: function eventTest(uint256 event_id) returns()
func (_Eth *EthTransactorSession) EventTest(event_id *big.Int) (*types.Transaction, error) {
	return _Eth.Contract.EventTest(&_Eth.TransactOpts, event_id)
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
	EventId     *big.Int
	Unimportant *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTest is a free log retrieval operation binding the contract event 0x91916a5e2c96453ddf6b585497262675140eb9f7a774095fb003d93e6dc69216.
//
// Solidity: event Test(uint256 indexed event_id, uint256 unimportant)
func (_Eth *EthFilterer) FilterTest(opts *bind.FilterOpts, event_id []*big.Int) (*EthTestIterator, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Eth.contract.FilterLogs(opts, "Test", event_idRule)
	if err != nil {
		return nil, err
	}
	return &EthTestIterator{contract: _Eth.contract, event: "Test", logs: logs, sub: sub}, nil
}

// WatchTest is a free log subscription operation binding the contract event 0x91916a5e2c96453ddf6b585497262675140eb9f7a774095fb003d93e6dc69216.
//
// Solidity: event Test(uint256 indexed event_id, uint256 unimportant)
func (_Eth *EthFilterer) WatchTest(opts *bind.WatchOpts, sink chan<- *EthTest, event_id []*big.Int) (event.Subscription, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Eth.contract.WatchLogs(opts, "Test", event_idRule)
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

// ParseTest is a log parse operation binding the contract event 0x91916a5e2c96453ddf6b585497262675140eb9f7a774095fb003d93e6dc69216.
//
// Solidity: event Test(uint256 indexed event_id, uint256 unimportant)
func (_Eth *EthFilterer) ParseTest(log types.Log) (*EthTest, error) {
	event := new(EthTest)
	if err := _Eth.contract.UnpackLog(event, "Test", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EthTransferEventIterator is returned from FilterTransferEvent and is used to iterate over the raw logs and unpacked data for TransferEvent events raised by the Eth contract.
type EthTransferEventIterator struct {
	Event *EthTransferEvent // Event containing the contract specifics and raw log

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
func (it *EthTransferEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthTransferEvent)
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
		it.Event = new(EthTransferEvent)
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
func (it *EthTransferEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthTransferEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FilterTransferEvent is a free log retrieval operation binding the contract event 0xf81cdc6472b697220bc91c185dc3a1e3faec8c316a022f46a5a284d94b4cd8ce.
//
// Solidity: event TransferEvent(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Eth *EthFilterer) FilterTransferEvent(opts *bind.FilterOpts, event_id []*big.Int) (*EthTransferEventIterator, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Eth.contract.FilterLogs(opts, "TransferEvent", event_idRule)
	if err != nil {
		return nil, err
	}
	return &EthTransferEventIterator{contract: _Eth.contract, event: "TransferEvent", logs: logs, sub: sub}, nil
}

// WatchTransferEvent is a free log subscription operation binding the contract event 0xf81cdc6472b697220bc91c185dc3a1e3faec8c316a022f46a5a284d94b4cd8ce.
//
// Solidity: event TransferEvent(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Eth *EthFilterer) WatchTransferEvent(opts *bind.WatchOpts, sink chan<- *EthTransferEvent, event_id []*big.Int) (event.Subscription, error) {

	var event_idRule []interface{}
	for _, event_idItem := range event_id {
		event_idRule = append(event_idRule, event_idItem)
	}

	logs, sub, err := _Eth.contract.WatchLogs(opts, "TransferEvent", event_idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthTransferEvent)
				if err := _Eth.contract.UnpackLog(event, "TransferEvent", log); err != nil {
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

// ParseTransferEvent is a log parse operation binding the contract event 0xf81cdc6472b697220bc91c185dc3a1e3faec8c316a022f46a5a284d94b4cd8ce.
//
// Solidity: event TransferEvent(uint256 indexed event_id, (address,address,uint256)[] queue)
func (_Eth *EthFilterer) ParseTransferEvent(log types.Log) (*EthTransferEvent, error) {
	event := new(EthTransferEvent)
	if err := _Eth.contract.UnpackLog(event, "TransferEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
