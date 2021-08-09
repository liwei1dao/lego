package registry

import (
	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type RegistryType uint8

const (
	Registry_Consul RegistryType = iota
	Registry_Nacos
)

type (
	Option  func(*Options)
	Options struct {
		RegistryType            RegistryType
		Service                 base.IClusterService
		Listener                IListener
		Consul_Addr             string
		Consul_Timeout          int
		Consul_RegisterInterval int //定期注册
		Consul_RegisterTTL      int
		Nacos_NamespaceId       string
		Nacos_NacosAddr         string
		Nacos_Port              uint64
		Nacos_TimeoutMs         uint64 //连接超时 ms
		Nacos_BeatInterval      int64  //心跳间隔 ms
	}
)

func SetRegistryType(v RegistryType) Option {
	return func(o *Options) {
		o.RegistryType = v
	}
}
func SetService(v base.IClusterService) Option {
	return func(o *Options) {
		o.Service = v
	}
}
func SetListener(v IListener) Option {
	return func(o *Options) {
		o.Listener = v
	}
}
func SetConsul_Addr(v string) Option {
	return func(o *Options) {
		o.Consul_Addr = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{}
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
