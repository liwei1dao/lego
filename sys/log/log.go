package log

import (
	"fmt"
	"lego/core"

	"go.uber.org/zap"
)

var (
	service core.IService
	log     Ilog
)

func init() {
	tlog, _ := zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(2))
	log = &Logger{
		tlog: tlog,
		log:  tlog.Sugar(),
	}
}
func OnInit(s core.IService, opt ...Option) (err error) {
	service = s
	opt = append(opt, SetFileName(fmt.Sprintf("%sbin/log/%s/%s.log", s.GetWorkPath(), s.GetId(), s.GetType())))
	log, err = newLog(opt...)
	return
}
func NewILog(opt ...Option) (log Ilog, err error) {
	log, err = newLog(opt...)
	return
}

func Debug(msg string, fields ...Field)      { log.Debug(msg, fields...) }
func Info(msg string, fields ...Field)       { log.Info(msg, fields...) }
func Warn(msg string, fields ...Field)       { log.Warn(msg, fields...) }
func Error(msg string, fields ...Field)      { log.Error(msg, fields...) }
func Panic(msg string, fields ...Field)      { log.Panic(msg, fields...) }
func Fatal(msg string, fields ...Field)      { log.Fatal(msg, fields...) }
func Debugf(format string, a ...interface{}) { log.Debugf(format, a...) }
func Infof(format string, a ...interface{})  { log.Infof(format, a...) }
func Warnf(format string, a ...interface{})  { log.Warnf(format, a...) }
func Errorf(format string, a ...interface{}) { log.Errorf(format, a...) }
func Panicf(format string, a ...interface{}) { log.Panicf(format, a...) }
func Fatalf(format string, a ...interface{}) { log.Fatalf(format, a...) }
