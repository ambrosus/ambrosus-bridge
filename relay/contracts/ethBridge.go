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

// EthBridgeBlock is an auto generated low-level Go binding around an user-defined struct.
type EthBridgeBlock struct {
	P1                    []byte
	PrevHashOrReceiptRoot [32]byte
	P2                    []byte
	Timestamp             []byte
	P3                    []byte
	Signature             []byte
}

// EthBridgeWithdraw is an auto generated low-level Go binding around an user-defined struct.
type EthBridgeWithdraw struct {
	FromAddress common.Address
	ToAddress   common.Address
	Amount      *big.Int
}

// EthMetaData contains all meta data concerning the Eth contract.
var EthMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"validator_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawEvent\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"p1\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"prevHashOrReceiptRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"p2\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"timestamp\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"p3\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structEthBridge.Block[]\",\"name\":\"blocks\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"fromAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"toAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"internalType\":\"structEthBridge.Withdraw[]\",\"name\":\"events\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes[]\",\"name\":\"proof\",\"type\":\"bytes[]\"}],\"name\":\"TestAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"bloom\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"topicHash\",\"type\":\"bytes\"}],\"name\":\"TestBloom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"proof\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"eventToSearch\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"receiptsRoot\",\"type\":\"bytes32\"}],\"name\":\"TestReceiptsProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"b\",\"type\":\"bytes\"}],\"name\":\"bytesToUint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"proof\",\"type\":\"bytes[]\"},{\"internalType\":\"bytes\",\"name\":\"eventToSearch\",\"type\":\"bytes\"}],\"name\":\"calcReceiptsRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// TestAll is a paid mutator transaction binding the contract method 0x3315de1f.
//
// Solidity: function TestAll((bytes,bytes32,bytes,bytes,bytes,bytes)[] blocks, (address,address,uint256)[] events, bytes[] proof) returns()
func (_Eth *EthTransactor) TestAll(opts *bind.TransactOpts, blocks []EthBridgeBlock, events []EthBridgeWithdraw, proof [][]byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "TestAll", blocks, events, proof)
}

// TestAll is a paid mutator transaction binding the contract method 0x3315de1f.
//
// Solidity: function TestAll((bytes,bytes32,bytes,bytes,bytes,bytes)[] blocks, (address,address,uint256)[] events, bytes[] proof) returns()
func (_Eth *EthSession) TestAll(blocks []EthBridgeBlock, events []EthBridgeWithdraw, proof [][]byte) (*types.Transaction, error) {
	return _Eth.Contract.TestAll(&_Eth.TransactOpts, blocks, events, proof)
}

// TestAll is a paid mutator transaction binding the contract method 0x3315de1f.
//
// Solidity: function TestAll((bytes,bytes32,bytes,bytes,bytes,bytes)[] blocks, (address,address,uint256)[] events, bytes[] proof) returns()
func (_Eth *EthTransactorSession) TestAll(blocks []EthBridgeBlock, events []EthBridgeWithdraw, proof [][]byte) (*types.Transaction, error) {
	return _Eth.Contract.TestAll(&_Eth.TransactOpts, blocks, events, proof)
}

// TestBloom is a paid mutator transaction binding the contract method 0xb11a8408.
//
// Solidity: function TestBloom(bytes bloom, bytes topicHash) returns(bool)
func (_Eth *EthTransactor) TestBloom(opts *bind.TransactOpts, bloom []byte, topicHash []byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "TestBloom", bloom, topicHash)
}

// TestBloom is a paid mutator transaction binding the contract method 0xb11a8408.
//
// Solidity: function TestBloom(bytes bloom, bytes topicHash) returns(bool)
func (_Eth *EthSession) TestBloom(bloom []byte, topicHash []byte) (*types.Transaction, error) {
	return _Eth.Contract.TestBloom(&_Eth.TransactOpts, bloom, topicHash)
}

// TestBloom is a paid mutator transaction binding the contract method 0xb11a8408.
//
// Solidity: function TestBloom(bytes bloom, bytes topicHash) returns(bool)
func (_Eth *EthTransactorSession) TestBloom(bloom []byte, topicHash []byte) (*types.Transaction, error) {
	return _Eth.Contract.TestBloom(&_Eth.TransactOpts, bloom, topicHash)
}

// TestReceiptsProof is a paid mutator transaction binding the contract method 0x98f4333f.
//
// Solidity: function TestReceiptsProof(bytes[] proof, bytes eventToSearch, bytes32 receiptsRoot) returns()
func (_Eth *EthTransactor) TestReceiptsProof(opts *bind.TransactOpts, proof [][]byte, eventToSearch []byte, receiptsRoot [32]byte) (*types.Transaction, error) {
	return _Eth.contract.Transact(opts, "TestReceiptsProof", proof, eventToSearch, receiptsRoot)
}

// TestReceiptsProof is a paid mutator transaction binding the contract method 0x98f4333f.
//
// Solidity: function TestReceiptsProof(bytes[] proof, bytes eventToSearch, bytes32 receiptsRoot) returns()
func (_Eth *EthSession) TestReceiptsProof(proof [][]byte, eventToSearch []byte, receiptsRoot [32]byte) (*types.Transaction, error) {
	return _Eth.Contract.TestReceiptsProof(&_Eth.TransactOpts, proof, eventToSearch, receiptsRoot)
}

// TestReceiptsProof is a paid mutator transaction binding the contract method 0x98f4333f.
//
// Solidity: function TestReceiptsProof(bytes[] proof, bytes eventToSearch, bytes32 receiptsRoot) returns()
func (_Eth *EthTransactorSession) TestReceiptsProof(proof [][]byte, eventToSearch []byte, receiptsRoot [32]byte) (*types.Transaction, error) {
	return _Eth.Contract.TestReceiptsProof(&_Eth.TransactOpts, proof, eventToSearch, receiptsRoot)
}

// EthWithdrawEventIterator is returned from FilterWithdrawEvent and is used to iterate over the raw logs and unpacked data for WithdrawEvent events raised by the Eth contract.
type EthWithdrawEventIterator struct {
	Event *EthWithdrawEvent // Event containing the contract specifics and raw log

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
func (it *EthWithdrawEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EthWithdrawEvent)
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
		it.Event = new(EthWithdrawEvent)
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
func (it *EthWithdrawEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EthWithdrawEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EthWithdrawEvent represents a WithdrawEvent event raised by the Eth contract.
type EthWithdrawEvent struct {
	From   common.Address
	To     common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawEvent is a free log retrieval operation binding the contract event 0x93cb7b4ba12c5bb07f02e52c4e43788d8f4db1e66e9d30aaaeffc5ab325b810c.
//
// Solidity: event WithdrawEvent(address indexed from, address indexed to, uint256 amount)
func (_Eth *EthFilterer) FilterWithdrawEvent(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*EthWithdrawEventIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Eth.contract.FilterLogs(opts, "WithdrawEvent", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &EthWithdrawEventIterator{contract: _Eth.contract, event: "WithdrawEvent", logs: logs, sub: sub}, nil
}

// WatchWithdrawEvent is a free log subscription operation binding the contract event 0x93cb7b4ba12c5bb07f02e52c4e43788d8f4db1e66e9d30aaaeffc5ab325b810c.
//
// Solidity: event WithdrawEvent(address indexed from, address indexed to, uint256 amount)
func (_Eth *EthFilterer) WatchWithdrawEvent(opts *bind.WatchOpts, sink chan<- *EthWithdrawEvent, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Eth.contract.WatchLogs(opts, "WithdrawEvent", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EthWithdrawEvent)
				if err := _Eth.contract.UnpackLog(event, "WithdrawEvent", log); err != nil {
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

// ParseWithdrawEvent is a log parse operation binding the contract event 0x93cb7b4ba12c5bb07f02e52c4e43788d8f4db1e66e9d30aaaeffc5ab325b810c.
//
// Solidity: event WithdrawEvent(address indexed from, address indexed to, uint256 amount)
func (_Eth *EthFilterer) ParseWithdrawEvent(log types.Log) (*EthWithdrawEvent, error) {
	event := new(EthWithdrawEvent)
	if err := _Eth.contract.UnpackLog(event, "WithdrawEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
