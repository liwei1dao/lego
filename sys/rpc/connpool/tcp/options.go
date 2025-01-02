package tcp

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	ListterAddr       string //监听地址
	KeepAlivePeriod   int    //保持连接时间
	ConnectionTimeout int    //连接超时时间
	ReadTimeout       int    //读取超时时间
	WriteTimeout      int    //写入超时时间
	Debug             bool   //日志是否开启
	Log               log.ILogger
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
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.connpool_tcp", 3))
	}

	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{}
	for _, o := range opts {
		o(options)
	}
	return
}
