package blockcache

import (
	"errors"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	CacheMaxSzie int64
	Debug        bool //日志是否开启
	Log          log.Ilogf
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
func SetLog(v log.Ilogf) Option {
	return func(o *Options) {
		o.Log = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		CacheMaxSzie: 1024 * 1024 * 100,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}

	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.Blockcache", 2)); options.Log == nil {
		err = errors.New("log is nil")
	}

	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		CacheMaxSzie: 1024 * 1024 * 100,
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.Blockcache", 2)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return
}
