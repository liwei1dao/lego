package log

import (
	"fmt"
	"github.com/liwei1dao/lego/utils/flietools"
	"os"
	"strings"
)

type (
	LogStrut interface {
		ToString() (str string)
	}
	Field struct {
		Key   string
		Value interface{}
	}
	Ilog interface {
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
