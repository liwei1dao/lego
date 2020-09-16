package token

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/sdks/token/solidity"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	opt        *Options
	client     *ethclient.Client
	tokeninstance 
	closesignal chan bool
}

func (this *Token) Init() (err error) {
	this.client, err = ethclient.Dial(this.opt.EthPoolAdrr)
	if err != nil {
		return
	}
	// Golem (GNT) Address
	tokenAddress := common.HexToAddress(this.opt.TokenAddr)
	this.tokeninstance, err = solidity.NewToken(tokenAddress, this.client)
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

func (this *Token)BalanceOf(_address string){



	address := common.HexToAddress(_address)
    bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
    if err != nil {
        log.Fatal(err)
    }
}


//监听代币事件
func (this *Token) MonitorTokenEvent(){
	contractAbi, err := abi.JSON(strings.NewReader(string(solidity.HiToolCoinABI)))
	if err != nil {
		log.Fatalf("Eth Token MonitorTokenEvent Fatal:%s",err.Error())
		return 
	}

	contractAddress := common.HexToAddress(this.opt.TokenAddr)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Eth Token MonitorTokenEvent Fatal:%s",err.Error())
		return 
	}
	ctx, cancel := context.WithDeadline(context.Background(), timeout)

	transferevent := struct {
		from common.Address, 
		to common.Address, 
		value *big.Int,
	}{}
	approvalevent := struct {
		from common.Address, 
		to common.Address, 
		value *big.Int,
	}{}

	for {
		select {
		case err := <-sub.Err():
			log.Errorf("Eth Token MonitorTokenEvent Fatal:%s",err.Error())
			return
		case <- this.opt.closesignal:
			return
		case vLog := <-logs:
			//交易事件
			err := contractAbi.Unpack(&transferevent, "Transfer", vLog.Data)
			if err == nil && this.opt.TransferEvent != nil {
				this.opt.TransferEvent(transferevent.from,transferevent.to,transferevent.value)
			}
			//授权事件
			err := contractAbi.Unpack(&approvalevent, "Approval", vLog.Data)
			if err == nil && this.opt.ApprovalEvent != nil{
				this.opt.ApprovalEvent(transferevent.from,transferevent.to,transferevent.value)
			}
		}
	}
}
