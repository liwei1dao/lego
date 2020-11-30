package timewheel

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
)

type Option func(*Options)
type Options struct {
	Tick       int //单位毫秒
	BucketsNum int
	IsSyncPool bool
}

func SetTick(v time.Duration) Option {
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


func newOptions(config map[string]interface{},opts ...Option) Options {
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
	if opt.Tick < 100 {
		log.Errorf("创建时间轮参数异常 Tick 必须大于 100 ms ")
		opt.Tick = 100
	}
	if opt.BucketsNum < 0 {
		log.Errorf("创建时间轮参数异常 BucketsNum 必须大于 0 ")
		opt.BucketsNum = 1
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
	if opt.Tick < 100 {
		log.Errorf("创建时间轮参数异常 Tick 必须大于 100 ms ")
		opt.Tick = 100
	}
	if opt.BucketsNum < 0 {
		log.Errorf("创建时间轮参数异常 BucketsNum 必须大于 0 ")
		opt.BucketsNum = 1
	}
	return options
}