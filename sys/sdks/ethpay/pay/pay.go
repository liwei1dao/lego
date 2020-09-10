// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package pays

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

// AccountABI is the input ABI used to generate the binding from.
const AccountABI = "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_reciever\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Flush\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Recharge\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"flush\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reciever\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AccountFuncSigs maps the 4-byte function signature to its string representation.
var AccountFuncSigs = map[string]string{
	"6b9f96ea": "flush()",
	"f4b0b756": "reciever()",
}

// AccountBin is the compiled bytecode used for deploying new contracts.
var AccountBin = "0x608060405234801561001057600080fd5b506040516101e13803806101e18339818101604052602081101561003357600080fd5b5051600080546001600160a01b039092166001600160a01b031990921691909117905561017c806100656000396000f3fe6080604052600436106100295760003560e01c80636b9f96ea1461005e578063f4b0b75614610075575b6040805134815290517f3b47a73ca305dcdfb3a69bdd0305a75776fbe85525d92cf22644e47e16eb5c1b9181900360200190a1005b34801561006a57600080fd5b506100736100a6565b005b34801561008157600080fd5b5061008a610137565b604080516001600160a01b039092168252519081900360200190f35b47806100b25750610135565b600080546040516001600160a01b039091169183156108fc02918491818181858888f193505050501580156100eb573d6000803e3d6000fd5b50600054604080516001600160a01b0390921682526020820183905280517f12b2a0ee977e74c33898f8be30fde7ae3a32ac7409a3666da55ce77e9bc32e879281900390910190a1505b565b6000546001600160a01b03168156fea264697066735822122066d9e560c66006a9376726110cda2f35329700f29dd9ae7cdba69191bb56b40864736f6c63430007000033"

// DeployAccount deploys a new Ethereum contract, binding an instance of Account to it.
func DeployAccount(auth *bind.TransactOpts, backend bind.ContractBackend, _reciever common.Address) (common.Address, *types.Transaction, *Account, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccountBin), backend, _reciever)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Account{AccountCaller: AccountCaller{contract: contract}, AccountTransactor: AccountTransactor{contract: contract}, AccountFilterer: AccountFilterer{contract: contract}}, nil
}

// Account is an auto generated Go binding around an Ethereum contract.
type Account struct {
	AccountCaller     // Read-only binding to the contract
	AccountTransactor // Write-only binding to the contract
	AccountFilterer   // Log filterer for contract events
}

// AccountCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccountCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccountTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccountFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccountSession struct {
	Contract     *Account          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccountCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccountCallerSession struct {
	Contract *AccountCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// AccountTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccountTransactorSession struct {
	Contract     *AccountTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// AccountRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccountRaw struct {
	Contract *Account // Generic contract binding to access the raw methods on
}

// AccountCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccountCallerRaw struct {
	Contract *AccountCaller // Generic read-only contract binding to access the raw methods on
}

// AccountTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccountTransactorRaw struct {
	Contract *AccountTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccount creates a new instance of Account, bound to a specific deployed contract.
func NewAccount(address common.Address, backend bind.ContractBackend) (*Account, error) {
	contract, err := bindAccount(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Account{AccountCaller: AccountCaller{contract: contract}, AccountTransactor: AccountTransactor{contract: contract}, AccountFilterer: AccountFilterer{contract: contract}}, nil
}

// NewAccountCaller creates a new read-only instance of Account, bound to a specific deployed contract.
func NewAccountCaller(address common.Address, caller bind.ContractCaller) (*AccountCaller, error) {
	contract, err := bindAccount(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccountCaller{contract: contract}, nil
}

// NewAccountTransactor creates a new write-only instance of Account, bound to a specific deployed contract.
func NewAccountTransactor(address common.Address, transactor bind.ContractTransactor) (*AccountTransactor, error) {
	contract, err := bindAccount(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccountTransactor{contract: contract}, nil
}

// NewAccountFilterer creates a new log filterer instance of Account, bound to a specific deployed contract.
func NewAccountFilterer(address common.Address, filterer bind.ContractFilterer) (*AccountFilterer, error) {
	contract, err := bindAccount(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccountFilterer{contract: contract}, nil
}

// bindAccount binds a generic wrapper to an already deployed contract.
func bindAccount(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Account *AccountRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Account.Contract.AccountCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Account *AccountRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Account.Contract.AccountTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Account *AccountRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Account.Contract.AccountTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Account *AccountCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Account.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Account *AccountTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Account.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Account *AccountTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Account.Contract.contract.Transact(opts, method, params...)
}

// Reciever is a free data retrieval call binding the contract method 0xf4b0b756.
//
// Solidity: function reciever() view returns(address)
func (_Account *AccountCaller) Reciever(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Account.contract.Call(opts, out, "reciever")
	return *ret0, err
}

// Reciever is a free data retrieval call binding the contract method 0xf4b0b756.
//
// Solidity: function reciever() view returns(address)
func (_Account *AccountSession) Reciever() (common.Address, error) {
	return _Account.Contract.Reciever(&_Account.CallOpts)
}

// Reciever is a free data retrieval call binding the contract method 0xf4b0b756.
//
// Solidity: function reciever() view returns(address)
func (_Account *AccountCallerSession) Reciever() (common.Address, error) {
	return _Account.Contract.Reciever(&_Account.CallOpts)
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_Account *AccountTransactor) Flush(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Account.contract.Transact(opts, "flush")
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_Account *AccountSession) Flush() (*types.Transaction, error) {
	return _Account.Contract.Flush(&_Account.TransactOpts)
}

// Flush is a paid mutator transaction binding the contract method 0x6b9f96ea.
//
// Solidity: function flush() returns()
func (_Account *AccountTransactorSession) Flush() (*types.Transaction, error) {
	return _Account.Contract.Flush(&_Account.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Account *AccountTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Account.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Account *AccountSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Account.Contract.Fallback(&_Account.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Account *AccountTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Account.Contract.Fallback(&_Account.TransactOpts, calldata)
}

// AccountFlushIterator is returned from FilterFlush and is used to iterate over the raw logs and unpacked data for Flush events raised by the Account contract.
type AccountFlushIterator struct {
	Event *AccountFlush // Event containing the contract specifics and raw log

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
func (it *AccountFlushIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountFlush)
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
		it.Event = new(AccountFlush)
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
func (it *AccountFlushIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountFlushIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountFlush represents a Flush event raised by the Account contract.
type AccountFlush struct {
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterFlush is a free log retrieval operation binding the contract event 0x12b2a0ee977e74c33898f8be30fde7ae3a32ac7409a3666da55ce77e9bc32e87.
//
// Solidity: event Flush(address to, uint256 value)
func (_Account *AccountFilterer) FilterFlush(opts *bind.FilterOpts) (*AccountFlushIterator, error) {

	logs, sub, err := _Account.contract.FilterLogs(opts, "Flush")
	if err != nil {
		return nil, err
	}
	return &AccountFlushIterator{contract: _Account.contract, event: "Flush", logs: logs, sub: sub}, nil
}

// WatchFlush is a free log subscription operation binding the contract event 0x12b2a0ee977e74c33898f8be30fde7ae3a32ac7409a3666da55ce77e9bc32e87.
//
// Solidity: event Flush(address to, uint256 value)
func (_Account *AccountFilterer) WatchFlush(opts *bind.WatchOpts, sink chan<- *AccountFlush) (event.Subscription, error) {

	logs, sub, err := _Account.contract.WatchLogs(opts, "Flush")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountFlush)
				if err := _Account.contract.UnpackLog(event, "Flush", log); err != nil {
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

// ParseFlush is a log parse operation binding the contract event 0x12b2a0ee977e74c33898f8be30fde7ae3a32ac7409a3666da55ce77e9bc32e87.
//
// Solidity: event Flush(address to, uint256 value)
func (_Account *AccountFilterer) ParseFlush(log types.Log) (*AccountFlush, error) {
	event := new(AccountFlush)
	if err := _Account.contract.UnpackLog(event, "Flush", log); err != nil {
		return nil, err
	}
	return event, nil
}

// AccountRechargeIterator is returned from FilterRecharge and is used to iterate over the raw logs and unpacked data for Recharge events raised by the Account contract.
type AccountRechargeIterator struct {
	Event *AccountRecharge // Event containing the contract specifics and raw log

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
func (it *AccountRechargeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountRecharge)
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
		it.Event = new(AccountRecharge)
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
func (it *AccountRechargeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountRechargeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountRecharge represents a Recharge event raised by the Account contract.
type AccountRecharge struct {
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRecharge is a free log retrieval operation binding the contract event 0x3b47a73ca305dcdfb3a69bdd0305a75776fbe85525d92cf22644e47e16eb5c1b.
//
// Solidity: event Recharge(uint256 value)
func (_Account *AccountFilterer) FilterRecharge(opts *bind.FilterOpts) (*AccountRechargeIterator, error) {

	logs, sub, err := _Account.contract.FilterLogs(opts, "Recharge")
	if err != nil {
		return nil, err
	}
	return &AccountRechargeIterator{contract: _Account.contract, event: "Recharge", logs: logs, sub: sub}, nil
}

// WatchRecharge is a free log subscription operation binding the contract event 0x3b47a73ca305dcdfb3a69bdd0305a75776fbe85525d92cf22644e47e16eb5c1b.
//
// Solidity: event Recharge(uint256 value)
func (_Account *AccountFilterer) WatchRecharge(opts *bind.WatchOpts, sink chan<- *AccountRecharge) (event.Subscription, error) {

	logs, sub, err := _Account.contract.WatchLogs(opts, "Recharge")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountRecharge)
				if err := _Account.contract.UnpackLog(event, "Recharge", log); err != nil {
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

// ParseRecharge is a log parse operation binding the contract event 0x3b47a73ca305dcdfb3a69bdd0305a75776fbe85525d92cf22644e47e16eb5c1b.
//
// Solidity: event Recharge(uint256 value)
func (_Account *AccountFilterer) ParseRecharge(log types.Log) (*AccountRecharge, error) {
	event := new(AccountRecharge)
	if err := _Account.contract.UnpackLog(event, "Recharge", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WalletABI is the input ABI used to generate the binding from.
const WalletABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"Create\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"accounts\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_salt\",\"type\":\"bytes32\"}],\"name\":\"create\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// WalletFuncSigs maps the 4-byte function signature to its string representation.
var WalletFuncSigs = map[string]string{
	"5e5c06e2": "accounts(address)",
	"f851a440": "admin()",
	"a3def923": "create(address,bytes32)",
}

// WalletBin is the compiled bytecode used for deploying new contracts.
var WalletBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556103e5806100326000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635e5c06e214610046578063a3def92314610080578063f851a440146100ae575b600080fd5b61006c6004803603602081101561005c57600080fd5b50356001600160a01b03166100d2565b604080519115158252519081900360200190f35b6100ac6004803603604081101561009657600080fd5b506001600160a01b0381351690602001356100e7565b005b6100b66101b2565b604080516001600160a01b039092168252519081900360200190f35b60016020526000908152604090205460ff1681565b6000546001600160a01b0316331461012c576040805162461bcd60e51b815260206004820152600360248201526234303360e81b604482015290519081900360640190fd5b6000818360405161013c906101c1565b6001600160a01b0390911681526040518291819003602001906000f590508015801561016c573d6000803e3d6000fd5b50604080516001600160a01b038316815290519192507fe3758539c1bd6726422843471b2886c2d2cefd3b4aead6778386283e20a32a80919081900360200190a1505050565b6000546001600160a01b031681565b6101e1806101cf8339019056fe608060405234801561001057600080fd5b506040516101e13803806101e18339818101604052602081101561003357600080fd5b5051600080546001600160a01b039092166001600160a01b031990921691909117905561017c806100656000396000f3fe6080604052600436106100295760003560e01c80636b9f96ea1461005e578063f4b0b75614610075575b6040805134815290517f3b47a73ca305dcdfb3a69bdd0305a75776fbe85525d92cf22644e47e16eb5c1b9181900360200190a1005b34801561006a57600080fd5b506100736100a6565b005b34801561008157600080fd5b5061008a610137565b604080516001600160a01b039092168252519081900360200190f35b47806100b25750610135565b600080546040516001600160a01b039091169183156108fc02918491818181858888f193505050501580156100eb573d6000803e3d6000fd5b50600054604080516001600160a01b0390921682526020820183905280517f12b2a0ee977e74c33898f8be30fde7ae3a32ac7409a3666da55ce77e9bc32e879281900390910190a1505b565b6000546001600160a01b03168156fea264697066735822122066d9e560c66006a9376726110cda2f35329700f29dd9ae7cdba69191bb56b40864736f6c63430007000033a2646970667358221220245cbf9737205ebb8eb8aad356036578ddcd254714581c96baf101e43e07712664736f6c63430007000033"

// DeployWallet deploys a new Ethereum contract, binding an instance of Wallet to it.
func DeployWallet(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Wallet, error) {
	parsed, err := abi.JSON(strings.NewReader(WalletABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(WalletBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Wallet{WalletCaller: WalletCaller{contract: contract}, WalletTransactor: WalletTransactor{contract: contract}, WalletFilterer: WalletFilterer{contract: contract}}, nil
}

// Wallet is an auto generated Go binding around an Ethereum contract.
type Wallet struct {
	WalletCaller     // Read-only binding to the contract
	WalletTransactor // Write-only binding to the contract
	WalletFilterer   // Log filterer for contract events
}

// WalletCaller is an auto generated read-only Go binding around an Ethereum contract.
type WalletCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WalletTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WalletTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WalletFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WalletFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WalletSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WalletSession struct {
	Contract     *Wallet           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WalletCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WalletCallerSession struct {
	Contract *WalletCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// WalletTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WalletTransactorSession struct {
	Contract     *WalletTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WalletRaw is an auto generated low-level Go binding around an Ethereum contract.
type WalletRaw struct {
	Contract *Wallet // Generic contract binding to access the raw methods on
}

// WalletCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WalletCallerRaw struct {
	Contract *WalletCaller // Generic read-only contract binding to access the raw methods on
}

// WalletTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WalletTransactorRaw struct {
	Contract *WalletTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWallet creates a new instance of Wallet, bound to a specific deployed contract.
func NewWallet(address common.Address, backend bind.ContractBackend) (*Wallet, error) {
	contract, err := bindWallet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Wallet{WalletCaller: WalletCaller{contract: contract}, WalletTransactor: WalletTransactor{contract: contract}, WalletFilterer: WalletFilterer{contract: contract}}, nil
}

// NewWalletCaller creates a new read-only instance of Wallet, bound to a specific deployed contract.
func NewWalletCaller(address common.Address, caller bind.ContractCaller) (*WalletCaller, error) {
	contract, err := bindWallet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WalletCaller{contract: contract}, nil
}

// NewWalletTransactor creates a new write-only instance of Wallet, bound to a specific deployed contract.
func NewWalletTransactor(address common.Address, transactor bind.ContractTransactor) (*WalletTransactor, error) {
	contract, err := bindWallet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WalletTransactor{contract: contract}, nil
}

// NewWalletFilterer creates a new log filterer instance of Wallet, bound to a specific deployed contract.
func NewWalletFilterer(address common.Address, filterer bind.ContractFilterer) (*WalletFilterer, error) {
	contract, err := bindWallet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WalletFilterer{contract: contract}, nil
}

// bindWallet binds a generic wrapper to an already deployed contract.
func bindWallet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WalletABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Wallet *WalletRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Wallet.Contract.WalletCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Wallet *WalletRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Wallet.Contract.WalletTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Wallet *WalletRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Wallet.Contract.WalletTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Wallet *WalletCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Wallet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Wallet *WalletTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Wallet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Wallet *WalletTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Wallet.Contract.contract.Transact(opts, method, params...)
}

// Accounts is a free data retrieval call binding the contract method 0x5e5c06e2.
//
// Solidity: function accounts(address ) view returns(bool)
func (_Wallet *WalletCaller) Accounts(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Wallet.contract.Call(opts, out, "accounts", arg0)
	return *ret0, err
}

// Accounts is a free data retrieval call binding the contract method 0x5e5c06e2.
//
// Solidity: function accounts(address ) view returns(bool)
func (_Wallet *WalletSession) Accounts(arg0 common.Address) (bool, error) {
	return _Wallet.Contract.Accounts(&_Wallet.CallOpts, arg0)
}

// Accounts is a free data retrieval call binding the contract method 0x5e5c06e2.
//
// Solidity: function accounts(address ) view returns(bool)
func (_Wallet *WalletCallerSession) Accounts(arg0 common.Address) (bool, error) {
	return _Wallet.Contract.Accounts(&_Wallet.CallOpts, arg0)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Wallet *WalletCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Wallet.contract.Call(opts, out, "admin")
	return *ret0, err
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Wallet *WalletSession) Admin() (common.Address, error) {
	return _Wallet.Contract.Admin(&_Wallet.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Wallet *WalletCallerSession) Admin() (common.Address, error) {
	return _Wallet.Contract.Admin(&_Wallet.CallOpts)
}

// Create is a paid mutator transaction binding the contract method 0xa3def923.
//
// Solidity: function create(address _to, bytes32 _salt) returns()
func (_Wallet *WalletTransactor) Create(opts *bind.TransactOpts, _to common.Address, _salt [32]byte) (*types.Transaction, error) {
	return _Wallet.contract.Transact(opts, "create", _to, _salt)
}

// Create is a paid mutator transaction binding the contract method 0xa3def923.
//
// Solidity: function create(address _to, bytes32 _salt) returns()
func (_Wallet *WalletSession) Create(_to common.Address, _salt [32]byte) (*types.Transaction, error) {
	return _Wallet.Contract.Create(&_Wallet.TransactOpts, _to, _salt)
}

// Create is a paid mutator transaction binding the contract method 0xa3def923.
//
// Solidity: function create(address _to, bytes32 _salt) returns()
func (_Wallet *WalletTransactorSession) Create(_to common.Address, _salt [32]byte) (*types.Transaction, error) {
	return _Wallet.Contract.Create(&_Wallet.TransactOpts, _to, _salt)
}

// WalletCreateIterator is returned from FilterCreate and is used to iterate over the raw logs and unpacked data for Create events raised by the Wallet contract.
type WalletCreateIterator struct {
	Event *WalletCreate // Event containing the contract specifics and raw log

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
func (it *WalletCreateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WalletCreate)
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
		it.Event = new(WalletCreate)
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
func (it *WalletCreateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WalletCreateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WalletCreate represents a Create event raised by the Wallet contract.
type WalletCreate struct {
	Arg0 common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterCreate is a free log retrieval operation binding the contract event 0xe3758539c1bd6726422843471b2886c2d2cefd3b4aead6778386283e20a32a80.
//
// Solidity: event Create(address arg0)
func (_Wallet *WalletFilterer) FilterCreate(opts *bind.FilterOpts) (*WalletCreateIterator, error) {

	logs, sub, err := _Wallet.contract.FilterLogs(opts, "Create")
	if err != nil {
		return nil, err
	}
	return &WalletCreateIterator{contract: _Wallet.contract, event: "Create", logs: logs, sub: sub}, nil
}

// WatchCreate is a free log subscription operation binding the contract event 0xe3758539c1bd6726422843471b2886c2d2cefd3b4aead6778386283e20a32a80.
//
// Solidity: event Create(address arg0)
func (_Wallet *WalletFilterer) WatchCreate(opts *bind.WatchOpts, sink chan<- *WalletCreate) (event.Subscription, error) {

	logs, sub, err := _Wallet.contract.WatchLogs(opts, "Create")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WalletCreate)
				if err := _Wallet.contract.UnpackLog(event, "Create", log); err != nil {
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

// ParseCreate is a log parse operation binding the contract event 0xe3758539c1bd6726422843471b2886c2d2cefd3b4aead6778386283e20a32a80.
//
// Solidity: event Create(address arg0)
func (_Wallet *WalletFilterer) ParseCreate(log types.Log) (*WalletCreate, error) {
	event := new(WalletCreate)
	if err := _Wallet.contract.UnpackLog(event, "Create", log); err != nil {
		return nil, err
	}
	return event, nil
}
