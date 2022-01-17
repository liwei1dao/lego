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
		Consul_RegisterInterval int //定期注册
		Consul_RegisterTTL      int
		Nacos_NamespaceId       string
		Nacos_NacosAddr         string
		Nacos_Port              uint64
		Nacos_TimeoutMs         uint64 //连接超时 ms
		Nacos_BeatInterval      int64  //心跳间隔 ms
		Nacos_RegisterTTL       int
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
func SetConsul_RegisterInterval(v int) Option {
	return func(o *Options) {
		o.Consul_RegisterInterval = v
	}
}
func SetConsul_RegisterTTL(v int) Option {
	return func(o *Options) {
		o.Consul_RegisterTTL = v
	}
}
func SetNacos_NamespaceId(v string) Option {
	return func(o *Options) {
		o.Nacos_NamespaceId = v
	}
}
func SetNacos_NacosAddr(v string) Option {
	return func(o *Options) {
		o.Nacos_NacosAddr = v
	}
}

func SetNacos_Port(v uint64) Option {
	return func(o *Options) {
		o.Nacos_Port = v
	}
}

func SetNacos_TimeoutMs(v uint64) Option {
	return func(o *Options) {
		o.Nacos_TimeoutMs = v
	}
}
func SetNacos_BeatInterval(v int64) Option {
	return func(o *Options) {
		o.Nacos_BeatInterval = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Consul_RegisterInterval: 10,
		Consul_RegisterTTL:      15,
		Nacos_TimeoutMs:         10000,
		Nacos_BeatInterval:      5000,
		Nacos_RegisterTTL:       8,
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
		Consul_RegisterInterval: 10,
		Consul_RegisterTTL:      30,
		Nacos_TimeoutMs:         10000,
		Nacos_BeatInterval:      5000,
		Nacos_RegisterTTL:       8,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
