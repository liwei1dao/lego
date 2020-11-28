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

func newOptionsByConfig(config map[string]interface{}) Options {
	options := Options{
		RedisUrl: "redis://127.0.0.1:6379/1",
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		RedisUrl: "redis://127.0.0.1:6379/1",
	}
	for _, o := range opts {
		o(&options)
	}
	return options
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
