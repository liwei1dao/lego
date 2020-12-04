package gate

import (
	"fmt"
	"sync"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/proto"
	"github.com/liwei1dao/lego/sys/registry"
)

type RemoteRouteMgrComp struct {
	cbase.ModuleCompBase
	service    base.IClusterService
	routs      map[uint16]map[string]*RemoteRoute
	routslock  sync.RWMutex
	NewSession func(service base.IClusterService, data map[string]interface{}) (s core.IUserSession, err error)
}

func (this *RemoteRouteMgrComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	if s, ok := service.(base.IClusterService); !ok {
		return fmt.Errorf("RemoteRouteMgrComp Init Service is no core.IClusterService")
	} else {
		this.service = s
	}
	if this.NewSession == nil {
		return fmt.Errorf("RemoteRouteMgrComp Init is no install NewSession")
	}
	this.ModuleCompBase.Init(service, module, comp, options)
	this.routs = make(map[uint16]map[string]*RemoteRoute)
	return
}

func (this *RemoteRouteMgrComp) Start() (err error) {
	this.ModuleCompBase.Start()
	this.service.Subscribe(Rpc_GateRouteRegister, this.RegisterRoute)
	return
}

func (this *RemoteRouteMgrComp) SetNewSession(f func(service base.IClusterService, data map[string]interface{}) (s core.IUserSession, err error)) {
	this.NewSession = f
}

func (this *RemoteRouteMgrComp) RegisterRoute(comId uint16, sId, sType string) (result string, err string) {
	// log.Infof("注册远程服务 comId：%d sId:%s", comId, sId)
	// snode, e := registry.GetServiceById(sId)
	// if e != nil {
	// 	return "", fmt.Sprintf("未发现目标服务【%s】", sId)
	// }
	snode := &registry.ServiceNode{
		Id:   sId,
		Type: sType,
	}
	this.routslock.Lock()
	defer this.routslock.Unlock()
	if _, ok := this.routs[comId]; ok {
		if r, ok1 := this.routs[comId][snode.Type]; ok1 {
			e := r.RegisterRoute(*snode)
			if e != nil {
				return "", e.Error()
			}
		} else {
			route := NewRemoteRoute(this.service, comId, this.NewSession, *snode)
			this.routs[comId][snode.Type] = route
		}

	} else {
		route := NewRemoteRoute(this.service, comId, this.NewSession, *snode)
		this.routs[comId] = map[string]*RemoteRoute{snode.Type: route}
	}
	return
}

func (this *RemoteRouteMgrComp) UnRegisterRoute(comId uint16, sType, sId string) {
	this.routslock.Lock()
	defer this.routslock.Unlock()
	if r, ok := this.routs[comId]; ok {
		if r1, ok1 := r[sType]; ok1 {
			r1.UnRegisterRoute(sId)
			if r1.Count() == 0 {
				delete(r, sType)
				if len(r) == 0 {
					delete(this.routs, comId)
				}
			}
		}
	}
	return
}

func (this *RemoteRouteMgrComp) OnRoute(agent IAgent, msg proto.IMessage) (code core.ErrorCode, err error) {
	this.routslock.RLock()
	routes, ok := this.routs[msg.GetComId()]
	this.routslock.RUnlock()
	if ok {
		for _, v := range routes {
			if cd, e := v.OnRoute(agent, msg); cd != core.ErrorCode_Success || e != "" {
				return cd, fmt.Errorf("LocalRoute err:%s", e)
			}
		}
		return core.ErrorCode_Success, nil
	} else {
		return core.ErrorCode_NoRoute, fmt.Errorf("NoRoute ComId:%d", msg.GetComId())
	}
}
