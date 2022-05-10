package rpc

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	RPCConnType   RPCConnType
	Listener      RPCListener
	ClusterTag    string
	ServiceId     string
	MaxCoroutine  int
	RpcExpired    int
	Nats_Addr     string
	Kafka_Host    []string
	Kafka_Version string
	Debug         bool //日志是否开启
	Log           log.ILog
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

func SetDebug(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}
func SetLog(v log.ILog) Option {
	return func(o *Options) {
		o.Log = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		RPCConnType:   Nats,
		MaxCoroutine:  2000,
		RpcExpired:    5,
		Kafka_Version: "1.0.0",
		Debug:         true,
		Log:           log.Clone(log.SetLoglayer(2)),
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
	options := Options{
		RPCConnType:   Nats,
		MaxCoroutine:  2000,
		RpcExpired:    5,
		Kafka_Version: "1.0.0",
		Debug:         true,
		Log:           log.Clone(log.SetLoglayer(2)),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
