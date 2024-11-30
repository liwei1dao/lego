package tcppool

import (
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Debug             bool //日志是否开启
	Log               log.ILogger
	Prot              int32         //监听端口
	ConnectionTimeout time.Duration //连接超时
	ReadTimeout       time.Duration //读取超时
	WriteTimeout      time.Duration //写入超时
	KeepAlivePeriod   time.Duration //保持活跃时期
	Username          string        //用户名
	Password          string        //密码
	Vsersion          string        //版本
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.tcppool", 3))
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.tcppool", 3))
	}
	return options
}
