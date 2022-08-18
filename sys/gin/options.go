package gin

import (
	"errors"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	ListenPort int    //监听端口
	CertFile   string //tls文件
	KeyFile    string //tls文件
	Debug      bool   //日志是否开启
	Log        log.ILogger
}

func SetListenPort(v int) Option {
	return func(o *Options) {
		o.ListenPort = v
	}
}

func SetCertFile(v string) Option {
	return func(o *Options) {
		o.CertFile = v
	}
}

func SetKeyFile(v string) Option {
	return func(o *Options) {
		o.KeyFile = v
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
	options = &Options{
		ListenPort: 8080,
		CertFile:   "",
		KeyFile:    "",
		Debug:      true,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.gin", 2)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		ListenPort: 8080,
		CertFile:   "",
		KeyFile:    "",
		Debug:      true,
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.gin", 2)); options.Log == nil {
		err = errors.New("log is nil")
	}
	return
}
