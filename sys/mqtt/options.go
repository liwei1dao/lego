package mqtt

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	MqttUrl     string
	MqttPort    int
	ClientID    string
	UserName    string
	Password    string
	MqttFactory IMqttFactory
}

func SetMqttUrl(v string) Option {
	return func(o *Options) {
		o.MqttUrl = v
	}
}

func SetMqttPort(v int) Option {
	return func(o *Options) {
		o.MqttPort = v
	}
}

func SetClientID(v string) Option {
	return func(o *Options) {
		o.ClientID = v
	}
}

func SetUserName(v string) Option {
	return func(o *Options) {
		o.UserName = v
	}
}
func SetMPassword(v string) Option {
	return func(o *Options) {
		o.Password = v
	}
}

func SetMMqttFactory(v IMqttFactory) Option {
	return func(o *Options) {
		o.MqttFactory = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		MqttFactory: new(DefMqttFactory),
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
		MqttFactory: new(DefMqttFactory),
	}
	for _, o := range opts {
		o(&options)
	}

	return options
}
