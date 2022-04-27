package livego

import "github.com/liwei1dao/lego/utils/mapstructure"

type Option func(*Options)
type Options struct {
	Appname         string   `mapstructure:"appname"`
	Live            bool     `mapstructure:"live"`
	Hls             bool     `mapstructure:"hls"`
	Flv             bool     `mapstructure:"flv"`
	Api             bool     `mapstructure:"api"`
	StaticPush      []string `mapstructure:"static_push"`
	RTMPNoAuth      bool     `mapstructure:"rtmp_noauth"`
	EnableTLSVerify bool     `mapstructure:"enable_tls_verify"`
	Timeout         int      //单位秒
	ConnBuffSzie    int      //连接对象读写缓存
	Debug           bool     //日志是否开启
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Timeout:      5,
		ConnBuffSzie: 4 * 1024,
		Debug:        true,
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
		Timeout:      5,
		ConnBuffSzie: 4 * 1024,
		Debug:        true,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
