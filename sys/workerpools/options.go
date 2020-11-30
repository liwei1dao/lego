package workerpools

import (
	"time"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	DefWrokers     int
	MaxWorkers     int
	Tasktimeout    time.Duration //任务执行操超时间
	IdleTimeoutSec time.Duration //超时释放空闲工作人员
}

func SetDefWorkers(v int) Option {
	return func(o *Options) {
		o.DefWrokers = v
	}
}

func SetMaxWorkers(v int) Option {
	return func(o *Options) {
		o.MaxWorkers = v
	}
}

func SetTaskTimeOut(v time.Duration) Option {
	return func(o *Options) {
		o.Tasktimeout = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		DefWrokers:  10,
		MaxWorkers:  20,
		Tasktimeout: time.Second * 3,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if options.DefWrokers < 1 {
		options.DefWrokers = 10
	}
	if options.MaxWorkers < 1 {
		options.MaxWorkers = 20
	}
	if options.Tasktimeout < time.Millisecond {
		options.Tasktimeout = time.Second * 3
	}
	if options.IdleTimeoutSec < time.Second*5 {
		options.IdleTimeoutSec = time.Second * 5
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		DefWrokers:  10,
		MaxWorkers:  20,
		Tasktimeout: time.Second * 3,
	}
	for _, o := range opts {
		o(&options)
	}
	if options.DefWrokers < 1 {
		options.DefWrokers = 10
	}
	if options.MaxWorkers < 1 {
		options.MaxWorkers = 20
	}
	if options.Tasktimeout < time.Millisecond {
		options.Tasktimeout = time.Second * 3
	}
	if options.IdleTimeoutSec < time.Second*5 {
		options.IdleTimeoutSec = time.Second * 5
	}
	return options
}
