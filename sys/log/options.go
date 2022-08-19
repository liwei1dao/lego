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
	IsDebug      bool       //是否是开发模式
	ReportCaller Loglevel   //是否输出堆栈信息
	CallerSkip   int        //堆栈深度
	Encoder      LogEncoder //日志输出样式
	CupTimeTime  int        //日志分割时间 单位 小时
	MaxAgeTime   int        //日志最大保存时间 单位天
	MaxBackups   int        //最大备份日志个数
	Compress     bool       //是否压缩备份日志
}

///日志文件名包含
func SetFileName(v string) Option {
	return func(o *Options) {
		o.FileName = v
	}
}

//是否是开发模式
func SetIsDebug(v bool) Option {
	return func(o *Options) {
		o.IsDebug = v
	}
}

///日志输出级别 debug info warning error fatal panic
func SetLoglevel(v Loglevel) Option {
	return func(o *Options) {
		o.Loglevel = v
	}
}

func SetReportCaller(v Loglevel) Option {
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
func SetRotationTime(v int) Option {
	return func(o *Options) {
		o.CupTimeTime = v
	}
}

///日志最大保存时间 单位天
func SetMaxAgeTime(v int) Option {
	return func(o *Options) {
		o.MaxAgeTime = v
	}
}

///日志备份最大文件数
func SetMaxBackups(v int) Option {
	return func(o *Options) {
		o.MaxBackups = v
	}
}

///是否压缩备份日志
func SetCompress(v bool) Option {
	return func(o *Options) {
		o.Compress = v
	}
}
func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{
		FileName:     "log",
		Loglevel:     DebugLevel,
		CupTimeTime:  24,
		MaxAgeTime:   7,
		MaxBackups:   250,
		Compress:     false,
		Encoder:      TextEncoder,
		CallerSkip:   3,
		ReportCaller: ErrorLevel,
		IsDebug:      true,
	}
	if config != nil {
		mapstructure.Decode(config, options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.CupTimeTime <= 0 || options.MaxAgeTime <= 0 {
		err = errors.New("log options RotationTime or MaxAgeTime is zero!")
		return
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{
		FileName:     "log.log",
		Loglevel:     DebugLevel,
		CupTimeTime:  24,
		MaxAgeTime:   7,
		MaxBackups:   250,
		Compress:     false,
		Encoder:      TextEncoder,
		CallerSkip:   3,
		ReportCaller: ErrorLevel,
		IsDebug:      true,
	}
	for _, o := range opts {
		o(options)
	}
	return
}
