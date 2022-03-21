package storage

import "github.com/liwei1dao/lego/utils/mapstructure"

type Option func(*Options)
type Options struct {
	ServiceAccountPath string
	BucketName         string
	TimeOut            int
}

func SetServiceAccountPath(v string) Option {
	return func(o *Options) {
		o.ServiceAccountPath = v
	}
}
func SetBucketName(v string) Option {
	return func(o *Options) {
		o.BucketName = v
	}
}
func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		ServiceAccountPath: "service-account.json",
		TimeOut:            5,
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
		ServiceAccountPath: "service-account.json",
		TimeOut:            5,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
