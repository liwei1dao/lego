package pay

import (
	"fmt"
)

type Option func(*Options)
type Options struct {
	EthPoolAdrr          string //以太坊币池地址
	WalletAdrr           string //字符系统钱包合约地址 需提前部署
	ControllerPrivateKey string //系统控制者ETH 私钥 确保账号下有充足的以太币
	FundRecoveryAddr     string //资金回收地址	最后收钱的账号 请确保安全
}

func SetEthPoolAdrr(v string) Option {
	return func(o *Options) {
		o.EthPoolAdrr = v
	}
}

func SetWalletAdrr(v string) Option {
	return func(o *Options) {
		o.WalletAdrr = v
	}
}

func SetControllerPrivateKey(v string) Option {
	return func(o *Options) {
		o.ControllerPrivateKey = v
	}
}

func SetFundRecoveryAddr(v string) Option {
	return func(o *Options) {
		o.FundRecoveryAddr = v
	}
}

func newOptions(opts ...Option) (opt *Options, err error) {
	opt = &Options{}
	for _, o := range opts {
		o(opt)
	}
	if opt.ControllerPrivateKey == "" || opt.FundRecoveryAddr == "" {
		return nil, fmt.Errorf("")
	}

	return opt, nil
}
