package sitemap

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	DefaultHost string //网址
	Filename    string //文件地址
	Compress    bool   //是否压缩
	Pretty      bool
	MaxLinks    int  //最大url数量
	Debug       bool //日志是否开启
	Log         log.ILogger
}

///网址
func SetDefaultHost(v string) Option {
	return func(o *Options) {
		o.DefaultHost = v
	}
}

///文件地址
func SetFilename(v string) Option {
	return func(o *Options) {
		o.Filename = v
	}
}

///是否压缩
func SetCompress(v bool) Option {
	return func(o *Options) {
		o.Compress = v
	}
}

func SetPretty(v bool) Option {
	return func(o *Options) {
		o.Pretty = v
	}
}

func SetMaxLinks(v int) Option {
	return func(o *Options) {
		o.MaxLinks = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		MaxLinks: 50000,
		Filename: "./sitemap.xml",
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.sitemap", 3))
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		MaxLinks: 50000,
		Filename: "./sitemap.xml",
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.sitemap", 3))
	}

	return
}
