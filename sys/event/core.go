package event

import (
	"reflect"

	"github.com/liwei1dao/lego/core"
)

type (
	FunctionInfo struct {
		Function  reflect.Value
		Goroutine bool
	}
	IEventSys interface {
		Register(eId core.Event_Key, f interface{}) (err error)
		RegisterGO(eId core.Event_Key, f interface{}) (err error)
		RemoveEvent(eId core.Event_Key, f interface{}) (err error)
		TriggerEvent(eId core.Event_Key, agr ...interface{})
	}
)

var (
	defsys IEventSys
)

func OnInit(config map[string]interface{}) (err error) {
	defsys, err = newSys(newOptionsByConfig(config))
	return
}

func NewSys(opts ...Option) (err error) {
	defsys, err = newSys(newOptionsByOption(opts...))
	return
}

func Register(eId core.Event_Key, f interface{}) (err error) {
	return defsys.Register(eId, f)
}

func RegisterGO(eId core.Event_Key, f interface{}) (err error) {
	return defsys.Register(eId, f)
}

func RemoveEvent(eId core.Event_Key, f interface{}) (err error) {
	return defsys.RemoveEvent(eId, f)
}

func TriggerEvent(eId core.Event_Key, agr ...interface{}) {
	defsys.TriggerEvent(eId, agr...)
}
