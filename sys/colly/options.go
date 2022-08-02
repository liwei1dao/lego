package colly

import (
	"errors"
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	DomainGlob  string
	Parallelism int
	Delay       time.Duration
	Proxy       []string
	UserAgent   string
	Debug       bool //日志是否开启
	Log         log.ILog
}

func SetDomainGlob(v string) Option {
	return func(o *Options) {
		o.DomainGlob = v
	}
}

func SetParallelism(v int) Option {
	return func(o *Options) {
		o.Parallelism = v
	}
}
func SetDelay(v time.Duration) Option {
	return func(o *Options) {
		o.Delay = v
	}
}
func SetProxy(v []string) Option {
	return func(o *Options) {
		o.Proxy = v
	}
}
func SetUserAgent(v string) Option {
	return func(o *Options) {
		o.UserAgent = v
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
	if options.Debug && options.Log == nil {
		if options.Log = log.Clone(2); options.Log == nil {
			err = errors.New("log is nil")
		}
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{}
	for _, o := range opts {
		o(options)
	}
	if options.Debug && options.Log == nil {
		if options.Log = log.Clone(2); options.Log == nil {
			err = errors.New("log is nil")
		}
	}
	return
}
