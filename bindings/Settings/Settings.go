// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Settings

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

// SettingsABI is the input ABI used to generate the binding from.
const SettingsABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stafiStorageAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"getNodeConsensusThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeRefundRatio\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getNodeTrustedRefundRatio\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getPlatformFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProcessWithdrawalsEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getSubmitBalancesEnabled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getWithdrawalCredentials\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"setNodeConsensusThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"setNodeFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"setNodeRefundRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"setNodeTrustedRefundRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"setPlatformFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_value\",\"type\":\"bool\"}],\"name\":\"setProcessWithdrawalsEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"_value\",\"type\":\"bool\"}],\"name\":\"setSubmitBalancesEnabled\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_value\",\"type\":\"bytes\"}],\"name\":\"setWithdrawalCredentials\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Settings is an auto generated Go binding around an Ethereum contract.
type Settings struct {
	SettingsCaller     // Read-only binding to the contract
	SettingsTransactor // Write-only binding to the contract
	SettingsFilterer   // Log filterer for contract events
}

// SettingsCaller is an auto generated read-only Go binding around an Ethereum contract.
type SettingsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettingsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SettingsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettingsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SettingsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SettingsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SettingsSession struct {
	Contract     *Settings         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SettingsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SettingsCallerSession struct {
	Contract *SettingsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SettingsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SettingsTransactorSession struct {
	Contract     *SettingsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SettingsRaw is an auto generated low-level Go binding around an Ethereum contract.
type SettingsRaw struct {
	Contract *Settings // Generic contract binding to access the raw methods on
}

// SettingsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SettingsCallerRaw struct {
	Contract *SettingsCaller // Generic read-only contract binding to access the raw methods on
}

// SettingsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SettingsTransactorRaw struct {
	Contract *SettingsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSettings creates a new instance of Settings, bound to a specific deployed contract.
func NewSettings(address common.Address, backend bind.ContractBackend) (*Settings, error) {
	contract, err := bindSettings(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Settings{SettingsCaller: SettingsCaller{contract: contract}, SettingsTransactor: SettingsTransactor{contract: contract}, SettingsFilterer: SettingsFilterer{contract: contract}}, nil
}

// NewSettingsCaller creates a new read-only instance of Settings, bound to a specific deployed contract.
func NewSettingsCaller(address common.Address, caller bind.ContractCaller) (*SettingsCaller, error) {
	contract, err := bindSettings(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SettingsCaller{contract: contract}, nil
}

// NewSettingsTransactor creates a new write-only instance of Settings, bound to a specific deployed contract.
func NewSettingsTransactor(address common.Address, transactor bind.ContractTransactor) (*SettingsTransactor, error) {
	contract, err := bindSettings(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SettingsTransactor{contract: contract}, nil
}

// NewSettingsFilterer creates a new log filterer instance of Settings, bound to a specific deployed contract.
func NewSettingsFilterer(address common.Address, filterer bind.ContractFilterer) (*SettingsFilterer, error) {
	contract, err := bindSettings(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SettingsFilterer{contract: contract}, nil
}

// bindSettings binds a generic wrapper to an already deployed contract.
func bindSettings(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SettingsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Settings *SettingsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Settings.Contract.SettingsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Settings *SettingsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settings.Contract.SettingsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Settings *SettingsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Settings.Contract.SettingsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Settings *SettingsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Settings.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Settings *SettingsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Settings.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Settings *SettingsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Settings.Contract.contract.Transact(opts, method, params...)
}

// GetNodeConsensusThreshold is a free data retrieval call binding the contract method 0x1f66e8ed.
//
// Solidity: function getNodeConsensusThreshold() view returns(uint256)
func (_Settings *SettingsCaller) GetNodeConsensusThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Settings.contract.Call(opts, &out, "getNodeConsensusThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeConsensusThreshold is a free data retrieval call binding the contract method 0x1f66e8ed.
//
// Solidity: function getNodeConsensusThreshold() view returns(uint256)
func (_Settings *SettingsSession) GetNodeConsensusThreshold() (*big.Int, error) {
	return _Settings.Contract.GetNodeConsensusThreshold(&_Settings.CallOpts)
}

// GetNodeConsensusThreshold is a free data retrieval call binding the contract method 0x1f66e8ed.
//
// Solidity: function getNodeConsensusThreshold() view returns(uint256)
func (_Settings *SettingsCallerSession) GetNodeConsensusThreshold() (*big.Int, error) {
	return _Settings.Contract.GetNodeConsensusThreshold(&_Settings.CallOpts)
}

// GetNodeFee is a free data retrieval call binding the contract method 0xe7150134.
//
// Solidity: function getNodeFee() view returns(uint256)
func (_Settings *SettingsCaller) GetNodeFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Settings.contract.Call(opts, &out, "getNodeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeFee is a free data retrieval call binding the contract method 0xe7150134.
//
// Solidity: function getNodeFee() view returns(uint256)
func (_Settings *SettingsSession) GetNodeFee() (*big.Int, error) {
	return _Settings.Contract.GetNodeFee(&_Settings.CallOpts)
}

// GetNodeFee is a free data retrieval call binding the contract method 0xe7150134.
//
// Solidity: function getNodeFee() view returns(uint256)
func (_Settings *SettingsCallerSession) GetNodeFee() (*big.Int, error) {
	return _Settings.Contract.GetNodeFee(&_Settings.CallOpts)
}

// GetNodeRefundRatio is a free data retrieval call binding the contract method 0xeb364204.
//
// Solidity: function getNodeRefundRatio() view returns(uint256)
func (_Settings *SettingsCaller) GetNodeRefundRatio(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Settings.contract.Call(opts, &out, "getNodeRefundRatio")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeRefundRatio is a free data retrieval call binding the contract method 0xeb364204.
//
// Solidity: function getNodeRefundRatio() view returns(uint256)
func (_Settings *SettingsSession) GetNodeRefundRatio() (*big.Int, error) {
	return _Settings.Contract.GetNodeRefundRatio(&_Settings.CallOpts)
}

// GetNodeRefundRatio is a free data retrieval call binding the contract method 0xeb364204.
//
// Solidity: function getNodeRefundRatio() view returns(uint256)
func (_Settings *SettingsCallerSession) GetNodeRefundRatio() (*big.Int, error) {
	return _Settings.Contract.GetNodeRefundRatio(&_Settings.CallOpts)
}

// GetNodeTrustedRefundRatio is a free data retrieval call binding the contract method 0xf4cf0830.
//
// Solidity: function getNodeTrustedRefundRatio() view returns(uint256)
func (_Settings *SettingsCaller) GetNodeTrustedRefundRatio(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Settings.contract.Call(opts, &out, "getNodeTrustedRefundRatio")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeTrustedRefundRatio is a free data retrieval call binding the contract method 0xf4cf0830.
//
// Solidity: function getNodeTrustedRefundRatio() view returns(uint256)
func (_Settings *SettingsSession) GetNodeTrustedRefundRatio() (*big.Int, error) {
	return _Settings.Contract.GetNodeTrustedRefundRatio(&_Settings.CallOpts)
}

// GetNodeTrustedRefundRatio is a free data retrieval call binding the contract method 0xf4cf0830.
//
// Solidity: function getNodeTrustedRefundRatio() view returns(uint256)
func (_Settings *SettingsCallerSession) GetNodeTrustedRefundRatio() (*big.Int, error) {
	return _Settings.Contract.GetNodeTrustedRefundRatio(&_Settings.CallOpts)
}

// GetPlatformFee is a free data retrieval call binding the contract method 0x6ea8bc10.
//
// Solidity: function getPlatformFee() view returns(uint256)
func (_Settings *SettingsCaller) GetPlatformFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Settings.contract.Call(opts, &out, "getPlatformFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPlatformFee is a free data retrieval call binding the contract method 0x6ea8bc10.
//
// Solidity: function getPlatformFee() view returns(uint256)
func (_Settings *SettingsSession) GetPlatformFee() (*big.Int, error) {
	return _Settings.Contract.GetPlatformFee(&_Settings.CallOpts)
}

// GetPlatformFee is a free data retrieval call binding the contract method 0x6ea8bc10.
//
// Solidity: function getPlatformFee() view returns(uint256)
func (_Settings *SettingsCallerSession) GetPlatformFee() (*big.Int, error) {
	return _Settings.Contract.GetPlatformFee(&_Settings.CallOpts)
}

// GetProcessWithdrawalsEnabled is a free data retrieval call binding the contract method 0x28ea8d06.
//
// Solidity: function getProcessWithdrawalsEnabled() view returns(bool)
func (_Settings *SettingsCaller) GetProcessWithdrawalsEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Settings.contract.Call(opts, &out, "getProcessWithdrawalsEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetProcessWithdrawalsEnabled is a free data retrieval call binding the contract method 0x28ea8d06.
//
// Solidity: function getProcessWithdrawalsEnabled() view returns(bool)
func (_Settings *SettingsSession) GetProcessWithdrawalsEnabled() (bool, error) {
	return _Settings.Contract.GetProcessWithdrawalsEnabled(&_Settings.CallOpts)
}

// GetProcessWithdrawalsEnabled is a free data retrieval call binding the contract method 0x28ea8d06.
//
// Solidity: function getProcessWithdrawalsEnabled() view returns(bool)
func (_Settings *SettingsCallerSession) GetProcessWithdrawalsEnabled() (bool, error) {
	return _Settings.Contract.GetProcessWithdrawalsEnabled(&_Settings.CallOpts)
}

// GetSubmitBalancesEnabled is a free data retrieval call binding the contract method 0xfcdb60db.
//
// Solidity: function getSubmitBalancesEnabled() view returns(bool)
func (_Settings *SettingsCaller) GetSubmitBalancesEnabled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Settings.contract.Call(opts, &out, "getSubmitBalancesEnabled")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetSubmitBalancesEnabled is a free data retrieval call binding the contract method 0xfcdb60db.
//
// Solidity: function getSubmitBalancesEnabled() view returns(bool)
func (_Settings *SettingsSession) GetSubmitBalancesEnabled() (bool, error) {
	return _Settings.Contract.GetSubmitBalancesEnabled(&_Settings.CallOpts)
}

// GetSubmitBalancesEnabled is a free data retrieval call binding the contract method 0xfcdb60db.
//
// Solidity: function getSubmitBalancesEnabled() view returns(bool)
func (_Settings *SettingsCallerSession) GetSubmitBalancesEnabled() (bool, error) {
	return _Settings.Contract.GetSubmitBalancesEnabled(&_Settings.CallOpts)
}

// GetWithdrawalCredentials is a free data retrieval call binding the contract method 0x56396715.
//
// Solidity: function getWithdrawalCredentials() view returns(bytes)
func (_Settings *SettingsCaller) GetWithdrawalCredentials(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Settings.contract.Call(opts, &out, "getWithdrawalCredentials")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetWithdrawalCredentials is a free data retrieval call binding the contract method 0x56396715.
//
// Solidity: function getWithdrawalCredentials() view returns(bytes)
func (_Settings *SettingsSession) GetWithdrawalCredentials() ([]byte, error) {
	return _Settings.Contract.GetWithdrawalCredentials(&_Settings.CallOpts)
}

// GetWithdrawalCredentials is a free data retrieval call binding the contract method 0x56396715.
//
// Solidity: function getWithdrawalCredentials() view returns(bytes)
func (_Settings *SettingsCallerSession) GetWithdrawalCredentials() ([]byte, error) {
	return _Settings.Contract.GetWithdrawalCredentials(&_Settings.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_Settings *SettingsCaller) Version(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Settings.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_Settings *SettingsSession) Version() (uint8, error) {
	return _Settings.Contract.Version(&_Settings.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_Settings *SettingsCallerSession) Version() (uint8, error) {
	return _Settings.Contract.Version(&_Settings.CallOpts)
}

// SetNodeConsensusThreshold is a paid mutator transaction binding the contract method 0x977886ea.
//
// Solidity: function setNodeConsensusThreshold(uint256 _value) returns()
func (_Settings *SettingsTransactor) SetNodeConsensusThreshold(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _Settings.contract.Transact(opts, "setNodeConsensusThreshold", _value)
}

// SetNodeConsensusThreshold is a paid mutator transaction binding the contract method 0x977886ea.
//
// Solidity: function setNodeConsensusThreshold(uint256 _value) returns()
func (_Settings *SettingsSession) SetNodeConsensusThreshold(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetNodeConsensusThreshold(&_Settings.TransactOpts, _value)
}

// SetNodeConsensusThreshold is a paid mutator transaction binding the contract method 0x977886ea.
//
// Solidity: function setNodeConsensusThreshold(uint256 _value) returns()
func (_Settings *SettingsTransactorSession) SetNodeConsensusThreshold(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetNodeConsensusThreshold(&_Settings.TransactOpts, _value)
}

// SetNodeFee is a paid mutator transaction binding the contract method 0x89093310.
//
// Solidity: function setNodeFee(uint256 _value) returns()
func (_Settings *SettingsTransactor) SetNodeFee(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _Settings.contract.Transact(opts, "setNodeFee", _value)
}

// SetNodeFee is a paid mutator transaction binding the contract method 0x89093310.
//
// Solidity: function setNodeFee(uint256 _value) returns()
func (_Settings *SettingsSession) SetNodeFee(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetNodeFee(&_Settings.TransactOpts, _value)
}

// SetNodeFee is a paid mutator transaction binding the contract method 0x89093310.
//
// Solidity: function setNodeFee(uint256 _value) returns()
func (_Settings *SettingsTransactorSession) SetNodeFee(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetNodeFee(&_Settings.TransactOpts, _value)
}

// SetNodeRefundRatio is a paid mutator transaction binding the contract method 0xb67b446a.
//
// Solidity: function setNodeRefundRatio(uint256 _value) returns()
func (_Settings *SettingsTransactor) SetNodeRefundRatio(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _Settings.contract.Transact(opts, "setNodeRefundRatio", _value)
}

// SetNodeRefundRatio is a paid mutator transaction binding the contract method 0xb67b446a.
//
// Solidity: function setNodeRefundRatio(uint256 _value) returns()
func (_Settings *SettingsSession) SetNodeRefundRatio(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetNodeRefundRatio(&_Settings.TransactOpts, _value)
}

// SetNodeRefundRatio is a paid mutator transaction binding the contract method 0xb67b446a.
//
// Solidity: function setNodeRefundRatio(uint256 _value) returns()
func (_Settings *SettingsTransactorSession) SetNodeRefundRatio(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetNodeRefundRatio(&_Settings.TransactOpts, _value)
}

// SetNodeTrustedRefundRatio is a paid mutator transaction binding the contract method 0x89799aa1.
//
// Solidity: function setNodeTrustedRefundRatio(uint256 _value) returns()
func (_Settings *SettingsTransactor) SetNodeTrustedRefundRatio(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _Settings.contract.Transact(opts, "setNodeTrustedRefundRatio", _value)
}

// SetNodeTrustedRefundRatio is a paid mutator transaction binding the contract method 0x89799aa1.
//
// Solidity: function setNodeTrustedRefundRatio(uint256 _value) returns()
func (_Settings *SettingsSession) SetNodeTrustedRefundRatio(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetNodeTrustedRefundRatio(&_Settings.TransactOpts, _value)
}

// SetNodeTrustedRefundRatio is a paid mutator transaction binding the contract method 0x89799aa1.
//
// Solidity: function setNodeTrustedRefundRatio(uint256 _value) returns()
func (_Settings *SettingsTransactorSession) SetNodeTrustedRefundRatio(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetNodeTrustedRefundRatio(&_Settings.TransactOpts, _value)
}

// SetPlatformFee is a paid mutator transaction binding the contract method 0x12e8e2c3.
//
// Solidity: function setPlatformFee(uint256 _value) returns()
func (_Settings *SettingsTransactor) SetPlatformFee(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _Settings.contract.Transact(opts, "setPlatformFee", _value)
}

// SetPlatformFee is a paid mutator transaction binding the contract method 0x12e8e2c3.
//
// Solidity: function setPlatformFee(uint256 _value) returns()
func (_Settings *SettingsSession) SetPlatformFee(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetPlatformFee(&_Settings.TransactOpts, _value)
}

// SetPlatformFee is a paid mutator transaction binding the contract method 0x12e8e2c3.
//
// Solidity: function setPlatformFee(uint256 _value) returns()
func (_Settings *SettingsTransactorSession) SetPlatformFee(_value *big.Int) (*types.Transaction, error) {
	return _Settings.Contract.SetPlatformFee(&_Settings.TransactOpts, _value)
}

// SetProcessWithdrawalsEnabled is a paid mutator transaction binding the contract method 0xc6fbae5e.
//
// Solidity: function setProcessWithdrawalsEnabled(bool _value) returns()
func (_Settings *SettingsTransactor) SetProcessWithdrawalsEnabled(opts *bind.TransactOpts, _value bool) (*types.Transaction, error) {
	return _Settings.contract.Transact(opts, "setProcessWithdrawalsEnabled", _value)
}

// SetProcessWithdrawalsEnabled is a paid mutator transaction binding the contract method 0xc6fbae5e.
//
// Solidity: function setProcessWithdrawalsEnabled(bool _value) returns()
func (_Settings *SettingsSession) SetProcessWithdrawalsEnabled(_value bool) (*types.Transaction, error) {
	return _Settings.Contract.SetProcessWithdrawalsEnabled(&_Settings.TransactOpts, _value)
}

// SetProcessWithdrawalsEnabled is a paid mutator transaction binding the contract method 0xc6fbae5e.
//
// Solidity: function setProcessWithdrawalsEnabled(bool _value) returns()
func (_Settings *SettingsTransactorSession) SetProcessWithdrawalsEnabled(_value bool) (*types.Transaction, error) {
	return _Settings.Contract.SetProcessWithdrawalsEnabled(&_Settings.TransactOpts, _value)
}

// SetSubmitBalancesEnabled is a paid mutator transaction binding the contract method 0xa37a9e61.
//
// Solidity: function setSubmitBalancesEnabled(bool _value) returns()
func (_Settings *SettingsTransactor) SetSubmitBalancesEnabled(opts *bind.TransactOpts, _value bool) (*types.Transaction, error) {
	return _Settings.contract.Transact(opts, "setSubmitBalancesEnabled", _value)
}

// SetSubmitBalancesEnabled is a paid mutator transaction binding the contract method 0xa37a9e61.
//
// Solidity: function setSubmitBalancesEnabled(bool _value) returns()
func (_Settings *SettingsSession) SetSubmitBalancesEnabled(_value bool) (*types.Transaction, error) {
	return _Settings.Contract.SetSubmitBalancesEnabled(&_Settings.TransactOpts, _value)
}

// SetSubmitBalancesEnabled is a paid mutator transaction binding the contract method 0xa37a9e61.
//
// Solidity: function setSubmitBalancesEnabled(bool _value) returns()
func (_Settings *SettingsTransactorSession) SetSubmitBalancesEnabled(_value bool) (*types.Transaction, error) {
	return _Settings.Contract.SetSubmitBalancesEnabled(&_Settings.TransactOpts, _value)
}

// SetWithdrawalCredentials is a paid mutator transaction binding the contract method 0xc4cf71c9.
//
// Solidity: function setWithdrawalCredentials(bytes _value) returns()
func (_Settings *SettingsTransactor) SetWithdrawalCredentials(opts *bind.TransactOpts, _value []byte) (*types.Transaction, error) {
	return _Settings.contract.Transact(opts, "setWithdrawalCredentials", _value)
}

// SetWithdrawalCredentials is a paid mutator transaction binding the contract method 0xc4cf71c9.
//
// Solidity: function setWithdrawalCredentials(bytes _value) returns()
func (_Settings *SettingsSession) SetWithdrawalCredentials(_value []byte) (*types.Transaction, error) {
	return _Settings.Contract.SetWithdrawalCredentials(&_Settings.TransactOpts, _value)
}

// SetWithdrawalCredentials is a paid mutator transaction binding the contract method 0xc4cf71c9.
//
// Solidity: function setWithdrawalCredentials(bytes _value) returns()
func (_Settings *SettingsTransactorSession) SetWithdrawalCredentials(_value []byte) (*types.Transaction, error) {
	return _Settings.Contract.SetWithdrawalCredentials(&_Settings.TransactOpts, _value)
}
