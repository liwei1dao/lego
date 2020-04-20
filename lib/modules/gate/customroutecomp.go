package gate

import (
	"time"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	cbase "github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/proto"
)

//自定义网关组件 特殊消息转发可以在自定义网关内实现与注册
type CustomRouteComp struct {
	cbase.ModuleCompBase
	service   base.IClusterService
	route     map[core.CustomRoute]map[uint16][]uint16
	routeFunc map[core.CustomRoute]func(a IAgent, msg proto.IMessage) (code core.ErrorCode, err string)
}

func (this *CustomRouteComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, settings map[string]interface{}) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, settings)
	this.service = service.(base.IClusterService)
	this.route = make(map[core.CustomRoute]map[uint16][]uint16)
	this.routeFunc = make(map[core.CustomRoute]func(a IAgent, msg proto.IMessage) (code core.ErrorCode, err string))
	return
}

func (this *CustomRouteComp) Start() (err error) {
	this.ModuleCompBase.Start()

locp:
	for { //保证rpc事件能成功写入
		err = this.service.Subscribe(Rpc_GateCustomRouteRegister, this.RegisterRoute) //订阅网关注册接口
		if err == nil {
			break locp
		}
		time.Sleep(time.Second)
	}
	return
}

//添加自定义网关组注册接口
func (this *CustomRouteComp) RegisterRoute(route core.CustomRoute, msgs map[uint16][]uint16) (result string, err string) {
	if _, ok := this.routeFunc[route]; !ok {
		log.Errorf("注册自定义路由错误 路由接口未注册 route：%d", route)
		return
	}
	var ishave bool
	if _, ok := this.route[route]; !ok {
		this.route[route] = make(map[uint16][]uint16)
	}
	for k, v := range msgs {
		if _, ok := this.route[route][k]; !ok {
			this.route[route][k] = make([]uint16, 0)
		}
		for _, id1 := range v {
			ishave = false
			for _, id2 := range this.route[route][k] {
				if id1 == id2 {
					ishave = true
				}
			}
			if !ishave {
				this.route[route][k] = append(this.route[route][k], id1)
			}
		}
	}
	return
}

//添加自定义网关组注册接口
func (this *CustomRouteComp) RegisterRouteFunc(route core.CustomRoute, f func(a IAgent, msg proto.IMessage) (code core.ErrorCode, err string)) {
	if _, ok := this.routeFunc[route]; !ok {
		this.routeFunc[route] = f
	} else {
		log.Errorf("重复注册自定义网关 route：%d", route)
	}
	return
}

func (this *CustomRouteComp) OnRoute(agent IAgent, msg proto.IMessage) (iscontinue bool) {
	var ishave bool
	var route core.CustomRoute
	for k, v1 := range this.route {
		if v2, ok := v1[msg.GetComId()]; ok {
			for _, id := range v2 {
				if id == msg.GetMsgId() {
					route = k
					ishave = true
					break
				}
			}
		}
	}
	if ishave {
		if f, ok := this.routeFunc[route]; ok {
			f(agent, msg)
			return false
		} else {
			log.Errorf("自定义网关消息注册但是处理函数没有注册 route：%d", route)
		}
	}
	return true
}
