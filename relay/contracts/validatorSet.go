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

// VsMetaData contains all meta data concerning the Vs contract.
var VsMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"parentHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address[]\",\"name\":\"newSet\",\"type\":\"address[]\"}],\"name\":\"InitiateChange\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"addValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ambBlockNumber\",\"type\":\"uint256\"}],\"name\":\"finalizeChange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"finalizeChange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPendingValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ambBlockNumber\",\"type\":\"uint256\"}],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getValidators\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"pendingValidators\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ambBlockNumber\",\"type\":\"uint256\"}],\"name\":\"removeValidators\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// GetValidators is a free data retrieval call binding the contract method 0x471f40fb.
//
// Solidity: function getValidators(uint256 ambBlockNumber) view returns(address[])
func (_Vs *VsCaller) GetValidators(opts *bind.CallOpts, ambBlockNumber *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _Vs.contract.Call(opts, &out, "getValidators", ambBlockNumber)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetValidators is a free data retrieval call binding the contract method 0x471f40fb.
//
// Solidity: function getValidators(uint256 ambBlockNumber) view returns(address[])
func (_Vs *VsSession) GetValidators(ambBlockNumber *big.Int) ([]common.Address, error) {
	return _Vs.Contract.GetValidators(&_Vs.CallOpts, ambBlockNumber)
}

// GetValidators is a free data retrieval call binding the contract method 0x471f40fb.
//
// Solidity: function getValidators(uint256 ambBlockNumber) view returns(address[])
func (_Vs *VsCallerSession) GetValidators(ambBlockNumber *big.Int) ([]common.Address, error) {
	return _Vs.Contract.GetValidators(&_Vs.CallOpts, ambBlockNumber)
}

// GetValidators0 is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Vs *VsCaller) GetValidators0(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Vs.contract.Call(opts, &out, "getValidators0")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetValidators0 is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Vs *VsSession) GetValidators0() ([]common.Address, error) {
	return _Vs.Contract.GetValidators0(&_Vs.CallOpts)
}

// GetValidators0 is a free data retrieval call binding the contract method 0xb7ab4db5.
//
// Solidity: function getValidators() view returns(address[])
func (_Vs *VsCallerSession) GetValidators0() ([]common.Address, error) {
	return _Vs.Contract.GetValidators0(&_Vs.CallOpts)
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

// Validators is a free data retrieval call binding the contract method 0xdcf2793a.
//
// Solidity: function validators(uint256 , uint256 ) view returns(address)
func (_Vs *VsCaller) Validators(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Vs.contract.Call(opts, &out, "validators", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Validators is a free data retrieval call binding the contract method 0xdcf2793a.
//
// Solidity: function validators(uint256 , uint256 ) view returns(address)
func (_Vs *VsSession) Validators(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	return _Vs.Contract.Validators(&_Vs.CallOpts, arg0, arg1)
}

// Validators is a free data retrieval call binding the contract method 0xdcf2793a.
//
// Solidity: function validators(uint256 , uint256 ) view returns(address)
func (_Vs *VsCallerSession) Validators(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	return _Vs.Contract.Validators(&_Vs.CallOpts, arg0, arg1)
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

// FinalizeChange is a paid mutator transaction binding the contract method 0x55a6b050.
//
// Solidity: function finalizeChange(uint256 ambBlockNumber) returns()
func (_Vs *VsTransactor) FinalizeChange(opts *bind.TransactOpts, ambBlockNumber *big.Int) (*types.Transaction, error) {
	return _Vs.contract.Transact(opts, "finalizeChange", ambBlockNumber)
}

// FinalizeChange is a paid mutator transaction binding the contract method 0x55a6b050.
//
// Solidity: function finalizeChange(uint256 ambBlockNumber) returns()
func (_Vs *VsSession) FinalizeChange(ambBlockNumber *big.Int) (*types.Transaction, error) {
	return _Vs.Contract.FinalizeChange(&_Vs.TransactOpts, ambBlockNumber)
}

// FinalizeChange is a paid mutator transaction binding the contract method 0x55a6b050.
//
// Solidity: function finalizeChange(uint256 ambBlockNumber) returns()
func (_Vs *VsTransactorSession) FinalizeChange(ambBlockNumber *big.Int) (*types.Transaction, error) {
	return _Vs.Contract.FinalizeChange(&_Vs.TransactOpts, ambBlockNumber)
}

// FinalizeChange0 is a paid mutator transaction binding the contract method 0x75286211.
//
// Solidity: function finalizeChange() returns()
func (_Vs *VsTransactor) FinalizeChange0(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vs.contract.Transact(opts, "finalizeChange0")
}

// FinalizeChange0 is a paid mutator transaction binding the contract method 0x75286211.
//
// Solidity: function finalizeChange() returns()
func (_Vs *VsSession) FinalizeChange0() (*types.Transaction, error) {
	return _Vs.Contract.FinalizeChange0(&_Vs.TransactOpts)
}

// FinalizeChange0 is a paid mutator transaction binding the contract method 0x75286211.
//
// Solidity: function finalizeChange() returns()
func (_Vs *VsTransactorSession) FinalizeChange0() (*types.Transaction, error) {
	return _Vs.Contract.FinalizeChange0(&_Vs.TransactOpts)
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

// RemoveValidators is a paid mutator transaction binding the contract method 0x47ff07fc.
//
// Solidity: function removeValidators(uint256 ambBlockNumber) returns()
func (_Vs *VsTransactor) RemoveValidators(opts *bind.TransactOpts, ambBlockNumber *big.Int) (*types.Transaction, error) {
	return _Vs.contract.Transact(opts, "removeValidators", ambBlockNumber)
}

// RemoveValidators is a paid mutator transaction binding the contract method 0x47ff07fc.
//
// Solidity: function removeValidators(uint256 ambBlockNumber) returns()
func (_Vs *VsSession) RemoveValidators(ambBlockNumber *big.Int) (*types.Transaction, error) {
	return _Vs.Contract.RemoveValidators(&_Vs.TransactOpts, ambBlockNumber)
}

// RemoveValidators is a paid mutator transaction binding the contract method 0x47ff07fc.
//
// Solidity: function removeValidators(uint256 ambBlockNumber) returns()
func (_Vs *VsTransactorSession) RemoveValidators(ambBlockNumber *big.Int) (*types.Transaction, error) {
	return _Vs.Contract.RemoveValidators(&_Vs.TransactOpts, ambBlockNumber)
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
