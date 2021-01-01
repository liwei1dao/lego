package cron

type Option func(*Options)
type Options struct {
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
