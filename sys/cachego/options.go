package cachego

import (
	"errors"
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Expiration      time.Duration
	CleanupInterval time.Duration
	Debug           bool //日志是否开启
	Log             log.ILogger
}

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		Expiration:      -1,
		CleanupInterval: 0,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}

	if options.Debug && options.Log == nil {
		options.Log = log.Clone("sys.cachego", 2)
	}
	if options.Log = log.NewTurnlog(options.Debug, options.Log); options.Log == nil {
		err = errors.New("log is nil")
	}

	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		Expiration:      -1,
		CleanupInterval: 0,
	}
	for _, o := range opts {
		o(options)
	}
	if options.Debug && options.Log == nil {
		options.Log = log.Clone("sys.cachego", 2)
	}
	if options.Log = log.NewTurnlog(options.Debug, options.Log); options.Log == nil {
		err = errors.New("log is nil")
	}

	return
}
