package log

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Loglevel int8

const (
	DebugLevel Loglevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

type Option func(*Options)
type Options struct {
	FileName  string   //日志文件名包含
	Loglevel  Loglevel //日志输出级别
	Debugmode bool     //是否debug模式
	Loglayer  int      //日志堆栈信息打印层级
}

func SetFileName(v string) Option {
	return func(o *Options) {
		o.FileName = v
	}
}

func SetLoglevel(v Loglevel) Option {
	return func(o *Options) {
		o.Loglevel = v
	}
}

func SetDebugMode(v bool) Option {
	return func(o *Options) {
		o.Debugmode = v
	}
}

func SetLoglayer(v int) Option {
	return func(o *Options) {
		o.Loglayer = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Loglevel:  WarnLevel,
		Debugmode: false,
		Loglayer:  2,
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
		Loglevel:  WarnLevel,
		Debugmode: false,
		Loglayer:  2,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
