package bot

import (
	"errors"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Debug    bool //日志是否开启
	Log      log.ILogger
	ApiToken string
}

func SetApiToken(v string) Option {
	return func(o *Options) {
		o.ApiToken = v
	}
}

func SetDebug(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}
func SetLog(v log.ILogger) Option {
	return func(o *Options) {
		o.Log = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}

	if options.Log == nil {
		if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.telegram.bot", 2)); options.Log == nil {
			err = errors.New("log is nil")
		}
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{}
	for _, o := range opts {
		o(options)
	}
	if options.Debug && options.Log == nil {
		if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.telegram.bot", 2)); options.Log == nil {
			err = errors.New("log is nil")
		}
	} else if !options.Debug && options.Log == nil {
		if options.Log = log.NewTurnlog(options.Debug, nil); options.Log == nil {
			err = errors.New("log is nil")
		}
	}
	return
}
