package event

import (
	"reflect"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
)

/*
系统描述:进程级别的事件系统
*/
type (
	FunctionInfo struct {
		Function  reflect.Value
		Goroutine bool
	}
	ISys interface {
		Register(eId core.Event_Key, f interface{}) (err error)
		RegisterGO(eId core.Event_Key, f interface{}) (err error)
		RemoveEvent(eId core.Event_Key, f interface{}) (err error)
		TriggerEvent(eId core.Event_Key, agr ...interface{})
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

//注册同步事件处理函数
func Register(eId core.Event_Key, f interface{}) (err error) {
	return defsys.Register(eId, f)
}

//注册异步事件处理函数
func RegisterGO(eId core.Event_Key, f interface{}) (err error) {
	return defsys.Register(eId, f)
}

//移除事件
func RemoveEvent(eId core.Event_Key, f interface{}) (err error) {
	return defsys.RemoveEvent(eId, f)
}

//触发事件
func TriggerEvent(eId core.Event_Key, agr ...interface{}) {
	if defsys != nil {
		defsys.TriggerEvent(eId, agr...)
	} else {
		log.Warnf("event no start")
	}
}
