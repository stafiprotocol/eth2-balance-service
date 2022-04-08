// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PoolManager

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

// PoolManagerMetaData contains all meta data concerning the PoolManager contract.
var PoolManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getStakingPoolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getStakingPoolAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"}],\"name\":\"getNodeStakingPoolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getNodeStakingPoolAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"}],\"name\":\"getNodeValidatingStakingPoolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getNodeValidatingStakingPoolAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_pubkey\",\"type\":\"bytes\"}],\"name\":\"getStakingPoolByPubkey\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingPoolAddress\",\"type\":\"address\"}],\"name\":\"getStakingPoolExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingPoolAddress\",\"type\":\"address\"}],\"name\":\"getStakingPoolPubkey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingPoolAddress\",\"type\":\"address\"}],\"name\":\"getStakingPoolWithdrawalProcessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getPrelaunchStakingpools\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"enumDepositType\",\"name\":\"_depositType\",\"type\":\"uint8\"}],\"name\":\"createStakingPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"destroyStakingPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_pubkey\",\"type\":\"bytes\"}],\"name\":\"setStakingPoolPubkey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingPoolAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_processed\",\"type\":\"bool\"}],\"name\":\"setStakingPoolWithdrawalProcessed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// PoolManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use PoolManagerMetaData.ABI instead.
var PoolManagerABI = PoolManagerMetaData.ABI

// PoolManager is an auto generated Go binding around an Ethereum contract.
type PoolManager struct {
	PoolManagerCaller     // Read-only binding to the contract
	PoolManagerTransactor // Write-only binding to the contract
	PoolManagerFilterer   // Log filterer for contract events
}

// PoolManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolManagerSession struct {
	Contract     *PoolManager      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolManagerCallerSession struct {
	Contract *PoolManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PoolManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolManagerTransactorSession struct {
	Contract     *PoolManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PoolManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolManagerRaw struct {
	Contract *PoolManager // Generic contract binding to access the raw methods on
}

// PoolManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolManagerCallerRaw struct {
	Contract *PoolManagerCaller // Generic read-only contract binding to access the raw methods on
}

// PoolManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolManagerTransactorRaw struct {
	Contract *PoolManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoolManager creates a new instance of PoolManager, bound to a specific deployed contract.
func NewPoolManager(address common.Address, backend bind.ContractBackend) (*PoolManager, error) {
	contract, err := bindPoolManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PoolManager{PoolManagerCaller: PoolManagerCaller{contract: contract}, PoolManagerTransactor: PoolManagerTransactor{contract: contract}, PoolManagerFilterer: PoolManagerFilterer{contract: contract}}, nil
}

// NewPoolManagerCaller creates a new read-only instance of PoolManager, bound to a specific deployed contract.
func NewPoolManagerCaller(address common.Address, caller bind.ContractCaller) (*PoolManagerCaller, error) {
	contract, err := bindPoolManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolManagerCaller{contract: contract}, nil
}

// NewPoolManagerTransactor creates a new write-only instance of PoolManager, bound to a specific deployed contract.
func NewPoolManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolManagerTransactor, error) {
	contract, err := bindPoolManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolManagerTransactor{contract: contract}, nil
}

// NewPoolManagerFilterer creates a new log filterer instance of PoolManager, bound to a specific deployed contract.
func NewPoolManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolManagerFilterer, error) {
	contract, err := bindPoolManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolManagerFilterer{contract: contract}, nil
}

// bindPoolManager binds a generic wrapper to an already deployed contract.
func bindPoolManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PoolManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolManager *PoolManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolManager.Contract.PoolManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolManager *PoolManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolManager.Contract.PoolManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolManager *PoolManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolManager.Contract.PoolManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolManager *PoolManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolManager *PoolManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolManager *PoolManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolManager.Contract.contract.Transact(opts, method, params...)
}

// GetNodeStakingPoolAt is a free data retrieval call binding the contract method 0xc0f27c38.
//
// Solidity: function getNodeStakingPoolAt(address _nodeAddress, uint256 _index) view returns(address)
func (_PoolManager *PoolManagerCaller) GetNodeStakingPoolAt(opts *bind.CallOpts, _nodeAddress common.Address, _index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getNodeStakingPoolAt", _nodeAddress, _index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetNodeStakingPoolAt is a free data retrieval call binding the contract method 0xc0f27c38.
//
// Solidity: function getNodeStakingPoolAt(address _nodeAddress, uint256 _index) view returns(address)
func (_PoolManager *PoolManagerSession) GetNodeStakingPoolAt(_nodeAddress common.Address, _index *big.Int) (common.Address, error) {
	return _PoolManager.Contract.GetNodeStakingPoolAt(&_PoolManager.CallOpts, _nodeAddress, _index)
}

// GetNodeStakingPoolAt is a free data retrieval call binding the contract method 0xc0f27c38.
//
// Solidity: function getNodeStakingPoolAt(address _nodeAddress, uint256 _index) view returns(address)
func (_PoolManager *PoolManagerCallerSession) GetNodeStakingPoolAt(_nodeAddress common.Address, _index *big.Int) (common.Address, error) {
	return _PoolManager.Contract.GetNodeStakingPoolAt(&_PoolManager.CallOpts, _nodeAddress, _index)
}

// GetNodeStakingPoolCount is a free data retrieval call binding the contract method 0xc3d8073b.
//
// Solidity: function getNodeStakingPoolCount(address _nodeAddress) view returns(uint256)
func (_PoolManager *PoolManagerCaller) GetNodeStakingPoolCount(opts *bind.CallOpts, _nodeAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getNodeStakingPoolCount", _nodeAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeStakingPoolCount is a free data retrieval call binding the contract method 0xc3d8073b.
//
// Solidity: function getNodeStakingPoolCount(address _nodeAddress) view returns(uint256)
func (_PoolManager *PoolManagerSession) GetNodeStakingPoolCount(_nodeAddress common.Address) (*big.Int, error) {
	return _PoolManager.Contract.GetNodeStakingPoolCount(&_PoolManager.CallOpts, _nodeAddress)
}

// GetNodeStakingPoolCount is a free data retrieval call binding the contract method 0xc3d8073b.
//
// Solidity: function getNodeStakingPoolCount(address _nodeAddress) view returns(uint256)
func (_PoolManager *PoolManagerCallerSession) GetNodeStakingPoolCount(_nodeAddress common.Address) (*big.Int, error) {
	return _PoolManager.Contract.GetNodeStakingPoolCount(&_PoolManager.CallOpts, _nodeAddress)
}

// GetNodeValidatingStakingPoolAt is a free data retrieval call binding the contract method 0x0fed8b63.
//
// Solidity: function getNodeValidatingStakingPoolAt(address _nodeAddress, uint256 _index) view returns(address)
func (_PoolManager *PoolManagerCaller) GetNodeValidatingStakingPoolAt(opts *bind.CallOpts, _nodeAddress common.Address, _index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getNodeValidatingStakingPoolAt", _nodeAddress, _index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetNodeValidatingStakingPoolAt is a free data retrieval call binding the contract method 0x0fed8b63.
//
// Solidity: function getNodeValidatingStakingPoolAt(address _nodeAddress, uint256 _index) view returns(address)
func (_PoolManager *PoolManagerSession) GetNodeValidatingStakingPoolAt(_nodeAddress common.Address, _index *big.Int) (common.Address, error) {
	return _PoolManager.Contract.GetNodeValidatingStakingPoolAt(&_PoolManager.CallOpts, _nodeAddress, _index)
}

// GetNodeValidatingStakingPoolAt is a free data retrieval call binding the contract method 0x0fed8b63.
//
// Solidity: function getNodeValidatingStakingPoolAt(address _nodeAddress, uint256 _index) view returns(address)
func (_PoolManager *PoolManagerCallerSession) GetNodeValidatingStakingPoolAt(_nodeAddress common.Address, _index *big.Int) (common.Address, error) {
	return _PoolManager.Contract.GetNodeValidatingStakingPoolAt(&_PoolManager.CallOpts, _nodeAddress, _index)
}

// GetNodeValidatingStakingPoolCount is a free data retrieval call binding the contract method 0x7be535ad.
//
// Solidity: function getNodeValidatingStakingPoolCount(address _nodeAddress) view returns(uint256)
func (_PoolManager *PoolManagerCaller) GetNodeValidatingStakingPoolCount(opts *bind.CallOpts, _nodeAddress common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getNodeValidatingStakingPoolCount", _nodeAddress)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetNodeValidatingStakingPoolCount is a free data retrieval call binding the contract method 0x7be535ad.
//
// Solidity: function getNodeValidatingStakingPoolCount(address _nodeAddress) view returns(uint256)
func (_PoolManager *PoolManagerSession) GetNodeValidatingStakingPoolCount(_nodeAddress common.Address) (*big.Int, error) {
	return _PoolManager.Contract.GetNodeValidatingStakingPoolCount(&_PoolManager.CallOpts, _nodeAddress)
}

// GetNodeValidatingStakingPoolCount is a free data retrieval call binding the contract method 0x7be535ad.
//
// Solidity: function getNodeValidatingStakingPoolCount(address _nodeAddress) view returns(uint256)
func (_PoolManager *PoolManagerCallerSession) GetNodeValidatingStakingPoolCount(_nodeAddress common.Address) (*big.Int, error) {
	return _PoolManager.Contract.GetNodeValidatingStakingPoolCount(&_PoolManager.CallOpts, _nodeAddress)
}

// GetPrelaunchStakingpools is a free data retrieval call binding the contract method 0xcc677617.
//
// Solidity: function getPrelaunchStakingpools(uint256 offset, uint256 limit) view returns(address[])
func (_PoolManager *PoolManagerCaller) GetPrelaunchStakingpools(opts *bind.CallOpts, offset *big.Int, limit *big.Int) ([]common.Address, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getPrelaunchStakingpools", offset, limit)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetPrelaunchStakingpools is a free data retrieval call binding the contract method 0xcc677617.
//
// Solidity: function getPrelaunchStakingpools(uint256 offset, uint256 limit) view returns(address[])
func (_PoolManager *PoolManagerSession) GetPrelaunchStakingpools(offset *big.Int, limit *big.Int) ([]common.Address, error) {
	return _PoolManager.Contract.GetPrelaunchStakingpools(&_PoolManager.CallOpts, offset, limit)
}

// GetPrelaunchStakingpools is a free data retrieval call binding the contract method 0xcc677617.
//
// Solidity: function getPrelaunchStakingpools(uint256 offset, uint256 limit) view returns(address[])
func (_PoolManager *PoolManagerCallerSession) GetPrelaunchStakingpools(offset *big.Int, limit *big.Int) ([]common.Address, error) {
	return _PoolManager.Contract.GetPrelaunchStakingpools(&_PoolManager.CallOpts, offset, limit)
}

// GetStakingPoolAt is a free data retrieval call binding the contract method 0x26a71b13.
//
// Solidity: function getStakingPoolAt(uint256 _index) view returns(address)
func (_PoolManager *PoolManagerCaller) GetStakingPoolAt(opts *bind.CallOpts, _index *big.Int) (common.Address, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getStakingPoolAt", _index)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetStakingPoolAt is a free data retrieval call binding the contract method 0x26a71b13.
//
// Solidity: function getStakingPoolAt(uint256 _index) view returns(address)
func (_PoolManager *PoolManagerSession) GetStakingPoolAt(_index *big.Int) (common.Address, error) {
	return _PoolManager.Contract.GetStakingPoolAt(&_PoolManager.CallOpts, _index)
}

// GetStakingPoolAt is a free data retrieval call binding the contract method 0x26a71b13.
//
// Solidity: function getStakingPoolAt(uint256 _index) view returns(address)
func (_PoolManager *PoolManagerCallerSession) GetStakingPoolAt(_index *big.Int) (common.Address, error) {
	return _PoolManager.Contract.GetStakingPoolAt(&_PoolManager.CallOpts, _index)
}

// GetStakingPoolByPubkey is a free data retrieval call binding the contract method 0x02ec7481.
//
// Solidity: function getStakingPoolByPubkey(bytes _pubkey) view returns(address)
func (_PoolManager *PoolManagerCaller) GetStakingPoolByPubkey(opts *bind.CallOpts, _pubkey []byte) (common.Address, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getStakingPoolByPubkey", _pubkey)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetStakingPoolByPubkey is a free data retrieval call binding the contract method 0x02ec7481.
//
// Solidity: function getStakingPoolByPubkey(bytes _pubkey) view returns(address)
func (_PoolManager *PoolManagerSession) GetStakingPoolByPubkey(_pubkey []byte) (common.Address, error) {
	return _PoolManager.Contract.GetStakingPoolByPubkey(&_PoolManager.CallOpts, _pubkey)
}

// GetStakingPoolByPubkey is a free data retrieval call binding the contract method 0x02ec7481.
//
// Solidity: function getStakingPoolByPubkey(bytes _pubkey) view returns(address)
func (_PoolManager *PoolManagerCallerSession) GetStakingPoolByPubkey(_pubkey []byte) (common.Address, error) {
	return _PoolManager.Contract.GetStakingPoolByPubkey(&_PoolManager.CallOpts, _pubkey)
}

// GetStakingPoolCount is a free data retrieval call binding the contract method 0x9bf00f18.
//
// Solidity: function getStakingPoolCount() view returns(uint256)
func (_PoolManager *PoolManagerCaller) GetStakingPoolCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getStakingPoolCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetStakingPoolCount is a free data retrieval call binding the contract method 0x9bf00f18.
//
// Solidity: function getStakingPoolCount() view returns(uint256)
func (_PoolManager *PoolManagerSession) GetStakingPoolCount() (*big.Int, error) {
	return _PoolManager.Contract.GetStakingPoolCount(&_PoolManager.CallOpts)
}

// GetStakingPoolCount is a free data retrieval call binding the contract method 0x9bf00f18.
//
// Solidity: function getStakingPoolCount() view returns(uint256)
func (_PoolManager *PoolManagerCallerSession) GetStakingPoolCount() (*big.Int, error) {
	return _PoolManager.Contract.GetStakingPoolCount(&_PoolManager.CallOpts)
}

// GetStakingPoolExists is a free data retrieval call binding the contract method 0x506419ca.
//
// Solidity: function getStakingPoolExists(address _stakingPoolAddress) view returns(bool)
func (_PoolManager *PoolManagerCaller) GetStakingPoolExists(opts *bind.CallOpts, _stakingPoolAddress common.Address) (bool, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getStakingPoolExists", _stakingPoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetStakingPoolExists is a free data retrieval call binding the contract method 0x506419ca.
//
// Solidity: function getStakingPoolExists(address _stakingPoolAddress) view returns(bool)
func (_PoolManager *PoolManagerSession) GetStakingPoolExists(_stakingPoolAddress common.Address) (bool, error) {
	return _PoolManager.Contract.GetStakingPoolExists(&_PoolManager.CallOpts, _stakingPoolAddress)
}

// GetStakingPoolExists is a free data retrieval call binding the contract method 0x506419ca.
//
// Solidity: function getStakingPoolExists(address _stakingPoolAddress) view returns(bool)
func (_PoolManager *PoolManagerCallerSession) GetStakingPoolExists(_stakingPoolAddress common.Address) (bool, error) {
	return _PoolManager.Contract.GetStakingPoolExists(&_PoolManager.CallOpts, _stakingPoolAddress)
}

// GetStakingPoolPubkey is a free data retrieval call binding the contract method 0x0d4354c6.
//
// Solidity: function getStakingPoolPubkey(address _stakingPoolAddress) view returns(bytes)
func (_PoolManager *PoolManagerCaller) GetStakingPoolPubkey(opts *bind.CallOpts, _stakingPoolAddress common.Address) ([]byte, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getStakingPoolPubkey", _stakingPoolAddress)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetStakingPoolPubkey is a free data retrieval call binding the contract method 0x0d4354c6.
//
// Solidity: function getStakingPoolPubkey(address _stakingPoolAddress) view returns(bytes)
func (_PoolManager *PoolManagerSession) GetStakingPoolPubkey(_stakingPoolAddress common.Address) ([]byte, error) {
	return _PoolManager.Contract.GetStakingPoolPubkey(&_PoolManager.CallOpts, _stakingPoolAddress)
}

// GetStakingPoolPubkey is a free data retrieval call binding the contract method 0x0d4354c6.
//
// Solidity: function getStakingPoolPubkey(address _stakingPoolAddress) view returns(bytes)
func (_PoolManager *PoolManagerCallerSession) GetStakingPoolPubkey(_stakingPoolAddress common.Address) ([]byte, error) {
	return _PoolManager.Contract.GetStakingPoolPubkey(&_PoolManager.CallOpts, _stakingPoolAddress)
}

// GetStakingPoolWithdrawalProcessed is a free data retrieval call binding the contract method 0xf3e0ce85.
//
// Solidity: function getStakingPoolWithdrawalProcessed(address _stakingPoolAddress) view returns(bool)
func (_PoolManager *PoolManagerCaller) GetStakingPoolWithdrawalProcessed(opts *bind.CallOpts, _stakingPoolAddress common.Address) (bool, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "getStakingPoolWithdrawalProcessed", _stakingPoolAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// GetStakingPoolWithdrawalProcessed is a free data retrieval call binding the contract method 0xf3e0ce85.
//
// Solidity: function getStakingPoolWithdrawalProcessed(address _stakingPoolAddress) view returns(bool)
func (_PoolManager *PoolManagerSession) GetStakingPoolWithdrawalProcessed(_stakingPoolAddress common.Address) (bool, error) {
	return _PoolManager.Contract.GetStakingPoolWithdrawalProcessed(&_PoolManager.CallOpts, _stakingPoolAddress)
}

// GetStakingPoolWithdrawalProcessed is a free data retrieval call binding the contract method 0xf3e0ce85.
//
// Solidity: function getStakingPoolWithdrawalProcessed(address _stakingPoolAddress) view returns(bool)
func (_PoolManager *PoolManagerCallerSession) GetStakingPoolWithdrawalProcessed(_stakingPoolAddress common.Address) (bool, error) {
	return _PoolManager.Contract.GetStakingPoolWithdrawalProcessed(&_PoolManager.CallOpts, _stakingPoolAddress)
}

// CreateStakingPool is a paid mutator transaction binding the contract method 0x6e5da826.
//
// Solidity: function createStakingPool(address _nodeAddress, uint8 _depositType) returns(address)
func (_PoolManager *PoolManagerTransactor) CreateStakingPool(opts *bind.TransactOpts, _nodeAddress common.Address, _depositType uint8) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "createStakingPool", _nodeAddress, _depositType)
}

// CreateStakingPool is a paid mutator transaction binding the contract method 0x6e5da826.
//
// Solidity: function createStakingPool(address _nodeAddress, uint8 _depositType) returns(address)
func (_PoolManager *PoolManagerSession) CreateStakingPool(_nodeAddress common.Address, _depositType uint8) (*types.Transaction, error) {
	return _PoolManager.Contract.CreateStakingPool(&_PoolManager.TransactOpts, _nodeAddress, _depositType)
}

// CreateStakingPool is a paid mutator transaction binding the contract method 0x6e5da826.
//
// Solidity: function createStakingPool(address _nodeAddress, uint8 _depositType) returns(address)
func (_PoolManager *PoolManagerTransactorSession) CreateStakingPool(_nodeAddress common.Address, _depositType uint8) (*types.Transaction, error) {
	return _PoolManager.Contract.CreateStakingPool(&_PoolManager.TransactOpts, _nodeAddress, _depositType)
}

// DestroyStakingPool is a paid mutator transaction binding the contract method 0x6b16c49d.
//
// Solidity: function destroyStakingPool() returns()
func (_PoolManager *PoolManagerTransactor) DestroyStakingPool(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "destroyStakingPool")
}

// DestroyStakingPool is a paid mutator transaction binding the contract method 0x6b16c49d.
//
// Solidity: function destroyStakingPool() returns()
func (_PoolManager *PoolManagerSession) DestroyStakingPool() (*types.Transaction, error) {
	return _PoolManager.Contract.DestroyStakingPool(&_PoolManager.TransactOpts)
}

// DestroyStakingPool is a paid mutator transaction binding the contract method 0x6b16c49d.
//
// Solidity: function destroyStakingPool() returns()
func (_PoolManager *PoolManagerTransactorSession) DestroyStakingPool() (*types.Transaction, error) {
	return _PoolManager.Contract.DestroyStakingPool(&_PoolManager.TransactOpts)
}

// SetStakingPoolPubkey is a paid mutator transaction binding the contract method 0xaf40d55d.
//
// Solidity: function setStakingPoolPubkey(bytes _pubkey) returns()
func (_PoolManager *PoolManagerTransactor) SetStakingPoolPubkey(opts *bind.TransactOpts, _pubkey []byte) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "setStakingPoolPubkey", _pubkey)
}

// SetStakingPoolPubkey is a paid mutator transaction binding the contract method 0xaf40d55d.
//
// Solidity: function setStakingPoolPubkey(bytes _pubkey) returns()
func (_PoolManager *PoolManagerSession) SetStakingPoolPubkey(_pubkey []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.SetStakingPoolPubkey(&_PoolManager.TransactOpts, _pubkey)
}

// SetStakingPoolPubkey is a paid mutator transaction binding the contract method 0xaf40d55d.
//
// Solidity: function setStakingPoolPubkey(bytes _pubkey) returns()
func (_PoolManager *PoolManagerTransactorSession) SetStakingPoolPubkey(_pubkey []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.SetStakingPoolPubkey(&_PoolManager.TransactOpts, _pubkey)
}

// SetStakingPoolWithdrawalProcessed is a paid mutator transaction binding the contract method 0x1189504e.
//
// Solidity: function setStakingPoolWithdrawalProcessed(address _stakingPoolAddress, bool _processed) returns()
func (_PoolManager *PoolManagerTransactor) SetStakingPoolWithdrawalProcessed(opts *bind.TransactOpts, _stakingPoolAddress common.Address, _processed bool) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "setStakingPoolWithdrawalProcessed", _stakingPoolAddress, _processed)
}

// SetStakingPoolWithdrawalProcessed is a paid mutator transaction binding the contract method 0x1189504e.
//
// Solidity: function setStakingPoolWithdrawalProcessed(address _stakingPoolAddress, bool _processed) returns()
func (_PoolManager *PoolManagerSession) SetStakingPoolWithdrawalProcessed(_stakingPoolAddress common.Address, _processed bool) (*types.Transaction, error) {
	return _PoolManager.Contract.SetStakingPoolWithdrawalProcessed(&_PoolManager.TransactOpts, _stakingPoolAddress, _processed)
}

// SetStakingPoolWithdrawalProcessed is a paid mutator transaction binding the contract method 0x1189504e.
//
// Solidity: function setStakingPoolWithdrawalProcessed(address _stakingPoolAddress, bool _processed) returns()
func (_PoolManager *PoolManagerTransactorSession) SetStakingPoolWithdrawalProcessed(_stakingPoolAddress common.Address, _processed bool) (*types.Transaction, error) {
	return _PoolManager.Contract.SetStakingPoolWithdrawalProcessed(&_PoolManager.TransactOpts, _stakingPoolAddress, _processed)
}
