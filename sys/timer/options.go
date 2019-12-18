package timer

import (
	"time"
)

type Option func(*Options)
type Options struct {
	Inteval time.Duration
}

func Inteval(v time.Duration) Option {
	return func(o *Options) {
		o.Inteval = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	if opt.Inteval == 0 {
		opt.Inteval = time.Second
	}
	return opt
}
