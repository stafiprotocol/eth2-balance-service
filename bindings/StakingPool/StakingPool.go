// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package StakingPool

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

// StakingPoolMetaData contains all meta data concerning the StakingPool contract.
var StakingPoolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"enumDepositType\",\"name\":\"_depositType\",\"type\":\"uint8\"}],\"name\":\"initialise\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStatus\",\"outputs\":[{\"internalType\":\"enumStakingPoolStatus\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStatusBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStatusTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWithdrawalCredentialsMatch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getDepositType\",\"outputs\":[{\"internalType\":\"enumDepositType\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeRefundBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeDepositAssigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeCommonlyRefunded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeTrustedRefunded\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUserDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUserDepositAssigned\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getUserDepositAssignedTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPlatformDepositBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_validatorPubkey\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_validatorSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"_depositDataRoot\",\"type\":\"bytes32\"}],\"name\":\"nodeDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"userDeposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_validatorSignature\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"_depositDataRoot\",\"type\":\"bytes32\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dissolve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"close\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"voteWithdrawCredentials\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// StakingPoolABI is the input ABI used to generate the binding from.
// Deprecated: Use StakingPoolMetaData.ABI instead.
var StakingPoolABI = StakingPoolMetaData.ABI

// StakingPool is an auto generated Go binding around an Ethereum contract.
type StakingPool struct {
	StakingPoolCaller     // Read-only binding to the contract
	StakingPoolTransactor // Write-only binding to the contract
	StakingPoolFilterer   // Log filterer for contract events
}

// StakingPoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakingPoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakingPoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakingPoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakingPoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakingPoolSession struct {
	Contract     *StakingPool      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakingPoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakingPoolCallerSession struct {
	Contract *StakingPoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StakingPoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakingPoolTransactorSession struct {
	Contract     *StakingPoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StakingPoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakingPoolRaw struct {
	Contract *StakingPool // Generic contract binding to access the raw methods on
}

// StakingPoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakingPoolCallerRaw struct {
	Contract *StakingPoolCaller // Generic read-only contract binding to access the raw methods on
}

// StakingPoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakingPoolTransactorRaw struct {
	Contract *StakingPoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakingPool creates a new instance of StakingPool, bound to a specific deployed contract.
func NewStakingPool(address common.Address, backend bind.ContractBackend) (*StakingPool, error) {
	contract, err := bindStakingPool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StakingPool{StakingPoolCaller: StakingPoolCaller{contract: contract}, StakingPoolTransactor: StakingPoolTransactor{contract: contract}, StakingPoolFilterer: StakingPoolFilterer{contract: contract}}, nil
}

// NewStakingPoolCaller creates a new read-only instance of StakingPool, bound to a specific deployed contract.
func NewStakingPoolCaller(address common.Address, caller bind.ContractCaller) (*StakingPoolCaller, error) {
	contract, err := bindStakingPool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakingPoolCaller{contract: contract}, nil
}

// NewStakingPoolTransactor creates a new write-only instance of StakingPool, bound to a specific deployed contract.
func NewStakingPoolTransactor(address common.Address, transactor bind.ContractTransactor) (*StakingPoolTransactor, error) {
	contract, err := bindStakingPool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakingPoolTransactor{contract: contract}, nil
}

// NewStakingPoolFilterer creates a new log filterer instance of StakingPool, bound to a specific deployed contract.
func NewStakingPoolFilterer(address common.Address, filterer bind.ContractFilterer) (*StakingPoolFilterer, error) {
	contract, err := bindStakingPool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakingPoolFilterer{contract: contract}, nil
}

// bindStakingPool binds a generic wrapper to an already deployed contract.
func bindStakingPool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakingPoolABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingPool *StakingPoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingPool.Contract.StakingPoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingPool *StakingPoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.Contract.StakingPoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingPool *StakingPoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingPool.Contract.StakingPoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StakingPool *StakingPoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StakingPool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StakingPool *StakingPoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StakingPool *StakingPoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StakingPool.Contract.contract.Transact(opts, method, params...)
}

// GetDepositType is a free data retrieval call binding the contract method 0x5abd37e4.
//
// Solidity: function getDepositType() view returns(uint8)
func (_StakingPool *StakingPoolCaller) GetDepositType(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getDepositType")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetDepositType is a free data retrieval call binding the contract method 0x5abd37e4.
//
// Solidity: function getDepositType() view returns(uint8)
func (_StakingPool *StakingPoolSession) GetDepositType() (uint8, error) {
	return _StakingPool.Contract.GetDepositType(&_StakingPool.CallOpts)
}

// GetDepositType is a free data retrieval call binding the contract method 0x5abd37e4.
//
// Solidity: function getDepositType() view returns(uint8)
func (_StakingPool *StakingPoolCallerSession) GetDepositType() (uint8, error) {
	return _StakingPool.Contract.GetDepositType(&_StakingPool.CallOpts)
}

// GetNodeAddress is a free data retrieval call binding the contract method 0x70dabc9e.
//
// Solidity: function getNodeAddress() view returns(address)
func (_StakingPool *StakingPoolCaller) GetNodeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getNodeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetNodeAddress is a free data retrieval call binding the contract method 0x70dabc9e.
//
// Solidity: function getNodeAddress() view returns(address)
func (_StakingPool *StakingPoolSession) GetNodeAddress() (common.Address, error) {
	return _StakingPool.Contract.GetNodeAddress(&_StakingPool.CallOpts)
}

// GetNodeAddress is a free data retrieval call binding the contract method 0x70dabc9e.
//
// Solidity: function getNodeAddress() view returns(address)
func (_StakingPool *StakingPoolCallerSession) GetNodeAddress() (common.Address, error) {
	return _StakingPool.Contract.GetNodeAddress(&_StakingPool.CallOpts)
}

// GetNodeCommonlyRefunded is a free data retrieval call binding the contract method 0xa0a8440f.
//
// Solidity: function getNodeCommonlyRefunded() view returns(bool)
func (_StakingPool *StakingPoolCaller) GetNodeCommonlyRefunded(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getNodeCommonlyRefunded")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetNodeCommonlyRefunded is a free data retrieval call binding the contract method 0xa0a8440f.
//
// Solidity: function getNodeCommonlyRefunded() view returns(bool)
func (_StakingPool *StakingPoolSession) GetNodeCommonlyRefunded() (bool, error) {
	return _StakingPool.Contract.GetNodeCommonlyRefunded(&_StakingPool.CallOpts)
}

// GetNodeCommonlyRefunded is a free data retrieval call binding the contract method 0xa0a8440f.
//
// Solidity: function getNodeCommonlyRefunded() view returns(bool)
func (_StakingPool *StakingPoolCallerSession) GetNodeCommonlyRefunded() (bool, error) {
	return _StakingPool.Contract.GetNodeCommonlyRefunded(&_StakingPool.CallOpts)
}

// GetNodeDepositAssigned is a free data retrieval call binding the contract method 0x69c089ea.
//
// Solidity: function getNodeDepositAssigned() view returns(bool)
func (_StakingPool *StakingPoolCaller) GetNodeDepositAssigned(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getNodeDepositAssigned")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetNodeDepositAssigned is a free data retrieval call binding the contract method 0x69c089ea.
//
// Solidity: function getNodeDepositAssigned() view returns(bool)
func (_StakingPool *StakingPoolSession) GetNodeDepositAssigned() (bool, error) {
	return _StakingPool.Contract.GetNodeDepositAssigned(&_StakingPool.CallOpts)
}

// GetNodeDepositAssigned is a free data retrieval call binding the contract method 0x69c089ea.
//
// Solidity: function getNodeDepositAssigned() view returns(bool)
func (_StakingPool *StakingPoolCallerSession) GetNodeDepositAssigned() (bool, error) {
	return _StakingPool.Contract.GetNodeDepositAssigned(&_StakingPool.CallOpts)
}

// GetNodeDepositBalance is a free data retrieval call binding the contract method 0x74ca6bf2.
//
// Solidity: function getNodeDepositBalance() view returns(uint256)
func (_StakingPool *StakingPoolCaller) GetNodeDepositBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getNodeDepositBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeDepositBalance is a free data retrieval call binding the contract method 0x74ca6bf2.
//
// Solidity: function getNodeDepositBalance() view returns(uint256)
func (_StakingPool *StakingPoolSession) GetNodeDepositBalance() (*big.Int, error) {
	return _StakingPool.Contract.GetNodeDepositBalance(&_StakingPool.CallOpts)
}

// GetNodeDepositBalance is a free data retrieval call binding the contract method 0x74ca6bf2.
//
// Solidity: function getNodeDepositBalance() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) GetNodeDepositBalance() (*big.Int, error) {
	return _StakingPool.Contract.GetNodeDepositBalance(&_StakingPool.CallOpts)
}

// GetNodeFee is a free data retrieval call binding the contract method 0xe7150134.
//
// Solidity: function getNodeFee() view returns(uint256)
func (_StakingPool *StakingPoolCaller) GetNodeFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getNodeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeFee is a free data retrieval call binding the contract method 0xe7150134.
//
// Solidity: function getNodeFee() view returns(uint256)
func (_StakingPool *StakingPoolSession) GetNodeFee() (*big.Int, error) {
	return _StakingPool.Contract.GetNodeFee(&_StakingPool.CallOpts)
}

// GetNodeFee is a free data retrieval call binding the contract method 0xe7150134.
//
// Solidity: function getNodeFee() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) GetNodeFee() (*big.Int, error) {
	return _StakingPool.Contract.GetNodeFee(&_StakingPool.CallOpts)
}

// GetNodeRefundBalance is a free data retrieval call binding the contract method 0xfbc02c42.
//
// Solidity: function getNodeRefundBalance() view returns(uint256)
func (_StakingPool *StakingPoolCaller) GetNodeRefundBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getNodeRefundBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeRefundBalance is a free data retrieval call binding the contract method 0xfbc02c42.
//
// Solidity: function getNodeRefundBalance() view returns(uint256)
func (_StakingPool *StakingPoolSession) GetNodeRefundBalance() (*big.Int, error) {
	return _StakingPool.Contract.GetNodeRefundBalance(&_StakingPool.CallOpts)
}

// GetNodeRefundBalance is a free data retrieval call binding the contract method 0xfbc02c42.
//
// Solidity: function getNodeRefundBalance() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) GetNodeRefundBalance() (*big.Int, error) {
	return _StakingPool.Contract.GetNodeRefundBalance(&_StakingPool.CallOpts)
}

// GetNodeTrustedRefunded is a free data retrieval call binding the contract method 0x542bac20.
//
// Solidity: function getNodeTrustedRefunded() view returns(bool)
func (_StakingPool *StakingPoolCaller) GetNodeTrustedRefunded(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getNodeTrustedRefunded")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetNodeTrustedRefunded is a free data retrieval call binding the contract method 0x542bac20.
//
// Solidity: function getNodeTrustedRefunded() view returns(bool)
func (_StakingPool *StakingPoolSession) GetNodeTrustedRefunded() (bool, error) {
	return _StakingPool.Contract.GetNodeTrustedRefunded(&_StakingPool.CallOpts)
}

// GetNodeTrustedRefunded is a free data retrieval call binding the contract method 0x542bac20.
//
// Solidity: function getNodeTrustedRefunded() view returns(bool)
func (_StakingPool *StakingPoolCallerSession) GetNodeTrustedRefunded() (bool, error) {
	return _StakingPool.Contract.GetNodeTrustedRefunded(&_StakingPool.CallOpts)
}

// GetPlatformDepositBalance is a free data retrieval call binding the contract method 0xdd58cc23.
//
// Solidity: function getPlatformDepositBalance() view returns(uint256)
func (_StakingPool *StakingPoolCaller) GetPlatformDepositBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getPlatformDepositBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPlatformDepositBalance is a free data retrieval call binding the contract method 0xdd58cc23.
//
// Solidity: function getPlatformDepositBalance() view returns(uint256)
func (_StakingPool *StakingPoolSession) GetPlatformDepositBalance() (*big.Int, error) {
	return _StakingPool.Contract.GetPlatformDepositBalance(&_StakingPool.CallOpts)
}

// GetPlatformDepositBalance is a free data retrieval call binding the contract method 0xdd58cc23.
//
// Solidity: function getPlatformDepositBalance() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) GetPlatformDepositBalance() (*big.Int, error) {
	return _StakingPool.Contract.GetPlatformDepositBalance(&_StakingPool.CallOpts)
}

// GetStatus is a free data retrieval call binding the contract method 0x4e69d560.
//
// Solidity: function getStatus() view returns(uint8)
func (_StakingPool *StakingPoolCaller) GetStatus(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getStatus")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetStatus is a free data retrieval call binding the contract method 0x4e69d560.
//
// Solidity: function getStatus() view returns(uint8)
func (_StakingPool *StakingPoolSession) GetStatus() (uint8, error) {
	return _StakingPool.Contract.GetStatus(&_StakingPool.CallOpts)
}

// GetStatus is a free data retrieval call binding the contract method 0x4e69d560.
//
// Solidity: function getStatus() view returns(uint8)
func (_StakingPool *StakingPoolCallerSession) GetStatus() (uint8, error) {
	return _StakingPool.Contract.GetStatus(&_StakingPool.CallOpts)
}

// GetStatusBlock is a free data retrieval call binding the contract method 0xe67cd5b0.
//
// Solidity: function getStatusBlock() view returns(uint256)
func (_StakingPool *StakingPoolCaller) GetStatusBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getStatusBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStatusBlock is a free data retrieval call binding the contract method 0xe67cd5b0.
//
// Solidity: function getStatusBlock() view returns(uint256)
func (_StakingPool *StakingPoolSession) GetStatusBlock() (*big.Int, error) {
	return _StakingPool.Contract.GetStatusBlock(&_StakingPool.CallOpts)
}

// GetStatusBlock is a free data retrieval call binding the contract method 0xe67cd5b0.
//
// Solidity: function getStatusBlock() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) GetStatusBlock() (*big.Int, error) {
	return _StakingPool.Contract.GetStatusBlock(&_StakingPool.CallOpts)
}

// GetStatusTime is a free data retrieval call binding the contract method 0x3e0a56b0.
//
// Solidity: function getStatusTime() view returns(uint256)
func (_StakingPool *StakingPoolCaller) GetStatusTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getStatusTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStatusTime is a free data retrieval call binding the contract method 0x3e0a56b0.
//
// Solidity: function getStatusTime() view returns(uint256)
func (_StakingPool *StakingPoolSession) GetStatusTime() (*big.Int, error) {
	return _StakingPool.Contract.GetStatusTime(&_StakingPool.CallOpts)
}

// GetStatusTime is a free data retrieval call binding the contract method 0x3e0a56b0.
//
// Solidity: function getStatusTime() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) GetStatusTime() (*big.Int, error) {
	return _StakingPool.Contract.GetStatusTime(&_StakingPool.CallOpts)
}

// GetUserDepositAssigned is a free data retrieval call binding the contract method 0xd91eda62.
//
// Solidity: function getUserDepositAssigned() view returns(bool)
func (_StakingPool *StakingPoolCaller) GetUserDepositAssigned(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getUserDepositAssigned")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetUserDepositAssigned is a free data retrieval call binding the contract method 0xd91eda62.
//
// Solidity: function getUserDepositAssigned() view returns(bool)
func (_StakingPool *StakingPoolSession) GetUserDepositAssigned() (bool, error) {
	return _StakingPool.Contract.GetUserDepositAssigned(&_StakingPool.CallOpts)
}

// GetUserDepositAssigned is a free data retrieval call binding the contract method 0xd91eda62.
//
// Solidity: function getUserDepositAssigned() view returns(bool)
func (_StakingPool *StakingPoolCallerSession) GetUserDepositAssigned() (bool, error) {
	return _StakingPool.Contract.GetUserDepositAssigned(&_StakingPool.CallOpts)
}

// GetUserDepositAssignedTime is a free data retrieval call binding the contract method 0xa2940a90.
//
// Solidity: function getUserDepositAssignedTime() view returns(uint256)
func (_StakingPool *StakingPoolCaller) GetUserDepositAssignedTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getUserDepositAssignedTime")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserDepositAssignedTime is a free data retrieval call binding the contract method 0xa2940a90.
//
// Solidity: function getUserDepositAssignedTime() view returns(uint256)
func (_StakingPool *StakingPoolSession) GetUserDepositAssignedTime() (*big.Int, error) {
	return _StakingPool.Contract.GetUserDepositAssignedTime(&_StakingPool.CallOpts)
}

// GetUserDepositAssignedTime is a free data retrieval call binding the contract method 0xa2940a90.
//
// Solidity: function getUserDepositAssignedTime() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) GetUserDepositAssignedTime() (*big.Int, error) {
	return _StakingPool.Contract.GetUserDepositAssignedTime(&_StakingPool.CallOpts)
}

// GetUserDepositBalance is a free data retrieval call binding the contract method 0xe7e04aba.
//
// Solidity: function getUserDepositBalance() view returns(uint256)
func (_StakingPool *StakingPoolCaller) GetUserDepositBalance(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getUserDepositBalance")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserDepositBalance is a free data retrieval call binding the contract method 0xe7e04aba.
//
// Solidity: function getUserDepositBalance() view returns(uint256)
func (_StakingPool *StakingPoolSession) GetUserDepositBalance() (*big.Int, error) {
	return _StakingPool.Contract.GetUserDepositBalance(&_StakingPool.CallOpts)
}

// GetUserDepositBalance is a free data retrieval call binding the contract method 0xe7e04aba.
//
// Solidity: function getUserDepositBalance() view returns(uint256)
func (_StakingPool *StakingPoolCallerSession) GetUserDepositBalance() (*big.Int, error) {
	return _StakingPool.Contract.GetUserDepositBalance(&_StakingPool.CallOpts)
}

// GetWithdrawalCredentialsMatch is a free data retrieval call binding the contract method 0x05dab351.
//
// Solidity: function getWithdrawalCredentialsMatch() view returns(bool)
func (_StakingPool *StakingPoolCaller) GetWithdrawalCredentialsMatch(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _StakingPool.contract.Call(opts, &out, "getWithdrawalCredentialsMatch")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetWithdrawalCredentialsMatch is a free data retrieval call binding the contract method 0x05dab351.
//
// Solidity: function getWithdrawalCredentialsMatch() view returns(bool)
func (_StakingPool *StakingPoolSession) GetWithdrawalCredentialsMatch() (bool, error) {
	return _StakingPool.Contract.GetWithdrawalCredentialsMatch(&_StakingPool.CallOpts)
}

// GetWithdrawalCredentialsMatch is a free data retrieval call binding the contract method 0x05dab351.
//
// Solidity: function getWithdrawalCredentialsMatch() view returns(bool)
func (_StakingPool *StakingPoolCallerSession) GetWithdrawalCredentialsMatch() (bool, error) {
	return _StakingPool.Contract.GetWithdrawalCredentialsMatch(&_StakingPool.CallOpts)
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() returns()
func (_StakingPool *StakingPoolTransactor) Close(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "close")
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() returns()
func (_StakingPool *StakingPoolSession) Close() (*types.Transaction, error) {
	return _StakingPool.Contract.Close(&_StakingPool.TransactOpts)
}

// Close is a paid mutator transaction binding the contract method 0x43d726d6.
//
// Solidity: function close() returns()
func (_StakingPool *StakingPoolTransactorSession) Close() (*types.Transaction, error) {
	return _StakingPool.Contract.Close(&_StakingPool.TransactOpts)
}

// Dissolve is a paid mutator transaction binding the contract method 0x3bef8a3a.
//
// Solidity: function dissolve() returns()
func (_StakingPool *StakingPoolTransactor) Dissolve(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "dissolve")
}

// Dissolve is a paid mutator transaction binding the contract method 0x3bef8a3a.
//
// Solidity: function dissolve() returns()
func (_StakingPool *StakingPoolSession) Dissolve() (*types.Transaction, error) {
	return _StakingPool.Contract.Dissolve(&_StakingPool.TransactOpts)
}

// Dissolve is a paid mutator transaction binding the contract method 0x3bef8a3a.
//
// Solidity: function dissolve() returns()
func (_StakingPool *StakingPoolTransactorSession) Dissolve() (*types.Transaction, error) {
	return _StakingPool.Contract.Dissolve(&_StakingPool.TransactOpts)
}

// Initialise is a paid mutator transaction binding the contract method 0xdd0ddfcf.
//
// Solidity: function initialise(address _nodeAddress, uint8 _depositType) returns()
func (_StakingPool *StakingPoolTransactor) Initialise(opts *bind.TransactOpts, _nodeAddress common.Address, _depositType uint8) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "initialise", _nodeAddress, _depositType)
}

// Initialise is a paid mutator transaction binding the contract method 0xdd0ddfcf.
//
// Solidity: function initialise(address _nodeAddress, uint8 _depositType) returns()
func (_StakingPool *StakingPoolSession) Initialise(_nodeAddress common.Address, _depositType uint8) (*types.Transaction, error) {
	return _StakingPool.Contract.Initialise(&_StakingPool.TransactOpts, _nodeAddress, _depositType)
}

// Initialise is a paid mutator transaction binding the contract method 0xdd0ddfcf.
//
// Solidity: function initialise(address _nodeAddress, uint8 _depositType) returns()
func (_StakingPool *StakingPoolTransactorSession) Initialise(_nodeAddress common.Address, _depositType uint8) (*types.Transaction, error) {
	return _StakingPool.Contract.Initialise(&_StakingPool.TransactOpts, _nodeAddress, _depositType)
}

// NodeDeposit is a paid mutator transaction binding the contract method 0x7476a6c3.
//
// Solidity: function nodeDeposit(bytes _validatorPubkey, bytes _validatorSignature, bytes32 _depositDataRoot) payable returns()
func (_StakingPool *StakingPoolTransactor) NodeDeposit(opts *bind.TransactOpts, _validatorPubkey []byte, _validatorSignature []byte, _depositDataRoot [32]byte) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "nodeDeposit", _validatorPubkey, _validatorSignature, _depositDataRoot)
}

// NodeDeposit is a paid mutator transaction binding the contract method 0x7476a6c3.
//
// Solidity: function nodeDeposit(bytes _validatorPubkey, bytes _validatorSignature, bytes32 _depositDataRoot) payable returns()
func (_StakingPool *StakingPoolSession) NodeDeposit(_validatorPubkey []byte, _validatorSignature []byte, _depositDataRoot [32]byte) (*types.Transaction, error) {
	return _StakingPool.Contract.NodeDeposit(&_StakingPool.TransactOpts, _validatorPubkey, _validatorSignature, _depositDataRoot)
}

// NodeDeposit is a paid mutator transaction binding the contract method 0x7476a6c3.
//
// Solidity: function nodeDeposit(bytes _validatorPubkey, bytes _validatorSignature, bytes32 _depositDataRoot) payable returns()
func (_StakingPool *StakingPoolTransactorSession) NodeDeposit(_validatorPubkey []byte, _validatorSignature []byte, _depositDataRoot [32]byte) (*types.Transaction, error) {
	return _StakingPool.Contract.NodeDeposit(&_StakingPool.TransactOpts, _validatorPubkey, _validatorSignature, _depositDataRoot)
}

// Refund is a paid mutator transaction binding the contract method 0x590e1ae3.
//
// Solidity: function refund() returns()
func (_StakingPool *StakingPoolTransactor) Refund(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "refund")
}

// Refund is a paid mutator transaction binding the contract method 0x590e1ae3.
//
// Solidity: function refund() returns()
func (_StakingPool *StakingPoolSession) Refund() (*types.Transaction, error) {
	return _StakingPool.Contract.Refund(&_StakingPool.TransactOpts)
}

// Refund is a paid mutator transaction binding the contract method 0x590e1ae3.
//
// Solidity: function refund() returns()
func (_StakingPool *StakingPoolTransactorSession) Refund() (*types.Transaction, error) {
	return _StakingPool.Contract.Refund(&_StakingPool.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0xf7ae36d1.
//
// Solidity: function stake(bytes _validatorSignature, bytes32 _depositDataRoot) returns()
func (_StakingPool *StakingPoolTransactor) Stake(opts *bind.TransactOpts, _validatorSignature []byte, _depositDataRoot [32]byte) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "stake", _validatorSignature, _depositDataRoot)
}

// Stake is a paid mutator transaction binding the contract method 0xf7ae36d1.
//
// Solidity: function stake(bytes _validatorSignature, bytes32 _depositDataRoot) returns()
func (_StakingPool *StakingPoolSession) Stake(_validatorSignature []byte, _depositDataRoot [32]byte) (*types.Transaction, error) {
	return _StakingPool.Contract.Stake(&_StakingPool.TransactOpts, _validatorSignature, _depositDataRoot)
}

// Stake is a paid mutator transaction binding the contract method 0xf7ae36d1.
//
// Solidity: function stake(bytes _validatorSignature, bytes32 _depositDataRoot) returns()
func (_StakingPool *StakingPoolTransactorSession) Stake(_validatorSignature []byte, _depositDataRoot [32]byte) (*types.Transaction, error) {
	return _StakingPool.Contract.Stake(&_StakingPool.TransactOpts, _validatorSignature, _depositDataRoot)
}

// UserDeposit is a paid mutator transaction binding the contract method 0x48146113.
//
// Solidity: function userDeposit() payable returns()
func (_StakingPool *StakingPoolTransactor) UserDeposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "userDeposit")
}

// UserDeposit is a paid mutator transaction binding the contract method 0x48146113.
//
// Solidity: function userDeposit() payable returns()
func (_StakingPool *StakingPoolSession) UserDeposit() (*types.Transaction, error) {
	return _StakingPool.Contract.UserDeposit(&_StakingPool.TransactOpts)
}

// UserDeposit is a paid mutator transaction binding the contract method 0x48146113.
//
// Solidity: function userDeposit() payable returns()
func (_StakingPool *StakingPoolTransactorSession) UserDeposit() (*types.Transaction, error) {
	return _StakingPool.Contract.UserDeposit(&_StakingPool.TransactOpts)
}

// VoteWithdrawCredentials is a paid mutator transaction binding the contract method 0x6e63103f.
//
// Solidity: function voteWithdrawCredentials() returns()
func (_StakingPool *StakingPoolTransactor) VoteWithdrawCredentials(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StakingPool.contract.Transact(opts, "voteWithdrawCredentials")
}

// VoteWithdrawCredentials is a paid mutator transaction binding the contract method 0x6e63103f.
//
// Solidity: function voteWithdrawCredentials() returns()
func (_StakingPool *StakingPoolSession) VoteWithdrawCredentials() (*types.Transaction, error) {
	return _StakingPool.Contract.VoteWithdrawCredentials(&_StakingPool.TransactOpts)
}

// VoteWithdrawCredentials is a paid mutator transaction binding the contract method 0x6e63103f.
//
// Solidity: function voteWithdrawCredentials() returns()
func (_StakingPool *StakingPoolTransactorSession) VoteWithdrawCredentials() (*types.Transaction, error) {
	return _StakingPool.Contract.VoteWithdrawCredentials(&_StakingPool.TransactOpts)
}
