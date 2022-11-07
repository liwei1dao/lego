package log

var AllLevels = []Loglevel{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
}

type (
	Field struct {
		Key   string
		Value interface{}
	}
	Fields []Field
	Ilogf  interface {
		Debugf(format string, args ...interface{})
		Infof(format string, args ...interface{})
		Printf(format string, args ...interface{})
		Warnf(format string, args ...interface{})
		Errorf(format string, args ...interface{})
		Fatalf(format string, args ...interface{})
		Panicf(format string, args ...interface{})
	}
	IlogIn interface {
		Debugln(args ...interface{})
		Infoln(args ...interface{})
		Println(args ...interface{})
		Warnln(args ...interface{})
		Errorln(args ...interface{})
		Fatalln(args ...interface{})
		Panicln(args ...interface{})
	}
	ILog interface {
		Debug(msg string, args ...Field)
		Info(msg string, args ...Field)
		Print(msg string, args ...Field)
		Warn(msg string, args ...Field)
		Error(msg string, args ...Field)
		Fatal(msg string, args ...Field)
		Panic(msg string, args ...Field)
	}

	ILogger interface {
		SetName(name string)
		Enabled(lvl Loglevel) bool
		Ilogf
		IlogIn
		ILog
	}
	ISys interface {
		Clone(name string, skip int) ILogger
		ILogger
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}
func Clone(name string, skip int) ILogger {
	if defsys != nil {
		return defsys.Clone(name, skip)
	}
	return nil
}
func Debug(msg string, args ...Field) {
	defsys.Debug(msg, args...)
}
func Info(msg string, args ...Field) {
	defsys.Info(msg, args...)
}
func Warn(msg string, args ...Field) {
	defsys.Warn(msg, args...)
}
func Error(msg string, args ...Field) {
	defsys.Error(msg, args...)
}
func Fatal(msg string, args ...Field) {
	defsys.Fatal(msg, args...)
}
func Panic(msg string, args ...Field) {
	defsys.Panic(msg, args...)
}
func Debugf(format string, args ...interface{}) {
	defsys.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	defsys.Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	defsys.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	defsys.Errorf(format, args...)
}
func Fatalf(format string, args ...interface{}) {
	defsys.Fatalf(format, args...)
}
func Panicf(format string, args ...interface{}) {
	defsys.Panicf(format, args...)
}
func Debugln(args ...interface{}) {
	defsys.Debugln(args...)
}
func Infoln(args ...interface{}) {
	defsys.Infoln(args...)
}
func Println(args ...interface{}) {
	defsys.Println(args...)
}
func Warnln(args ...interface{}) {
	defsys.Warnln(args...)
}
func Errorln(args ...interface{}) {
	defsys.Errorln(args...)
}
func Fatalln(args ...interface{}) {
	defsys.Fatalln(args...)
}
func Panicln(args ...interface{}) {
	defsys.Panicln(args...)
}
