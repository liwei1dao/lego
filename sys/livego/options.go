package livego

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	Appname         string   `mapstructure:"appname"`
	Hls             bool     `mapstructure:"hls"`
	Flv             bool     `mapstructure:"flv"`
	Api             bool     `mapstructure:"api"`
	APIAddr         string   `mapstructure:"api_addr"`
	JWTSecret       string   `mapstructure:"jwt_secret"`
	JWTAlgorithm    string   `mapstructure:"jwt_algorithm"`
	StaticPush      []string `mapstructure:"static_push"`
	RTMPAddr        string   `mapstructure:"rtmp_addr"`
	RTMPNoAuth      bool     `mapstructure:"rtmp_noauth"`
	FLVDir          string   `mapstructure:"flv_dir"`
	FLVArchive      bool     `mapstructure:"flv_archive"`
	EnableTLSVerify bool     `mapstructure:"enable_tls_verify"`
	Timeout         int      //单位秒
	ConnBuffSzie    int      //连接对象读写缓存
	Debug           bool     //日志是否开启
	Log             log.ILog
}

func SetAppname(v string) Option {
	return func(o *Options) {
		o.Appname = v
	}
}
func SetHls(v bool) Option {
	return func(o *Options) {
		o.Hls = v
	}
}
func SetFlv(v bool) Option {
	return func(o *Options) {
		o.Flv = v
	}
}
func SetApi(v bool) Option {
	return func(o *Options) {
		o.Flv = v
	}
}
func SetStaticPush(v []string) Option {
	return func(o *Options) {
		o.StaticPush = v
	}
}

func SetRTMPAddr(v string) Option {
	return func(o *Options) {
		o.RTMPAddr = v
	}
}

func SetRTMPNoAuth(v bool) Option {
	return func(o *Options) {
		o.RTMPNoAuth = v
	}
}

func SetFLVDir(v string) Option {
	return func(o *Options) {
		o.FLVDir = v
	}
}

func SetFLVArchive(v bool) Option {
	return func(o *Options) {
		o.FLVArchive = v
	}
}

func SetEnableTLSVerify(v bool) Option {
	return func(o *Options) {
		o.EnableTLSVerify = v
	}
}

func SetTimeout(v int) Option {
	return func(o *Options) {
		o.Timeout = v
	}
}

func SetConnBuffSzie(v int) Option {
	return func(o *Options) {
		o.ConnBuffSzie = v
	}
}

func SetDebug(v bool) Option {
	return func(o *Options) {
		o.Debug = v
	}
}
func SetLog(v log.ILog) Option {
	return func(o *Options) {
		o.Log = v
	}
}
func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Appname:      "live",
		RTMPAddr:     ":1935",
		Timeout:      5,
		ConnBuffSzie: 4 * 1024,
		Debug:        true,
		Log:          log.Clone(log.SetLoglayer(2)),
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		Appname:      "live",
		RTMPAddr:     ":1935",
		Timeout:      5,
		ConnBuffSzie: 4 * 1024,
		Debug:        true,
		Log:          log.Clone(log.SetLoglayer(2)),
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
