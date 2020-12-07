package gate

import (
	"fmt"
	"sync"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/proto"
)

type LocalRouteMgrComp struct {
	cbase.ModuleCompBase
	module     IGateModule
	routs      map[uint16]map[*func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string)]*LocalRoute
	routslock  sync.RWMutex
	NewSession func(module IGateModule, data map[string]interface{}) (s core.IUserSession, err error)
}

func (this *LocalRouteMgrComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	if m, ok := module.(IGateModule); !ok {
		return fmt.Errorf("LocalRouteMgrComp Init module is no IGateModule")
	} else {
		this.module = m
	}
	if this.NewSession == nil {
		return fmt.Errorf("LocalRouteMgrComp Init is no install NewLocalSession")
	}
	this.ModuleCompBase.Init(service, module, comp, options)
	this.routs = make(map[uint16]map[*func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string)]*LocalRoute)
	return
}

func (this *LocalRouteMgrComp) SetNewSession(f func(module IGateModule, data map[string]interface{}) (s core.IUserSession, err error)) {
	this.NewSession = f
}

func (this *LocalRouteMgrComp) RegisterRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string)) {
	log.Infof("注册本地服务 comId：%d f:%v", comId, f)
	this.routslock.Lock()
	defer this.routslock.Unlock()
	route := NewLocalRoute(this.module, comId, this.NewSession, f)
	if _, ok := this.routs[comId]; !ok {
		this.routs[comId] = make(map[*func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string)]*LocalRoute)
	}
	this.routs[comId][&f] = route
	return
}

func (this *LocalRouteMgrComp) UnRegisterRoute(comId uint16, f func(session core.IUserSession, msg proto.IMessage) (code core.ErrorCode, err string)) {
	this.routslock.Lock()
	defer this.routslock.Unlock()
	if r, ok := this.routs[comId]; ok {
		delete(r, &f)
		if len(r) == 0 {
			delete(this.routs, comId)
		}
	}
	return
}

func (this *LocalRouteMgrComp) OnRoute(agent IAgent, msg proto.IMessage) (code core.ErrorCode, err error) {
	this.routslock.RLock()
	routes, ok := this.routs[msg.GetComId()]
	this.routslock.RUnlock()
	if ok {
		for _, v := range routes {
			if cd, e := v.OnRoute(agent, msg); e != "" {
				return cd, fmt.Errorf("LocalRoute err:%s", e)
			}
		}
		return core.ErrorCode_Success, nil
	} else {
		return core.ErrorCode_NoRoute, fmt.Errorf("NoRoute ComId:%d", msg.GetComId())
	}
}
