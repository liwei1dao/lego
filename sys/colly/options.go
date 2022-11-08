package colly

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	DomainGlob      string
	Parallelism     int
	Delay           time.Duration
	Proxy           []string
	SkipCertificate bool //是否跳过证书验证
	UserAgent       string
	Debug           bool //日志是否开启
	Log             log.ILogger
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
func SetSkipCertificate(v bool) Option {
	return func(o *Options) {
		o.SkipCertificate = v
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
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.colly", 3))
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{}
	for _, o := range opts {
		o(options)
	}

	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.colly", 3))
	}
	return
}
