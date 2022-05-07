package blockcache

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	CacheMaxSzie int64
	Debug        bool //日志是否开启
	Log          log.ILog
}

func SetCacheMaxSzie(v int64) Option {
	return func(o *Options) {
		o.CacheMaxSzie = v
	}
}

func SetDebug(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}
func SetLog(v log.ILog) Option {
	return func(o *Options) {
		o.Log = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		CacheMaxSzie: 1024 * 1024 * 100,
		Debug:        true,
		Log:          log.Clone(log.SetLoglayer(2)),
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		CacheMaxSzie: 1024 * 1024 * 100,
		Debug:        true,
		Log:          log.Clone(log.SetLoglayer(2)),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
