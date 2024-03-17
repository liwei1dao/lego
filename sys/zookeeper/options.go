package zookeeper

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Adds    []string
	TimeOut int
}

func SetAdds(v []string) Option {
	return func(o *Options) {
		o.Adds = v
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
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	return options
}
