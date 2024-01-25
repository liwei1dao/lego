package timewheel

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Tick       time.Duration //不小于 10毫秒
	BucketsNum int
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
func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Tick:       time.Second,
		BucketsNum: 1024,
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
		Tick:       time.Second,
		BucketsNum: 1024,
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
