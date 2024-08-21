package httppool

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	MaxIdleConns        int //最大待机链接数量
	MaxIdleConnsPerHost int //每个主机的最大活跃待机链接数
	IdleConnTimeout     int //待机链接超时时间
	RequestTimeout      int //请求超时时间
	//日志是否开启
	Debug bool
	Log   log.ILogger
}

func SetMaxIdleConns(v int) Option {
	return func(o *Options) {
		o.MaxIdleConns = v
	}
}

func SetMaxIdleConnsPerHost(v int) Option {
	return func(o *Options) {
		o.MaxIdleConnsPerHost = v
	}
}

func SetIdleConnTimeout(v int) Option {
	return func(o *Options) {
		o.IdleConnTimeout = v
	}
}

func SetRequestTimeout(v int) Option {
	return func(o *Options) {
		o.RequestTimeout = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.httppool", 3))
	}

	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.httppool", 3))
	}
	return options
}
