package timer

import (
	"fmt"
	"lego/core"
	"lego/sys/timer/wtimer"
)

var (
	timer ITimer
)

func OnInit(s core.IService, opt ...Option) (err error) {
	option := newOptions(opt...)
	timer = wtimer.NewTimer(option.Inteval)
	return
}

func Add(inteval uint32, handler func(string, ...interface{}), args ...interface{}) (tkey string, err error) {
	if timer == nil {
		return "", fmt.Errorf("timer 系统未初始化")
	}
	return timer.Add(inteval, handler, args...), nil
}

func Remove(key string) (err error) {
	if timer == nil {
		return fmt.Errorf("timer 系统未初始化")
	}
	timer.Remove(key)
	return
}
