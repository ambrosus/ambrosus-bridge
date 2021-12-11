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

// AmbBridgeWithdraw is an auto generated low-level Go binding around an user-defined struct.
type AmbBridgeWithdraw struct {
	FromAddress common.Address
	ToAddress   common.Address
	Amount      *big.Int
}

// AmbMetaData contains all meta data concerning the Amb contract.
var AmbMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"withdraws_hash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"fromAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structAmbBridge.Withdraw[]\",\"name\":\"withdraws\",\"type\":\"tuple[]\"}],\"name\":\"Test\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"eventTest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"toAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// EventTest is a paid mutator transaction binding the contract method 0xac002055.
//
// Solidity: function eventTest() returns()
func (_Amb *AmbTransactor) EventTest(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "eventTest")
}

// EventTest is a paid mutator transaction binding the contract method 0xac002055.
//
// Solidity: function eventTest() returns()
func (_Amb *AmbSession) EventTest() (*types.Transaction, error) {
	return _Amb.Contract.EventTest(&_Amb.TransactOpts)
}

// EventTest is a paid mutator transaction binding the contract method 0xac002055.
//
// Solidity: function eventTest() returns()
func (_Amb *AmbTransactorSession) EventTest() (*types.Transaction, error) {
	return _Amb.Contract.EventTest(&_Amb.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address toAddr, uint256 amount) returns()
func (_Amb *AmbTransactor) Withdraw(opts *bind.TransactOpts, toAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Amb.contract.Transact(opts, "withdraw", toAddr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address toAddr, uint256 amount) returns()
func (_Amb *AmbSession) Withdraw(toAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.Withdraw(&_Amb.TransactOpts, toAddr, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xf3fef3a3.
//
// Solidity: function withdraw(address toAddr, uint256 amount) returns()
func (_Amb *AmbTransactorSession) Withdraw(toAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Amb.Contract.Withdraw(&_Amb.TransactOpts, toAddr, amount)
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
	WithdrawsHash [32]byte
	Withdraws     []AmbBridgeWithdraw
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTest is a free log retrieval operation binding the contract event 0x88d88376da47bd16f6e8e45064996475d1536ab348a4f812de2f22518e06ee2d.
//
// Solidity: event Test(bytes32 indexed withdraws_hash, (address,address,uint256)[] withdraws)
func (_Amb *AmbFilterer) FilterTest(opts *bind.FilterOpts, withdraws_hash [][32]byte) (*AmbTestIterator, error) {

	var withdraws_hashRule []interface{}
	for _, withdraws_hashItem := range withdraws_hash {
		withdraws_hashRule = append(withdraws_hashRule, withdraws_hashItem)
	}

	logs, sub, err := _Amb.contract.FilterLogs(opts, "Test", withdraws_hashRule)
	if err != nil {
		return nil, err
	}
	return &AmbTestIterator{contract: _Amb.contract, event: "Test", logs: logs, sub: sub}, nil
}

// WatchTest is a free log subscription operation binding the contract event 0x88d88376da47bd16f6e8e45064996475d1536ab348a4f812de2f22518e06ee2d.
//
// Solidity: event Test(bytes32 indexed withdraws_hash, (address,address,uint256)[] withdraws)
func (_Amb *AmbFilterer) WatchTest(opts *bind.WatchOpts, sink chan<- *AmbTest, withdraws_hash [][32]byte) (event.Subscription, error) {

	var withdraws_hashRule []interface{}
	for _, withdraws_hashItem := range withdraws_hash {
		withdraws_hashRule = append(withdraws_hashRule, withdraws_hashItem)
	}

	logs, sub, err := _Amb.contract.WatchLogs(opts, "Test", withdraws_hashRule)
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

// ParseTest is a log parse operation binding the contract event 0x88d88376da47bd16f6e8e45064996475d1536ab348a4f812de2f22518e06ee2d.
//
// Solidity: event Test(bytes32 indexed withdraws_hash, (address,address,uint256)[] withdraws)
func (_Amb *AmbFilterer) ParseTest(log types.Log) (*AmbTest, error) {
	event := new(AmbTest)
	if err := _Amb.contract.UnpackLog(event, "Test", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
