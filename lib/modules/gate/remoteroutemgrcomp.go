package gate

import (
	"fmt"
	"lego/base"
	"lego/core"
	"lego/core/cbase"
	"lego/sys/log"
	"lego/sys/proto"
	"lego/sys/registry"
	"sync"
	"time"
)

type RemoteRouteMgrComp struct {
	cbase.ModuleCompBase
	service    base.IClusterService
	routs      map[uint16]*RemoteRoute
	routslock  sync.RWMutex
	NewSession func(service base.IClusterService, data map[string]interface{}) (s core.IUserSession, err error)
}

func (this *RemoteRouteMgrComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, settings map[string]interface{}) (err error) {
	if s, ok := service.(base.IClusterService); !ok {
		return fmt.Errorf("RemoteRouteMgrComp Init Service is no core.IClusterService")
	} else {
		this.service = s
	}
	if this.NewSession == nil {
		return fmt.Errorf("RemoteRouteMgrComp Init is no install NewSession")
	}
	this.ModuleCompBase.Init(service, module, comp, settings)
	this.routs = make(map[uint16]*RemoteRoute)
	return
}

func (this *RemoteRouteMgrComp) Start() (err error) {
	this.ModuleCompBase.Start()

locp:
	for { //保证rpc事件能成功写入
		err = this.service.Subscribe(Rpc_GateRouteRegister, this.RegisterRoute) //订阅网关注册接口
		if err == nil {
			break locp
		}
		time.Sleep(time.Second)
	}
	return
}

func (this *RemoteRouteMgrComp) RegisterRoute(comId uint16, sId string) (result string, err string) {
	log.Infof("注册远程服务 comId：%d sId:%s", comId, sId)
	snode, e := registry.GetServiceById(sId)
	if e != nil {
		return "", fmt.Sprintf("未发现目标服务【%s】", sId)
	}
	this.routslock.Lock()
	defer this.routslock.Unlock()
	if r, ok := this.routs[comId]; ok {
		e := r.RegisterRoute(*snode)
		if e != nil {
			return "", e.Error()
		}
	} else {
		route := NewRemoteRoute(this.service, comId, this.NewSession, *snode)
		this.routs[comId] = route
	}
	return
}

func (this *RemoteRouteMgrComp) UnRegisterRoute(comId uint16, sId string) {
	this.routslock.Lock()
	defer this.routslock.Unlock()
	if r, ok := this.routs[comId]; ok {
		r.UnRegisterRoute(sId)
		if r.Count() == 0 {
			delete(this.routs, comId)
		}
	}
	return
}

func (this *RemoteRouteMgrComp) OnRoute(agent IAgent, msg proto.IMessage) (code int, err string) {
	this.routslock.RLock()
	route, ok := this.routs[msg.GetComId()]
	this.routslock.RUnlock()
	if ok {
		return route.OnRoute(agent, msg)
	} else {
		return 0, fmt.Sprintf("网关没有注册Comid【%d】路由", msg.GetComId())
	}
}
