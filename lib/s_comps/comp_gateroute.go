package s_comps

import (
	"fmt"
	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib"
	"github.com/liwei1dao/lego/lib/modules/gate"
	"github.com/liwei1dao/lego/sys/event"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/proto"
	"github.com/liwei1dao/lego/sys/registry"
)

type SComp_GateRouteComp struct {
	cbase.ServiceCompBase
	Service base.IClusterService
	Routes  map[uint16]func(s core.IUserSession, msg proto.IMessage) (code int, err string)
}

func (this *SComp_GateRouteComp) GetName() core.S_Comps {
	return lib.SC_ServiceGateRouteComp
}

func (this *SComp_GateRouteComp) Init(service core.IService, comp core.IServiceComp) (err error) {
	if s, ok := service.(base.IClusterService); !ok {
		return fmt.Errorf("SC_GateRouteComp Init service is no IClusterService")
	} else {
		this.Service = s
	}
	err = this.ServiceCompBase.Init(service, comp)
	this.Routes = make(map[uint16]func(s core.IUserSession, msg proto.IMessage) (code int, err string))
	return err
}

func (this *SComp_GateRouteComp) Start() (err error) {
	err = this.ServiceCompBase.Start()
	this.Service.RegisterGO(gate.Rpc_GateRoute, this.ReceiveMsg) //注册网关路由接收接口
	event.RegisterGO(core.Event_ServiceStartEnd, this.registergateroute)
	event.RegisterGO(core.Event_FindNewService, this.findnewservice)
	return
}

func (this *SComp_GateRouteComp) registergateroute() {
	if ss := this.Service.GetSessionsByCategory("GateService"); ss != nil {
		for _, s := range ss {
			for k, _ := range this.Routes {
				log.Infof("向网关【%s】服务器注册服务comId:%d", s.GetId(), k)
				this.Service.RpcInvokeById(s.GetId(), gate.Rpc_GateRouteRegister, false, k, this.Service.GetId()) //向所有网关服务注册
			}
		}
	}
}

/*
注意
”GateService“ 这里是判断发现的服务是否是提供网关注册服务的服务
*/
func (this *SComp_GateRouteComp) findnewservice(node registry.ServiceNode) {
	if len(this.Routes) > 0 && node.Category == "GateService" {
		for k, _ := range this.Routes {
			log.Infof("向网关【%s】服务器注册服务comId:%d", node.Id, k)
			this.Service.RpcInvokeById(node.Id, gate.Rpc_GateRouteRegister, false, k, this.Service.GetId()) //向所有网关服务注册
		}
	}
}

func (this *SComp_GateRouteComp) ReceiveMsg(s core.IUserSession, msg proto.IMessage) (result string, err string) {
	log.Infof("收到网关消息comd:%d msg:%d", msg.GetComId(), msg.GetMsgId())
	if v, ok := this.Routes[msg.GetComId()]; ok {
		v(s, msg)
	}
	return
}

func (this *SComp_GateRouteComp) RegisterRoute(comId uint16, f func(s core.IUserSession, msg proto.IMessage) (code int, err string)) (err error) {
	if _, ok := this.Routes[comId]; !ok {
		this.Routes[comId] = f
		return
	} else {
		return fmt.Errorf("SC_GateRouteComp 重复注册路由")
	}
}
