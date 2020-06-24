package event

import (
	"fmt"
	"reflect"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
)

var (
	service   core.IService
	functions map[core.Event_Key][]*FunctionInfo
)

func OnInit(s core.IService, opt ...Option) (err error) {
	service = s
	functions = make(map[core.Event_Key][]*FunctionInfo)
	return
}

func Register(eId core.Event_Key, f interface{}) (err error) {
	if functions == nil {
		return fmt.Errorf("event 系统未初始化")
	}
	if _, ok := functions[eId]; !ok {
		functions[eId] = []*FunctionInfo{}
	}
	if checkIsRegister(eId, f) {
		return fmt.Errorf("重复注册相同事件【%s】方法", eId)
	}
	functions[eId] = append(functions[eId], &FunctionInfo{
		Function:  reflect.ValueOf(f),
		Goroutine: false,
	})
	return
}

func RegisterGO(eId core.Event_Key, f interface{}) (err error) {
	if functions == nil {
		return fmt.Errorf("event 系统未初始化")
	}
	if _, ok := functions[eId]; !ok {
		functions[eId] = []*FunctionInfo{}
	}
	if checkIsRegister(eId, f) {
		return fmt.Errorf("重复注册相同事件【%s】方法", eId)
	}
	functions[eId] = append(functions[eId], &FunctionInfo{
		Function:  reflect.ValueOf(f),
		Goroutine: true,
	})
	return
}

func checkIsRegister(eId core.Event_Key, f interface{}) bool {
	if _, ok := functions[eId]; !ok {
		return false
	}
	for _, v := range functions[eId] {
		if v.Function == reflect.ValueOf(f) {
			return true
		}
	}
	return false
}

//移除事件
func RemoveEvent(eId core.Event_Key, f interface{}) (err error) {
	if _, ok := functions[eId]; !ok {
		return fmt.Errorf("未注册【%s】事件", eId)
	}
	for i, v := range functions[eId] {
		if v.Function == reflect.ValueOf(f) {
			functions[eId] = append(functions[eId][0:i], functions[eId][i+1:]...)
			return
		}
	}
	return fmt.Errorf("未注册【%s】事件", eId)
}

//触发
func TriggerEvent(eId core.Event_Key, agr ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			var rn = ""
			switch r.(type) {

			case string:
				rn = r.(string)
			case error:
				rn = r.(error).Error()
			}
			log.Errorf("Event:[%d] recover:[%s]", eId, rn)
		}
	}()
	if v, ok := functions[eId]; ok {
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
