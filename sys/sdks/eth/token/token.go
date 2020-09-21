package token

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/sdks/eth/token/solidity"
)

func newToken(opt ...Option) (IToken, error) {
	opts, err := newOptions(opt...)
	if err != nil {
		return nil, err
	}
	token := &Token{
		opt: opts,
	}
	err = token.Init()
	return token, err
}

type Token struct {
	opt            *Options
	client         *ethclient.Client
	privateKey     *ecdsa.PrivateKey
	ControllerAddr common.Address //控制账号钱包地址
	tokenInstance  *solidity.HiToolCoin
	closesignal    chan bool
}

func (this *Token) Init() (err error) {
	this.client, err = ethclient.Dial(this.opt.EthPoolAdrr)
	if err != nil {
		return
	}
	this.privateKey, err = crypto.HexToECDSA(this.opt.ControllerPrivateKey)
	if err != nil {
		return
	}
	publicKey := this.privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	this.ControllerAddr = crypto.PubkeyToAddress(*publicKeyECDSA)
	//代币合约
	tokenaddress := common.HexToAddress(this.opt.TokenAddr)
	this.tokenInstance, err = solidity.NewHiToolCoin(tokenaddress, this.client)
	if err != nil {
		return
	}
	return
}

func (this *Token) Start() {
	go this.MonitorTokenEvent()
}

func (this *Token) Stop() {
	this.closesignal <- true
}

//获取目标地址的代币存量
func (this *Token) BalanceOf(_address string) (uint64, error) {
	address := common.HexToAddress(_address)
	bal, err := this.tokenInstance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return 0, err
	}
	return bal.Uint64(), nil
}

//设置代币汇率
func (this *Token) SetTokenExchangeRate(exchange uint32) (string, error) {
	gasPrice, err := this.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	nonce, err := this.client.PendingNonceAt(context.Background(), this.ControllerAddr)
	if err != nil {
		return "", err
	}
	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(1000000) // in units
	auth.GasPrice = gasPrice
	tx, err := this.tokenInstance.SetTokenExchangeRate(auth, big.NewInt(int64(exchange)))
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

//设置合约eth接收地址
func (this *Token) ChangeEthFundDeposit(newFundDeposit string) (string, error) {
	gasPrice, err := this.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	nonce, err := this.client.PendingNonceAt(context.Background(), this.ControllerAddr)
	if err != nil {
		return "", err
	}
	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(1000000) // in units
	auth.GasPrice = gasPrice
	address := common.HexToAddress(newFundDeposit)
	tx, err := this.tokenInstance.ChangeEthFundDeposit(auth, address)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

//设置合约eth接收地址
func (this *Token) TransferETH() (string, error) {
	gasPrice, err := this.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	nonce, err := this.client.PendingNonceAt(context.Background(), this.ControllerAddr)
	if err != nil {
		return "", err
	}
	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(1000000) // in units
	auth.GasPrice = gasPrice
	tx, err := this.tokenInstance.TransferETH(auth)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

//设置合约eth接收地址
func (this *Token) Addmint(amount uint32) (string, error) {
	gasPrice, err := this.client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	nonce, err := this.client.PendingNonceAt(context.Background(), this.ControllerAddr)
	if err != nil {
		return "", err
	}
	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(1000000) // in units
	auth.GasPrice = gasPrice
	tx, err := this.tokenInstance.Addmint(auth, big.NewInt(int64(amount)))
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

//监听代币事件
func (this *Token) MonitorTokenEvent() {
	contractAbi, err := abi.JSON(strings.NewReader(string(solidity.HiToolCoinABI)))
	if err != nil {
		log.Fatalf("Eth Token MonitorTokenEvent Fatal:%s", err.Error())
		return
	}

	contractAddress := common.HexToAddress(this.opt.TokenAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := this.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Eth Token MonitorTokenEvent Fatal:%s", err.Error())
		return
	}

	transferevent := struct {
		from  common.Address
		to    common.Address
		value *big.Int
	}{}
	approvalevent := struct {
		from  common.Address
		to    common.Address
		value *big.Int
	}{}

	for {
		select {
		case err := <-sub.Err():
			log.Errorf("Eth Token MonitorTokenEvent Fatal:%s", err.Error())
			return
		case <-this.closesignal:
			return
		case vLog := <-logs:
			//交易事件
			err := contractAbi.Unpack(&transferevent, "Transfer", vLog.Data)
			if err == nil && this.opt.TransferEvent != nil {
				this.opt.TransferEvent(transferevent.from, transferevent.to, transferevent.value)
			}
			//授权事件
			err = contractAbi.Unpack(&approvalevent, "Approval", vLog.Data)
			if err == nil && this.opt.ApprovalEvent != nil {
				this.opt.ApprovalEvent(approvalevent.from, approvalevent.to, approvalevent.value)
			}
		}
	}
}
