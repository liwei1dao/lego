package oss

type Option func(*Options)
type Options struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
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

func SetBucketName(v string) Option {
	return func(o *Options) {
		o.BucketName = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
