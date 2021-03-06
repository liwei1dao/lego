package cluster

import (
	"fmt"
	"runtime"
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
	ServiceMonitor core.IServiceMonitor
	IsInClustered  bool
}

func (this *ClusterService) GetTag() string {
	return this.opts.Setting.Tag
}
func (this *ClusterService) GetId() string {
	return this.opts.Setting.Id
}
func (this *ClusterService) GetType() string {
	return this.opts.Setting.Type
}

func (this *ClusterService) GetCategory() core.S_Category {
	return this.opts.Setting.Category
}

func (this *ClusterService) GetVersion() float32 {
	return this.opts.Setting.Version
}
func (this *ClusterService) GetSettings() core.ServiceSttings {
	return this.opts.Setting
}
func (this *ClusterService) GetRpcId() string {
	return rpc.RpcId()
}
func (this *ClusterService) GetPreWeight() float64 {
	return float64(1) / float64(runtime.NumGoroutine())
}

func (this *ClusterService) SetPreWeight(weight int32) {

}

func (this *ClusterService) GetServiceMonitor() core.IServiceMonitor {
	return this.ServiceMonitor
}

func (this *ClusterService) Options() *Options {
	return this.opts
}
func (this *ClusterService) Configure(opts ...Option) {
	this.opts = newOptions(opts...)
}

func (this *ClusterService) Init(service core.IService) (err error) {
	this.ClusterService = service.(base.IClusterService)
	return this.ServiceBase.Init(service)
}

func (this *ClusterService) InitSys() {
	if err := log.OnInit(this.opts.Setting.Sys["log"]); err != nil {
		panic(fmt.Sprintf("初始化log系统失败 err:%v", err))
	} else {
		log.Infof("Sys log Init success !")
	}
	if err := event.OnInit(this.opts.Setting.Sys["event"]); err != nil {
		panic(fmt.Sprintf("初始化event系统失败 err:%v", err))
	} else {
		log.Infof("Sys event Init success !")
	}
	if err := cron.OnInit(this.opts.Setting.Sys["cron"]); err != nil {
		panic(fmt.Sprintf("初始化cron系统 err:%v", err))
	} else {
		log.Infof("Sys cron Init success !")
	}
	if err := registry.OnInit(this.opts.Setting.Sys["registry"], registry.SetService(this.ClusterService)); err != nil {
		panic(fmt.Sprintf("初始化registry系统失败 err:%v", err))
	} else {
		log.Infof("Sys registry Init success !")
	}
	if err := rpc.OnInit(this.opts.Setting.Sys["rpc"]); err != nil {
		panic(fmt.Sprintf("初始化rpc系统 err:%v", err))
	} else {
		log.Infof("Sys rpc Init success !")
	}
	event.Register(core.Event_ServiceStartEnd, func() { //阻塞 先注册服务集群 保证其他服务能及时发现
		if err := registry.Start(); err != nil {
			panic(fmt.Sprintf("加入集群服务失败 err:%v", err))
		}
	})
}

func (this *ClusterService) Start() (err error) {
	if err = this.ServiceBase.Start(); err != nil {
		return
	}
	return
}

func (this *ClusterService) Run(mod ...core.IModule) {
	this.ServiceMonitor = monitor.NewModule()
	modules := make([]core.IModule, 0)
	modules = append(modules, this.ServiceMonitor)
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
			log.Errorf("创建服务会话失败【%s】 err:%v", node.Id, err)
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
				log.Errorf("更新服务会话失败【%s】 err:%v", node.Id, err)
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
					log.Errorf("创建服务会话失败【%s】 err:%v", v.Id, err)
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
					log.Errorf("创建服务会话失败【%s】 err:%v", v.Id, err)
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
					if as.GetPreWeight() > bs.GetPreWeight() {
						return 1
					} else if as.GetPreWeight() < bs.GetPreWeight() {
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
			log.Errorf("未找到目标服务【%s】节点 err:%v", sId, err)
			return nil, fmt.Errorf("No Found " + sId)
		} else {
			ss, err = cbase.NewServiceSession(node)
			if err != nil {
				return nil, fmt.Errorf(fmt.Sprintf("创建服务会话失败【%s】 err:%v", sId, err))
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
		log.Errorf("未找到目标服务【%s】节点 err:%v", sType, err)
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
func (this *ClusterService) Register(id core.Rpc_Key, f interface{}) {
	rpc.Register(id, f)
}
func (this *ClusterService) RegisterGO(id core.Rpc_Key, f interface{}) {
	rpc.RegisterGO(id, f)
}
func (this *ClusterService) Subscribe(id core.Rpc_Key, f interface{}) (err error) {
	rpc.RegisterGO(id, f)
	err = registry.PushServiceInfo()
	return
}
func (this *ClusterService) UnSubscribe(id core.Rpc_Key, f interface{}) (err error) {
	rpc.UnRegister(id, this.GetId())
	err = registry.PushServiceInfo()
	return
}
