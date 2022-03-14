package log_test

import (
	"testing"

	"github.com/liwei1dao/lego/sys/log"
)

func Test_sys(t *testing.T) {
	if err := log.OnInit(map[string]interface{}{
		"FileName":  "./test.log",
		"Loglevel":  log.FatalLevel,
		"Debugmode": true,
		"Loglayer":  2,
	}); err == nil {
		log.Infof("测试日志接口代码!")
	}
}
