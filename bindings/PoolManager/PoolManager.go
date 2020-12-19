// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package PoolManager

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

// PoolManagerABI is the input ABI used to generate the binding from.
const PoolManagerABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stafiStorageAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"stakingPool\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"StakingPoolCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"stakingPool\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"node\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"time\",\"type\":\"uint256\"}],\"name\":\"StakingPoolDestroyed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"enumDepositType\",\"name\":\"_depositType\",\"type\":\"uint8\"}],\"name\":\"createStakingPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"destroyStakingPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getNodeStakingPoolAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"}],\"name\":\"getNodeStakingPoolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getNodeValidatingStakingPoolAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_nodeAddress\",\"type\":\"address\"}],\"name\":\"getNodeValidatingStakingPoolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getStakingPoolAt\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_pubkey\",\"type\":\"bytes\"}],\"name\":\"getStakingPoolByPubkey\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getStakingPoolCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingPoolAddress\",\"type\":\"address\"}],\"name\":\"getStakingPoolExists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingPoolAddress\",\"type\":\"address\"}],\"name\":\"getStakingPoolPubkey\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingPoolAddress\",\"type\":\"address\"}],\"name\":\"getStakingPoolWithdrawalProcessed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_pubkey\",\"type\":\"bytes\"}],\"name\":\"setStakingPoolPubkey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_stakingPoolAddress\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_processed\",\"type\":\"bool\"}],\"name\":\"setStakingPoolWithdrawalProcessed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

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

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_PoolManager *PoolManagerCaller) Version(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "version")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_PoolManager *PoolManagerSession) Version() (uint8, error) {
	return _PoolManager.Contract.Version(&_PoolManager.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() view returns(uint8)
func (_PoolManager *PoolManagerCallerSession) Version() (uint8, error) {
	return _PoolManager.Contract.Version(&_PoolManager.CallOpts)
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

// PoolManagerStakingPoolCreatedIterator is returned from FilterStakingPoolCreated and is used to iterate over the raw logs and unpacked data for StakingPoolCreated events raised by the PoolManager contract.
type PoolManagerStakingPoolCreatedIterator struct {
	Event *PoolManagerStakingPoolCreated // Event containing the contract specifics and raw log

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
func (it *PoolManagerStakingPoolCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerStakingPoolCreated)
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
		it.Event = new(PoolManagerStakingPoolCreated)
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
func (it *PoolManagerStakingPoolCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerStakingPoolCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerStakingPoolCreated represents a StakingPoolCreated event raised by the PoolManager contract.
type PoolManagerStakingPoolCreated struct {
	StakingPool common.Address
	Node        common.Address
	Time        *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStakingPoolCreated is a free log retrieval operation binding the contract event 0xdfba889e07ceb4b33e759a1015aa4ae7b45c3881d9cf96dbf09b8971a27fc4b7.
//
// Solidity: event StakingPoolCreated(address indexed stakingPool, address indexed node, uint256 time)
func (_PoolManager *PoolManagerFilterer) FilterStakingPoolCreated(opts *bind.FilterOpts, stakingPool []common.Address, node []common.Address) (*PoolManagerStakingPoolCreatedIterator, error) {

	var stakingPoolRule []interface{}
	for _, stakingPoolItem := range stakingPool {
		stakingPoolRule = append(stakingPoolRule, stakingPoolItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "StakingPoolCreated", stakingPoolRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerStakingPoolCreatedIterator{contract: _PoolManager.contract, event: "StakingPoolCreated", logs: logs, sub: sub}, nil
}

// WatchStakingPoolCreated is a free log subscription operation binding the contract event 0xdfba889e07ceb4b33e759a1015aa4ae7b45c3881d9cf96dbf09b8971a27fc4b7.
//
// Solidity: event StakingPoolCreated(address indexed stakingPool, address indexed node, uint256 time)
func (_PoolManager *PoolManagerFilterer) WatchStakingPoolCreated(opts *bind.WatchOpts, sink chan<- *PoolManagerStakingPoolCreated, stakingPool []common.Address, node []common.Address) (event.Subscription, error) {

	var stakingPoolRule []interface{}
	for _, stakingPoolItem := range stakingPool {
		stakingPoolRule = append(stakingPoolRule, stakingPoolItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "StakingPoolCreated", stakingPoolRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerStakingPoolCreated)
				if err := _PoolManager.contract.UnpackLog(event, "StakingPoolCreated", log); err != nil {
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

// ParseStakingPoolCreated is a log parse operation binding the contract event 0xdfba889e07ceb4b33e759a1015aa4ae7b45c3881d9cf96dbf09b8971a27fc4b7.
//
// Solidity: event StakingPoolCreated(address indexed stakingPool, address indexed node, uint256 time)
func (_PoolManager *PoolManagerFilterer) ParseStakingPoolCreated(log types.Log) (*PoolManagerStakingPoolCreated, error) {
	event := new(PoolManagerStakingPoolCreated)
	if err := _PoolManager.contract.UnpackLog(event, "StakingPoolCreated", log); err != nil {
		return nil, err
	}
	return event, nil
}

// PoolManagerStakingPoolDestroyedIterator is returned from FilterStakingPoolDestroyed and is used to iterate over the raw logs and unpacked data for StakingPoolDestroyed events raised by the PoolManager contract.
type PoolManagerStakingPoolDestroyedIterator struct {
	Event *PoolManagerStakingPoolDestroyed // Event containing the contract specifics and raw log

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
func (it *PoolManagerStakingPoolDestroyedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerStakingPoolDestroyed)
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
		it.Event = new(PoolManagerStakingPoolDestroyed)
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
func (it *PoolManagerStakingPoolDestroyedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerStakingPoolDestroyedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerStakingPoolDestroyed represents a StakingPoolDestroyed event raised by the PoolManager contract.
type PoolManagerStakingPoolDestroyed struct {
	StakingPool common.Address
	Node        common.Address
	Time        *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStakingPoolDestroyed is a free log retrieval operation binding the contract event 0x458a7280951da092252cc4009cfcf959d019fef8696665c814ebff7a8d4750ea.
//
// Solidity: event StakingPoolDestroyed(address indexed stakingPool, address indexed node, uint256 time)
func (_PoolManager *PoolManagerFilterer) FilterStakingPoolDestroyed(opts *bind.FilterOpts, stakingPool []common.Address, node []common.Address) (*PoolManagerStakingPoolDestroyedIterator, error) {

	var stakingPoolRule []interface{}
	for _, stakingPoolItem := range stakingPool {
		stakingPoolRule = append(stakingPoolRule, stakingPoolItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "StakingPoolDestroyed", stakingPoolRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerStakingPoolDestroyedIterator{contract: _PoolManager.contract, event: "StakingPoolDestroyed", logs: logs, sub: sub}, nil
}

// WatchStakingPoolDestroyed is a free log subscription operation binding the contract event 0x458a7280951da092252cc4009cfcf959d019fef8696665c814ebff7a8d4750ea.
//
// Solidity: event StakingPoolDestroyed(address indexed stakingPool, address indexed node, uint256 time)
func (_PoolManager *PoolManagerFilterer) WatchStakingPoolDestroyed(opts *bind.WatchOpts, sink chan<- *PoolManagerStakingPoolDestroyed, stakingPool []common.Address, node []common.Address) (event.Subscription, error) {

	var stakingPoolRule []interface{}
	for _, stakingPoolItem := range stakingPool {
		stakingPoolRule = append(stakingPoolRule, stakingPoolItem)
	}
	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "StakingPoolDestroyed", stakingPoolRule, nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerStakingPoolDestroyed)
				if err := _PoolManager.contract.UnpackLog(event, "StakingPoolDestroyed", log); err != nil {
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

// ParseStakingPoolDestroyed is a log parse operation binding the contract event 0x458a7280951da092252cc4009cfcf959d019fef8696665c814ebff7a8d4750ea.
//
// Solidity: event StakingPoolDestroyed(address indexed stakingPool, address indexed node, uint256 time)
func (_PoolManager *PoolManagerFilterer) ParseStakingPoolDestroyed(log types.Log) (*PoolManagerStakingPoolDestroyed, error) {
	event := new(PoolManagerStakingPoolDestroyed)
	if err := _PoolManager.contract.UnpackLog(event, "StakingPoolDestroyed", log); err != nil {
		return nil, err
	}
	return event, nil
}
