package nacos

import "github.com/liwei1dao/lego/utils/mapstructure"

type Option func(*Options)
type Options struct {
	NamespaceId string
	NacosAddr   string
	Port        uint64
	TimeoutMs   uint64
}

func SetNamespaceId(v string) Option {
	return func(o *Options) {
		o.NamespaceId = v
	}
}
func SetNacosAddr(v string) Option {
	return func(o *Options) {
		o.NacosAddr = v
	}
}

func SetPort(v uint64) Option {
	return func(o *Options) {
		o.Port = v
	}
}

func SetTimeoutMs(v uint64) Option {
	return func(o *Options) {
		o.TimeoutMs = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		TimeoutMs: 5000,
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
		TimeoutMs: 5000,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
