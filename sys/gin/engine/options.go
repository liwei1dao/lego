package engine

import (
	"github.com/liwei1dao/lego/sys/log"
)

type Option func(*Options)
type Options struct {
	MultipartMemory       int64 //文件上传最大尺寸
	RedirectTrailingSlash bool
	Log                   log.ILogger
}

func SetMultipartMemory(v int64) Option {
	return func(o *Options) {
		o.MultipartMemory = v
	}
}

/*
	如果当前路由无法匹配，但启用自动重定向
	带有（不带）尾部斜杠的路径的处理程序存在。
	例如，如果 /foo/ 被请求，但路由只存在于 /foo，则
	对于 GET 请求，客户端被重定向到 /foo，http 状态码为 301
	对于所有其他请求方法，则为 307。
*/
func SetRedirectTrailingSlash(v bool) Option {
	return func(o *Options) {
		o.RedirectTrailingSlash = v
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
