package cachego

import (
	"time"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Expiration      time.Duration
	CleanupInterval time.Duration
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Expiration:      -1,
		CleanupInterval: 0,
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
		Expiration:      -1,
		CleanupInterval: 0,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
