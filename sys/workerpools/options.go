package workerpools

import (
	"time"

	"github.com/liwei1dao/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	DefWrokers     int
	MaxWorkers     int
	Tasktimeout    time.Duration //任务执行操超时间
	IdleTimeoutSec time.Duration  //超时释放空闲工作人员
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

func newOptionsByConfig(config map[string]interface{}) Options {
	options := Options{
		DefWrokers:  10,
		MaxWorkers:  20,
		Tasktimeout: time.Second * 3,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	if opt.DefWrokers < 1 {
		opt.DefWrokers = 10
	}
	if opt.MaxWorkers < 1 {
		opt.MaxWorkers = 20
	}
	if opt.Tasktimeout < time.Millisecond {
		opt.Tasktimeout = time.Second * 3
	}
	if opt.IdleTimeoutSec < time.Second*5 {
		opt.IdleTimeoutSec = time.Second*5
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
	if opt.DefWrokers < 1 {
		opt.DefWrokers = 10
	}
	if opt.MaxWorkers < 1 {
		opt.MaxWorkers = 20
	}
	if opt.Tasktimeout < time.Millisecond {
		opt.Tasktimeout = time.Second * 3
	}
	if opt.IdleTimeoutSec < time.Second*5 {
		opt.IdleTimeoutSec = time.Second*5
	}
	return options
}
