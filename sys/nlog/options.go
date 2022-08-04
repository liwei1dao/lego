package log

import (
	"errors"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type LogEncoder int8

const (
	TextEncoder LogEncoder = iota
	JSONEncoder
)

type Option func(*Options)
type Options struct {
	FileName     string     //日志文件名包含
	Loglevel     Loglevel   //日志输出级别
	ReportCaller bool       //是否输出堆栈信息
	CallerSkip   int        //堆栈深度
	Encoder      LogEncoder //日志输出样式
	RotationTime int32      //日志分割时间 单位 小时
	MaxAgeTime   int32      //日志最大保存时间 单位天
}

///日志文件名包含
func SetFileName(v string) Option {
	return func(o *Options) {
		o.FileName = v
	}
}

///日志输出级别 debug info warning error fatal panic
func SetLoglevel(v Loglevel) Option {
	return func(o *Options) {
		o.Loglevel = v
	}
}

func SetReportCaller(v bool) Option {
	return func(o *Options) {
		o.ReportCaller = v
	}
}
func SetCallerSkip(v int) Option {
	return func(o *Options) {
		o.CallerSkip = v
	}
}

///日志输出样式
func SetEncoder(v LogEncoder) Option {
	return func(o *Options) {
		o.Encoder = v
	}
}

///日志分割时间 单位 小时
func SetRotationTime(v int32) Option {
	return func(o *Options) {
		o.RotationTime = v
	}
}

///日志最大保存时间 单位天
func SetMaxAgeTime(v int32) Option {
	return func(o *Options) {
		o.RotationTime = v
	}
}
func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		FileName:     "log",
		Loglevel:     DebugLevel,
		RotationTime: 24,
		MaxAgeTime:   7,
		Encoder:      TextEncoder,
		CallerSkip:   0,
		ReportCaller: true,
	}
	if config != nil {
		mapstructure.Decode(config, options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.RotationTime <= 0 || options.MaxAgeTime <= 0 {
		err = errors.New("log options RotationTime or MaxAgeTime is zero!")
		return
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		FileName:     "log.log",
		Loglevel:     DebugLevel,
		RotationTime: 24,
		MaxAgeTime:   7,
		Encoder:      TextEncoder,
		CallerSkip:   0,
		ReportCaller: true,
	}
	for _, o := range opts {
		o(options)
	}
	return
}
