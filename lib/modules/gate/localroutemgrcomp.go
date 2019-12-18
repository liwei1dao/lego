package gate

import (
	"fmt"
	"lego/core"
	"lego/core/cbase"
	"lego/sys/log"
	"lego/sys/proto"
	"sync"
)

type LocalRouteMgrComp struct {
	cbase.ModuleCompBase
	module     IGateModule
	routs      map[uint16]*LocalRoute
	routslock  sync.RWMutex
	NewSession func(module IGateModule, data map[string]interface{}) (s core.IUserSession, err error)
}

func (this *LocalRouteMgrComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, settings map[string]interface{}) (err error) {
	if m, ok := module.(IGateModule); !ok {
		return fmt.Errorf("LocalRouteMgrComp Init module is no IGateModule")
	} else {
		this.module = m
	}
	if this.NewSession == nil {
		return fmt.Errorf("LocalRouteMgrComp Init is no install NewLocalSession")
	}
	this.ModuleCompBase.Init(service, module, comp, settings)
	this.routs = make(map[uint16]*LocalRoute)
	return
}

func (this *LocalRouteMgrComp) RegisterRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code int, err string)) {
	log.Infof("注册本地服务 comId：%d ", comId)
	this.routslock.Lock()
	defer this.routslock.Unlock()
	if _, ok := this.routs[comId]; ok {
		log.Errorf("重复注册网关路由【%d】", comId)
	} else {
		route := NewLocalRoute(this.module, comId, this.NewSession, f)
		this.routs[comId] = route
	}
	return
}

func (this *LocalRouteMgrComp) UnRegisterRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage)) {
	this.routslock.Lock()
	defer this.routslock.Unlock()
	delete(this.routs, comId)
	return
}

func (this *LocalRouteMgrComp) OnRoute(agent IAgent, msg proto.IMessage) (code int, err string) {
	this.routslock.RLock()
	route, ok := this.routs[msg.GetComId()]
	this.routslock.RUnlock()
	if ok {
		return route.OnRoute(agent, msg)
	} else {
		return 0, fmt.Sprintf("网关没有注册Comid【%d】路由", msg.GetComId())
	}
}
