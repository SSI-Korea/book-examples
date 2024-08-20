// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package didregistry

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
	_ = abi.ConvertType
)

// DidregistryMetaData contains all meta data concerning the Didregistry contract.
var DidregistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_did\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_document\",\"type\":\"string\"}],\"name\":\"RegisterDid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_did\",\"type\":\"string\"}],\"name\":\"ResolveDid\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DidregistryABI is the input ABI used to generate the binding from.
// Deprecated: Use DidregistryMetaData.ABI instead.
var DidregistryABI = DidregistryMetaData.ABI

// Didregistry is an auto generated Go binding around an Ethereum contract.
type Didregistry struct {
	DidregistryCaller     // Read-only binding to the contract
	DidregistryTransactor // Write-only binding to the contract
	DidregistryFilterer   // Log filterer for contract events
}

// DidregistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type DidregistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DidregistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DidregistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DidregistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DidregistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DidregistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DidregistrySession struct {
	Contract     *Didregistry      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DidregistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DidregistryCallerSession struct {
	Contract *DidregistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// DidregistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DidregistryTransactorSession struct {
	Contract     *DidregistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DidregistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type DidregistryRaw struct {
	Contract *Didregistry // Generic contract binding to access the raw methods on
}

// DidregistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DidregistryCallerRaw struct {
	Contract *DidregistryCaller // Generic read-only contract binding to access the raw methods on
}

// DidregistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DidregistryTransactorRaw struct {
	Contract *DidregistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDidregistry creates a new instance of Didregistry, bound to a specific deployed contract.
func NewDidregistry(address common.Address, backend bind.ContractBackend) (*Didregistry, error) {
	contract, err := bindDidregistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Didregistry{DidregistryCaller: DidregistryCaller{contract: contract}, DidregistryTransactor: DidregistryTransactor{contract: contract}, DidregistryFilterer: DidregistryFilterer{contract: contract}}, nil
}

// NewDidregistryCaller creates a new read-only instance of Didregistry, bound to a specific deployed contract.
func NewDidregistryCaller(address common.Address, caller bind.ContractCaller) (*DidregistryCaller, error) {
	contract, err := bindDidregistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DidregistryCaller{contract: contract}, nil
}

// NewDidregistryTransactor creates a new write-only instance of Didregistry, bound to a specific deployed contract.
func NewDidregistryTransactor(address common.Address, transactor bind.ContractTransactor) (*DidregistryTransactor, error) {
	contract, err := bindDidregistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DidregistryTransactor{contract: contract}, nil
}

// NewDidregistryFilterer creates a new log filterer instance of Didregistry, bound to a specific deployed contract.
func NewDidregistryFilterer(address common.Address, filterer bind.ContractFilterer) (*DidregistryFilterer, error) {
	contract, err := bindDidregistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DidregistryFilterer{contract: contract}, nil
}

// bindDidregistry binds a generic wrapper to an already deployed contract.
func bindDidregistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DidregistryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Didregistry *DidregistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Didregistry.Contract.DidregistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Didregistry *DidregistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Didregistry.Contract.DidregistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Didregistry *DidregistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Didregistry.Contract.DidregistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Didregistry *DidregistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Didregistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Didregistry *DidregistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Didregistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Didregistry *DidregistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Didregistry.Contract.contract.Transact(opts, method, params...)
}

// ResolveDid is a free data retrieval call binding the contract method 0x15f0f802.
//
// Solidity: function ResolveDid(string _did) view returns(string)
func (_Didregistry *DidregistryCaller) ResolveDid(opts *bind.CallOpts, _did string) (string, error) {
	var out []interface{}
	err := _Didregistry.contract.Call(opts, &out, "ResolveDid", _did)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ResolveDid is a free data retrieval call binding the contract method 0x15f0f802.
//
// Solidity: function ResolveDid(string _did) view returns(string)
func (_Didregistry *DidregistrySession) ResolveDid(_did string) (string, error) {
	return _Didregistry.Contract.ResolveDid(&_Didregistry.CallOpts, _did)
}

// ResolveDid is a free data retrieval call binding the contract method 0x15f0f802.
//
// Solidity: function ResolveDid(string _did) view returns(string)
func (_Didregistry *DidregistryCallerSession) ResolveDid(_did string) (string, error) {
	return _Didregistry.Contract.ResolveDid(&_Didregistry.CallOpts, _did)
}

// RegisterDid is a paid mutator transaction binding the contract method 0x841e2cc7.
//
// Solidity: function RegisterDid(string _did, string _document) returns()
func (_Didregistry *DidregistryTransactor) RegisterDid(opts *bind.TransactOpts, _did string, _document string) (*types.Transaction, error) {
	return _Didregistry.contract.Transact(opts, "RegisterDid", _did, _document)
}

// RegisterDid is a paid mutator transaction binding the contract method 0x841e2cc7.
//
// Solidity: function RegisterDid(string _did, string _document) returns()
func (_Didregistry *DidregistrySession) RegisterDid(_did string, _document string) (*types.Transaction, error) {
	return _Didregistry.Contract.RegisterDid(&_Didregistry.TransactOpts, _did, _document)
}

// RegisterDid is a paid mutator transaction binding the contract method 0x841e2cc7.
//
// Solidity: function RegisterDid(string _did, string _document) returns()
func (_Didregistry *DidregistryTransactorSession) RegisterDid(_did string, _document string) (*types.Transaction, error) {
	return _Didregistry.Contract.RegisterDid(&_Didregistry.TransactOpts, _did, _document)
}
