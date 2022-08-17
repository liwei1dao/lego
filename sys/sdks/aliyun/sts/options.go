package sts

import "github.com/liwei1dao/lego/utils/mapstructure"

type Option func(*Options)
type Options struct {
	RegionId        string
	AccessKeyId     string
	AccessKeySecret string
}

//注册点 cn-shenzhen
func SetRegionId(v string) Option {
	return func(o *Options) {
		o.RegionId = v
	}
}

//头像->访问控制->用户->用户 AccessKey
func SetAccessKeyId(v string) Option {
	return func(o *Options) {
		o.AccessKeyId = v
	}
}

//头像->访问控制->用户->用户 AccessKeySecret
func SetAccessKeySecret(v string) Option {
	return func(o *Options) {
		o.AccessKeySecret = v
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
