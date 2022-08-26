package paypal

import (
	"errors"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	ClientID  string //paypal clientID
	SecretID  string //paypal secretID
	IsSandBox bool   //paypal 是否是沙盒环境
	Debug     bool   //日志是否开启
	Log       log.ILogger
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
	options = &Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}

	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.paypal", 2)); options.Log == nil {
		err = errors.New("log is nil")
	}

	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{}
	for _, o := range opts {
		o(options)
	}
	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.paypal", 2)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return
}
