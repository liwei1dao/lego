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
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		RPCConnType:  Nats,
		MaxCoroutine: 2000,
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
