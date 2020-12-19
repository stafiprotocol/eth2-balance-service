// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package NetworkBalance

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// NetworkBalanceABI is the input ABI used to generate the binding from.
const NetworkBalanceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stafiStorageAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalEth\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakingEth\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rethSupply\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"BalancesSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"block\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalEth\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"stakingEth\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"rethSupply\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"BalancesUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getBalancesBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getETHStakingRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStakingETHBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalETHBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalRETHSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_block\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_totalEth\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_stakingEth\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_rethSupply\",\"type\":\"uint256\"}],\"name\":\"submitBalances\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// NetworkBalance is an auto generated Go binding around an Ethereum contract.
type NetworkBalance struct {
	NetworkBalanceCaller     // Read-only binding to the contract
	NetworkBalanceTransactor // Write-only binding to the contract
	NetworkBalanceFilterer   // Log filterer for contract events
}

// NetworkBalanceCaller is an auto generated read-only Go binding around an Ethereum contract.
type NetworkBalanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkBalanceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NetworkBalanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkBalanceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NetworkBalanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NetworkBalanceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NetworkBalanceSession struct {
	Contract     *NetworkBalance   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NetworkBalanceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NetworkBalanceCallerSession struct {
	Contract *NetworkBalanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// NetworkBalanceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NetworkBalanceTransactorSession struct {
	Contract     *NetworkBalanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// NetworkBalanceRaw is an auto generated low-level Go binding around an Ethereum contract.
type NetworkBalanceRaw struct {
	Contract *NetworkBalance // Generic contract binding to access the raw methods on
}

// NetworkBalanceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NetworkBalanceCallerRaw struct {
	Contract *NetworkBalanceCaller // Generic read-only contract binding to access the raw methods on
}

// NetworkBalanceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NetworkBalanceTransactorRaw struct {
	Contract *NetworkBalanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNetworkBalance creates a new instance of NetworkBalance, bound to a specific deployed contract.
func NewNetworkBalance(address common.Address, backend bind.ContractBackend) (*NetworkBalance, error) {
	contract, err := bindNetworkBalance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NetworkBalance{NetworkBalanceCaller: NetworkBalanceCaller{contract: contract}, NetworkBalanceTransactor: NetworkBalanceTransactor{contract: contract}, NetworkBalanceFilterer: NetworkBalanceFilterer{contract: contract}}, nil
}

// NewNetworkBalanceCaller creates a new read-only instance of NetworkBalance, bound to a specific deployed contract.
func NewNetworkBalanceCaller(address common.Address, caller bind.ContractCaller) (*NetworkBalanceCaller, error) {
	contract, err := bindNetworkBalance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkBalanceCaller{contract: contract}, nil
}

// NewNetworkBalanceTransactor creates a new write-only instance of NetworkBalance, bound to a specific deployed contract.
func NewNetworkBalanceTransactor(address common.Address, transactor bind.ContractTransactor) (*NetworkBalanceTransactor, error) {
	contract, err := bindNetworkBalance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NetworkBalanceTransactor{contract: contract}, nil
}

// NewNetworkBalanceFilterer creates a new log filterer instance of NetworkBalance, bound to a specific deployed contract.
func NewNetworkBalanceFilterer(address common.Address, filterer bind.ContractFilterer) (*NetworkBalanceFilterer, error) {
	contract, err := bindNetworkBalance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NetworkBalanceFilterer{contract: contract}, nil
}

// bindNetworkBalance binds a generic wrapper to an already deployed contract.
func bindNetworkBalance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NetworkBalanceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkBalance *NetworkBalanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkBalance.Contract.NetworkBalanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkBalance *NetworkBalanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkBalance.Contract.NetworkBalanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkBalance *NetworkBalanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkBalance.Contract.NetworkBalanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NetworkBalance *NetworkBalanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NetworkBalance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NetworkBalance *NetworkBalanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NetworkBalance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NetworkBalance *NetworkBalanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NetworkBalance.Contract.contract.Transact(opts, method, params...)
}

// GetBalancesBlock is a free data retrieval call binding the contract method 0x9100c13d.
//
// Solidity: function getBalancesBlock() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCaller) GetBalancesBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkBalance.contract.Call(opts, &out, "getBalancesBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalancesBlock is a free data retrieval call binding the contract method 0x9100c13d.
//
// Solidity: function getBalancesBlock() view returns(uint256)
func (_NetworkBalance *NetworkBalanceSession) GetBalancesBlock() (*big.Int, error) {
	return _NetworkBalance.Contract.GetBalancesBlock(&_NetworkBalance.CallOpts)
}

// GetBalancesBlock is a free data retrieval call binding the contract method 0x9100c13d.
//
// Solidity: function getBalancesBlock() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCallerSession) GetBalancesBlock() (*big.Int, error) {
	return _NetworkBalance.Contract.GetBalancesBlock(&_NetworkBalance.CallOpts)
}

// GetETHStakingRate is a free data retrieval call binding the contract method 0xcd482567.
//
// Solidity: function getETHStakingRate() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCaller) GetETHStakingRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkBalance.contract.Call(opts, &out, "getETHStakingRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetETHStakingRate is a free data retrieval call binding the contract method 0xcd482567.
//
// Solidity: function getETHStakingRate() view returns(uint256)
func (_NetworkBalance *NetworkBalanceSession) GetETHStakingRate() (*big.Int, error) {
	return _NetworkBalance.Contract.GetETHStakingRate(&_NetworkBalance.CallOpts)
}

// GetETHStakingRate is a free data retrieval call binding the contract method 0xcd482567.
//
// Solidity: function getETHStakingRate() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCallerSession) GetETHStakingRate() (*big.Int, error) {
	return _NetworkBalance.Contract.GetETHStakingRate(&_NetworkBalance.CallOpts)
}

// GetStakingETHBalance is a free data retrieval call binding the contract method 0xf1eda634.
//
// Solidity: function getStakingETHBalance() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCaller) GetStakingETHBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkBalance.contract.Call(opts, &out, "getStakingETHBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakingETHBalance is a free data retrieval call binding the contract method 0xf1eda634.
//
// Solidity: function getStakingETHBalance() view returns(uint256)
func (_NetworkBalance *NetworkBalanceSession) GetStakingETHBalance() (*big.Int, error) {
	return _NetworkBalance.Contract.GetStakingETHBalance(&_NetworkBalance.CallOpts)
}

// GetStakingETHBalance is a free data retrieval call binding the contract method 0xf1eda634.
//
// Solidity: function getStakingETHBalance() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCallerSession) GetStakingETHBalance() (*big.Int, error) {
	return _NetworkBalance.Contract.GetStakingETHBalance(&_NetworkBalance.CallOpts)
}

// GetTotalETHBalance is a free data retrieval call binding the contract method 0x964d042c.
//
// Solidity: function getTotalETHBalance() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCaller) GetTotalETHBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkBalance.contract.Call(opts, &out, "getTotalETHBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalETHBalance is a free data retrieval call binding the contract method 0x964d042c.
//
// Solidity: function getTotalETHBalance() view returns(uint256)
func (_NetworkBalance *NetworkBalanceSession) GetTotalETHBalance() (*big.Int, error) {
	return _NetworkBalance.Contract.GetTotalETHBalance(&_NetworkBalance.CallOpts)
}

// GetTotalETHBalance is a free data retrieval call binding the contract method 0x964d042c.
//
// Solidity: function getTotalETHBalance() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCallerSession) GetTotalETHBalance() (*big.Int, error) {
	return _NetworkBalance.Contract.GetTotalETHBalance(&_NetworkBalance.CallOpts)
}

// GetTotalRETHSupply is a free data retrieval call binding the contract method 0xc4c8d0ad.
//
// Solidity: function getTotalRETHSupply() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCaller) GetTotalRETHSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NetworkBalance.contract.Call(opts, &out, "getTotalRETHSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalRETHSupply is a free data retrieval call binding the contract method 0xc4c8d0ad.
//
// Solidity: function getTotalRETHSupply() view returns(uint256)
func (_NetworkBalance *NetworkBalanceSession) GetTotalRETHSupply() (*big.Int, error) {
	return _NetworkBalance.Contract.GetTotalRETHSupply(&_NetworkBalance.CallOpts)
}

// GetTotalRETHSupply is a free data retrieval call binding the contract method 0xc4c8d0ad.
//
// Solidity: function getTotalRETHSupply() view returns(uint256)
func (_NetworkBalance *NetworkBalanceCallerSession) GetTotalRETHSupply() (*big.Int, error) {
	return _NetworkBalance.Contract.GetTotalRETHSupply(&_NetworkBalance.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_NetworkBalance *NetworkBalanceCaller) Version(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _NetworkBalance.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_NetworkBalance *NetworkBalanceSession) Version() (uint8, error) {
	return _NetworkBalance.Contract.Version(&_NetworkBalance.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_NetworkBalance *NetworkBalanceCallerSession) Version() (uint8, error) {
	return _NetworkBalance.Contract.Version(&_NetworkBalance.CallOpts)
}

// SubmitBalances is a paid mutator transaction binding the contract method 0x32db5470.
//
// Solidity: function submitBalances(uint256 _block, uint256 _totalEth, uint256 _stakingEth, uint256 _rethSupply) returns()
func (_NetworkBalance *NetworkBalanceTransactor) SubmitBalances(opts *bind.TransactOpts, _block *big.Int, _totalEth *big.Int, _stakingEth *big.Int, _rethSupply *big.Int) (*types.Transaction, error) {
	return _NetworkBalance.contract.Transact(opts, "submitBalances", _block, _totalEth, _stakingEth, _rethSupply)
}

// SubmitBalances is a paid mutator transaction binding the contract method 0x32db5470.
//
// Solidity: function submitBalances(uint256 _block, uint256 _totalEth, uint256 _stakingEth, uint256 _rethSupply) returns()
func (_NetworkBalance *NetworkBalanceSession) SubmitBalances(_block *big.Int, _totalEth *big.Int, _stakingEth *big.Int, _rethSupply *big.Int) (*types.Transaction, error) {
	return _NetworkBalance.Contract.SubmitBalances(&_NetworkBalance.TransactOpts, _block, _totalEth, _stakingEth, _rethSupply)
}

// SubmitBalances is a paid mutator transaction binding the contract method 0x32db5470.
//
// Solidity: function submitBalances(uint256 _block, uint256 _totalEth, uint256 _stakingEth, uint256 _rethSupply) returns()
func (_NetworkBalance *NetworkBalanceTransactorSession) SubmitBalances(_block *big.Int, _totalEth *big.Int, _stakingEth *big.Int, _rethSupply *big.Int) (*types.Transaction, error) {
	return _NetworkBalance.Contract.SubmitBalances(&_NetworkBalance.TransactOpts, _block, _totalEth, _stakingEth, _rethSupply)
}

// NetworkBalanceBalancesSubmittedIterator is returned from FilterBalancesSubmitted and is used to iterate over the raw logs and unpacked data for BalancesSubmitted events raised by the NetworkBalance contract.
type NetworkBalanceBalancesSubmittedIterator struct {
	Event *NetworkBalanceBalancesSubmitted // Event containing the contract specifics and raw log

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
func (it *NetworkBalanceBalancesSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkBalanceBalancesSubmitted)
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
		it.Event = new(NetworkBalanceBalancesSubmitted)
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
func (it *NetworkBalanceBalancesSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkBalanceBalancesSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkBalanceBalancesSubmitted represents a BalancesSubmitted event raised by the NetworkBalance contract.
type NetworkBalanceBalancesSubmitted struct {
	From       common.Address
	Block      *big.Int
	TotalEth   *big.Int
	StakingEth *big.Int
	RethSupply *big.Int
	Time       *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBalancesSubmitted is a free log retrieval operation binding the contract event 0xe657a6d6957f4fabb37b86d4d6571e82df061bd2d8a3ede5d197b0b98a5a1bdf.
//
// Solidity: event BalancesSubmitted(address indexed from, uint256 block, uint256 totalEth, uint256 stakingEth, uint256 rethSupply, uint256 time)
func (_NetworkBalance *NetworkBalanceFilterer) FilterBalancesSubmitted(opts *bind.FilterOpts, from []common.Address) (*NetworkBalanceBalancesSubmittedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _NetworkBalance.contract.FilterLogs(opts, "BalancesSubmitted", fromRule)
	if err != nil {
		return nil, err
	}
	return &NetworkBalanceBalancesSubmittedIterator{contract: _NetworkBalance.contract, event: "BalancesSubmitted", logs: logs, sub: sub}, nil
}

// WatchBalancesSubmitted is a free log subscription operation binding the contract event 0xe657a6d6957f4fabb37b86d4d6571e82df061bd2d8a3ede5d197b0b98a5a1bdf.
//
// Solidity: event BalancesSubmitted(address indexed from, uint256 block, uint256 totalEth, uint256 stakingEth, uint256 rethSupply, uint256 time)
func (_NetworkBalance *NetworkBalanceFilterer) WatchBalancesSubmitted(opts *bind.WatchOpts, sink chan<- *NetworkBalanceBalancesSubmitted, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _NetworkBalance.contract.WatchLogs(opts, "BalancesSubmitted", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkBalanceBalancesSubmitted)
				if err := _NetworkBalance.contract.UnpackLog(event, "BalancesSubmitted", log); err != nil {
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

// ParseBalancesSubmitted is a log parse operation binding the contract event 0xe657a6d6957f4fabb37b86d4d6571e82df061bd2d8a3ede5d197b0b98a5a1bdf.
//
// Solidity: event BalancesSubmitted(address indexed from, uint256 block, uint256 totalEth, uint256 stakingEth, uint256 rethSupply, uint256 time)
func (_NetworkBalance *NetworkBalanceFilterer) ParseBalancesSubmitted(log types.Log) (*NetworkBalanceBalancesSubmitted, error) {
	event := new(NetworkBalanceBalancesSubmitted)
	if err := _NetworkBalance.contract.UnpackLog(event, "BalancesSubmitted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// NetworkBalanceBalancesUpdatedIterator is returned from FilterBalancesUpdated and is used to iterate over the raw logs and unpacked data for BalancesUpdated events raised by the NetworkBalance contract.
type NetworkBalanceBalancesUpdatedIterator struct {
	Event *NetworkBalanceBalancesUpdated // Event containing the contract specifics and raw log

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
func (it *NetworkBalanceBalancesUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NetworkBalanceBalancesUpdated)
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
		it.Event = new(NetworkBalanceBalancesUpdated)
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
func (it *NetworkBalanceBalancesUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NetworkBalanceBalancesUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NetworkBalanceBalancesUpdated represents a BalancesUpdated event raised by the NetworkBalance contract.
type NetworkBalanceBalancesUpdated struct {
	Block      *big.Int
	TotalEth   *big.Int
	StakingEth *big.Int
	RethSupply *big.Int
	Time       *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBalancesUpdated is a free log retrieval operation binding the contract event 0x7bbbb137fdad433d6168b1c75c714c72b8abe8d07460f0c0b433063e7bf1f394.
//
// Solidity: event BalancesUpdated(uint256 block, uint256 totalEth, uint256 stakingEth, uint256 rethSupply, uint256 time)
func (_NetworkBalance *NetworkBalanceFilterer) FilterBalancesUpdated(opts *bind.FilterOpts) (*NetworkBalanceBalancesUpdatedIterator, error) {

	logs, sub, err := _NetworkBalance.contract.FilterLogs(opts, "BalancesUpdated")
	if err != nil {
		return nil, err
	}
	return &NetworkBalanceBalancesUpdatedIterator{contract: _NetworkBalance.contract, event: "BalancesUpdated", logs: logs, sub: sub}, nil
}

// WatchBalancesUpdated is a free log subscription operation binding the contract event 0x7bbbb137fdad433d6168b1c75c714c72b8abe8d07460f0c0b433063e7bf1f394.
//
// Solidity: event BalancesUpdated(uint256 block, uint256 totalEth, uint256 stakingEth, uint256 rethSupply, uint256 time)
func (_NetworkBalance *NetworkBalanceFilterer) WatchBalancesUpdated(opts *bind.WatchOpts, sink chan<- *NetworkBalanceBalancesUpdated) (event.Subscription, error) {

	logs, sub, err := _NetworkBalance.contract.WatchLogs(opts, "BalancesUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NetworkBalanceBalancesUpdated)
				if err := _NetworkBalance.contract.UnpackLog(event, "BalancesUpdated", log); err != nil {
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

// ParseBalancesUpdated is a log parse operation binding the contract event 0x7bbbb137fdad433d6168b1c75c714c72b8abe8d07460f0c0b433063e7bf1f394.
//
// Solidity: event BalancesUpdated(uint256 block, uint256 totalEth, uint256 stakingEth, uint256 rethSupply, uint256 time)
func (_NetworkBalance *NetworkBalanceFilterer) ParseBalancesUpdated(log types.Log) (*NetworkBalanceBalancesUpdated, error) {
	event := new(NetworkBalanceBalancesUpdated)
	if err := _NetworkBalance.contract.UnpackLog(event, "BalancesUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
