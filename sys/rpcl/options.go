package rpcl

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/liwei1dao/lego/sys/rpcl/selector"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	UpDateInterval time.Duration
	Selector       core.ISelector
	Debug          bool //日志是否开启
	Log            log.ILog
}

func Set_UpDateInterval(v time.Duration) Option {
	return func(o *Options) {
		o.UpDateInterval = v
	}
}

func Set_Selector(v core.ISelector) Option {
	return func(o *Options) {
		o.Selector = v
	}
}

func Set_Debug(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}

func Set_Log(v log.ILog) Option {
	return func(o *Options) {
		o.Log = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Selector: selector.NewSelector(core.RandomSelect, map[string]string{}),
		Debug:    true,
		Log:      log.Clone(log.SetLoglayer(2)),
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
		Selector: selector.NewSelector(core.RandomSelect, map[string]string{}),
		Debug:    true,
		Log:      log.Clone(log.SetLoglayer(2)),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
