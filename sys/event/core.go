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
