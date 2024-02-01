package openai

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
	openai "github.com/sashabaranov/go-openai"
)

type Option func(*Options)
type Options struct {
	Token string //密钥
	Model string //使用模型
	Debug bool   //日志是否开启
	Log   log.ILogger
}

//设置密钥
func SetToken(v string) Option {
	return func(o *Options) {
		o.Token = v
	}
}

//设置模型
func SetModel(v string) Option {
	return func(o *Options) {
		o.Model = v
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

func newOptions(config map[string]interface{}, opts ...Option) *Options {
	options := &Options{
		Model: openai.GPT3Dot5Turbo,
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.openai", 3))
	}
	return options
}

func newOptionsByOption(opts ...Option) *Options {
	options := &Options{
		Model: openai.GPT3Dot5Turbo,
	}
	for _, o := range opts {
		o(options)
	}
	if options.Log == nil {
		options.Log = log.NewTurnlog(options.Debug, log.Clone("sys.openai", 3))
	}
	return options
}
