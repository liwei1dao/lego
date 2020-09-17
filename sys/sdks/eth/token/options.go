package token

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Option func(*Options)
type Options struct {
	EthPoolAdrr          string                                                       //以太坊币池地址
	TokenAddr            string                                                       //代币合约地址
	ControllerPrivateKey string                                                       //系统控制者ETH 私钥 确保账号下有充足的以太币
	TransferEvent        func(form common.Address, to common.Address, value *big.Int) //代币交易事件
	ApprovalEvent        func(form common.Address, to common.Address, value *big.Int) //代币授权事件
}

func SetEthPoolAdrr(v string) Option {
	return func(o *Options) {
		o.EthPoolAdrr = v
	}
}

func SetTokenAddr(v string) Option {
	return func(o *Options) {
		o.TokenAddr = v
	}
}

func SetControllerPrivateKey(v string) Option {
	return func(o *Options) {
		o.ControllerPrivateKey = v
	}
}

func SetTransferEvent(v func(form common.Address, to common.Address, value *big.Int)) Option {
	return func(o *Options) {
		o.TransferEvent = v
	}
}

func SetApprovalEvent(v func(form common.Address, to common.Address, value *big.Int)) Option {
	return func(o *Options) {
		o.ApprovalEvent = v
	}
}

func newOptions(opts ...Option) (opt *Options, err error) {
	opt = &Options{}
	for _, o := range opts {
		o(opt)
	}

	if opt.EthPoolAdrr == "" || opt.ControllerPrivateKey == "" || opt.TokenAddr == "" {
		return nil, fmt.Errorf("Eth Token options Error")
	}
	return opt, nil
}
