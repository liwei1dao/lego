package cluster

import (
	"fmt"
	"sync"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/monitor"
	"github.com/liwei1dao/lego/sys/cron"
	"github.com/liwei1dao/lego/sys/event"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/registry"
	"github.com/liwei1dao/lego/sys/rpc"
	"github.com/liwei1dao/lego/utils"
)

type ClusterService struct {
	cbase.ServiceBase
	opts           *Options
	serverList     sync.Map
	ClusterService base.IClusterService
	IsInClustered  bool
	preweight      int32
}

func (this *ClusterService) GetTag() string {
	return this.opts.Tag
}
func (this *ClusterService) GetId() string {
	return this.opts.Id
}
func (this *ClusterService) GetType() string {
	return this.opts.Type
}

func (this *ClusterService) GetCategory() core.S_Category {
	return this.opts.Category
}

func (this *ClusterService) GetVersion() int32 {
	return this.opts.Version
}
func (this *ClusterService) GetSettings() core.ServiceSttings {
	return this.opts.Setting
}
func (this *ClusterService) GetWorkPath() string {
	return this.opts.WorkPath
}
func (this *ClusterService) GetRpcId() string {
	if rpcid, err := rpc.RpcId(); err != nil {
		return ""
	} else {
		return rpcid
	}
}
func (this *ClusterService) GetPreWeight() int32 {
	return this.preweight
}

func (this *ClusterService) SetPreWeight(weight int32) {
	this.preweight = weight
}

func (this *ClusterService) Options() *Options {
	return this.opts
}
func (this *ClusterService) Configure(opts ...Option) {
	this.opts = newOptions(opts...)
}

func (this *ClusterService) Init(service core.IService) (err error) {
	if s, ok := service.(base.IClusterService); !ok {
		panic("service No Is  ClusterService !")
	} else {
		this.ClusterService = s
	}
	return this.ServiceBase.Init(service)
}

func (this *ClusterService) InitSys() {
	// this.ServiceBase.InitSys()
	if err := log.OnInit(this.Service, log.SetLoglevel(this.opts.LogLvel), log.SetDebugMode(this.opts.Debugmode)); err != nil {
		panic(fmt.Sprintf("初始化log系统失败 %s", err.Error()))
	}
	if err := event.OnInit(this.Service); err != nil {
		panic(fmt.Sprintf("初始化event系统失败 %s", err.Error()))
	}
	if err := registry.OnInit(this.ClusterService,
		registry.SetAddress(this.Service.GetSettings().Settings["ConsulAddr"].(string)),
		registry.SetTag(this.opts.Tag),
		registry.SetListener(this)); err != nil {
		panic(fmt.Sprintf("初始化registry系统失败 %s", err.Error()))
	} else {
		log.Debugf("初始化registry系统完成!")
	}
	if err := rpc.OnInit(this.Service, rpc.NatsAddr(this.Service.GetSettings().Settings["NatsAddr"].(string))); err != nil {
		panic(fmt.Sprintf("初始化rpc系统【%s】失败%s", this.Service.GetSettings().Settings["NatsAddr"].(string), err.Error()))
	} else {
		log.Debugf("初始化rpc系统完成!")
	}
	event.Register(core.Event_ServiceStartEnd, func() { //阻塞 先注册服务集群 保证其他服务能及时发现
		if err := registry.Start(); err != nil {
			panic(fmt.Sprintf("加入集群服务失败 err:%s", err.Error()))
		}
	})
	if err := cron.OnInit(this.Service); err != nil {
		panic(fmt.Sprintf("初始化cron系统【%s】失败%s", err.Error()))
	}
}

func (this *ClusterService) Start() (err error) {
	if err = this.ServiceBase.Start(); err != nil {
		return
	}
	return
}

func (this *ClusterService) Run(mod ...core.IModule) {
	modules := make([]core.IModule, 0)
	modules = append(modules, monitor.NewModule())
	modules = append(modules, mod...)
	this.ServiceBase.Run(modules...)
}

func (this *ClusterService) Destroy() (err error) {
	if err = registry.Stop(); err != nil {
		return
	}
	err = this.ServiceBase.Destroy()
	return
}

//注册服务会话 当有新的服务加入时
func (this *ClusterService) FindServiceHandlefunc(node registry.ServiceNode) {
	if _, ok := this.serverList.Load(node.Id); ok { //已经在缓存中 需要更新节点信息
		if s, err := cbase.NewServiceSession(&node); err != nil {
			log.Error("创建服务会话失败【%s】 err = %s")
		} else {
			this.serverList.Store(node.Id, s)
		}
	}
	if this.IsInClustered {
		event.TriggerEvent(core.Event_FindNewService, node) //触发发现新的服务事件
	} else {
		if node.Id == this.opts.Id { //发现自己 加入集群成功
			this.IsInClustered = true
			event.TriggerEvent(core.Event_RegistryStart)
		}
	}
}

//更新服务会话 当有新的服务加入时
func (this *ClusterService) UpDataServiceHandlefunc(node registry.ServiceNode) {
	if ss, ok := this.serverList.Load(node.Id); ok { //已经在缓存中 需要更新节点信息
		session := ss.(core.IServiceSession)
		if session.GetRpcId() != node.RpcId {
			if s, err := cbase.NewServiceSession(&node); err != nil {
				log.Error("创建服务会话失败【%s】 err = %s")
			} else {
				this.serverList.Store(node.Id, s)
			}
			event.TriggerEvent(core.Event_FindNewService, node) //触发发现新的服务事件
		} else {
			if session.GetVersion() != node.Version {
				session.SetVersion(node.Version)
			}
			if session.GetPreWeight() != node.PreWeight {
				session.SetPreWeight(node.PreWeight)
			}
			event.TriggerEvent(core.Event_UpDataOldService, node) //触发发现新的服务事件
		}
	}

}

//注销服务会话
func (this *ClusterService) LoseServiceHandlefunc(sId string) {
	session, ok := this.serverList.Load(sId)
	if ok && session != nil {
		session.(core.IServiceSession).Done()
		this.serverList.Delete(sId)
	}
	event.TriggerEvent(core.Event_LoseService, sId) //触发发现新的服务事件
}
func (this *ClusterService) getServiceSessionByType(sType string) (ss []core.IServiceSession, err error) {
	ss = make([]core.IServiceSession, 0)
	if nodes := registry.GetServiceByType(sType); nodes == nil {
		log.Errorf("获取目标类型【%s】服务集失败", sType)
		return nil, err
	} else {
		for _, v := range nodes {
			if s, ok := this.serverList.Load(v.Id); ok {
				ss = append(ss, s.(core.IServiceSession))
			} else {
				s, err = cbase.NewServiceSession(v)
				if err != nil {
					log.Errorf("创建服务会话失败【%s】 err = %s", v.Id, err.Error())
					continue
				} else {
					this.serverList.Store(v.Id, s)
					ss = append(ss, s.(core.IServiceSession))
				}
			}
		}
	}
	return
}

func (this *ClusterService) GetSessionsByCategory(category core.S_Category) (ss []core.IServiceSession) {
	ss = make([]core.IServiceSession, 0)
	if nodes := registry.GetServiceByCategory(category); nodes == nil {
		log.Errorf("获取目标类型【%s】服务集失败", category)
		return ss
	} else {
		for _, v := range nodes {
			if s, ok := this.serverList.Load(v.Id); ok {
				ss = append(ss, s.(core.IServiceSession))
			} else {
				s, err := cbase.NewServiceSession(v)
				if err != nil {
					log.Errorf("创建服务会话失败【%s】 err = %s", v.Id, err.Error())
					continue
				} else {
					this.serverList.Store(v.Id, s)
					ss = append(ss, s.(core.IServiceSession))
				}
			}
		}
	}
	return
}

//默认路由规则
func (this *ClusterService) DefauleRpcRouteRules(stype string) (ss core.IServiceSession, err error) {
	if s, e := this.getServiceSessionByType(stype); e != nil {
		return nil, e
	} else {
		ss := make([]interface{}, len(s))
		for i, v := range s {
			ss[i] = v
		}
		if len(ss) > 0 {
			//排序找到最优服务
			utils.Sort(ss, func(a interface{}, b interface{}) int8 {
				as := a.(core.IServiceSession)
				bs := b.(core.IServiceSession)
				if as.GetVersion() > bs.GetVersion() {
					return 1
				} else if as.GetVersion() == bs.GetVersion() {
					if as.GetPreWeight() < bs.GetPreWeight() {
						return 1
					} else if as.GetPreWeight() > bs.GetPreWeight() {
						return -1
					} else {
						return 0
					}
				} else {
					return -1
				}
			})
			return ss[0].(core.IServiceSession), nil
		} else {
			return nil, fmt.Errorf("未找到类型【%s】的服务信息", stype)
		}
	}
}
func (this *ClusterService) RpcInvokeById(sId string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error) {
	ss, ok := this.serverList.Load(sId)
	if !ok {
		if node, err := registry.GetServiceById(sId); err != nil {
			log.Errorf("未找到目标服务【%s】节点 err:%s", sId, err.Error())
			return nil, fmt.Errorf("No Found " + sId)
		} else {
			ss, err = cbase.NewServiceSession(node)
			if err != nil {
				return nil, fmt.Errorf(fmt.Sprintf("创建服务会话失败【%s】 err = %s", sId, err.Error()))
			} else {
				this.serverList.Store(node.Id, ss)
			}
		}
	}
	if iscall {
		result, err = ss.(core.IServiceSession).Call(rkey, arg...)
	} else {
		err = ss.(core.IServiceSession).CallNR(rkey, arg...)
	}
	return
}
func (this *ClusterService) RpcInvokeByType(sType string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error) {
	ss, err := this.ClusterService.DefauleRpcRouteRules(sType)
	if err != nil {
		return nil, err
	}
	if iscall {
		result, err = ss.Call(rkey, arg...)
	} else {
		err = ss.CallNR(rkey, arg...)
	}
	return
}
func (this *ClusterService) ReleaseRpc(rkey core.Rpc_Key, arg ...interface{}) {
	rpcf := registry.GetRpcSubById(rkey)
	for _, v := range rpcf {
		this.RpcInvokeById(v.Id, rkey, false, arg...)
	}
	return
}
func (this *ClusterService) Register(id core.Rpc_Key, f interface{}) (err error) {
	return rpc.Register(string(id), f)
}
func (this *ClusterService) RegisterGO(id core.Rpc_Key, f interface{}) (err error) {
	return rpc.RegisterGO(string(id), f)
}
func (this *ClusterService) Subscribe(id core.Rpc_Key, f interface{}) (err error) {
	if err = rpc.RegisterGO(string(id), f); err != nil {
		log.Errorf("注册rpc消息【%s】失败err:%s ", id, err.Error())
		return
	} else {
		err = registry.PushServiceInfo()
	}
	return
}
func (this *ClusterService) UnSubscribe(id core.Rpc_Key, f interface{}) (err error) {
	if err = rpc.UnRegister(string(id), this.GetId()); err != nil {
		log.Errorf("取消订阅rpc消息【%s】失败err:%s ", id, err.Error())
		return
	} else {
		err = registry.PushServiceInfo()
	}
	return
}
