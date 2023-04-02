package engine

import (
	"github.com/liwei1dao/lego/sys/log"
)

type Option func(*Options)
type Options struct {
	MultipartMemory int64 //文件上传最大尺寸
	Log             log.ILogger
}

func SetMultipartMemory(v int64) Option {
	return func(o *Options) {
		o.MultipartMemory = v
	}
}

func SetLog(log log.ILogger) Option {
	return func(o *Options) {
		o.Log = log
	}
}

func newOptions(opts ...Option) (options *Options, err error) {
	options = &Options{
		MultipartMemory: defaultMultipartMemory,
	}
	for _, o := range opts {
		o(options)
	}
	return
}
