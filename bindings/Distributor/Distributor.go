// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package distributor

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

// DistributorMetaData contains all meta data concerning the Distributor contract.
var DistributorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stafiStorageAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"round\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Claimed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"claimRoundList\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indexList\",\"type\":\"uint256[]\"},{\"internalType\":\"address[]\",\"name\":\"accountList\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amountList\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"merkleProofList\",\"type\":\"bytes32[][]\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"distributeFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"distributeSuperNodeFee\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCurrentNodeDepositAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"claimRound\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"isClaimed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"receiveEtherWithdrawal\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"claimRound\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"merkleRoot\",\"type\":\"bytes32\"}],\"name\":\"setMerkleRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// DistributorABI is the input ABI used to generate the binding from.
// Deprecated: Use DistributorMetaData.ABI instead.
var DistributorABI = DistributorMetaData.ABI

// Distributor is an auto generated Go binding around an Ethereum contract.
type Distributor struct {
	DistributorCaller     // Read-only binding to the contract
	DistributorTransactor // Write-only binding to the contract
	DistributorFilterer   // Log filterer for contract events
}

// DistributorCaller is an auto generated read-only Go binding around an Ethereum contract.
type DistributorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DistributorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DistributorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DistributorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DistributorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DistributorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DistributorSession struct {
	Contract     *Distributor      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DistributorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DistributorCallerSession struct {
	Contract *DistributorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// DistributorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DistributorTransactorSession struct {
	Contract     *DistributorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DistributorRaw is an auto generated low-level Go binding around an Ethereum contract.
type DistributorRaw struct {
	Contract *Distributor // Generic contract binding to access the raw methods on
}

// DistributorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DistributorCallerRaw struct {
	Contract *DistributorCaller // Generic read-only contract binding to access the raw methods on
}

// DistributorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DistributorTransactorRaw struct {
	Contract *DistributorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDistributor creates a new instance of Distributor, bound to a specific deployed contract.
func NewDistributor(address common.Address, backend bind.ContractBackend) (*Distributor, error) {
	contract, err := bindDistributor(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Distributor{DistributorCaller: DistributorCaller{contract: contract}, DistributorTransactor: DistributorTransactor{contract: contract}, DistributorFilterer: DistributorFilterer{contract: contract}}, nil
}

// NewDistributorCaller creates a new read-only instance of Distributor, bound to a specific deployed contract.
func NewDistributorCaller(address common.Address, caller bind.ContractCaller) (*DistributorCaller, error) {
	contract, err := bindDistributor(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DistributorCaller{contract: contract}, nil
}

// NewDistributorTransactor creates a new write-only instance of Distributor, bound to a specific deployed contract.
func NewDistributorTransactor(address common.Address, transactor bind.ContractTransactor) (*DistributorTransactor, error) {
	contract, err := bindDistributor(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DistributorTransactor{contract: contract}, nil
}

// NewDistributorFilterer creates a new log filterer instance of Distributor, bound to a specific deployed contract.
func NewDistributorFilterer(address common.Address, filterer bind.ContractFilterer) (*DistributorFilterer, error) {
	contract, err := bindDistributor(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DistributorFilterer{contract: contract}, nil
}

// bindDistributor binds a generic wrapper to an already deployed contract.
func bindDistributor(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DistributorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Distributor *DistributorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Distributor.Contract.DistributorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Distributor *DistributorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Distributor.Contract.DistributorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Distributor *DistributorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Distributor.Contract.DistributorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Distributor *DistributorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Distributor.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Distributor *DistributorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Distributor.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Distributor *DistributorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Distributor.Contract.contract.Transact(opts, method, params...)
}

// GetCurrentNodeDepositAmount is a free data retrieval call binding the contract method 0x83f3f086.
//
// Solidity: function getCurrentNodeDepositAmount() view returns(uint256)
func (_Distributor *DistributorCaller) GetCurrentNodeDepositAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Distributor.contract.Call(opts, &out, "getCurrentNodeDepositAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCurrentNodeDepositAmount is a free data retrieval call binding the contract method 0x83f3f086.
//
// Solidity: function getCurrentNodeDepositAmount() view returns(uint256)
func (_Distributor *DistributorSession) GetCurrentNodeDepositAmount() (*big.Int, error) {
	return _Distributor.Contract.GetCurrentNodeDepositAmount(&_Distributor.CallOpts)
}

// GetCurrentNodeDepositAmount is a free data retrieval call binding the contract method 0x83f3f086.
//
// Solidity: function getCurrentNodeDepositAmount() view returns(uint256)
func (_Distributor *DistributorCallerSession) GetCurrentNodeDepositAmount() (*big.Int, error) {
	return _Distributor.Contract.GetCurrentNodeDepositAmount(&_Distributor.CallOpts)
}

// IsClaimed is a free data retrieval call binding the contract method 0xf364c90c.
//
// Solidity: function isClaimed(uint256 claimRound, uint256 index) view returns(bool)
func (_Distributor *DistributorCaller) IsClaimed(opts *bind.CallOpts, claimRound *big.Int, index *big.Int) (bool, error) {
	var out []interface{}
	err := _Distributor.contract.Call(opts, &out, "isClaimed", claimRound, index)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsClaimed is a free data retrieval call binding the contract method 0xf364c90c.
//
// Solidity: function isClaimed(uint256 claimRound, uint256 index) view returns(bool)
func (_Distributor *DistributorSession) IsClaimed(claimRound *big.Int, index *big.Int) (bool, error) {
	return _Distributor.Contract.IsClaimed(&_Distributor.CallOpts, claimRound, index)
}

// IsClaimed is a free data retrieval call binding the contract method 0xf364c90c.
//
// Solidity: function isClaimed(uint256 claimRound, uint256 index) view returns(bool)
func (_Distributor *DistributorCallerSession) IsClaimed(claimRound *big.Int, index *big.Int) (bool, error) {
	return _Distributor.Contract.IsClaimed(&_Distributor.CallOpts, claimRound, index)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_Distributor *DistributorCaller) Version(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Distributor.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_Distributor *DistributorSession) Version() (uint8, error) {
	return _Distributor.Contract.Version(&_Distributor.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_Distributor *DistributorCallerSession) Version() (uint8, error) {
	return _Distributor.Contract.Version(&_Distributor.CallOpts)
}

// Claim is a paid mutator transaction binding the contract method 0x7fb4e918.
//
// Solidity: function claim(uint256[] claimRoundList, uint256[] indexList, address[] accountList, uint256[] amountList, bytes32[][] merkleProofList) returns()
func (_Distributor *DistributorTransactor) Claim(opts *bind.TransactOpts, claimRoundList []*big.Int, indexList []*big.Int, accountList []common.Address, amountList []*big.Int, merkleProofList [][][32]byte) (*types.Transaction, error) {
	return _Distributor.contract.Transact(opts, "claim", claimRoundList, indexList, accountList, amountList, merkleProofList)
}

// Claim is a paid mutator transaction binding the contract method 0x7fb4e918.
//
// Solidity: function claim(uint256[] claimRoundList, uint256[] indexList, address[] accountList, uint256[] amountList, bytes32[][] merkleProofList) returns()
func (_Distributor *DistributorSession) Claim(claimRoundList []*big.Int, indexList []*big.Int, accountList []common.Address, amountList []*big.Int, merkleProofList [][][32]byte) (*types.Transaction, error) {
	return _Distributor.Contract.Claim(&_Distributor.TransactOpts, claimRoundList, indexList, accountList, amountList, merkleProofList)
}

// Claim is a paid mutator transaction binding the contract method 0x7fb4e918.
//
// Solidity: function claim(uint256[] claimRoundList, uint256[] indexList, address[] accountList, uint256[] amountList, bytes32[][] merkleProofList) returns()
func (_Distributor *DistributorTransactorSession) Claim(claimRoundList []*big.Int, indexList []*big.Int, accountList []common.Address, amountList []*big.Int, merkleProofList [][][32]byte) (*types.Transaction, error) {
	return _Distributor.Contract.Claim(&_Distributor.TransactOpts, claimRoundList, indexList, accountList, amountList, merkleProofList)
}

// DistributeFee is a paid mutator transaction binding the contract method 0x05cc49dd.
//
// Solidity: function distributeFee(uint256 amount) payable returns()
func (_Distributor *DistributorTransactor) DistributeFee(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Distributor.contract.Transact(opts, "distributeFee", amount)
}

// DistributeFee is a paid mutator transaction binding the contract method 0x05cc49dd.
//
// Solidity: function distributeFee(uint256 amount) payable returns()
func (_Distributor *DistributorSession) DistributeFee(amount *big.Int) (*types.Transaction, error) {
	return _Distributor.Contract.DistributeFee(&_Distributor.TransactOpts, amount)
}

// DistributeFee is a paid mutator transaction binding the contract method 0x05cc49dd.
//
// Solidity: function distributeFee(uint256 amount) payable returns()
func (_Distributor *DistributorTransactorSession) DistributeFee(amount *big.Int) (*types.Transaction, error) {
	return _Distributor.Contract.DistributeFee(&_Distributor.TransactOpts, amount)
}

// DistributeSuperNodeFee is a paid mutator transaction binding the contract method 0x4fdcc7cd.
//
// Solidity: function distributeSuperNodeFee(uint256 amount) payable returns()
func (_Distributor *DistributorTransactor) DistributeSuperNodeFee(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Distributor.contract.Transact(opts, "distributeSuperNodeFee", amount)
}

// DistributeSuperNodeFee is a paid mutator transaction binding the contract method 0x4fdcc7cd.
//
// Solidity: function distributeSuperNodeFee(uint256 amount) payable returns()
func (_Distributor *DistributorSession) DistributeSuperNodeFee(amount *big.Int) (*types.Transaction, error) {
	return _Distributor.Contract.DistributeSuperNodeFee(&_Distributor.TransactOpts, amount)
}

// DistributeSuperNodeFee is a paid mutator transaction binding the contract method 0x4fdcc7cd.
//
// Solidity: function distributeSuperNodeFee(uint256 amount) payable returns()
func (_Distributor *DistributorTransactorSession) DistributeSuperNodeFee(amount *big.Int) (*types.Transaction, error) {
	return _Distributor.Contract.DistributeSuperNodeFee(&_Distributor.TransactOpts, amount)
}

// ReceiveEtherWithdrawal is a paid mutator transaction binding the contract method 0x0a019eaf.
//
// Solidity: function receiveEtherWithdrawal() payable returns()
func (_Distributor *DistributorTransactor) ReceiveEtherWithdrawal(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Distributor.contract.Transact(opts, "receiveEtherWithdrawal")
}

// ReceiveEtherWithdrawal is a paid mutator transaction binding the contract method 0x0a019eaf.
//
// Solidity: function receiveEtherWithdrawal() payable returns()
func (_Distributor *DistributorSession) ReceiveEtherWithdrawal() (*types.Transaction, error) {
	return _Distributor.Contract.ReceiveEtherWithdrawal(&_Distributor.TransactOpts)
}

// ReceiveEtherWithdrawal is a paid mutator transaction binding the contract method 0x0a019eaf.
//
// Solidity: function receiveEtherWithdrawal() payable returns()
func (_Distributor *DistributorTransactorSession) ReceiveEtherWithdrawal() (*types.Transaction, error) {
	return _Distributor.Contract.ReceiveEtherWithdrawal(&_Distributor.TransactOpts)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x18712c21.
//
// Solidity: function setMerkleRoot(uint256 claimRound, bytes32 merkleRoot) returns()
func (_Distributor *DistributorTransactor) SetMerkleRoot(opts *bind.TransactOpts, claimRound *big.Int, merkleRoot [32]byte) (*types.Transaction, error) {
	return _Distributor.contract.Transact(opts, "setMerkleRoot", claimRound, merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x18712c21.
//
// Solidity: function setMerkleRoot(uint256 claimRound, bytes32 merkleRoot) returns()
func (_Distributor *DistributorSession) SetMerkleRoot(claimRound *big.Int, merkleRoot [32]byte) (*types.Transaction, error) {
	return _Distributor.Contract.SetMerkleRoot(&_Distributor.TransactOpts, claimRound, merkleRoot)
}

// SetMerkleRoot is a paid mutator transaction binding the contract method 0x18712c21.
//
// Solidity: function setMerkleRoot(uint256 claimRound, bytes32 merkleRoot) returns()
func (_Distributor *DistributorTransactorSession) SetMerkleRoot(claimRound *big.Int, merkleRoot [32]byte) (*types.Transaction, error) {
	return _Distributor.Contract.SetMerkleRoot(&_Distributor.TransactOpts, claimRound, merkleRoot)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Distributor *DistributorTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Distributor.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Distributor *DistributorSession) Receive() (*types.Transaction, error) {
	return _Distributor.Contract.Receive(&_Distributor.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Distributor *DistributorTransactorSession) Receive() (*types.Transaction, error) {
	return _Distributor.Contract.Receive(&_Distributor.TransactOpts)
}

// DistributorClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the Distributor contract.
type DistributorClaimedIterator struct {
	Event *DistributorClaimed // Event containing the contract specifics and raw log

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
func (it *DistributorClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DistributorClaimed)
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
		it.Event = new(DistributorClaimed)
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
func (it *DistributorClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DistributorClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DistributorClaimed represents a Claimed event raised by the Distributor contract.
type DistributorClaimed struct {
	Round   *big.Int
	Index   *big.Int
	Account common.Address
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0xb94bf7f9302edf52a596286915a69b4b0685574cffdedd0712e3c62f2550f0ba.
//
// Solidity: event Claimed(uint256 round, uint256 index, address account, uint256 amount)
func (_Distributor *DistributorFilterer) FilterClaimed(opts *bind.FilterOpts) (*DistributorClaimedIterator, error) {

	logs, sub, err := _Distributor.contract.FilterLogs(opts, "Claimed")
	if err != nil {
		return nil, err
	}
	return &DistributorClaimedIterator{contract: _Distributor.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0xb94bf7f9302edf52a596286915a69b4b0685574cffdedd0712e3c62f2550f0ba.
//
// Solidity: event Claimed(uint256 round, uint256 index, address account, uint256 amount)
func (_Distributor *DistributorFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *DistributorClaimed) (event.Subscription, error) {

	logs, sub, err := _Distributor.contract.WatchLogs(opts, "Claimed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DistributorClaimed)
				if err := _Distributor.contract.UnpackLog(event, "Claimed", log); err != nil {
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

// ParseClaimed is a log parse operation binding the contract event 0xb94bf7f9302edf52a596286915a69b4b0685574cffdedd0712e3c62f2550f0ba.
//
// Solidity: event Claimed(uint256 round, uint256 index, address account, uint256 amount)
func (_Distributor *DistributorFilterer) ParseClaimed(log types.Log) (*DistributorClaimed, error) {
	event := new(DistributorClaimed)
	if err := _Distributor.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
