package rpc

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	RPCConnType  RPCConnType
	Listener     RPCListener
	ClusterTag   string
	ServiceId    string
	MaxCoroutine int
	RpcExpired   int
	Nats_Addr    string
	Kafka_Host   []string
}

func SetClusterTag(v string) Option {
	return func(o *Options) {
		o.ClusterTag = v
	}
}
func SetServiceId(v string) Option {
	return func(o *Options) {
		o.ServiceId = v
	}
}
func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		RPCConnType:  Nats,
		MaxCoroutine: 2000,
		RpcExpired:   5,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return options
}
