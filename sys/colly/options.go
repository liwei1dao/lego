package colly

import (
	"time"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	DomainGlob  string
	Parallelism int
	Delay       time.Duration
	Proxy       []string
	UserAgent   string
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

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return options
}
