package codec

import (
	"errors"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

func newOptions(config map[string]interface{}, opts ...core.Option) (options *core.Options, err error) {
	options = &core.Options{
		IndentionStep: 2,
		TagKey:        "json",
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.Debug && options.Log == nil {
		options.Log = log.Clone()
	}
	if options.Log = log.Clone(log.SetLoglayer(2)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return
}

func newOptionsByOption(opts ...core.Option) (options *core.Options, err error) {
	options = &core.Options{
		IndentionStep: 2,
		TagKey:        "json",
	}
	for _, o := range opts {
		o(options)
	}
	if options.Debug && options.Log == nil {
		options.Log = log.Clone()
	}
	if options.Log = log.Clone(log.SetLoglayer(2)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return
}
