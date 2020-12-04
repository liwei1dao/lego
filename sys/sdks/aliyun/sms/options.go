package sms

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	key       string
	KeySecret string
	SignName  string
}

func Setkey(v string) Option {
	return func(o *Options) {
		o.key = v
	}
}

func KeySecret(v string) Option {
	return func(o *Options) {
		o.KeySecret = v
	}
}

func SetSignName(v string) Option {
	return func(o *Options) {
		o.SignName = v
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
	if options.key == "" || options.KeySecret == "" || options.SignName == "" {
		log.Panicf("start smsverification Missing necessary configuration : Serverhost:%s Fromemail:%s Fompasswd:%s", options.key, options.KeySecret, options.SignName)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	if options.key == "" || options.KeySecret == "" || options.SignName == "" {
		log.Panicf("start smsverification Missing necessary configuration : Serverhost:%s Fromemail:%s Fompasswd:%s", options.key, options.KeySecret, options.SignName)
	}
	return options
}
