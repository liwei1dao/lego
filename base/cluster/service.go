package cluster

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/liwei1dao/lego"
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
func (this *ClusterService) GetIp() string {
	return this.opts.Setting.Ip
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
		log.Panicf(fmt.Sprintf("初始化event系统失败 err:%v", err))
	} else {
		log.Infof("Sys event Init success !")
	}
	if err := cron.OnInit(this.opts.Setting.Sys["cron"]); err != nil {
		log.Panicf(fmt.Sprintf("初始化cron系统 err:%v", err))
	} else {
		log.Infof("Sys cron Init success !")
	}
	if err := registry.OnInit(this.opts.Setting.Sys["registry"], registry.SetService(this.ClusterService), registry.SetListener(this.ClusterService.(registry.IListener))); err != nil {
		log.Panicf(fmt.Sprintf("初始化registry系统失败 err:%v", err))
	} else {
		log.Infof("Sys registry Init success !")
	}
	if err := rpc.OnInit(this.opts.Setting.Sys["rpc"], rpc.SetClusterTag(this.GetTag()), rpc.SetServiceId(this.GetId())); err != nil {
		log.Panicf(fmt.Sprintf("初始化rpc系统 err:%v", err))
	} else {
		log.Infof("Sys rpc Init success !")
	}
	event.Register(core.Event_ServiceStartEnd, func() { //阻塞 先注册服务集群 保证其他服务能及时发现
		if err := rpc.Start(); err != nil {
			log.Panicf(fmt.Sprintf("启动RPC失败 err:%v", err))
		}
		if err := registry.Start(); err != nil {
			log.Panicf(fmt.Sprintf("加入集群服务失败 err:%v", err))
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
	if err = rpc.Stop(); err != nil {
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
func (this *ClusterService) FindServiceHandlefunc(node registry.ServiceNode) {
	if _, ok := this.serverList.Load(node.Id); !ok {
		if s, err := cbase.NewServiceSession(&node); err != nil {
			log.Errorf("创建服务会话失败【%s】 err:%v", node.Id, err)
		} else {
			this.serverList.Store(node.Id, s)
		}
	}
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

func (this *ClusterService) getServiceSessionByType(sType string, sIp string) (ss []core.IServiceSession, err error) {
	ss = make([]core.IServiceSession, 0)
	if nodes := registry.GetServiceByType(sType); nodes == nil {
		log.Errorf("获取目标类型 type【%s】ip [%s] 服务集失败", sType, sIp)
		return nil, err
	} else {
		if sIp == core.AutoIp {
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
		} else {
			for _, v := range nodes {
				if v.IP == sIp {
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
		}
	}
	return
}
func (this *ClusterService) getServiceSessionByIds(Ids []string) (result []core.IServiceSession, err error) {

	var (
		ss    = make(map[string]core.IServiceSession)
		nodes *registry.ServiceNode
	)
	for _, v := range Ids {
		if nodes, err = registry.GetServiceById(v); err != nil {
			log.Errorf("获取目标类型 ip[%v] 服务集失败:%v", v, err)
			return
		} else {
			if s, ok := this.serverList.Load(v); ok {
				ss[v] = s.(core.IServiceSession)
			} else {
				s, err = cbase.NewServiceSession(nodes)
				if err != nil {
					log.Errorf("创建服务会话失败【%s】 err:%v", v, err)
					continue
				} else {
					this.serverList.Store(v, s)
					ss[v] = s.(core.IServiceSession)
				}
			}
		}
	}
	result = make([]core.IServiceSession, len(ss))
	n := 0
	for _, v := range ss {
		result[n] = v
		n++
	}
	return
}
func (this *ClusterService) getServiceSessionByIps(sType string, sIps []string) (result []core.IServiceSession, err error) {
	var (
		ss    map[string]core.IServiceSession = make(map[string]core.IServiceSession)
		s     core.IServiceSession
		found bool
	)
	//容错处理
	if len(sIps) == 1 && sIps[0] == core.AutoIp {
		s, err = this.ClusterService.DefauleRpcRouteRules(sType, core.AutoIp)
		if err == nil {
			result = []core.IServiceSession{s}
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
					ss[v.Id] = s.(core.IServiceSession)
				} else {
					s, err = cbase.NewServiceSession(v)
					if err != nil {
						log.Errorf("创建服务会话失败【%s】 err:%v", v.Id, err)
						continue
					} else {
						this.serverList.Store(v.Id, s)
						ss[v.Id] = s.(core.IServiceSession)
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

	result = make([]core.IServiceSession, len(ss))
	n := 0
	for _, v := range ss {
		result[n] = v
		n++
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
func (this *ClusterService) DefauleRpcRouteRules(stype string, sip string) (ss core.IServiceSession, err error) {
	if s, e := this.getServiceSessionByType(stype, sip); e != nil {
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
			return nil, fmt.Errorf("未找到IP[%s]类型%s】的服务信息", sip, stype)
		}
	}
}
func (this *ClusterService) RpcInvokeById(sId string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error) {
	defer lego.Recover(fmt.Sprintf("RpcInvokeById sId:%s rkey:%v iscall %v arg %v", sId, rkey, iscall, arg))
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

func (this *ClusterService) RpcInvokeByIds(sId []string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (results map[string]*base.Result, err error) {
	var (
		ss   []core.IServiceSession
		lock sync.Mutex
		wg   *sync.WaitGroup
	)
	defer lego.Recover(fmt.Sprintf("RpcInvokeByIds sId:%v rkey:%v iscall %v arg %v", sId, rkey, iscall, arg))
	if ss, err = this.getServiceSessionByIds(sId); err != nil {
		log.Errorf("未找到目标服务 ip:%s 节点 err:%v", sId, err)
		return
	}
	results = make(map[string]*base.Result)
	if !iscall {
		for _, v := range ss {
			results[v.GetId()] = &base.Result{
				Index: v.GetId(),
				Err:   v.CallNR(rkey, arg...),
			}
		}
	} else {
		wg = new(sync.WaitGroup)
		wg.Add(len(ss))
		for _, v := range ss {
			go func(ss core.IServiceSession, wg *sync.WaitGroup) {
				result, e := ss.Call(rkey, arg...)
				lock.Lock()
				results[ss.GetId()] = &base.Result{Index: ss.GetId(), Result: result, Err: e}
				lock.Unlock()
				wg.Done()
			}(v, wg)
		}
		wg.Wait()
	}
	return
}

func (this *ClusterService) RpcInvokeByType(sType string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error) {
	defer lego.Recover(fmt.Sprintf("RpcInvokeByType sType:%s rkey:%v iscall %v arg %v", sType, rkey, iscall, arg))
	ss, err := this.ClusterService.DefauleRpcRouteRules(sType, core.AutoIp)
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
func (this *ClusterService) RpcInvokeByIp(sIp, sType string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error) {
	defer lego.Recover(fmt.Sprintf("RpcInvokeByIp sIp:%s sType:%s rkey:%v iscall %v arg %v", sIp, sType, rkey, iscall, arg))
	ss, err := this.ClusterService.DefauleRpcRouteRules(sType, sIp)
	if err != nil {
		log.Errorf("未找到目标服务 ip:%s type:%s 节点 err:%v", sIp, sType, err)
		return nil, err
	}
	if iscall {
		result, err = ss.Call(rkey, arg...)
	} else {
		err = ss.CallNR(rkey, arg...)
	}
	return
}

func (this *ClusterService) RpcInvokeByIps(sIp []string, sType string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (results map[string]*base.Result, err error) {
	var (
		ss   []core.IServiceSession
		lock sync.Mutex
		wg   *sync.WaitGroup
	)
	defer lego.Recover(fmt.Sprintf("RpcInvokeByIps sIp:%v sType:%s rkey:%v iscall %v arg %v", sIp, sType, rkey, iscall, arg))
	if ss, err = this.getServiceSessionByIps(sType, sIp); err != nil {
		log.Errorf("未找到目标服务 ip:%s type:%s 节点 err:%v", sIp, sType, err)
		return
	}
	if len(ss) == 0 {
		log.Errorf("未找到目标服务 ip:%s type:%s 节点", sIp, sType)
		err = fmt.Errorf("未找到目标服务 ip:%s type:%s 节点", sIp, sType)
		return
	}
	results = make(map[string]*base.Result)
	if !iscall {
		for _, v := range ss {
			results[v.GetIp()] = &base.Result{
				Err: v.CallNR(rkey, arg...),
			}
		}
	} else {
		wg = new(sync.WaitGroup)
		wg.Add(len(ss))
		for _, v := range ss {
			go func(ss core.IServiceSession, wg *sync.WaitGroup) {
				result, e := ss.Call(rkey, arg...)
				lock.Lock()
				results[ss.GetIp()] = &base.Result{Index: ss.GetIp(), Result: result, Err: e}
				lock.Unlock()
				wg.Done()
			}(v, wg)
		}
		wg.Wait()
	}
	return
}

func (this *ClusterService) ReleaseRpc(rkey core.Rpc_Key, arg ...interface{}) {
	defer lego.Recover(fmt.Sprintf("ReleaseRpc rkey:%v arg %v", rkey, arg))
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
