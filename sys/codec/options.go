package codec

import (
	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

func newOptions(config map[string]interface{}, opts ...core.Option) core.Options {
	options := core.Options{
		IndentionStep: 2,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...core.Option) core.Options {
	options := core.Options{
		IndentionStep: 2,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
