package influxdb

import "github.com/liwei1dao/lego/utils/mapstructure"

type Option func(*Options)
type Options struct {
	Addr    string
	Token   string
	TimeOut int
}

///地址
func SetAddr(v string) Option {
	return func(o *Options) {
		o.Addr = v
	}
}

func SetToken(v string) Option {
	return func(o *Options) {
		o.Token = v
	}
}
func SetTimeOut(v int) Option {
	return func(o *Options) {
		o.TimeOut = v
	}
}
func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		TimeOut: 5,
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
		TimeOut: 5,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
