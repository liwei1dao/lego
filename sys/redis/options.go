package redis

import "time"

type Option func(*Options)
type Options struct {
	RedisUrl string
}

func RedisUrl(v string) Option {
	return func(o *Options) {
		o.RedisUrl = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{
		RedisUrl: "redis://127.0.0.1:6379/1",
	}
	for _, o := range opts {
		o(&opt)
	}
	if len(opt.RedisUrl) == 0 {
		opt.RedisUrl = "redis://127.0.0.1:6379/1"
	}
	return opt
}

type RMutexOption func(*RMutexOptions)
type RMutexOptions struct {
	expiry int
	delay  time.Duration
}

func SetExpiry(v int) RMutexOption {
	return func(o *RMutexOptions) {
		o.expiry = v
	}
}
func Setdelay(v time.Duration) RMutexOption {
	return func(o *RMutexOptions) {
		o.delay = v
	}
}

func newRMutexOptions(opts ...RMutexOption) RMutexOptions {
	opt := RMutexOptions{
		expiry: 3,
		delay:  time.Millisecond * 8,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
