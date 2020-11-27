package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/liwei1dao/lego/utils/flietools"
)

type (
	LogStrut interface {
		ToString() (str string)
	}
	Field struct {
		Key   string
		Value interface{}
	}
	ILog interface {
		Debug(msg string, fields ...Field)
		Info(msg string, fields ...Field)
		Warn(msg string, fields ...Field)
		Error(msg string, fields ...Field)
		Panic(msg string, fields ...Field)
		Fatal(msg string, fields ...Field)
		Debugf(format string, a ...interface{})
		Infof(format string, a ...interface{})
		Warnf(format string, a ...interface{})
		Errorf(format string, a ...interface{})
		Panicf(format string, a ...interface{})
		Fatalf(format string, a ...interface{})
	}
)

var (
	defsys ILog
)

func OnInit(config map[string]interface{}) (err error) {
	defsys, err = newSys(newOptionsByConfig(config))
	return
}
func NewSys(option ...Option) (sys Ilog, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func Debug(msg string, fields ...Field)      { defsys.Debug(msg, fields...) }
func Info(msg string, fields ...Field)       { defsys.Info(msg, fields...) }
func Warn(msg string, fields ...Field)       { defsys.Warn(msg, fields...) }
func Error(msg string, fields ...Field)      { defsys.Error(msg, fields...) }
func Panic(msg string, fields ...Field)      { defsys.Panic(msg, fields...) }
func Fatal(msg string, fields ...Field)      { defsys.Fatal(msg, fields...) }
func Debugf(format string, a ...interface{}) { defsys.Debugf(format, a...) }
func Infof(format string, a ...interface{})  { defsys.Infof(format, a...) }
func Warnf(format string, a ...interface{})  { defsys.Warnf(format, a...) }
func Errorf(format string, a ...interface{}) { defsys.Errorf(format, a...) }
func Panicf(format string, a ...interface{}) { defsys.Panicf(format, a...) }
func Fatalf(format string, a ...interface{}) { defsys.Fatalf(format, a...) }

//创建日志文件
func createlogfile(logpath string) error {
	logdir := string(logpath[0:strings.LastIndex(logpath, "/")])
	if !flietools.IsExist(logdir) {
		err := os.MkdirAll(logdir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("创建日志路径失败 1" + err.Error())
		}
	}
	return nil
}
