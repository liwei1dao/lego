package rpc

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
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

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		NatsAddr:     nats.DefaultURL,
		MaxCoroutine: 10000,
		RpcExpired:   5,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.NatsAddr) == 0 {
		options.NatsAddr = nats.DefaultURL
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		NatsAddr:     nats.DefaultURL,
		MaxCoroutine: 10000,
		RpcExpired:   5,
	}
	for _, o := range opts {
		o(&options)
	}
	if len(options.NatsAddr) == 0 {
		options.NatsAddr = nats.DefaultURL
	}
	return options
}
