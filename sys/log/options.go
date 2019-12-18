package log

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
	fileName  string
	loglevel  Loglevel
	debugmode bool
	loglayer  int //日志堆栈信息打印层级
}

func SetFileName(v string) Option {
	return func(o *Options) {
		o.fileName = v
	}
}

func SetLoglevel(v Loglevel) Option {
	return func(o *Options) {
		o.loglevel = v
	}
}

func SetDebugMode(v bool) Option {
	return func(o *Options) {
		o.debugmode = v
	}
}

func SetLoglayer(v int) Option {
	return func(o *Options) {
		o.loglayer = v
	}
}
func newOptions(opts ...Option) Options {
	opt := Options{
		loglevel:  WarnLevel,
		debugmode: false,
		loglayer:  2,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
