package workerpools

import (
	"time"
)

type Option func(*Options)
type Options struct {
	defWrokers  int
	maxWorkers  int
	tasktimeout time.Duration
}

func SetDefWorkers(v int) Option {
	return func(o *Options) {
		o.defWrokers = v
	}
}

func SetMaxWorkers(v int) Option {
	return func(o *Options) {
		o.maxWorkers = v
	}
}

func SetTaskTimeOut(v time.Duration) Option {
	return func(o *Options) {
		o.tasktimeout = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		defWrokers:  10,
		maxWorkers:  20,
		tasktimeout: time.Second * 3,
	}
	for _, o := range opts {
		o(&opt)
	}
	if opt.maxWorkers < 1 {
		opt.maxWorkers = 10
	}
	if opt.maxWorkers < 1 {
		opt.maxWorkers = 20
	}
	if opt.tasktimeout < time.Millisecond {
		opt.tasktimeout = time.Second * 3
	}
	return opt
}
