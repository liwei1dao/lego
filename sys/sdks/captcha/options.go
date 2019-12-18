package sdk_aliyun

type Option func(*Options)
type Options struct {
	Aliyun CaptchaService
	Email  CaptchaService
}

func SetAliyun(v CaptchaService) Option {
	return func(o *Options) {
		o.Aliyun = v
	}
}

func SetEmail(v CaptchaService) Option {
	return func(o *Options) {
		o.Email = v
	}
}

func newOptions(opts ...Option) *Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return &opt
}
