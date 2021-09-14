package colly

import "github.com/liwei1dao/lego/utils/mapstructure"

type Option func(*Options)
type Options struct {
	Proxy     []string
	UserAgent string
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
