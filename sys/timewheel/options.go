package timewheel

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Tick       int //单位毫秒
	BucketsNum int
	IsSyncPool bool
}

func SetTick(v int) Option {
	return func(o *Options) {
		o.Tick = v
	}
}

func SetBucketsNum(v int) Option {
	return func(o *Options) {
		o.BucketsNum = v
	}
}

func SetIsSyncPool(v bool) Option {
	return func(o *Options) {
		o.IsSyncPool = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Tick:       1000,
		BucketsNum: 1,
		IsSyncPool: true,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if options.Tick < 100 {
		log.Errorf("创建时间轮参数异常 Tick 必须大于 100 ms ")
		options.Tick = 100
	}
	if options.BucketsNum < 0 {
		log.Errorf("创建时间轮参数异常 BucketsNum 必须大于 0 ")
		options.BucketsNum = 1
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		Tick:       1000,
		BucketsNum: 1,
		IsSyncPool: true,
	}
	for _, o := range opts {
		o(&options)
	}
	if options.Tick < 100 {
		log.Warnf("创建时间轮参数异常 Tick 必须大于 100 ms ")
		options.Tick = 100
	}
	if options.BucketsNum < 0 {
		log.Warnf("创建时间轮参数异常 BucketsNum 必须大于 0 ")
		options.BucketsNum = 1
	}
	return options
}
