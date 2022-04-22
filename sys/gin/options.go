package gin

import "github.com/liwei1dao/lego/utils/mapstructure"

type Option func(*Options)
type Options struct {
	ListenPort int  //监听端口
	Debug      bool //日志是否开启
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		ListenPort: 8080,
		Debug:      false,
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
		ListenPort: 8080,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
