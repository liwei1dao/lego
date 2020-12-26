package gate

import (
	"fmt"
	"sync"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/proto"
	"github.com/liwei1dao/lego/sys/registry"
)

func NewRemoteRoute(service base.IClusterService, comId uint16, f func(service base.IClusterService, data map[string]interface{}) (s core.IUserSession, err error), sNode registry.ServiceNode) *RemoteRoute {
	r := &RemoteRoute{
		Service:     service,
		ComId:       comId,
		services:    []string{sNode.Id},
		ServiceType: sNode.Type,
		sessionFun:  f,
	}
	return r
}
func NewLocalRoute(module IGateModule, comId uint16, sf func(module IGateModule, data map[string]interface{}) (s core.IUserSession, err error), f func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string)) *LocalRoute {
	r := &LocalRoute{
		Module:     module,
		ComId:      comId,
		f:          f,
		sessionFun: sf,
	}
	return r
}

type RemoteRoute struct {
	Service     base.IClusterService
	ComId       uint16
	services    []string
	ServiceType string
	slock       sync.RWMutex
	sessionFun  func(service base.IClusterService, data map[string]interface{}) (s core.IUserSession, err error)
}

func (this *RemoteRoute) Count() int32 {
	return int32(len(this.services))
}

//注册远程路由
func (this *RemoteRoute) RegisterRoute(sNode registry.ServiceNode) (err error) {
	defer this.slock.Unlock()
	this.slock.Lock()
	if this.ServiceType != sNode.Type {
		return fmt.Errorf("网关ComId【%d】已经注册服务类型【%s】,不可以在注册其他服务【%s】", this.ComId, this.ServiceType, sNode.Type)
	}
	for _, v := range this.services {
		if v == sNode.Id {
			log.Warnf("网关ComId【%d】已经注册服务【%s】", this.ComId, sNode.Id)
			return
		}
	}
	this.services = append(this.services, sNode.Id)
	return nil
}

//注销远程路由
func (this *RemoteRoute) UnRegisterRoute(sId string) {
	defer this.slock.Unlock()
	this.slock.Lock()
	for i, v := range this.services {
		if v == sId {
			this.services = append(this.services[0:i], this.services[i+1:]...)
		}
	}
}

func (this *RemoteRoute) OnRoute(a IAgent, msg proto.IMessage) (code core.ErrorCode, err string) {
	if s, e := this.sessionFun(this.Service, a.GetSessionData()); e != nil {
		return 0, e.Error()
	} else {
		_, e := this.Service.RpcInvokeByType(this.ServiceType, Rpc_GateRoute, false, s, msg)
		//code = result.(int)
		if e != nil {
			err = e.Error()
		}
		return
	}
}

type LocalRoute struct {
	Module     IGateModule
	ComId      uint16
	f          func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string)
	sessionFun func(module IGateModule, data map[string]interface{}) (s core.IUserSession, err error)
}

func (this *LocalRoute) OnRoute(a IAgent, msg proto.IMessage) (code core.ErrorCode, err string) {
	if s, err := this.sessionFun(this.Module, a.GetSessionData()); err != nil {
		return core.ErrorCode_RpcFuncExecutionError, err.Error()
	} else {
		return this.f(s, msg)
	}
}
