package rpcx

import (
	"context"
	"fmt"
	"runtime"

	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/cron"
	"github.com/liwei1dao/lego/sys/event"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/registry"
	"github.com/liwei1dao/lego/sys/rpcx"
	"github.com/smallnest/rpcx/client"
)

type RPCXService struct {
	cbase.ServiceBase
	opts          *Options
	rpcxService   base.IRPCXService
	selector      base.ISelector //选择器
	isInClustered bool
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

func (this *RPCXService) Options() *Options {
	return this.opts
}
func (this *RPCXService) Configure(opts ...Option) {
	this.opts = newOptions(opts...)
}

func (this *RPCXService) Init(service core.IService) (err error) {
	this.rpcxService = service.(base.IRPCXService)
	this.selector = base.NewSelector(base.WeightedRoundRobin)
	return this.ServiceBase.Init(service)
}

func (this *RPCXService) InitSys() {
	if err := log.OnInit(this.opts.Setting.Sys["log"]); err != nil {
		panic(fmt.Sprintf("Sys log Init err:%v", err))
	} else {
		log.Infof("Sys log Init success !")
	}
	if err := event.OnInit(this.opts.Setting.Sys["event"]); err != nil {
		log.Panicf(fmt.Sprintf("Sys event Init err:%v", err))
	} else {
		log.Infof("Sys event Init success !")
	}
	if err := registry.OnInit(this.opts.Setting.Sys["registry"], registry.SetService(this.rpcxService), registry.SetListener(this.rpcxService.(registry.IListener))); err != nil {
		log.Panicf(fmt.Sprintf("Sys registry Init err:%v", err))
	} else {
		log.Infof("Sys registry Init success !")
	}
	if err := rpcx.OnInit(this.opts.Setting.Sys["rpcx"], rpcx.SetServiceId(this.GetId()), rpcx.SetPort(this.GetPort())); err != nil {
		log.Panicf(fmt.Sprintf("Sys rpcx Init err:%v", err))
	} else {
		log.Infof("Sys rpcx Init success !")
	}
	event.Register(core.Event_ServiceStartEnd, func() { //阻塞 先注册服务集群 保证其他服务能及时发现
		if err := rpcx.Start(); err != nil {
			log.Panicf(fmt.Sprintf("Sys rpcx Start err:%v", err))
		}
		if err := registry.Start(); err != nil {
			log.Panicf(fmt.Sprintf("Sys registry Start err:%v", err))
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
	modules := make([]core.IModule, 0)
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
func (this *RPCXService) FindServiceHandlefunc(node core.ServiceNode) {
	if !this.selector.IsHave(node.Id) {
		if s, err := NewServiceSession(&node); err != nil {
			log.Errorf("创建服务会话失败【%s】 err:%v", node.Id, err)
		} else {
			this.selector.AddServer(s)
		}
	}
	if this.isInClustered {
		event.TriggerEvent(core.Event_FindNewService, node) //触发发现新的服务事件
	} else {
		if node.Id == this.opts.Setting.Id { //发现自己 加入集群成功
			this.isInClustered = true
			event.TriggerEvent(core.Event_RegistryStart)
		}
	}
}

//更新服务会话 当有新的服务加入时
func (this *RPCXService) UpDataServiceHandlefunc(node core.ServiceNode) {
	if this.selector.IsHave(node.Id) {
		this.selector.UpdateServer(node)
		event.TriggerEvent(core.Event_UpDataOldService, node)
	}
}

//注销服务会话
func (this *RPCXService) LoseServiceHandlefunc(sId string) {
	if this.selector.IsHave(sId) {
		this.selector.RemoveServer(sId)
		event.TriggerEvent(core.Event_LoseService, sId)
	}
}

//注册服务对象
func (this *RPCXService) Register(rcvr interface{}) (err error) {
	err = rpcx.Register(rcvr)
	return
}

//注册服务方法
func (this *RPCXService) RegisterFunction(fn interface{}) (err error) {
	err = rpcx.RegisterFunction(fn)
	return
}

//注册服务方法 自定义方法名称
func (this *RPCXService) RegisterFunctionName(name string, fn interface{}) (err error) {
	err = rpcx.RegisterFunctionName(name, fn)
	return
}

//同步 执行目标远程服务方法
func (this *RPCXService) RpcCallById(sId string, serviceMethod string, ctx context.Context, args interface{}, reply interface{}) (err error) {
	defer lego.Recover(fmt.Sprintf("RpcCallById sId:%s rkey:%v arg %v", sId, serviceMethod, args))
	ss, ok := this.selector.Load(sId)
	if !ok {
		log.Errorf("未找到目标服务【%s】节点 ", sId)
		return fmt.Errorf("No Found " + sId)

	}
	err = ss.(base.IRPCXServiceSession).Call(ctx, serviceMethod, args, reply)
	return
}

//异步 执行目标远程服务方法
func (this *RPCXService) RpcGoById(sId string, serviceMethod string, ctx context.Context, args interface{}, reply interface{}) (call *client.Call, err error) {
	defer lego.Recover(fmt.Sprintf("RpcGoById sId:%s rkey:%v arg %v", sId, serviceMethod, args))
	ss, ok := this.selector.Load(sId)
	if !ok {
		log.Errorf("未找到目标服务【%s】节点", sId)
		return nil, fmt.Errorf("No Found " + sId)
	}
	call, err = ss.(base.IRPCXServiceSession).Go(ctx, serviceMethod, args, reply)
	return
}

func (this *RPCXService) RpcCallByType(sType string, serviceMethod string, ctx context.Context, args interface{}, reply interface{}) (err error) {
	defer lego.Recover(fmt.Sprintf("RpcCallByType sType:%s rkey:%s arg %v", sType, serviceMethod, args))

	ss, err := this.selector.Select(ctx, sType, core.AutoIp)
	if err != nil {
		log.Errorf("未找到目标服务【%s】节点 err:%v", sType, err)
		return err
	}
	err = ss.(base.IRPCXServiceSession).Call(ctx, serviceMethod, args, reply)
	return
}

func (this *RPCXService) RpcGoByType(sType string, serviceMethod string, ctx context.Context, args interface{}, reply interface{}) (call *client.Call, err error) {
	defer lego.Recover(fmt.Sprintf("RpcCallByType sType:%s rkey:%s arg %v", sType, serviceMethod, args))
	ss, err := this.selector.Select(ctx, sType, core.AutoIp)
	if err != nil {
		log.Errorf("未找到目标服务【%s】节点 err:%v", sType, err)
		return nil, err
	}
	call, err = ss.(base.IRPCXServiceSession).Go(ctx, serviceMethod, args, reply)
	return
}
