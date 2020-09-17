package token

import (
	"context"
	"crypto/ecdsa"
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

func newPay(opt ...Option) (IToken, error) {
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
	opt           *Options
	client        *ethclient.Client
	privateKey    *ecdsa.PrivateKey
	tokenInstance *solidity.HiToolCoin
	closesignal   chan bool
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
	//代币合约
	tokenaddress := common.HexToAddress(this.opt.TokenAddr)
	tokenInstance, err := solidity.NewHiToolCoin(tokenaddress, client)
	if err != nil {
		return
	}
	return
}

func (this *Token) Start() (err error) {
	go this.MonitorTokenEvent()
}

func (this *Token) Stop() (err error) {
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
func (this *Token) SetTokenExchangeRate(exchange uint32) error {
	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	bal, err := this.tokenInstance.SetTokenExchangeRate(auth, big.NewInt(exchange))
	if err != nil {
		return err
	}
	return nil
}

//设置合约eth接收地址
func (this *Token) ChangeEthFundDeposit(newFundDeposit string) (common.Address, error) {
	auth := bind.NewKeyedTransactor(this.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	address := common.HexToAddress(newFundDeposit)
	bal, err := this.tokenInstance.SetTokenExchangeRate(auth, address)
	if err != nil {
		return err
	}
	return nil
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
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Eth Token MonitorTokenEvent Fatal:%s", err.Error())
		return
	}
	ctx, cancel := context.WithDeadline(context.Background(), timeout)

	transferevent := struct {
		from  common.Address
		to    common.Address
		value *big.Ints
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
		case <-this.opt.closesignal:
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
				this.opt.ApprovalEvent(transferevent.from, transferevent.to, transferevent.value)
			}
		}
	}
}
