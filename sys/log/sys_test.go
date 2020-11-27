package log_test

import (
	"testing"
)

func Test_sys(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"FileName":  "./test.log",
		"Loglevel":  FatalLevel,
		"Debugmode": true,
		"Loglayer":  2,
	}); err == nil {
		Infof("测试日志接口代码!")
	}
}
