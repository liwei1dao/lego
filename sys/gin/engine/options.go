package engine

import (
	"github.com/liwei1dao/lego/sys/log"
)

type Option func(*Options)
type Options struct {
	Log                   log.ILogger
	RedirectTrailingSlash bool
}

func SeRedirectTrailingSlash(v bool) Option {
	return func(o *Options) {
		o.RedirectTrailingSlash = v
	}
}
func SetLog(v log.ILogger) Option {
	return func(o *Options) {
		o.Log = v
	}
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		RedirectTrailingSlash: true,
	}
	for _, o := range opts {
		o(options)
	}

	return
}
