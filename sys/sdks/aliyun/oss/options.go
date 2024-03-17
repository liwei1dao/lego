package oss

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	BucketName      string
}

func SetEndpoint(v string) Option {
	return func(o *Options) {
		o.Endpoint = v
	}
}

func SetAccessKeyId(v string) Option {
	return func(o *Options) {
		o.AccessKeyId = v
	}
}

func SetAccessKeySecret(v string) Option {
	return func(o *Options) {
		o.AccessKeySecret = v
	}
}
func SetSecurityToken(v string) Option {
	return func(o *Options) {
		o.SecurityToken = v
	}
}

func SetBucketName(v string) Option {
	return func(o *Options) {
		o.BucketName = v
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
	if options.Endpoint == "" || options.AccessKeyId == "" || options.AccessKeySecret == "" || options.BucketName == "" {
		log.Panicf("start oss Missing necessary configuration : Endpoint:%s AccessKeyId:%s AccessKeySecret:%s BucketName:%s", options.Endpoint, options.AccessKeyId, options.AccessKeySecret, options.BucketName)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	if options.Endpoint == "" || options.AccessKeyId == "" || options.AccessKeySecret == "" || options.BucketName == "" {
		log.Panicf("start oss Missing necessary configuration : Endpoint:%s AccessKeyId:%s AccessKeySecret:%s BucketName:%s", options.Endpoint, options.AccessKeyId, options.AccessKeySecret, options.BucketName)
	}
	return options
}
