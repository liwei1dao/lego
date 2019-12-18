package workerpools

import "time"

type Option func(*Options)
type Options struct {
	maxWorkers  int
	tasktimeout time.Duration
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
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	if opt.tasktimeout <= 0 {
		opt.tasktimeout = time.Second * 3
	}
	return opt
}
