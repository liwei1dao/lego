package log

import (
	"fmt"

	"github.com/liwei1dao/lego/core"
)

var (
	service    core.IService
	defaultlog Ilog
)

func OnInit(s core.IService, opt ...Option) (err error) {
	service = s
	opt = append(opt, SetFileName(fmt.Sprintf("%slog/%s/%s.log", s.GetWorkPath(), s.GetId(), s.GetType())))
	defaultlog, err = newLog(opt...)
	return
}
func NewILog(opt ...Option) (log Ilog, err error) {
	log, err = newLog(opt...)
	return
}

func Debug(msg string, fields ...Field)      { defaultlog.Debug(msg, fields...) }
func Info(msg string, fields ...Field)       { defaultlog.Info(msg, fields...) }
func Warn(msg string, fields ...Field)       { defaultlog.Warn(msg, fields...) }
func Error(msg string, fields ...Field)      { defaultlog.Error(msg, fields...) }
func Panic(msg string, fields ...Field)      { defaultlog.Panic(msg, fields...) }
func Fatal(msg string, fields ...Field)      { defaultlog.Fatal(msg, fields...) }
func Debugf(format string, a ...interface{}) { defaultlog.Debugf(format, a...) }
func Infof(format string, a ...interface{})  { defaultlog.Infof(format, a...) }
func Warnf(format string, a ...interface{})  { defaultlog.Warnf(format, a...) }
func Errorf(format string, a ...interface{}) { defaultlog.Errorf(format, a...) }
func Panicf(format string, a ...interface{}) { defaultlog.Panicf(format, a...) }
func Fatalf(format string, a ...interface{}) { defaultlog.Fatalf(format, a...) }
