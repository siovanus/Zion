// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package side_chain_manager_abi

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

// ISideChainManagerSideChain is an auto generated low-level Go binding around an user-defined struct.
type ISideChainManagerSideChain struct {
	ChainID     uint64
	Router      uint64
	Name        string
	CCMCAddress []byte
	ExtraInfo   []byte
}

var (
	MethodRegisterAsset = "registerAsset"

	MethodUpdateFee = "updateFee"

	MethodGetFee = "getFee"

	MethodGetSideChain = "getSideChain"
)

// ISideChainManagerABI is the input ABI used to generate the binding from.
const ISideChainManagerABI = "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"getFee\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"getSideChain\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"router\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"CCMCAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"extraInfo\",\"type\":\"bytes\"}],\"internalType\":\"structISideChainManager.SideChain\",\"name\":\"sidechain\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64[]\",\"name\":\"AssetMapKey\",\"type\":\"uint64[]\"},{\"internalType\":\"bytes[]\",\"name\":\"AssetMapValue\",\"type\":\"bytes[]\"},{\"internalType\":\"uint64[]\",\"name\":\"LockProxyMapKey\",\"type\":\"uint64[]\"},{\"internalType\":\"bytes[]\",\"name\":\"LockProxyMapValue\",\"type\":\"bytes[]\"}],\"name\":\"registerAsset\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"viewNum\",\"type\":\"uint64\"},{\"internalType\":\"int256\",\"name\":\"fee\",\"type\":\"int256\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"updateFee\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ISideChainManagerFuncSigs maps the 4-byte function signature to its string representation.
var ISideChainManagerFuncSigs = map[string]string{
	"1982b1d0": "getFee(uint64)",
	"84838fb8": "getSideChain(uint64)",
	"e171240f": "registerAsset(uint64,uint64[],bytes[],uint64[],bytes[])",
	"db5d3488": "updateFee(uint64,uint64,int256,bytes)",
}

// ISideChainManager is an auto generated Go binding around an Ethereum contract.
type ISideChainManager struct {
	ISideChainManagerCaller     // Read-only binding to the contract
	ISideChainManagerTransactor // Write-only binding to the contract
	ISideChainManagerFilterer   // Log filterer for contract events
}

// ISideChainManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ISideChainManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISideChainManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ISideChainManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISideChainManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ISideChainManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ISideChainManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ISideChainManagerSession struct {
	Contract     *ISideChainManager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ISideChainManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ISideChainManagerCallerSession struct {
	Contract *ISideChainManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// ISideChainManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ISideChainManagerTransactorSession struct {
	Contract     *ISideChainManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// ISideChainManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ISideChainManagerRaw struct {
	Contract *ISideChainManager // Generic contract binding to access the raw methods on
}

// ISideChainManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ISideChainManagerCallerRaw struct {
	Contract *ISideChainManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ISideChainManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ISideChainManagerTransactorRaw struct {
	Contract *ISideChainManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewISideChainManager creates a new instance of ISideChainManager, bound to a specific deployed contract.
func NewISideChainManager(address common.Address, backend bind.ContractBackend) (*ISideChainManager, error) {
	contract, err := bindISideChainManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ISideChainManager{ISideChainManagerCaller: ISideChainManagerCaller{contract: contract}, ISideChainManagerTransactor: ISideChainManagerTransactor{contract: contract}, ISideChainManagerFilterer: ISideChainManagerFilterer{contract: contract}}, nil
}

// NewISideChainManagerCaller creates a new read-only instance of ISideChainManager, bound to a specific deployed contract.
func NewISideChainManagerCaller(address common.Address, caller bind.ContractCaller) (*ISideChainManagerCaller, error) {
	contract, err := bindISideChainManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerCaller{contract: contract}, nil
}

// NewISideChainManagerTransactor creates a new write-only instance of ISideChainManager, bound to a specific deployed contract.
func NewISideChainManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ISideChainManagerTransactor, error) {
	contract, err := bindISideChainManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerTransactor{contract: contract}, nil
}

// NewISideChainManagerFilterer creates a new log filterer instance of ISideChainManager, bound to a specific deployed contract.
func NewISideChainManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ISideChainManagerFilterer, error) {
	contract, err := bindISideChainManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ISideChainManagerFilterer{contract: contract}, nil
}

// bindISideChainManager binds a generic wrapper to an already deployed contract.
func bindISideChainManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ISideChainManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISideChainManager *ISideChainManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISideChainManager.Contract.ISideChainManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISideChainManager *ISideChainManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ISideChainManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISideChainManager *ISideChainManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISideChainManager.Contract.ISideChainManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ISideChainManager *ISideChainManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ISideChainManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ISideChainManager *ISideChainManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ISideChainManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ISideChainManager *ISideChainManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ISideChainManager.Contract.contract.Transact(opts, method, params...)
}

// GetFee is a free data retrieval call binding the contract method 0x1982b1d0.
//
// Solidity: function getFee(uint64 chainID) view returns(bytes)
func (_ISideChainManager *ISideChainManagerCaller) GetFee(opts *bind.CallOpts, chainID uint64) ([]byte, error) {
	var out []interface{}
	err := _ISideChainManager.contract.Call(opts, &out, "getFee", chainID)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetFee is a free data retrieval call binding the contract method 0x1982b1d0.
//
// Solidity: function getFee(uint64 chainID) view returns(bytes)
func (_ISideChainManager *ISideChainManagerSession) GetFee(chainID uint64) ([]byte, error) {
	return _ISideChainManager.Contract.GetFee(&_ISideChainManager.CallOpts, chainID)
}

// GetFee is a free data retrieval call binding the contract method 0x1982b1d0.
//
// Solidity: function getFee(uint64 chainID) view returns(bytes)
func (_ISideChainManager *ISideChainManagerCallerSession) GetFee(chainID uint64) ([]byte, error) {
	return _ISideChainManager.Contract.GetFee(&_ISideChainManager.CallOpts, chainID)
}

// GetSideChain is a free data retrieval call binding the contract method 0x84838fb8.
//
// Solidity: function getSideChain(uint64 chainID) view returns((uint64,uint64,string,bytes,bytes) sidechain)
func (_ISideChainManager *ISideChainManagerCaller) GetSideChain(opts *bind.CallOpts, chainID uint64) (ISideChainManagerSideChain, error) {
	var out []interface{}
	err := _ISideChainManager.contract.Call(opts, &out, "getSideChain", chainID)

	if err != nil {
		return *new(ISideChainManagerSideChain), err
	}

	out0 := *abi.ConvertType(out[0], new(ISideChainManagerSideChain)).(*ISideChainManagerSideChain)

	return out0, err

}

// GetSideChain is a free data retrieval call binding the contract method 0x84838fb8.
//
// Solidity: function getSideChain(uint64 chainID) view returns((uint64,uint64,string,bytes,bytes) sidechain)
func (_ISideChainManager *ISideChainManagerSession) GetSideChain(chainID uint64) (ISideChainManagerSideChain, error) {
	return _ISideChainManager.Contract.GetSideChain(&_ISideChainManager.CallOpts, chainID)
}

// GetSideChain is a free data retrieval call binding the contract method 0x84838fb8.
//
// Solidity: function getSideChain(uint64 chainID) view returns((uint64,uint64,string,bytes,bytes) sidechain)
func (_ISideChainManager *ISideChainManagerCallerSession) GetSideChain(chainID uint64) (ISideChainManagerSideChain, error) {
	return _ISideChainManager.Contract.GetSideChain(&_ISideChainManager.CallOpts, chainID)
}

// RegisterAsset is a paid mutator transaction binding the contract method 0xe171240f.
//
// Solidity: function registerAsset(uint64 chainID, uint64[] AssetMapKey, bytes[] AssetMapValue, uint64[] LockProxyMapKey, bytes[] LockProxyMapValue) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactor) RegisterAsset(opts *bind.TransactOpts, chainID uint64, AssetMapKey []uint64, AssetMapValue [][]byte, LockProxyMapKey []uint64, LockProxyMapValue [][]byte) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "registerAsset", chainID, AssetMapKey, AssetMapValue, LockProxyMapKey, LockProxyMapValue)
}

// RegisterAsset is a paid mutator transaction binding the contract method 0xe171240f.
//
// Solidity: function registerAsset(uint64 chainID, uint64[] AssetMapKey, bytes[] AssetMapValue, uint64[] LockProxyMapKey, bytes[] LockProxyMapValue) returns(bool success)
func (_ISideChainManager *ISideChainManagerSession) RegisterAsset(chainID uint64, AssetMapKey []uint64, AssetMapValue [][]byte, LockProxyMapKey []uint64, LockProxyMapValue [][]byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.RegisterAsset(&_ISideChainManager.TransactOpts, chainID, AssetMapKey, AssetMapValue, LockProxyMapKey, LockProxyMapValue)
}

// RegisterAsset is a paid mutator transaction binding the contract method 0xe171240f.
//
// Solidity: function registerAsset(uint64 chainID, uint64[] AssetMapKey, bytes[] AssetMapValue, uint64[] LockProxyMapKey, bytes[] LockProxyMapValue) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactorSession) RegisterAsset(chainID uint64, AssetMapKey []uint64, AssetMapValue [][]byte, LockProxyMapKey []uint64, LockProxyMapValue [][]byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.RegisterAsset(&_ISideChainManager.TransactOpts, chainID, AssetMapKey, AssetMapValue, LockProxyMapKey, LockProxyMapValue)
}

// UpdateFee is a paid mutator transaction binding the contract method 0xdb5d3488.
//
// Solidity: function updateFee(uint64 chainID, uint64 viewNum, int256 fee, bytes signature) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactor) UpdateFee(opts *bind.TransactOpts, chainID uint64, viewNum uint64, fee *big.Int, signature []byte) (*types.Transaction, error) {
	return _ISideChainManager.contract.Transact(opts, "updateFee", chainID, viewNum, fee, signature)
}

// UpdateFee is a paid mutator transaction binding the contract method 0xdb5d3488.
//
// Solidity: function updateFee(uint64 chainID, uint64 viewNum, int256 fee, bytes signature) returns(bool success)
func (_ISideChainManager *ISideChainManagerSession) UpdateFee(chainID uint64, viewNum uint64, fee *big.Int, signature []byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.UpdateFee(&_ISideChainManager.TransactOpts, chainID, viewNum, fee, signature)
}

// UpdateFee is a paid mutator transaction binding the contract method 0xdb5d3488.
//
// Solidity: function updateFee(uint64 chainID, uint64 viewNum, int256 fee, bytes signature) returns(bool success)
func (_ISideChainManager *ISideChainManagerTransactorSession) UpdateFee(chainID uint64, viewNum uint64, fee *big.Int, signature []byte) (*types.Transaction, error) {
	return _ISideChainManager.Contract.UpdateFee(&_ISideChainManager.TransactOpts, chainID, viewNum, fee, signature)
}

