// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package proof

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ProofABI is the input ABI used to generate the binding from.
const ProofABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"fileHash\",\"type\":\"string\"}],\"name\":\"get\",\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint256\"},{\"name\":\"owner\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner\",\"type\":\"string\"},{\"name\":\"fileHash\",\"type\":\"string\"}],\"name\":\"set\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"status\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"fileHash\",\"type\":\"string\"}],\"name\":\"logFileAddedStatus\",\"type\":\"event\"}]"

// Proof is an auto generated Go binding around an Ethereum contract.
type Proof struct {
	ProofCaller     // Read-only binding to the contract
	ProofTransactor // Write-only binding to the contract
	ProofFilterer   // Log filterer for contract events
}

// ProofCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProofCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProofTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProofFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProofSession struct {
	Contract     *Proof            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProofCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProofCallerSession struct {
	Contract *ProofCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ProofTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProofTransactorSession struct {
	Contract     *ProofTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProofRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProofRaw struct {
	Contract *Proof // Generic contract binding to access the raw methods on
}

// ProofCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProofCallerRaw struct {
	Contract *ProofCaller // Generic read-only contract binding to access the raw methods on
}

// ProofTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProofTransactorRaw struct {
	Contract *ProofTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProof creates a new instance of Proof, bound to a specific deployed contract.
func NewProof(address common.Address, backend bind.ContractBackend) (*Proof, error) {
	contract, err := bindProof(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Proof{ProofCaller: ProofCaller{contract: contract}, ProofTransactor: ProofTransactor{contract: contract}, ProofFilterer: ProofFilterer{contract: contract}}, nil
}

// NewProofCaller creates a new read-only instance of Proof, bound to a specific deployed contract.
func NewProofCaller(address common.Address, caller bind.ContractCaller) (*ProofCaller, error) {
	contract, err := bindProof(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProofCaller{contract: contract}, nil
}

// NewProofTransactor creates a new write-only instance of Proof, bound to a specific deployed contract.
func NewProofTransactor(address common.Address, transactor bind.ContractTransactor) (*ProofTransactor, error) {
	contract, err := bindProof(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProofTransactor{contract: contract}, nil
}

// NewProofFilterer creates a new log filterer instance of Proof, bound to a specific deployed contract.
func NewProofFilterer(address common.Address, filterer bind.ContractFilterer) (*ProofFilterer, error) {
	contract, err := bindProof(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProofFilterer{contract: contract}, nil
}

// bindProof binds a generic wrapper to an already deployed contract.
func bindProof(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProofABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proof *ProofRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Proof.Contract.ProofCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proof *ProofRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proof.Contract.ProofTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proof *ProofRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proof.Contract.ProofTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Proof *ProofCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Proof.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Proof *ProofTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Proof.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Proof *ProofTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Proof.Contract.contract.Transact(opts, method, params...)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string fileHash) constant returns(uint256 timestamp, string owner)
func (_Proof *ProofCaller) Get(opts *bind.CallOpts, fileHash string) (struct {
	Timestamp *big.Int
	Owner     string
}, error) {
	ret := new(struct {
		Timestamp *big.Int
		Owner     string
	})
	out := ret
	err := _Proof.contract.Call(opts, out, "get", fileHash)
	return *ret, err
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string fileHash) constant returns(uint256 timestamp, string owner)
func (_Proof *ProofSession) Get(fileHash string) (struct {
	Timestamp *big.Int
	Owner     string
}, error) {
	return _Proof.Contract.Get(&_Proof.CallOpts, fileHash)
}

// Get is a free data retrieval call binding the contract method 0x693ec85e.
//
// Solidity: function get(string fileHash) constant returns(uint256 timestamp, string owner)
func (_Proof *ProofCallerSession) Get(fileHash string) (struct {
	Timestamp *big.Int
	Owner     string
}, error) {
	return _Proof.Contract.Get(&_Proof.CallOpts, fileHash)
}

// Set is a paid mutator transaction binding the contract method 0xe942b516.
//
// Solidity: function set(string owner, string fileHash) returns()
func (_Proof *ProofTransactor) Set(opts *bind.TransactOpts, owner string, fileHash string) (*types.Transaction, error) {
	return _Proof.contract.Transact(opts, "set", owner, fileHash)
}

// Set is a paid mutator transaction binding the contract method 0xe942b516.
//
// Solidity: function set(string owner, string fileHash) returns()
func (_Proof *ProofSession) Set(owner string, fileHash string) (*types.Transaction, error) {
	return _Proof.Contract.Set(&_Proof.TransactOpts, owner, fileHash)
}

// Set is a paid mutator transaction binding the contract method 0xe942b516.
//
// Solidity: function set(string owner, string fileHash) returns()
func (_Proof *ProofTransactorSession) Set(owner string, fileHash string) (*types.Transaction, error) {
	return _Proof.Contract.Set(&_Proof.TransactOpts, owner, fileHash)
}

// ProofLogFileAddedStatusIterator is returned from FilterLogFileAddedStatus and is used to iterate over the raw logs and unpacked data for LogFileAddedStatus events raised by the Proof contract.
type ProofLogFileAddedStatusIterator struct {
	Event *ProofLogFileAddedStatus // Event containing the contract specifics and raw log

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
func (it *ProofLogFileAddedStatusIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ProofLogFileAddedStatus)
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
		it.Event = new(ProofLogFileAddedStatus)
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
func (it *ProofLogFileAddedStatusIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ProofLogFileAddedStatusIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ProofLogFileAddedStatus represents a LogFileAddedStatus event raised by the Proof contract.
type ProofLogFileAddedStatus struct {
	Status    bool
	Timestamp *big.Int
	Owner     string
	FileHash  string
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLogFileAddedStatus is a free log retrieval operation binding the contract event 0x0d3bbc3c02da6ed436712ca1a0f626f1269df703a105f034e4637c7b10fb7ba5.
//
// Solidity: event logFileAddedStatus(bool status, uint256 timestamp, string owner, string fileHash)
func (_Proof *ProofFilterer) FilterLogFileAddedStatus(opts *bind.FilterOpts) (*ProofLogFileAddedStatusIterator, error) {

	logs, sub, err := _Proof.contract.FilterLogs(opts, "logFileAddedStatus")
	if err != nil {
		return nil, err
	}
	return &ProofLogFileAddedStatusIterator{contract: _Proof.contract, event: "logFileAddedStatus", logs: logs, sub: sub}, nil
}

// WatchLogFileAddedStatus is a free log subscription operation binding the contract event 0x0d3bbc3c02da6ed436712ca1a0f626f1269df703a105f034e4637c7b10fb7ba5.
//
// Solidity: event logFileAddedStatus(bool status, uint256 timestamp, string owner, string fileHash)
func (_Proof *ProofFilterer) WatchLogFileAddedStatus(opts *bind.WatchOpts, sink chan<- *ProofLogFileAddedStatus) (event.Subscription, error) {

	logs, sub, err := _Proof.contract.WatchLogs(opts, "logFileAddedStatus")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ProofLogFileAddedStatus)
				if err := _Proof.contract.UnpackLog(event, "logFileAddedStatus", log); err != nil {
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
