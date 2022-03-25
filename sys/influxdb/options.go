package influxdb

import "github.com/liwei1dao/lego/utils/mapstructure"

type Option func(*Options)
type Options struct {
	Addr     string
	Username string
	Password string
}

///地址
func SetAddr(v string) Option {
	return func(o *Options) {
		o.Addr = v
	}
}

func SetUsername(v string) Option {
	return func(o *Options) {
		o.Username = v
	}
}
func SetPassword(v string) Option {
	return func(o *Options) {
		o.Password = v
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
