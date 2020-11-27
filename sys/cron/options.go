package cron

import (
	"github.com/liwei1dao/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
}

func newOptionsByConfig(config map[string]interface{}) Options {
	options := Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
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
