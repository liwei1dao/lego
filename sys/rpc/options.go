package rpc

import (
	"github.com/nats-io/nats.go"
)

type Option func(*Options)
type Options struct {
	sId          string
	NatsAddr     string
	MaxCoroutine int
	RpcExpired   int
}

func Id(v string) Option {
	return func(o *Options) {
		o.sId = v
	}
}
func NatsAddr(v string) Option {
	return func(o *Options) {
		o.NatsAddr = v
	}
}

func MaxCoroutine(v int) Option {
	return func(o *Options) {
		o.MaxCoroutine = v
	}
}

func RpcExpired(v int) Option {
	return func(o *Options) {
		o.RpcExpired = v
	}
}

func newOptions(opts ...Option) *Options {
	opt := Options{
		MaxCoroutine: 10000,
		RpcExpired:   5,
	}
	for _, o := range opts {
		o(&opt)
	}
	if len(opt.NatsAddr) == 0 {
		opt.NatsAddr = nats.DefaultURL
	}
	return &opt
}
