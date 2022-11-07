package paypal

import (
	"errors"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	AppName   string
	ClientID  string //paypal clientID
	SecretID  string //paypal secretID
	IsSandBox bool   //paypal 是否是沙盒环境
	Currency  string //货币类型 "USD" 美元
	ReturnURL string //支付回调地址
	Debug     bool   //日志是否开启
	Log       log.ILogger
}

func SetAppName(v string) Option {
	return func(o *Options) {
		o.AppName = v
	}
}
func SetClientID(v string) Option {
	return func(o *Options) {
		o.ClientID = v
	}
}
func SetSecretID(v string) Option {
	return func(o *Options) {
		o.SecretID = v
	}
}
func SetIsSandBox(v bool) Option {
	return func(o *Options) {
		o.IsSandBox = v
	}
}
func SetCurrency(v string) Option {
	return func(o *Options) {
		o.Currency = v
	}
}
func SetReturnURL(v string) Option {
	return func(o *Options) {
		o.ReturnURL = v
	}
}
func SetDebug(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}
func SetLog(v log.ILogger) Option {
	return func(o *Options) {
		o.Log = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		AppName:  "lego",
		Currency: "USD",
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}

	if options.Log == nil {
		if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.paypal", 2)); options.Log == nil {
			err = errors.New("log is nil")
		}
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		AppName:  "lego",
		Currency: "USD",
	}
	for _, o := range opts {
		o(options)
	}
	if options.Debug && options.Log == nil {
		if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.paypal", 2)); options.Log == nil {
			err = errors.New("log is nil")
		}
	} else if !options.Debug && options.Log == nil {
		if options.Log = log.NewTurnlog(options.Debug, nil); options.Log == nil {
			err = errors.New("log is nil")
		}
	}
	return
}
