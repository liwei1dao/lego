package event

import (
	"fmt"
	"reflect"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
)

func newSys(options Options) (sys *EventSys, err error) {
	sys = &EventSys{
		functions: make(map[core.Event_Key][]*FunctionInfo),
	}
	return
}

type EventSys struct {
	functions map[core.Event_Key][]*FunctionInfo
}

func (this *EventSys) Register(eId core.Event_Key, f interface{}) (err error) {
	if _, ok := this.functions[eId]; !ok {
		this.functions[eId] = []*FunctionInfo{}
	}
	if this.checkIsRegister(eId, f) {
		return fmt.Errorf("Register the same event repeatedly [%s] method", eId)
	}
	this.functions[eId] = append(this.functions[eId], &FunctionInfo{
		Function:  reflect.ValueOf(f),
		Goroutine: false,
	})
	return
}

func (this *EventSys) RegisterGO(eId core.Event_Key, f interface{}) (err error) {
	if _, ok := this.functions[eId]; !ok {
		this.functions[eId] = []*FunctionInfo{}
	}
	if this.checkIsRegister(eId, f) {
		return fmt.Errorf("Register the same event repeatedly [%s] method", eId)
	}
	this.functions[eId] = append(this.functions[eId], &FunctionInfo{
		Function:  reflect.ValueOf(f),
		Goroutine: true,
	})
	return
}

func (this *EventSys) checkIsRegister(eId core.Event_Key, f interface{}) bool {
	if _, ok := this.functions[eId]; !ok {
		return false
	}
	for _, v := range this.functions[eId] {
		if v.Function == reflect.ValueOf(f) {
			return true
		}
	}
	return false
}

//移除事件
func (this *EventSys) RemoveEvent(eId core.Event_Key, f interface{}) (err error) {
	for i, v := range this.functions[eId] {
		if v.Function == reflect.ValueOf(f) {
			this.functions[eId] = append(this.functions[eId][0:i], this.functions[eId][i+1:]...)
			return
		}
	}
	return fmt.Errorf("Unregistered [%s] event", eId)
}

//触发
func (this *EventSys) TriggerEvent(eId core.Event_Key, agr ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			var rn = ""
			switch r.(type) {

			case string:
				rn = r.(string)
			case error:
				rn = r.(error).Error()
			}
			log.Errorf("Event:%v recover:%v", eId, rn)
		}
	}()
	if v, ok := this.functions[eId]; ok {
		for _, f := range v {
			in := make([]reflect.Value, len(agr))
			for j, a := range agr {
				in[j] = reflect.ValueOf(a)
			}
			if f.Goroutine {
				go f.Function.Call(in)
			} else {
				f.Function.Call(in)
			}
		}
	}
}
