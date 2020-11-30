package registry

import (
	"time"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type IListener interface {
	FindServiceHandlefunc(snode ServiceNode)
	UpDataServiceHandlefunc(snode ServiceNode)
	LoseServiceHandlefunc(sId string)
}

type Option func(*Options)
type Options struct {
	ConsulAddr       string
	Timeout          time.Duration
	RegisterInterval time.Duration
	RegisterTTL      time.Duration
	Service          base.IClusterService
	Listener         IListener
}

func SetAddress(v string) Option {
	return func(o *Options) {
		o.ConsulAddr = v
	}
}

func SetService(v base.IClusterService) Option {
	return func(o *Options) {
		o.Service = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		RegisterInterval: time.Second * 10,
		RegisterTTL:      time.Second * 30,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if options.Service == nil {
		log.Panicf("start registry Missing necessary configuration : Service is nul")
	} else {
		if listener, ok := options.Service.(IListener); !ok {
			log.Warnf("start registry No Configuration Listener")
		} else {
			options.Listener = listener
		}
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		RegisterInterval: time.Second * 10,
		RegisterTTL:      time.Second * 30,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
