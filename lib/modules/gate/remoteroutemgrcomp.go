package gate

import (
	"fmt"
	"sync"
	"time"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
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
	this.routs = make(map[uint16]map[string]*RemoteRoute)
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

func (this *RemoteRouteMgrComp) SetNewSession(f func(service base.IClusterService, data map[string]interface{}) (s core.IUserSession, err error)) {
	this.NewSession = f
}

func (this *RemoteRouteMgrComp) RegisterRoute(comId uint16, sId string) (result string, err string) {
	log.Infof("注册远程服务 comId：%d sId:%s", comId, sId)
	snode, e := registry.GetServiceById(sId)
	if e != nil {
		return "", fmt.Sprintf("未发现目标服务【%s】", sId)
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

func (this *RemoteRouteMgrComp) OnRoute(agent IAgent, msg proto.IMessage) (iscontinue bool) {
	this.routslock.RLock()
	routes, ok := this.routs[msg.GetComId()]
	this.routslock.RUnlock()
	if ok {
		for _, v := range routes {
			v.OnRoute(agent, msg)
		}
	} else {
		return true
	}
	return true
}
