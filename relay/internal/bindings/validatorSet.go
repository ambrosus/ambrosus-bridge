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

// VsMetaData contains all meta data concerning the Vs contract.
var VsMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pendingValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeChange\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPendingValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_initialValidators\",\"type\":\"address[]\"},{\"name\":\"_superUser\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"parentHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"newSet\",\"type\":\"address[]\"}],\"name\":\"InitiateChange\",\"type\":\"event\"}]",
}

// VsABI is the input ABI used to generate the binding from.
// Deprecated: Use VsMetaData.ABI instead.
var VsABI = VsMetaData.ABI

// Vs is an auto generated Go binding around an Ethereum contract.
type Vs struct {
	VsCaller     // Read-only binding to the contract
	VsTransactor // Write-only binding to the contract
	VsFilterer   // Log filterer for contract events
}

// VsCaller is an auto generated read-only Go binding around an Ethereum contract.
type VsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VsSession struct {
	Contract     *Vs               // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VsCallerSession struct {
	Contract *VsCaller     // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VsTransactorSession struct {
	Contract     *VsTransactor     // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VsRaw is an auto generated low-level Go binding around an Ethereum contract.
type VsRaw struct {
	Contract *Vs // Generic contract binding to access the raw methods on
}

// VsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VsCallerRaw struct {
	Contract *VsCaller // Generic read-only contract binding to access the raw methods on
}

// VsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VsTransactorRaw struct {
	Contract *VsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVs creates a new instance of Vs, bound to a specific deployed contract.
func NewVs(address common.Address, backend bind.ContractBackend) (*Vs, error) {
	contract, err := bindVs(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Vs{VsCaller: VsCaller{contract: contract}, VsTransactor: VsTransactor{contract: contract}, VsFilterer: VsFilterer{contract: contract}}, nil
}

// NewVsCaller creates a new read-only instance of Vs, bound to a specific deployed contract.
func NewVsCaller(address common.Address, caller bind.ContractCaller) (*VsCaller, error) {
	contract, err := bindVs(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VsCaller{contract: contract}, nil
}

// NewVsTransactor creates a new write-only instance of Vs, bound to a specific deployed contract.
func NewVsTransactor(address common.Address, transactor bind.ContractTransactor) (*VsTransactor, error) {
	contract, err := bindVs(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VsTransactor{contract: contract}, nil
}

// NewVsFilterer creates a new log filterer instance of Vs, bound to a specific deployed contract.
func NewVsFilterer(address common.Address, filterer bind.ContractFilterer) (*VsFilterer, error) {
	contract, err := bindVs(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VsFilterer{contract: contract}, nil
}

// bindVs binds a generic wrapper to an already deployed contract.
func bindVs(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vs *VsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vs.Contract.VsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vs *VsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vs.Contract.VsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vs *VsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vs.Contract.VsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vs *VsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vs.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vs *VsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vs.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vs *VsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vs.Contract.contract.Transact(opts, method, params...)
}

// GetPendingValidators is a free data retrieval call binding the contract method 0xeebc7a39.
//
// Solidity: function getPendingValidators() view returns(address[])
func (_Vs *VsCaller) GetPendingValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Vs.contract.Call(opts, &out, "getPendingValidators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetPendingValidators is a free data retrieval call binding the contract method 0xeebc7a39.
//
// Solidity: function getPendingValidators() view returns(address[])
func (_Vs *VsSession) GetPendingValidators() ([]common.Address, error) {
	return _Vs.Contract.GetPendingValidators(&_Vs.CallOpts)
}

// GetPendingValidators is a free data retrieval call binding the contract method 0xeebc7a39.
//
// Solidity: function getPendingValidators() view returns(address[])
func (_Vs *VsCallerSession) GetPendingValidators() ([]common.Address, error) {
	return _Vs.Contract.GetPendingValidators(&_Vs.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Vs *VsCaller) GetValidators(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Vs.contract.Call(opts, &out, "getValidators")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Vs *VsSession) GetValidators() ([]common.Address, error) {
	return _Vs.Contract.GetValidators(&_Vs.CallOpts)
}

// GetValidators is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Vs *VsCallerSession) GetValidators() ([]common.Address, error) {
	return _Vs.Contract.GetValidators(&_Vs.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vs *VsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vs.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vs *VsSession) Owner() (common.Address, error) {
	return _Vs.Contract.Owner(&_Vs.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vs *VsCallerSession) Owner() (common.Address, error) {
	return _Vs.Contract.Owner(&_Vs.CallOpts)
}

// PendingValidators is a free data retrieval call binding the contract method 0x28569e1f.
//
// Solidity: function pendingValidators(uint256 ) view returns(address)
func (_Vs *VsCaller) PendingValidators(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Vs.contract.Call(opts, &out, "pendingValidators", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingValidators is a free data retrieval call binding the contract method 0x28569e1f.
//
// Solidity: function pendingValidators(uint256 ) view returns(address)
func (_Vs *VsSession) PendingValidators(arg0 *big.Int) (common.Address, error) {
	return _Vs.Contract.PendingValidators(&_Vs.CallOpts, arg0)
}

// PendingValidators is a free data retrieval call binding the contract method 0x28569e1f.
//
// Solidity: function pendingValidators(uint256 ) view returns(address)
func (_Vs *VsCallerSession) PendingValidators(arg0 *big.Int) (common.Address, error) {
	return _Vs.Contract.PendingValidators(&_Vs.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) view returns(address)
func (_Vs *VsCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Vs.contract.Call(opts, &out, "validators", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) view returns(address)
func (_Vs *VsSession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Vs.Contract.Validators(&_Vs.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) view returns(address)
func (_Vs *VsCallerSession) Validators(arg0 *big.Int) (common.Address, error) {
	return _Vs.Contract.Validators(&_Vs.CallOpts, arg0)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(address _validator) returns()
func (_Vs *VsTransactor) AddValidator(opts *bind.TransactOpts, _validator common.Address) (*types.Transaction, error) {
	return _Vs.contract.Transact(opts, "addValidator", _validator)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(address _validator) returns()
func (_Vs *VsSession) AddValidator(_validator common.Address) (*types.Transaction, error) {
	return _Vs.Contract.AddValidator(&_Vs.TransactOpts, _validator)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(address _validator) returns()
func (_Vs *VsTransactorSession) AddValidator(_validator common.Address) (*types.Transaction, error) {
	return _Vs.Contract.AddValidator(&_Vs.TransactOpts, _validator)
}

// FinalizeChange is a paid mutator transaction binding the contract method 0x75286211.
//
// Solidity: function finalizeChange() returns()
func (_Vs *VsTransactor) FinalizeChange(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vs.contract.Transact(opts, "finalizeChange")
}

// FinalizeChange is a paid mutator transaction binding the contract method 0x75286211.
//
// Solidity: function finalizeChange() returns()
func (_Vs *VsSession) FinalizeChange() (*types.Transaction, error) {
	return _Vs.Contract.FinalizeChange(&_Vs.TransactOpts)
}

// FinalizeChange is a paid mutator transaction binding the contract method 0x75286211.
//
// Solidity: function finalizeChange() returns()
func (_Vs *VsTransactorSession) FinalizeChange() (*types.Transaction, error) {
	return _Vs.Contract.FinalizeChange(&_Vs.TransactOpts)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validator) returns()
func (_Vs *VsTransactor) RemoveValidator(opts *bind.TransactOpts, _validator common.Address) (*types.Transaction, error) {
	return _Vs.contract.Transact(opts, "removeValidator", _validator)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validator) returns()
func (_Vs *VsSession) RemoveValidator(_validator common.Address) (*types.Transaction, error) {
	return _Vs.Contract.RemoveValidator(&_Vs.TransactOpts, _validator)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(address _validator) returns()
func (_Vs *VsTransactorSession) RemoveValidator(_validator common.Address) (*types.Transaction, error) {
	return _Vs.Contract.RemoveValidator(&_Vs.TransactOpts, _validator)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vs *VsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Vs.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vs *VsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Vs.Contract.TransferOwnership(&_Vs.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Vs *VsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Vs.Contract.TransferOwnership(&_Vs.TransactOpts, newOwner)
}

// VsInitiateChangeIterator is returned from FilterInitiateChange and is used to iterate over the raw logs and unpacked data for InitiateChange events raised by the Vs contract.
type VsInitiateChangeIterator struct {
	Event *VsInitiateChange // Event containing the contract specifics and raw log

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
func (it *VsInitiateChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VsInitiateChange)
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
		it.Event = new(VsInitiateChange)
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
func (it *VsInitiateChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VsInitiateChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VsInitiateChange represents a InitiateChange event raised by the Vs contract.
type VsInitiateChange struct {
	ParentHash [32]byte
	NewSet     []common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterInitiateChange is a free log retrieval operation binding the contract event 0x55252fa6eee4741b4e24a74a70e9c11fd2c2281df8d6ea13126ff845f7825c89.
//
// Solidity: event InitiateChange(bytes32 indexed parentHash, address[] newSet)
func (_Vs *VsFilterer) FilterInitiateChange(opts *bind.FilterOpts, parentHash [][32]byte) (*VsInitiateChangeIterator, error) {

	var parentHashRule []interface{}
	for _, parentHashItem := range parentHash {
		parentHashRule = append(parentHashRule, parentHashItem)
	}

	logs, sub, err := _Vs.contract.FilterLogs(opts, "InitiateChange", parentHashRule)
	if err != nil {
		return nil, err
	}
	return &VsInitiateChangeIterator{contract: _Vs.contract, event: "InitiateChange", logs: logs, sub: sub}, nil
}

// WatchInitiateChange is a free log subscription operation binding the contract event 0x55252fa6eee4741b4e24a74a70e9c11fd2c2281df8d6ea13126ff845f7825c89.
//
// Solidity: event InitiateChange(bytes32 indexed parentHash, address[] newSet)
func (_Vs *VsFilterer) WatchInitiateChange(opts *bind.WatchOpts, sink chan<- *VsInitiateChange, parentHash [][32]byte) (event.Subscription, error) {

	var parentHashRule []interface{}
	for _, parentHashItem := range parentHash {
		parentHashRule = append(parentHashRule, parentHashItem)
	}

	logs, sub, err := _Vs.contract.WatchLogs(opts, "InitiateChange", parentHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VsInitiateChange)
				if err := _Vs.contract.UnpackLog(event, "InitiateChange", log); err != nil {
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

// ParseInitiateChange is a log parse operation binding the contract event 0x55252fa6eee4741b4e24a74a70e9c11fd2c2281df8d6ea13126ff845f7825c89.
//
// Solidity: event InitiateChange(bytes32 indexed parentHash, address[] newSet)
func (_Vs *VsFilterer) ParseInitiateChange(log types.Log) (*VsInitiateChange, error) {
	event := new(VsInitiateChange)
	if err := _Vs.contract.UnpackLog(event, "InitiateChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Vs contract.
type VsOwnershipTransferredIterator struct {
	Event *VsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *VsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VsOwnershipTransferred)
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
		it.Event = new(VsOwnershipTransferred)
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
func (it *VsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VsOwnershipTransferred represents a OwnershipTransferred event raised by the Vs contract.
type VsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Vs *VsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*VsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vs.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &VsOwnershipTransferredIterator{contract: _Vs.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Vs *VsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *VsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Vs.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VsOwnershipTransferred)
				if err := _Vs.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Vs *VsFilterer) ParseOwnershipTransferred(log types.Log) (*VsOwnershipTransferred, error) {
	event := new(VsOwnershipTransferred)
	if err := _Vs.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
