package rpcxcluster

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
	"github.com/liwei1dao/lego/sys/rpcx"
	"github.com/liwei1dao/lego/utils/container/sortslice"
	"github.com/liwei1dao/lego/utils/container/version"
)

type RPCXService struct {
	cbase.ServiceBase
	opts           *Options
	serverList     sync.Map
	RPCXService    base.IRPCXService
	ServiceMonitor core.IServiceMonitor
	IsInClustered  bool
	lock           sync.RWMutex //服务锁
}

func (this *RPCXService) GetTag() string {
	return this.opts.Setting.Tag
}
func (this *RPCXService) GetId() string {
	return this.opts.Setting.Id
}
func (this *RPCXService) GetType() string {
	return this.opts.Setting.Type
}
func (this *RPCXService) GetIp() string {
	return this.opts.Setting.Ip
}
func (this *RPCXService) GetPort() int {
	return this.opts.Setting.Port
}
func (this *RPCXService) GetCategory() core.S_Category {
	return this.opts.Setting.Category
}

func (this *RPCXService) GetVersion() string {
	return this.opts.Version
}
func (this *RPCXService) GetSettings() core.ServiceSttings {
	return this.opts.Setting
}
func (this *RPCXService) GetRpcId() string {
	return ""
}
func (this *RPCXService) GetPreWeight() float64 {
	return float64(1) / float64(runtime.NumGoroutine())
}

func (this *RPCXService) SetPreWeight(weight int32) {

}

func (this *RPCXService) GetServiceMonitor() core.IServiceMonitor {
	return this.ServiceMonitor
}

func (this *RPCXService) Options() *Options {
	return this.opts
}
func (this *RPCXService) Configure(opts ...Option) {
	this.opts = newOptions(opts...)
}

func (this *RPCXService) Init(service core.IService) (err error) {
	this.RPCXService = service.(base.IRPCXService)
	return this.ServiceBase.Init(service)
}

func (this *RPCXService) InitSys() {
	if err := log.OnInit(this.opts.Setting.Sys["log"]); err != nil {
		panic(fmt.Sprintf("初始化log系统失败 err:%v", err))
	} else {
		log.Infof("Sys log Init success !")
	}
	if err := event.OnInit(this.opts.Setting.Sys["event"]); err != nil {
		log.Panicf(fmt.Sprintf("初始化event系统失败 err:%v", err))
	} else {
		log.Infof("Sys event Init success !")
	}
	if err := cron.OnInit(this.opts.Setting.Sys["cron"]); err != nil {
		log.Panicf(fmt.Sprintf("初始化cron系统 err:%v", err))
	} else {
		log.Infof("Sys cron Init success !")
	}
	if err := registry.OnInit(this.opts.Setting.Sys["registry"], registry.SetService(this.RPCXService), registry.SetListener(this.RPCXService.(registry.IListener))); err != nil {
		log.Panicf(fmt.Sprintf("初始化registry系统失败 err:%v", err))
	} else {
		log.Infof("Sys registry Init success !")
	}
	if err := rpcx.OnInit(this.opts.Setting.Sys["rpcx"], rpcx.SetAddr(fmt.Sprintf("%s:%d", this.GetIp(), this.GetPort()))); err != nil {
		log.Panicf(fmt.Sprintf("初始化rpcx系统 err:%v", err))
	} else {
		log.Infof("Sys rpcx Init success !")
	}
	event.Register(core.Event_ServiceStartEnd, func() { //阻塞 先注册服务集群 保证其他服务能及时发现
		if err := rpcx.Start(); err != nil {
			log.Panicf(fmt.Sprintf("启动RPC失败 err:%v", err))
		}
		if err := registry.Start(); err != nil {
			log.Panicf(fmt.Sprintf("加入集群服务失败 err:%v", err))
		}
	})
}

func (this *RPCXService) Start() (err error) {
	if err = this.ServiceBase.Start(); err != nil {
		return
	}
	return
}

func (this *RPCXService) Run(mod ...core.IModule) {
	this.ServiceMonitor = monitor.NewModule()
	modules := make([]core.IModule, 0)
	modules = append(modules, this.ServiceMonitor)
	modules = append(modules, mod...)
	this.ServiceBase.Run(modules...)
}

func (this *RPCXService) Destroy() (err error) {
	if err = rpcx.Stop(); err != nil {
		return
	}
	if err = registry.Stop(); err != nil {
		return
	}
	cron.Stop()
	err = this.ServiceBase.Destroy()
	return
}

//注册服务会话 当有新的服务加入时
func (this *RPCXService) FindServiceHandlefunc(node registry.ServiceNode) {
	this.lock.Lock()
	if _, ok := this.serverList.Load(node.Id); !ok {
		if s, err := NewServiceSession(&node); err != nil {
			log.Errorf("创建服务会话失败【%s】 err:%v", node.Id, err)
		} else {
			this.serverList.Store(node.Id, s)
		}
	}
	this.lock.Unlock()
	if this.IsInClustered {
		event.TriggerEvent(core.Event_FindNewService, node) //触发发现新的服务事件
	} else {
		if node.Id == this.opts.Setting.Id { //发现自己 加入集群成功
			this.IsInClustered = true
			event.TriggerEvent(core.Event_RegistryStart)
		}
	}
}

//更新服务会话 当有新的服务加入时
func (this *RPCXService) UpDataServiceHandlefunc(node registry.ServiceNode) {
	if ss, ok := this.serverList.Load(node.Id); ok { //已经在缓存中 需要更新节点信息
		session := ss.(base.IRPCXServiceSession)
		if session.GetRpcId() != node.RpcId {
			if s, err := NewServiceSession(&node); err != nil {
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
func (this *RPCXService) LoseServiceHandlefunc(sId string) {
	this.lock.Lock()
	session, ok := this.serverList.Load(sId)
	if ok && session != nil {
		session.(base.IRPCXServiceSession).Done()
		this.serverList.Delete(sId)
	}
	this.lock.Unlock()
	event.TriggerEvent(core.Event_LoseService, sId) //触发发现新的服务事件
}

func (this *RPCXService) getServiceSessionByType(sType string, sIp string) (ss []base.IRPCXServiceSession, err error) {
	ss = make([]base.IRPCXServiceSession, 0)
	if nodes := registry.GetServiceByType(sType); nodes == nil {
		log.Errorf("获取目标类型 type【%s】ip [%s] 服务集失败", sType, sIp)
		return nil, err
	} else {
		if sIp == core.AutoIp {
			for _, v := range nodes {
				if s, ok := this.serverList.Load(v.Id); ok {
					ss = append(ss, s.(base.IRPCXServiceSession))
				} else {
					s, err = NewServiceSession(v)
					if err != nil {
						log.Errorf("创建服务会话失败【%s】 err:%v", v.Id, err)
						continue
					} else {
						this.serverList.Store(v.Id, s)
						ss = append(ss, s.(base.IRPCXServiceSession))
					}
				}
			}
		} else {
			for _, v := range nodes {
				if v.IP == sIp {
					if s, ok := this.serverList.Load(v.Id); ok {
						ss = append(ss, s.(base.IRPCXServiceSession))
					} else {
						s, err = NewServiceSession(v)
						if err != nil {
							log.Errorf("创建服务会话失败【%s】 err:%v", v.Id, err)
							continue
						} else {
							this.serverList.Store(v.Id, s)
							ss = append(ss, s.(base.IRPCXServiceSession))
						}
					}
				}
			}
		}
	}
	return
}
func (this *RPCXService) getServiceSessionByIds(Ids []string) (result []base.IRPCXServiceSession, err error) {

	var (
		ss    = make(map[string]base.IRPCXServiceSession)
		nodes *registry.ServiceNode
	)
	for _, v := range Ids {
		if nodes, err = registry.GetServiceById(v); err != nil {
			log.Errorf("获取目标类型 ip[%v] 服务集失败:%v", v, err)
			return
		} else {
			if s, ok := this.serverList.Load(v); ok {
				ss[v] = s.(base.IRPCXServiceSession)
			} else {
				s, err = NewServiceSession(nodes)
				if err != nil {
					log.Errorf("创建服务会话失败【%s】 err:%v", v, err)
					continue
				} else {
					this.serverList.Store(v, s)
					ss[v] = s.(base.IRPCXServiceSession)
				}
			}
		}
	}
	result = make([]base.IRPCXServiceSession, len(ss))
	n := 0
	for _, v := range ss {
		result[n] = v
		n++
	}
	return
}
func (this *RPCXService) getServiceSessionByIps(sType string, sIps []string) (result []base.IRPCXServiceSession, err error) {
	var (
		ss    map[string]base.IRPCXServiceSession = make(map[string]base.IRPCXServiceSession)
		s     base.IRPCXServiceSession
		found bool
	)
	//容错处理
	if len(sIps) == 1 && sIps[0] == core.AutoIp {
		s, err = this.RPCXService.DefauleRpcRouteRules(sType, core.AutoIp)
		if err == nil {
			result = []base.IRPCXServiceSession{s}
		} else {
			log.Errorf("未找到目标服务 ip:%v type:%s 节点 err:%v", sIps, sType, err)
		}
		return
	}

	Include := func(ip string, ips []string) bool {
		for _, v := range ips {
			if v == core.AllIp || ip == v {
				return true
			}
		}
		return false
	}
	if nodes := registry.GetServiceByType(sType); nodes == nil {
		log.Errorf("获取目标类型 type【%s】ips [%v] 服务集失败", sType, sIps)
		return nil, err
	} else {
		for _, v := range nodes {
			if Include(v.IP, sIps) {
				if s, ok := this.serverList.Load(v.Id); ok {
					ss[v.Id] = s.(base.IRPCXServiceSession)
				} else {
					s, err = NewServiceSession(v)
					if err != nil {
						log.Errorf("创建服务会话失败【%s】 err:%v", v.Id, err)
						continue
					} else {
						this.serverList.Store(v.Id, s)
						ss[v.Id] = s.(base.IRPCXServiceSession)
					}
				}
			}
		}
	}

	for _, v := range sIps {
		if v != core.AllIp {
			found = false
			for _, v1 := range ss {
				if v == v1.GetIp() {
					found = true
					break
				}
			}
			if !found { //未查询到目标服务节点
				log.Errorf("获取目标类型 type【%s】ips [%v] 服务集失败", sType, sIps)
				err = fmt.Errorf("未查询到服务 type【%s】ip [%v]", sType, v)
				return
			}
		}
	}

	result = make([]base.IRPCXServiceSession, len(ss))
	n := 0
	for _, v := range ss {
		result[n] = v
		n++
	}
	return
}

func (this *RPCXService) GetSessionsByCategory(category core.S_Category) (ss []base.IRPCXServiceSession) {
	ss = make([]base.IRPCXServiceSession, 0)
	if nodes := registry.GetServiceByCategory(category); nodes == nil {
		log.Errorf("获取目标类型【%s】服务集失败", category)
		return ss
	} else {
		for _, v := range nodes {
			if s, ok := this.serverList.Load(v.Id); ok {
				ss = append(ss, s.(base.IRPCXServiceSession))
			} else {
				s, err := NewServiceSession(v)
				if err != nil {
					log.Errorf("创建服务会话失败【%s】 err:%v", v.Id, err)
					continue
				} else {
					this.serverList.Store(v.Id, s)
					ss = append(ss, s.(base.IRPCXServiceSession))
				}
			}
		}
	}
	return
}

//默认路由规则
func (this *RPCXService) DefauleRpcRouteRules(stype string, sip string) (ss base.IRPCXServiceSession, err error) {
	if s, e := this.getServiceSessionByType(stype, sip); e != nil {
		return nil, e
	} else {
		ss := make([]interface{}, len(s))
		for i, v := range s {
			ss[i] = v
		}
		if len(ss) > 0 {
			//排序找到最优服务
			sortslice.Sort(ss, func(a interface{}, b interface{}) int8 {
				as := a.(base.IRPCXServiceSession)
				bs := b.(base.IRPCXServiceSession)
				if iscompare := version.CompareStrVer(as.GetVersion(), bs.GetVersion()); iscompare != 0 {
					return iscompare
				} else {
					if as.GetPreWeight() < bs.GetPreWeight() {
						return 1
					} else if as.GetPreWeight() > bs.GetPreWeight() {
						return -1
					} else {
						return 0
					}
				}
			})
			return ss[0].(base.IRPCXServiceSession), nil
		} else {
			return nil, fmt.Errorf("未找到IP[%s]类型%s】的服务信息", sip, stype)
		}
	}
}

func (this *RPCXService) Register(id core.Rpc_Key, f interface{}) {
	rpcx.Register(id, f)
}

func (this *RPCXService) Subscribe(id core.Rpc_Key, f interface{}) (err error) {
	rpcx.Register(id, f)
	err = registry.PushServiceInfo()
	return
}

func (this *RPCXService) UnSubscribe(id core.Rpc_Key, f interface{}) (err error) {
	rpcx.UnRegister(id)
	err = registry.PushServiceInfo()
	return
}
