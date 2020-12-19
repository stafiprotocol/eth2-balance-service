// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PoolBalance

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

// PoolBalanceABI is the input ABI used to generate the binding from.
const PoolBalanceABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stafiStorageAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"enumDepositType\",\"name\":\"_depositType\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"EtherDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"EtherWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"StatusUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"close\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dissolve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDepositType\",\"outputs\":[{\"internalType\":\"enumDepositType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeDepositAssigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeRefundBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPlatformDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStakingEndBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStakingStartBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStatus\",\"outputs\":[{\"internalType\":\"enumStakingPoolStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStatusBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStatusTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUserDepositAssigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUserDepositAssignedTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUserDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nodeDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_stakingStartBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_stakingEndBalance\",\"type\":\"uint256\"}],\"name\":\"setWithdrawn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_validatorPubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_validatorSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"_depositDataRoot\",\"type\":\"bytes32\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"userDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]"

// PoolBalance is an auto generated Go binding around an Ethereum contract.
type PoolBalance struct {
	PoolBalanceCaller     // Read-only binding to the contract
	PoolBalanceTransactor // Write-only binding to the contract
	PoolBalanceFilterer   // Log filterer for contract events
}

// PoolBalanceCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolBalanceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolBalanceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolBalanceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolBalanceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolBalanceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolBalanceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolBalanceSession struct {
	Contract     *PoolBalance      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolBalanceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolBalanceCallerSession struct {
	Contract *PoolBalanceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PoolBalanceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolBalanceTransactorSession struct {
	Contract     *PoolBalanceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PoolBalanceRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolBalanceRaw struct {
	Contract *PoolBalance // Generic contract binding to access the raw methods on
}

// PoolBalanceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolBalanceCallerRaw struct {
	Contract *PoolBalanceCaller // Generic read-only contract binding to access the raw methods on
}

// PoolBalanceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolBalanceTransactorRaw struct {
	Contract *PoolBalanceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoolBalance creates a new instance of PoolBalance, bound to a specific deployed contract.
func NewPoolBalance(address common.Address, backend bind.ContractBackend) (*PoolBalance, error) {
	contract, err := bindPoolBalance(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PoolBalance{PoolBalanceCaller: PoolBalanceCaller{contract: contract}, PoolBalanceTransactor: PoolBalanceTransactor{contract: contract}, PoolBalanceFilterer: PoolBalanceFilterer{contract: contract}}, nil
}

// NewPoolBalanceCaller creates a new read-only instance of PoolBalance, bound to a specific deployed contract.
func NewPoolBalanceCaller(address common.Address, caller bind.ContractCaller) (*PoolBalanceCaller, error) {
	contract, err := bindPoolBalance(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolBalanceCaller{contract: contract}, nil
}

// NewPoolBalanceTransactor creates a new write-only instance of PoolBalance, bound to a specific deployed contract.
func NewPoolBalanceTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolBalanceTransactor, error) {
	contract, err := bindPoolBalance(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolBalanceTransactor{contract: contract}, nil
}

// NewPoolBalanceFilterer creates a new log filterer instance of PoolBalance, bound to a specific deployed contract.
func NewPoolBalanceFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolBalanceFilterer, error) {
	contract, err := bindPoolBalance(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolBalanceFilterer{contract: contract}, nil
}

// bindPoolBalance binds a generic wrapper to an already deployed contract.
func bindPoolBalance(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PoolBalanceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolBalance *PoolBalanceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolBalance.Contract.PoolBalanceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolBalance *PoolBalanceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolBalance.Contract.PoolBalanceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolBalance *PoolBalanceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolBalance.Contract.PoolBalanceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolBalance *PoolBalanceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolBalance.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolBalance *PoolBalanceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolBalance.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolBalance *PoolBalanceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolBalance.Contract.contract.Transact(opts, method, params...)
}

// GetDepositType is a free data retrieval call binding the contract method 0x5abd37e4.
//
// Solidity: function getDepositType() view returns(uint8)
func (_PoolBalance *PoolBalanceCaller) GetDepositType(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getDepositType")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetDepositType is a free data retrieval call binding the contract method 0x5abd37e4.
//
// Solidity: function getDepositType() view returns(uint8)
func (_PoolBalance *PoolBalanceSession) GetDepositType() (uint8, error) {
	return _PoolBalance.Contract.GetDepositType(&_PoolBalance.CallOpts)
}

// GetDepositType is a free data retrieval call binding the contract method 0x5abd37e4.
//
// Solidity: function getDepositType() view returns(uint8)
func (_PoolBalance *PoolBalanceCallerSession) GetDepositType() (uint8, error) {
	return _PoolBalance.Contract.GetDepositType(&_PoolBalance.CallOpts)
}

// GetNodeAddress is a free data retrieval call binding the contract method 0x70dabc9e.
//
// Solidity: function getNodeAddress() view returns(address)
func (_PoolBalance *PoolBalanceCaller) GetNodeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getNodeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetNodeAddress is a free data retrieval call binding the contract method 0x70dabc9e.
//
// Solidity: function getNodeAddress() view returns(address)
func (_PoolBalance *PoolBalanceSession) GetNodeAddress() (common.Address, error) {
	return _PoolBalance.Contract.GetNodeAddress(&_PoolBalance.CallOpts)
}

// GetNodeAddress is a free data retrieval call binding the contract method 0x70dabc9e.
//
// Solidity: function getNodeAddress() view returns(address)
func (_PoolBalance *PoolBalanceCallerSession) GetNodeAddress() (common.Address, error) {
	return _PoolBalance.Contract.GetNodeAddress(&_PoolBalance.CallOpts)
}

// GetNodeDepositAssigned is a free data retrieval call binding the contract method 0x69c089ea.
//
// Solidity: function getNodeDepositAssigned() view returns(bool)
func (_PoolBalance *PoolBalanceCaller) GetNodeDepositAssigned(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getNodeDepositAssigned")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetNodeDepositAssigned is a free data retrieval call binding the contract method 0x69c089ea.
//
// Solidity: function getNodeDepositAssigned() view returns(bool)
func (_PoolBalance *PoolBalanceSession) GetNodeDepositAssigned() (bool, error) {
	return _PoolBalance.Contract.GetNodeDepositAssigned(&_PoolBalance.CallOpts)
}

// GetNodeDepositAssigned is a free data retrieval call binding the contract method 0x69c089ea.
//
// Solidity: function getNodeDepositAssigned() view returns(bool)
func (_PoolBalance *PoolBalanceCallerSession) GetNodeDepositAssigned() (bool, error) {
	return _PoolBalance.Contract.GetNodeDepositAssigned(&_PoolBalance.CallOpts)
}

// GetNodeDepositBalance is a free data retrieval call binding the contract method 0x74ca6bf2.
//
// Solidity: function getNodeDepositBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetNodeDepositBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getNodeDepositBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeDepositBalance is a free data retrieval call binding the contract method 0x74ca6bf2.
//
// Solidity: function getNodeDepositBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetNodeDepositBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetNodeDepositBalance(&_PoolBalance.CallOpts)
}

// GetNodeDepositBalance is a free data retrieval call binding the contract method 0x74ca6bf2.
//
// Solidity: function getNodeDepositBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetNodeDepositBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetNodeDepositBalance(&_PoolBalance.CallOpts)
}

// GetNodeFee is a free data retrieval call binding the contract method 0xe7150134.
//
// Solidity: function getNodeFee() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetNodeFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getNodeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeFee is a free data retrieval call binding the contract method 0xe7150134.
//
// Solidity: function getNodeFee() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetNodeFee() (*big.Int, error) {
	return _PoolBalance.Contract.GetNodeFee(&_PoolBalance.CallOpts)
}

// GetNodeFee is a free data retrieval call binding the contract method 0xe7150134.
//
// Solidity: function getNodeFee() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetNodeFee() (*big.Int, error) {
	return _PoolBalance.Contract.GetNodeFee(&_PoolBalance.CallOpts)
}

// GetNodeRefundBalance is a free data retrieval call binding the contract method 0xfbc02c42.
//
// Solidity: function getNodeRefundBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetNodeRefundBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getNodeRefundBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeRefundBalance is a free data retrieval call binding the contract method 0xfbc02c42.
//
// Solidity: function getNodeRefundBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetNodeRefundBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetNodeRefundBalance(&_PoolBalance.CallOpts)
}

// GetNodeRefundBalance is a free data retrieval call binding the contract method 0xfbc02c42.
//
// Solidity: function getNodeRefundBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetNodeRefundBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetNodeRefundBalance(&_PoolBalance.CallOpts)
}

// GetPlatformDepositBalance is a free data retrieval call binding the contract method 0xdd58cc23.
//
// Solidity: function getPlatformDepositBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetPlatformDepositBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getPlatformDepositBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPlatformDepositBalance is a free data retrieval call binding the contract method 0xdd58cc23.
//
// Solidity: function getPlatformDepositBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetPlatformDepositBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetPlatformDepositBalance(&_PoolBalance.CallOpts)
}

// GetPlatformDepositBalance is a free data retrieval call binding the contract method 0xdd58cc23.
//
// Solidity: function getPlatformDepositBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetPlatformDepositBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetPlatformDepositBalance(&_PoolBalance.CallOpts)
}

// GetStakingEndBalance is a free data retrieval call binding the contract method 0x4d9f376c.
//
// Solidity: function getStakingEndBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetStakingEndBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getStakingEndBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakingEndBalance is a free data retrieval call binding the contract method 0x4d9f376c.
//
// Solidity: function getStakingEndBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetStakingEndBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetStakingEndBalance(&_PoolBalance.CallOpts)
}

// GetStakingEndBalance is a free data retrieval call binding the contract method 0x4d9f376c.
//
// Solidity: function getStakingEndBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetStakingEndBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetStakingEndBalance(&_PoolBalance.CallOpts)
}

// GetStakingStartBalance is a free data retrieval call binding the contract method 0x6b101711.
//
// Solidity: function getStakingStartBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetStakingStartBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getStakingStartBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakingStartBalance is a free data retrieval call binding the contract method 0x6b101711.
//
// Solidity: function getStakingStartBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetStakingStartBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetStakingStartBalance(&_PoolBalance.CallOpts)
}

// GetStakingStartBalance is a free data retrieval call binding the contract method 0x6b101711.
//
// Solidity: function getStakingStartBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetStakingStartBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetStakingStartBalance(&_PoolBalance.CallOpts)
}

// GetStatus is a free data retrieval call binding the contract method 0x4e69d560.
//
// Solidity: function getStatus() view returns(uint8)
func (_PoolBalance *PoolBalanceCaller) GetStatus(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getStatus")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetStatus is a free data retrieval call binding the contract method 0x4e69d560.
//
// Solidity: function getStatus() view returns(uint8)
func (_PoolBalance *PoolBalanceSession) GetStatus() (uint8, error) {
	return _PoolBalance.Contract.GetStatus(&_PoolBalance.CallOpts)
}

// GetStatus is a free data retrieval call binding the contract method 0x4e69d560.
//
// Solidity: function getStatus() view returns(uint8)
func (_PoolBalance *PoolBalanceCallerSession) GetStatus() (uint8, error) {
	return _PoolBalance.Contract.GetStatus(&_PoolBalance.CallOpts)
}

// GetStatusBlock is a free data retrieval call binding the contract method 0xe67cd5b0.
//
// Solidity: function getStatusBlock() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetStatusBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getStatusBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStatusBlock is a free data retrieval call binding the contract method 0xe67cd5b0.
//
// Solidity: function getStatusBlock() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetStatusBlock() (*big.Int, error) {
	return _PoolBalance.Contract.GetStatusBlock(&_PoolBalance.CallOpts)
}

// GetStatusBlock is a free data retrieval call binding the contract method 0xe67cd5b0.
//
// Solidity: function getStatusBlock() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetStatusBlock() (*big.Int, error) {
	return _PoolBalance.Contract.GetStatusBlock(&_PoolBalance.CallOpts)
}

// GetStatusTime is a free data retrieval call binding the contract method 0x3e0a56b0.
//
// Solidity: function getStatusTime() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetStatusTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getStatusTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStatusTime is a free data retrieval call binding the contract method 0x3e0a56b0.
//
// Solidity: function getStatusTime() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetStatusTime() (*big.Int, error) {
	return _PoolBalance.Contract.GetStatusTime(&_PoolBalance.CallOpts)
}

// GetStatusTime is a free data retrieval call binding the contract method 0x3e0a56b0.
//
// Solidity: function getStatusTime() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetStatusTime() (*big.Int, error) {
	return _PoolBalance.Contract.GetStatusTime(&_PoolBalance.CallOpts)
}

// GetUserDepositAssigned is a free data retrieval call binding the contract method 0xd91eda62.
//
// Solidity: function getUserDepositAssigned() view returns(bool)
func (_PoolBalance *PoolBalanceCaller) GetUserDepositAssigned(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getUserDepositAssigned")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetUserDepositAssigned is a free data retrieval call binding the contract method 0xd91eda62.
//
// Solidity: function getUserDepositAssigned() view returns(bool)
func (_PoolBalance *PoolBalanceSession) GetUserDepositAssigned() (bool, error) {
	return _PoolBalance.Contract.GetUserDepositAssigned(&_PoolBalance.CallOpts)
}

// GetUserDepositAssigned is a free data retrieval call binding the contract method 0xd91eda62.
//
// Solidity: function getUserDepositAssigned() view returns(bool)
func (_PoolBalance *PoolBalanceCallerSession) GetUserDepositAssigned() (bool, error) {
	return _PoolBalance.Contract.GetUserDepositAssigned(&_PoolBalance.CallOpts)
}

// GetUserDepositAssignedTime is a free data retrieval call binding the contract method 0xa2940a90.
//
// Solidity: function getUserDepositAssignedTime() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetUserDepositAssignedTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getUserDepositAssignedTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserDepositAssignedTime is a free data retrieval call binding the contract method 0xa2940a90.
//
// Solidity: function getUserDepositAssignedTime() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetUserDepositAssignedTime() (*big.Int, error) {
	return _PoolBalance.Contract.GetUserDepositAssignedTime(&_PoolBalance.CallOpts)
}

// GetUserDepositAssignedTime is a free data retrieval call binding the contract method 0xa2940a90.
//
// Solidity: function getUserDepositAssignedTime() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetUserDepositAssignedTime() (*big.Int, error) {
	return _PoolBalance.Contract.GetUserDepositAssignedTime(&_PoolBalance.CallOpts)
}

// GetUserDepositBalance is a free data retrieval call binding the contract method 0xe7e04aba.
//
// Solidity: function getUserDepositBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCaller) GetUserDepositBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolBalance.contract.Call(opts, &out, "getUserDepositBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserDepositBalance is a free data retrieval call binding the contract method 0xe7e04aba.
//
// Solidity: function getUserDepositBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceSession) GetUserDepositBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetUserDepositBalance(&_PoolBalance.CallOpts)
}

// GetUserDepositBalance is a free data retrieval call binding the contract method 0xe7e04aba.
//
// Solidity: function getUserDepositBalance() view returns(uint256)
func (_PoolBalance *PoolBalanceCallerSession) GetUserDepositBalance() (*big.Int, error) {
	return _PoolBalance.Contract.GetUserDepositBalance(&_PoolBalance.CallOpts)
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() returns()
func (_PoolBalance *PoolBalanceTransactor) Close(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolBalance.contract.Transact(opts, "close")
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() returns()
func (_PoolBalance *PoolBalanceSession) Close() (*types.Transaction, error) {
	return _PoolBalance.Contract.Close(&_PoolBalance.TransactOpts)
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() returns()
func (_PoolBalance *PoolBalanceTransactorSession) Close() (*types.Transaction, error) {
	return _PoolBalance.Contract.Close(&_PoolBalance.TransactOpts)
}

// Dissolve is a paid mutator transaction binding the contract method 0x3bef8a3a.
//
// Solidity: function dissolve() returns()
func (_PoolBalance *PoolBalanceTransactor) Dissolve(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolBalance.contract.Transact(opts, "dissolve")
}

// Dissolve is a paid mutator transaction binding the contract method 0x3bef8a3a.
//
// Solidity: function dissolve() returns()
func (_PoolBalance *PoolBalanceSession) Dissolve() (*types.Transaction, error) {
	return _PoolBalance.Contract.Dissolve(&_PoolBalance.TransactOpts)
}

// Dissolve is a paid mutator transaction binding the contract method 0x3bef8a3a.
//
// Solidity: function dissolve() returns()
func (_PoolBalance *PoolBalanceTransactorSession) Dissolve() (*types.Transaction, error) {
	return _PoolBalance.Contract.Dissolve(&_PoolBalance.TransactOpts)
}

// NodeDeposit is a paid mutator transaction binding the contract method 0xaa6c5ace.
//
// Solidity: function nodeDeposit() payable returns()
func (_PoolBalance *PoolBalanceTransactor) NodeDeposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolBalance.contract.Transact(opts, "nodeDeposit")
}

// NodeDeposit is a paid mutator transaction binding the contract method 0xaa6c5ace.
//
// Solidity: function nodeDeposit() payable returns()
func (_PoolBalance *PoolBalanceSession) NodeDeposit() (*types.Transaction, error) {
	return _PoolBalance.Contract.NodeDeposit(&_PoolBalance.TransactOpts)
}

// NodeDeposit is a paid mutator transaction binding the contract method 0xaa6c5ace.
//
// Solidity: function nodeDeposit() payable returns()
func (_PoolBalance *PoolBalanceTransactorSession) NodeDeposit() (*types.Transaction, error) {
	return _PoolBalance.Contract.NodeDeposit(&_PoolBalance.TransactOpts)
}

// Refund is a paid mutator transaction binding the contract method 0x590e1ae3.
//
// Solidity: function refund() returns()
func (_PoolBalance *PoolBalanceTransactor) Refund(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolBalance.contract.Transact(opts, "refund")
}

// Refund is a paid mutator transaction binding the contract method 0x590e1ae3.
//
// Solidity: function refund() returns()
func (_PoolBalance *PoolBalanceSession) Refund() (*types.Transaction, error) {
	return _PoolBalance.Contract.Refund(&_PoolBalance.TransactOpts)
}

// Refund is a paid mutator transaction binding the contract method 0x590e1ae3.
//
// Solidity: function refund() returns()
func (_PoolBalance *PoolBalanceTransactorSession) Refund() (*types.Transaction, error) {
	return _PoolBalance.Contract.Refund(&_PoolBalance.TransactOpts)
}

// SetWithdrawn is a paid mutator transaction binding the contract method 0x77c16e3a.
//
// Solidity: function setWithdrawn(uint256 _stakingStartBalance, uint256 _stakingEndBalance) returns()
func (_PoolBalance *PoolBalanceTransactor) SetWithdrawn(opts *bind.TransactOpts, _stakingStartBalance *big.Int, _stakingEndBalance *big.Int) (*types.Transaction, error) {
	return _PoolBalance.contract.Transact(opts, "setWithdrawn", _stakingStartBalance, _stakingEndBalance)
}

// SetWithdrawn is a paid mutator transaction binding the contract method 0x77c16e3a.
//
// Solidity: function setWithdrawn(uint256 _stakingStartBalance, uint256 _stakingEndBalance) returns()
func (_PoolBalance *PoolBalanceSession) SetWithdrawn(_stakingStartBalance *big.Int, _stakingEndBalance *big.Int) (*types.Transaction, error) {
	return _PoolBalance.Contract.SetWithdrawn(&_PoolBalance.TransactOpts, _stakingStartBalance, _stakingEndBalance)
}

// SetWithdrawn is a paid mutator transaction binding the contract method 0x77c16e3a.
//
// Solidity: function setWithdrawn(uint256 _stakingStartBalance, uint256 _stakingEndBalance) returns()
func (_PoolBalance *PoolBalanceTransactorSession) SetWithdrawn(_stakingStartBalance *big.Int, _stakingEndBalance *big.Int) (*types.Transaction, error) {
	return _PoolBalance.Contract.SetWithdrawn(&_PoolBalance.TransactOpts, _stakingStartBalance, _stakingEndBalance)
}

// Stake is a paid mutator transaction binding the contract method 0x9b4e4634.
//
// Solidity: function stake(bytes _validatorPubkey, bytes _validatorSignature, bytes32 _depositDataRoot) returns()
func (_PoolBalance *PoolBalanceTransactor) Stake(opts *bind.TransactOpts, _validatorPubkey []byte, _validatorSignature []byte, _depositDataRoot [32]byte) (*types.Transaction, error) {
	return _PoolBalance.contract.Transact(opts, "stake", _validatorPubkey, _validatorSignature, _depositDataRoot)
}

// Stake is a paid mutator transaction binding the contract method 0x9b4e4634.
//
// Solidity: function stake(bytes _validatorPubkey, bytes _validatorSignature, bytes32 _depositDataRoot) returns()
func (_PoolBalance *PoolBalanceSession) Stake(_validatorPubkey []byte, _validatorSignature []byte, _depositDataRoot [32]byte) (*types.Transaction, error) {
	return _PoolBalance.Contract.Stake(&_PoolBalance.TransactOpts, _validatorPubkey, _validatorSignature, _depositDataRoot)
}

// Stake is a paid mutator transaction binding the contract method 0x9b4e4634.
//
// Solidity: function stake(bytes _validatorPubkey, bytes _validatorSignature, bytes32 _depositDataRoot) returns()
func (_PoolBalance *PoolBalanceTransactorSession) Stake(_validatorPubkey []byte, _validatorSignature []byte, _depositDataRoot [32]byte) (*types.Transaction, error) {
	return _PoolBalance.Contract.Stake(&_PoolBalance.TransactOpts, _validatorPubkey, _validatorSignature, _depositDataRoot)
}

// UserDeposit is a paid mutator transaction binding the contract method 0x48146113.
//
// Solidity: function userDeposit() payable returns()
func (_PoolBalance *PoolBalanceTransactor) UserDeposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolBalance.contract.Transact(opts, "userDeposit")
}

// UserDeposit is a paid mutator transaction binding the contract method 0x48146113.
//
// Solidity: function userDeposit() payable returns()
func (_PoolBalance *PoolBalanceSession) UserDeposit() (*types.Transaction, error) {
	return _PoolBalance.Contract.UserDeposit(&_PoolBalance.TransactOpts)
}

// UserDeposit is a paid mutator transaction binding the contract method 0x48146113.
//
// Solidity: function userDeposit() payable returns()
func (_PoolBalance *PoolBalanceTransactorSession) UserDeposit() (*types.Transaction, error) {
	return _PoolBalance.Contract.UserDeposit(&_PoolBalance.TransactOpts)
}

// PoolBalanceEtherDepositedIterator is returned from FilterEtherDeposited and is used to iterate over the raw logs and unpacked data for EtherDeposited events raised by the PoolBalance contract.
type PoolBalanceEtherDepositedIterator struct {
	Event *PoolBalanceEtherDeposited // Event containing the contract specifics and raw log

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
func (it *PoolBalanceEtherDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolBalanceEtherDeposited)
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
		it.Event = new(PoolBalanceEtherDeposited)
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
func (it *PoolBalanceEtherDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolBalanceEtherDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolBalanceEtherDeposited represents a EtherDeposited event raised by the PoolBalance contract.
type PoolBalanceEtherDeposited struct {
	From   common.Address
	Amount *big.Int
	Time   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherDeposited is a free log retrieval operation binding the contract event 0xef51b4c870b8b0100eae2072e91db01222a303072af3728e58c9d4d2da33127f.
//
// Solidity: event EtherDeposited(address indexed from, uint256 amount, uint256 time)
func (_PoolBalance *PoolBalanceFilterer) FilterEtherDeposited(opts *bind.FilterOpts, from []common.Address) (*PoolBalanceEtherDepositedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _PoolBalance.contract.FilterLogs(opts, "EtherDeposited", fromRule)
	if err != nil {
		return nil, err
	}
	return &PoolBalanceEtherDepositedIterator{contract: _PoolBalance.contract, event: "EtherDeposited", logs: logs, sub: sub}, nil
}

// WatchEtherDeposited is a free log subscription operation binding the contract event 0xef51b4c870b8b0100eae2072e91db01222a303072af3728e58c9d4d2da33127f.
//
// Solidity: event EtherDeposited(address indexed from, uint256 amount, uint256 time)
func (_PoolBalance *PoolBalanceFilterer) WatchEtherDeposited(opts *bind.WatchOpts, sink chan<- *PoolBalanceEtherDeposited, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _PoolBalance.contract.WatchLogs(opts, "EtherDeposited", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolBalanceEtherDeposited)
				if err := _PoolBalance.contract.UnpackLog(event, "EtherDeposited", log); err != nil {
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

// ParseEtherDeposited is a log parse operation binding the contract event 0xef51b4c870b8b0100eae2072e91db01222a303072af3728e58c9d4d2da33127f.
//
// Solidity: event EtherDeposited(address indexed from, uint256 amount, uint256 time)
func (_PoolBalance *PoolBalanceFilterer) ParseEtherDeposited(log types.Log) (*PoolBalanceEtherDeposited, error) {
	event := new(PoolBalanceEtherDeposited)
	if err := _PoolBalance.contract.UnpackLog(event, "EtherDeposited", log); err != nil {
		return nil, err
	}
	return event, nil
}

// PoolBalanceEtherWithdrawnIterator is returned from FilterEtherWithdrawn and is used to iterate over the raw logs and unpacked data for EtherWithdrawn events raised by the PoolBalance contract.
type PoolBalanceEtherWithdrawnIterator struct {
	Event *PoolBalanceEtherWithdrawn // Event containing the contract specifics and raw log

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
func (it *PoolBalanceEtherWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolBalanceEtherWithdrawn)
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
		it.Event = new(PoolBalanceEtherWithdrawn)
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
func (it *PoolBalanceEtherWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolBalanceEtherWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolBalanceEtherWithdrawn represents a EtherWithdrawn event raised by the PoolBalance contract.
type PoolBalanceEtherWithdrawn struct {
	To     common.Address
	Amount *big.Int
	Time   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherWithdrawn is a free log retrieval operation binding the contract event 0xd5ca65e1ec4f4864fea7b9c5cb1ec3087a0dbf9c74641db3f6458edf445c4051.
//
// Solidity: event EtherWithdrawn(address indexed to, uint256 amount, uint256 time)
func (_PoolBalance *PoolBalanceFilterer) FilterEtherWithdrawn(opts *bind.FilterOpts, to []common.Address) (*PoolBalanceEtherWithdrawnIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PoolBalance.contract.FilterLogs(opts, "EtherWithdrawn", toRule)
	if err != nil {
		return nil, err
	}
	return &PoolBalanceEtherWithdrawnIterator{contract: _PoolBalance.contract, event: "EtherWithdrawn", logs: logs, sub: sub}, nil
}

// WatchEtherWithdrawn is a free log subscription operation binding the contract event 0xd5ca65e1ec4f4864fea7b9c5cb1ec3087a0dbf9c74641db3f6458edf445c4051.
//
// Solidity: event EtherWithdrawn(address indexed to, uint256 amount, uint256 time)
func (_PoolBalance *PoolBalanceFilterer) WatchEtherWithdrawn(opts *bind.WatchOpts, sink chan<- *PoolBalanceEtherWithdrawn, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _PoolBalance.contract.WatchLogs(opts, "EtherWithdrawn", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolBalanceEtherWithdrawn)
				if err := _PoolBalance.contract.UnpackLog(event, "EtherWithdrawn", log); err != nil {
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

// ParseEtherWithdrawn is a log parse operation binding the contract event 0xd5ca65e1ec4f4864fea7b9c5cb1ec3087a0dbf9c74641db3f6458edf445c4051.
//
// Solidity: event EtherWithdrawn(address indexed to, uint256 amount, uint256 time)
func (_PoolBalance *PoolBalanceFilterer) ParseEtherWithdrawn(log types.Log) (*PoolBalanceEtherWithdrawn, error) {
	event := new(PoolBalanceEtherWithdrawn)
	if err := _PoolBalance.contract.UnpackLog(event, "EtherWithdrawn", log); err != nil {
		return nil, err
	}
	return event, nil
}

// PoolBalanceStatusUpdatedIterator is returned from FilterStatusUpdated and is used to iterate over the raw logs and unpacked data for StatusUpdated events raised by the PoolBalance contract.
type PoolBalanceStatusUpdatedIterator struct {
	Event *PoolBalanceStatusUpdated // Event containing the contract specifics and raw log

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
func (it *PoolBalanceStatusUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolBalanceStatusUpdated)
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
		it.Event = new(PoolBalanceStatusUpdated)
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
func (it *PoolBalanceStatusUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolBalanceStatusUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolBalanceStatusUpdated represents a StatusUpdated event raised by the PoolBalance contract.
type PoolBalanceStatusUpdated struct {
	Status uint8
	Time   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStatusUpdated is a free log retrieval operation binding the contract event 0x26725881c2a4290b02cd153d6599fd484f0d4e6062c361e740fbbe39e7ad6142.
//
// Solidity: event StatusUpdated(uint8 indexed status, uint256 time)
func (_PoolBalance *PoolBalanceFilterer) FilterStatusUpdated(opts *bind.FilterOpts, status []uint8) (*PoolBalanceStatusUpdatedIterator, error) {

	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _PoolBalance.contract.FilterLogs(opts, "StatusUpdated", statusRule)
	if err != nil {
		return nil, err
	}
	return &PoolBalanceStatusUpdatedIterator{contract: _PoolBalance.contract, event: "StatusUpdated", logs: logs, sub: sub}, nil
}

// WatchStatusUpdated is a free log subscription operation binding the contract event 0x26725881c2a4290b02cd153d6599fd484f0d4e6062c361e740fbbe39e7ad6142.
//
// Solidity: event StatusUpdated(uint8 indexed status, uint256 time)
func (_PoolBalance *PoolBalanceFilterer) WatchStatusUpdated(opts *bind.WatchOpts, sink chan<- *PoolBalanceStatusUpdated, status []uint8) (event.Subscription, error) {

	var statusRule []interface{}
	for _, statusItem := range status {
		statusRule = append(statusRule, statusItem)
	}

	logs, sub, err := _PoolBalance.contract.WatchLogs(opts, "StatusUpdated", statusRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolBalanceStatusUpdated)
				if err := _PoolBalance.contract.UnpackLog(event, "StatusUpdated", log); err != nil {
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

// ParseStatusUpdated is a log parse operation binding the contract event 0x26725881c2a4290b02cd153d6599fd484f0d4e6062c361e740fbbe39e7ad6142.
//
// Solidity: event StatusUpdated(uint8 indexed status, uint256 time)
func (_PoolBalance *PoolBalanceFilterer) ParseStatusUpdated(log types.Log) (*PoolBalanceStatusUpdated, error) {
	event := new(PoolBalanceStatusUpdated)
	if err := _PoolBalance.contract.UnpackLog(event, "StatusUpdated", log); err != nil {
		return nil, err
	}
	return event, nil
}
