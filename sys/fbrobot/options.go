package fbrobot

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	ListenPort int //监听端口
	//facebook app 配置项
	VerifyToken string
	AppSecret   string
	AccessToken string
	PageID      string

	//自定义消息处理函数
	MessageReceived  MessageReceivedHandler
	MessageDelivered MessageDeliveredHandler
	Postback         PostbackHandler
	Authentication   AuthenticationHandler

	//日志是否开启
	Debug bool
	Log   log.ILogger
}

func Set_ListenPort(v int) Option {
	return func(o *Options) {
		o.ListenPort = v
	}
}

func Set_VerifyToken(v string) Option {
	return func(o *Options) {
		o.VerifyToken = v
	}
}

func Set_AppSecret(v string) Option {
	return func(o *Options) {
		o.AppSecret = v
	}
}

func Set_AccessToken(v string) Option {
	return func(o *Options) {
		o.AccessToken = v
	}
}

func Set_PageID(v string) Option {
	return func(o *Options) {
		o.PageID = v
	}
}

func Set_MessageReceived(v MessageReceivedHandler) Option {
	return func(o *Options) {
		o.MessageReceived = v
	}
}

func Set_MessageDelivered(v MessageDeliveredHandler) Option {
	return func(o *Options) {
		o.MessageDelivered = v
	}
}
func Set_Postback(v PostbackHandler) Option {
	return func(o *Options) {
		o.Postback = v
	}
}
func Set_Authentication(v AuthenticationHandler) Option {
	return func(o *Options) {
		o.Authentication = v
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
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.fbrobot", 3))
	}

	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{}
	for _, o := range opts {
		o(options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.fbrobot", 3))
	}

	return
}
