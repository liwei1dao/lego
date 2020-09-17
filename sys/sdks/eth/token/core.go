package token

import (
	"github.com/liwei1dao/lego/core"
)

type (
	IToken interface {
		Start()
		Stop()
		ChangeEthFundDeposit(newFundDeposit string) (string, error) //更改以太币接受地址
		SetTokenExchangeRate(exchange uint32) (string, error)       //设置代币汇率
		BalanceOf(_address string) (uint64, error)                  //查询地址代币值
	}
)

var (
	tokne IToken
)

func OnInit(s core.IService, opt ...Option) (err error) {
	tokne, err = newToken(opt...)
	return
}

func Start() {
	tokne.Start()
}

func Stop() {
	tokne.Stop()
}

func ChangeEthFundDeposit(newFundDeposit string) (string, error) {
	return tokne.ChangeEthFundDeposit(newFundDeposit)
}

func SetTokenExchangeRate(exchange uint32) (string, error) {
	return tokne.SetTokenExchangeRate(exchange)
}

func BalanceOf(_address string) (uint64, error) {
	return tokne.BalanceOf(_address)
}
