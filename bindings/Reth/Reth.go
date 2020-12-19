// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Reth

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

// RethABI is the input ABI used to generate the binding from.
const RethABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stafiStorageAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"EtherDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"TokensBurned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ethAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"TokensMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rethAmount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositExcess\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositRewards\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBurnEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCollateralRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_rethAmount\",\"type\":\"uint256\"}],\"name\":\"getEthValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getExchangeRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_ethAmount\",\"type\":\"uint256\"}],\"name\":\"getRethValue\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalCollateral\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_ethAmount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_value\",\"type\":\"bool\"}],\"name\":\"setBurnEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Reth is an auto generated Go binding around an Ethereum contract.
type Reth struct {
	RethCaller     // Read-only binding to the contract
	RethTransactor // Write-only binding to the contract
	RethFilterer   // Log filterer for contract events
}

// RethCaller is an auto generated read-only Go binding around an Ethereum contract.
type RethCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RethTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RethTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RethFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RethFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RethSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RethSession struct {
	Contract     *Reth             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RethCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RethCallerSession struct {
	Contract *RethCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RethTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RethTransactorSession struct {
	Contract     *RethTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RethRaw is an auto generated low-level Go binding around an Ethereum contract.
type RethRaw struct {
	Contract *Reth // Generic contract binding to access the raw methods on
}

// RethCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RethCallerRaw struct {
	Contract *RethCaller // Generic read-only contract binding to access the raw methods on
}

// RethTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RethTransactorRaw struct {
	Contract *RethTransactor // Generic write-only contract binding to access the raw methods on
}

// NewReth creates a new instance of Reth, bound to a specific deployed contract.
func NewReth(address common.Address, backend bind.ContractBackend) (*Reth, error) {
	contract, err := bindReth(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Reth{RethCaller: RethCaller{contract: contract}, RethTransactor: RethTransactor{contract: contract}, RethFilterer: RethFilterer{contract: contract}}, nil
}

// NewRethCaller creates a new read-only instance of Reth, bound to a specific deployed contract.
func NewRethCaller(address common.Address, caller bind.ContractCaller) (*RethCaller, error) {
	contract, err := bindReth(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RethCaller{contract: contract}, nil
}

// NewRethTransactor creates a new write-only instance of Reth, bound to a specific deployed contract.
func NewRethTransactor(address common.Address, transactor bind.ContractTransactor) (*RethTransactor, error) {
	contract, err := bindReth(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RethTransactor{contract: contract}, nil
}

// NewRethFilterer creates a new log filterer instance of Reth, bound to a specific deployed contract.
func NewRethFilterer(address common.Address, filterer bind.ContractFilterer) (*RethFilterer, error) {
	contract, err := bindReth(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RethFilterer{contract: contract}, nil
}

// bindReth binds a generic wrapper to an already deployed contract.
func bindReth(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RethABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reth *RethRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reth.Contract.RethCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reth *RethRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reth.Contract.RethTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reth *RethRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reth.Contract.RethTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Reth *RethCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Reth.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Reth *RethTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reth.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Reth *RethTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Reth.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Reth *RethCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Reth *RethSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Reth.Contract.Allowance(&_Reth.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Reth *RethCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Reth.Contract.Allowance(&_Reth.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Reth *RethCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Reth *RethSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Reth.Contract.BalanceOf(&_Reth.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Reth *RethCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Reth.Contract.BalanceOf(&_Reth.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Reth *RethCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Reth *RethSession) Decimals() (uint8, error) {
	return _Reth.Contract.Decimals(&_Reth.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Reth *RethCallerSession) Decimals() (uint8, error) {
	return _Reth.Contract.Decimals(&_Reth.CallOpts)
}

// GetBurnEnabled is a free data retrieval call binding the contract method 0xf8e80db3.
//
// Solidity: function getBurnEnabled() view returns(bool)
func (_Reth *RethCaller) GetBurnEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "getBurnEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetBurnEnabled is a free data retrieval call binding the contract method 0xf8e80db3.
//
// Solidity: function getBurnEnabled() view returns(bool)
func (_Reth *RethSession) GetBurnEnabled() (bool, error) {
	return _Reth.Contract.GetBurnEnabled(&_Reth.CallOpts)
}

// GetBurnEnabled is a free data retrieval call binding the contract method 0xf8e80db3.
//
// Solidity: function getBurnEnabled() view returns(bool)
func (_Reth *RethCallerSession) GetBurnEnabled() (bool, error) {
	return _Reth.Contract.GetBurnEnabled(&_Reth.CallOpts)
}

// GetCollateralRate is a free data retrieval call binding the contract method 0x852185fc.
//
// Solidity: function getCollateralRate() view returns(uint256)
func (_Reth *RethCaller) GetCollateralRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "getCollateralRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCollateralRate is a free data retrieval call binding the contract method 0x852185fc.
//
// Solidity: function getCollateralRate() view returns(uint256)
func (_Reth *RethSession) GetCollateralRate() (*big.Int, error) {
	return _Reth.Contract.GetCollateralRate(&_Reth.CallOpts)
}

// GetCollateralRate is a free data retrieval call binding the contract method 0x852185fc.
//
// Solidity: function getCollateralRate() view returns(uint256)
func (_Reth *RethCallerSession) GetCollateralRate() (*big.Int, error) {
	return _Reth.Contract.GetCollateralRate(&_Reth.CallOpts)
}

// GetEthValue is a free data retrieval call binding the contract method 0x8b32fa23.
//
// Solidity: function getEthValue(uint256 _rethAmount) view returns(uint256)
func (_Reth *RethCaller) GetEthValue(opts *bind.CallOpts, _rethAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "getEthValue", _rethAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEthValue is a free data retrieval call binding the contract method 0x8b32fa23.
//
// Solidity: function getEthValue(uint256 _rethAmount) view returns(uint256)
func (_Reth *RethSession) GetEthValue(_rethAmount *big.Int) (*big.Int, error) {
	return _Reth.Contract.GetEthValue(&_Reth.CallOpts, _rethAmount)
}

// GetEthValue is a free data retrieval call binding the contract method 0x8b32fa23.
//
// Solidity: function getEthValue(uint256 _rethAmount) view returns(uint256)
func (_Reth *RethCallerSession) GetEthValue(_rethAmount *big.Int) (*big.Int, error) {
	return _Reth.Contract.GetEthValue(&_Reth.CallOpts, _rethAmount)
}

// GetExchangeRate is a free data retrieval call binding the contract method 0xe6aa216c.
//
// Solidity: function getExchangeRate() view returns(uint256)
func (_Reth *RethCaller) GetExchangeRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "getExchangeRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExchangeRate is a free data retrieval call binding the contract method 0xe6aa216c.
//
// Solidity: function getExchangeRate() view returns(uint256)
func (_Reth *RethSession) GetExchangeRate() (*big.Int, error) {
	return _Reth.Contract.GetExchangeRate(&_Reth.CallOpts)
}

// GetExchangeRate is a free data retrieval call binding the contract method 0xe6aa216c.
//
// Solidity: function getExchangeRate() view returns(uint256)
func (_Reth *RethCallerSession) GetExchangeRate() (*big.Int, error) {
	return _Reth.Contract.GetExchangeRate(&_Reth.CallOpts)
}

// GetRethValue is a free data retrieval call binding the contract method 0x4346f03e.
//
// Solidity: function getRethValue(uint256 _ethAmount) view returns(uint256)
func (_Reth *RethCaller) GetRethValue(opts *bind.CallOpts, _ethAmount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "getRethValue", _ethAmount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRethValue is a free data retrieval call binding the contract method 0x4346f03e.
//
// Solidity: function getRethValue(uint256 _ethAmount) view returns(uint256)
func (_Reth *RethSession) GetRethValue(_ethAmount *big.Int) (*big.Int, error) {
	return _Reth.Contract.GetRethValue(&_Reth.CallOpts, _ethAmount)
}

// GetRethValue is a free data retrieval call binding the contract method 0x4346f03e.
//
// Solidity: function getRethValue(uint256 _ethAmount) view returns(uint256)
func (_Reth *RethCallerSession) GetRethValue(_ethAmount *big.Int) (*big.Int, error) {
	return _Reth.Contract.GetRethValue(&_Reth.CallOpts, _ethAmount)
}

// GetTotalCollateral is a free data retrieval call binding the contract method 0xd6eb5910.
//
// Solidity: function getTotalCollateral() view returns(uint256)
func (_Reth *RethCaller) GetTotalCollateral(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "getTotalCollateral")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalCollateral is a free data retrieval call binding the contract method 0xd6eb5910.
//
// Solidity: function getTotalCollateral() view returns(uint256)
func (_Reth *RethSession) GetTotalCollateral() (*big.Int, error) {
	return _Reth.Contract.GetTotalCollateral(&_Reth.CallOpts)
}

// GetTotalCollateral is a free data retrieval call binding the contract method 0xd6eb5910.
//
// Solidity: function getTotalCollateral() view returns(uint256)
func (_Reth *RethCallerSession) GetTotalCollateral() (*big.Int, error) {
	return _Reth.Contract.GetTotalCollateral(&_Reth.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Reth *RethCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Reth *RethSession) Name() (string, error) {
	return _Reth.Contract.Name(&_Reth.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Reth *RethCallerSession) Name() (string, error) {
	return _Reth.Contract.Name(&_Reth.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Reth *RethCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Reth *RethSession) Symbol() (string, error) {
	return _Reth.Contract.Symbol(&_Reth.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Reth *RethCallerSession) Symbol() (string, error) {
	return _Reth.Contract.Symbol(&_Reth.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Reth *RethCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Reth *RethSession) TotalSupply() (*big.Int, error) {
	return _Reth.Contract.TotalSupply(&_Reth.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Reth *RethCallerSession) TotalSupply() (*big.Int, error) {
	return _Reth.Contract.TotalSupply(&_Reth.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_Reth *RethCaller) Version(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Reth.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_Reth *RethSession) Version() (uint8, error) {
	return _Reth.Contract.Version(&_Reth.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_Reth *RethCallerSession) Version() (uint8, error) {
	return _Reth.Contract.Version(&_Reth.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Reth *RethTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Reth *RethSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.Approve(&_Reth.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Reth *RethTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.Approve(&_Reth.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 _rethAmount) returns()
func (_Reth *RethTransactor) Burn(opts *bind.TransactOpts, _rethAmount *big.Int) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "burn", _rethAmount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 _rethAmount) returns()
func (_Reth *RethSession) Burn(_rethAmount *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.Burn(&_Reth.TransactOpts, _rethAmount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 _rethAmount) returns()
func (_Reth *RethTransactorSession) Burn(_rethAmount *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.Burn(&_Reth.TransactOpts, _rethAmount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Reth *RethTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Reth *RethSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.DecreaseAllowance(&_Reth.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Reth *RethTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.DecreaseAllowance(&_Reth.TransactOpts, spender, subtractedValue)
}

// DepositExcess is a paid mutator transaction binding the contract method 0x6c985a88.
//
// Solidity: function depositExcess() payable returns()
func (_Reth *RethTransactor) DepositExcess(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "depositExcess")
}

// DepositExcess is a paid mutator transaction binding the contract method 0x6c985a88.
//
// Solidity: function depositExcess() payable returns()
func (_Reth *RethSession) DepositExcess() (*types.Transaction, error) {
	return _Reth.Contract.DepositExcess(&_Reth.TransactOpts)
}

// DepositExcess is a paid mutator transaction binding the contract method 0x6c985a88.
//
// Solidity: function depositExcess() payable returns()
func (_Reth *RethTransactorSession) DepositExcess() (*types.Transaction, error) {
	return _Reth.Contract.DepositExcess(&_Reth.TransactOpts)
}

// DepositRewards is a paid mutator transaction binding the contract method 0x152111f7.
//
// Solidity: function depositRewards() payable returns()
func (_Reth *RethTransactor) DepositRewards(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "depositRewards")
}

// DepositRewards is a paid mutator transaction binding the contract method 0x152111f7.
//
// Solidity: function depositRewards() payable returns()
func (_Reth *RethSession) DepositRewards() (*types.Transaction, error) {
	return _Reth.Contract.DepositRewards(&_Reth.TransactOpts)
}

// DepositRewards is a paid mutator transaction binding the contract method 0x152111f7.
//
// Solidity: function depositRewards() payable returns()
func (_Reth *RethTransactorSession) DepositRewards() (*types.Transaction, error) {
	return _Reth.Contract.DepositRewards(&_Reth.TransactOpts)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Reth *RethTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Reth *RethSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.IncreaseAllowance(&_Reth.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Reth *RethTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.IncreaseAllowance(&_Reth.TransactOpts, spender, addedValue)
}

// Mint is a paid mutator transaction binding the contract method 0x94bf804d.
//
// Solidity: function mint(uint256 _ethAmount, address _to) returns()
func (_Reth *RethTransactor) Mint(opts *bind.TransactOpts, _ethAmount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "mint", _ethAmount, _to)
}

// Mint is a paid mutator transaction binding the contract method 0x94bf804d.
//
// Solidity: function mint(uint256 _ethAmount, address _to) returns()
func (_Reth *RethSession) Mint(_ethAmount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _Reth.Contract.Mint(&_Reth.TransactOpts, _ethAmount, _to)
}

// Mint is a paid mutator transaction binding the contract method 0x94bf804d.
//
// Solidity: function mint(uint256 _ethAmount, address _to) returns()
func (_Reth *RethTransactorSession) Mint(_ethAmount *big.Int, _to common.Address) (*types.Transaction, error) {
	return _Reth.Contract.Mint(&_Reth.TransactOpts, _ethAmount, _to)
}

// SetBurnEnabled is a paid mutator transaction binding the contract method 0x7b2c835f.
//
// Solidity: function setBurnEnabled(bool _value) returns()
func (_Reth *RethTransactor) SetBurnEnabled(opts *bind.TransactOpts, _value bool) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "setBurnEnabled", _value)
}

// SetBurnEnabled is a paid mutator transaction binding the contract method 0x7b2c835f.
//
// Solidity: function setBurnEnabled(bool _value) returns()
func (_Reth *RethSession) SetBurnEnabled(_value bool) (*types.Transaction, error) {
	return _Reth.Contract.SetBurnEnabled(&_Reth.TransactOpts, _value)
}

// SetBurnEnabled is a paid mutator transaction binding the contract method 0x7b2c835f.
//
// Solidity: function setBurnEnabled(bool _value) returns()
func (_Reth *RethTransactorSession) SetBurnEnabled(_value bool) (*types.Transaction, error) {
	return _Reth.Contract.SetBurnEnabled(&_Reth.TransactOpts, _value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Reth *RethTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Reth *RethSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.Transfer(&_Reth.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Reth *RethTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.Transfer(&_Reth.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Reth *RethTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reth.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Reth *RethSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.TransferFrom(&_Reth.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Reth *RethTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Reth.Contract.TransferFrom(&_Reth.TransactOpts, sender, recipient, amount)
}

// RethApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Reth contract.
type RethApprovalIterator struct {
	Event *RethApproval // Event containing the contract specifics and raw log

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
func (it *RethApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RethApproval)
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
		it.Event = new(RethApproval)
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
func (it *RethApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RethApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RethApproval represents a Approval event raised by the Reth contract.
type RethApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Reth *RethFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*RethApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Reth.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &RethApprovalIterator{contract: _Reth.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Reth *RethFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *RethApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Reth.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RethApproval)
				if err := _Reth.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Reth *RethFilterer) ParseApproval(log types.Log) (*RethApproval, error) {
	event := new(RethApproval)
	if err := _Reth.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RethEtherDepositedIterator is returned from FilterEtherDeposited and is used to iterate over the raw logs and unpacked data for EtherDeposited events raised by the Reth contract.
type RethEtherDepositedIterator struct {
	Event *RethEtherDeposited // Event containing the contract specifics and raw log

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
func (it *RethEtherDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RethEtherDeposited)
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
		it.Event = new(RethEtherDeposited)
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
func (it *RethEtherDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RethEtherDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RethEtherDeposited represents a EtherDeposited event raised by the Reth contract.
type RethEtherDeposited struct {
	From   common.Address
	Amount *big.Int
	Time   *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterEtherDeposited is a free log retrieval operation binding the contract event 0xef51b4c870b8b0100eae2072e91db01222a303072af3728e58c9d4d2da33127f.
//
// Solidity: event EtherDeposited(address indexed from, uint256 amount, uint256 time)
func (_Reth *RethFilterer) FilterEtherDeposited(opts *bind.FilterOpts, from []common.Address) (*RethEtherDepositedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Reth.contract.FilterLogs(opts, "EtherDeposited", fromRule)
	if err != nil {
		return nil, err
	}
	return &RethEtherDepositedIterator{contract: _Reth.contract, event: "EtherDeposited", logs: logs, sub: sub}, nil
}

// WatchEtherDeposited is a free log subscription operation binding the contract event 0xef51b4c870b8b0100eae2072e91db01222a303072af3728e58c9d4d2da33127f.
//
// Solidity: event EtherDeposited(address indexed from, uint256 amount, uint256 time)
func (_Reth *RethFilterer) WatchEtherDeposited(opts *bind.WatchOpts, sink chan<- *RethEtherDeposited, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Reth.contract.WatchLogs(opts, "EtherDeposited", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RethEtherDeposited)
				if err := _Reth.contract.UnpackLog(event, "EtherDeposited", log); err != nil {
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
func (_Reth *RethFilterer) ParseEtherDeposited(log types.Log) (*RethEtherDeposited, error) {
	event := new(RethEtherDeposited)
	if err := _Reth.contract.UnpackLog(event, "EtherDeposited", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RethTokensBurnedIterator is returned from FilterTokensBurned and is used to iterate over the raw logs and unpacked data for TokensBurned events raised by the Reth contract.
type RethTokensBurnedIterator struct {
	Event *RethTokensBurned // Event containing the contract specifics and raw log

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
func (it *RethTokensBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RethTokensBurned)
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
		it.Event = new(RethTokensBurned)
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
func (it *RethTokensBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RethTokensBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RethTokensBurned represents a TokensBurned event raised by the Reth contract.
type RethTokensBurned struct {
	From      common.Address
	Amount    *big.Int
	EthAmount *big.Int
	Time      *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTokensBurned is a free log retrieval operation binding the contract event 0x19783b34589160c168487dc7f9c51ae0bcefe67a47d6708fba90f6ce0366d3d1.
//
// Solidity: event TokensBurned(address indexed from, uint256 amount, uint256 ethAmount, uint256 time)
func (_Reth *RethFilterer) FilterTokensBurned(opts *bind.FilterOpts, from []common.Address) (*RethTokensBurnedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Reth.contract.FilterLogs(opts, "TokensBurned", fromRule)
	if err != nil {
		return nil, err
	}
	return &RethTokensBurnedIterator{contract: _Reth.contract, event: "TokensBurned", logs: logs, sub: sub}, nil
}

// WatchTokensBurned is a free log subscription operation binding the contract event 0x19783b34589160c168487dc7f9c51ae0bcefe67a47d6708fba90f6ce0366d3d1.
//
// Solidity: event TokensBurned(address indexed from, uint256 amount, uint256 ethAmount, uint256 time)
func (_Reth *RethFilterer) WatchTokensBurned(opts *bind.WatchOpts, sink chan<- *RethTokensBurned, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Reth.contract.WatchLogs(opts, "TokensBurned", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RethTokensBurned)
				if err := _Reth.contract.UnpackLog(event, "TokensBurned", log); err != nil {
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

// ParseTokensBurned is a log parse operation binding the contract event 0x19783b34589160c168487dc7f9c51ae0bcefe67a47d6708fba90f6ce0366d3d1.
//
// Solidity: event TokensBurned(address indexed from, uint256 amount, uint256 ethAmount, uint256 time)
func (_Reth *RethFilterer) ParseTokensBurned(log types.Log) (*RethTokensBurned, error) {
	event := new(RethTokensBurned)
	if err := _Reth.contract.UnpackLog(event, "TokensBurned", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RethTokensMintedIterator is returned from FilterTokensMinted and is used to iterate over the raw logs and unpacked data for TokensMinted events raised by the Reth contract.
type RethTokensMintedIterator struct {
	Event *RethTokensMinted // Event containing the contract specifics and raw log

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
func (it *RethTokensMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RethTokensMinted)
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
		it.Event = new(RethTokensMinted)
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
func (it *RethTokensMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RethTokensMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RethTokensMinted represents a TokensMinted event raised by the Reth contract.
type RethTokensMinted struct {
	To        common.Address
	Amount    *big.Int
	EthAmount *big.Int
	Time      *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTokensMinted is a free log retrieval operation binding the contract event 0x6155cfd0fd028b0ca77e8495a60cbe563e8bce8611f0aad6fedbdaafc05d44a2.
//
// Solidity: event TokensMinted(address indexed to, uint256 amount, uint256 ethAmount, uint256 time)
func (_Reth *RethFilterer) FilterTokensMinted(opts *bind.FilterOpts, to []common.Address) (*RethTokensMintedIterator, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Reth.contract.FilterLogs(opts, "TokensMinted", toRule)
	if err != nil {
		return nil, err
	}
	return &RethTokensMintedIterator{contract: _Reth.contract, event: "TokensMinted", logs: logs, sub: sub}, nil
}

// WatchTokensMinted is a free log subscription operation binding the contract event 0x6155cfd0fd028b0ca77e8495a60cbe563e8bce8611f0aad6fedbdaafc05d44a2.
//
// Solidity: event TokensMinted(address indexed to, uint256 amount, uint256 ethAmount, uint256 time)
func (_Reth *RethFilterer) WatchTokensMinted(opts *bind.WatchOpts, sink chan<- *RethTokensMinted, to []common.Address) (event.Subscription, error) {

	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Reth.contract.WatchLogs(opts, "TokensMinted", toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RethTokensMinted)
				if err := _Reth.contract.UnpackLog(event, "TokensMinted", log); err != nil {
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

// ParseTokensMinted is a log parse operation binding the contract event 0x6155cfd0fd028b0ca77e8495a60cbe563e8bce8611f0aad6fedbdaafc05d44a2.
//
// Solidity: event TokensMinted(address indexed to, uint256 amount, uint256 ethAmount, uint256 time)
func (_Reth *RethFilterer) ParseTokensMinted(log types.Log) (*RethTokensMinted, error) {
	event := new(RethTokensMinted)
	if err := _Reth.contract.UnpackLog(event, "TokensMinted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// RethTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Reth contract.
type RethTransferIterator struct {
	Event *RethTransfer // Event containing the contract specifics and raw log

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
func (it *RethTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RethTransfer)
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
		it.Event = new(RethTransfer)
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
func (it *RethTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RethTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RethTransfer represents a Transfer event raised by the Reth contract.
type RethTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Reth *RethFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*RethTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Reth.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &RethTransferIterator{contract: _Reth.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Reth *RethFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *RethTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Reth.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RethTransfer)
				if err := _Reth.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Reth *RethFilterer) ParseTransfer(log types.Log) (*RethTransfer, error) {
	event := new(RethTransfer)
	if err := _Reth.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	return event, nil
}
