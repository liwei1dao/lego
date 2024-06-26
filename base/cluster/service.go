package cluster

import (
	"context"
	"fmt"

	"github.com/liwei1dao/lego/base"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/cron"
	"github.com/liwei1dao/lego/sys/event"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc"
)

type ClusterService struct {
	cbase.ServiceBase
	option         *Options
	serviceNode    *core.ServiceNode
	clusterService base.IClusterService
}

func (this *ClusterService) GetTag() string {
	return this.option.Setting.Tag
}
func (this *ClusterService) GetId() string {
	return this.option.Setting.Id
}
func (this *ClusterService) GetType() string {
	return this.option.Setting.Type
}
func (this *ClusterService) GetVersion() string {
	return this.option.Version
}
func (this *ClusterService) GetSettings() core.ServiceSttings {
	return this.option.Setting
}
func (this *ClusterService) Options() *Options {
	return this.option
}
func (this *ClusterService) Configure(option ...Option) {
	this.option = newOptions(option...)
	this.serviceNode = &core.ServiceNode{
		Tag:     this.option.Setting.Tag,
		Id:      this.option.Setting.Id,
		Type:    this.option.Setting.Type,
		Version: this.option.Version,
		Meta:    make(map[string]string),
	}
}

func (this *ClusterService) Init(service core.IService) (err error) {
	this.clusterService = service.(base.IClusterService)
	return this.ServiceBase.Init(service)
}

func (this *ClusterService) InitSys() {
	if err := log.OnInit(this.option.Setting.Sys["log"]); err != nil {
		panic(fmt.Sprintf("sys log Init err:%v", err))
	} else {
		log.Infof("sys log Init success !")
	}
	if err := event.OnInit(this.option.Setting.Sys["event"]); err != nil {
		log.Panicf(fmt.Sprintf("sys event Init err:%v", err))
	} else {
		log.Infof("sys event Init success !")
	}
	if err := rpc.OnInit(this.option.Setting.Sys["rpc"], rpc.SetServiceNode(this.serviceNode)); err != nil {
		log.Panicf(fmt.Sprintf("初始化rpc系统 err:%v", err))
	} else {
		log.Infof("sys rpc Init success !")
	}
	event.Register(core.Event_ServiceStartEnd, func() { //阻塞 先注册服务集群 保证其他服务能及时发现
		if err := rpc.Start(); err != nil {
			log.Panicf(fmt.Sprintf("sys rpc Init err:%v", err))
		}
	})
}

func (this *ClusterService) Destroy() (err error) {
	if err = rpc.Close(); err != nil {
		return
	}
	cron.Close()
	err = this.ServiceBase.Destroy()
	return
}

// 注册服务对象
func (this *ClusterService) Register(rcvr interface{}) (err error) {
	return rpc.Register(rcvr)
}

// 注册服务方法
func (this *ClusterService) RegisterFunction(fn interface{}) (err error) {
	return rpc.RegisterFunction(fn)
}

// 注册服务方法 自定义服务名
func (this *ClusterService) RegisterFunctionName(name string, fn interface{}) (err error) {
	return rpc.RegisterFunctionName(name, fn)
}

// 调用远端服务接口 同步
func (this *ClusterService) RpcCall(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) (err error) {
	err = rpc.Call(ctx, servicePath, serviceMethod, args, reply)
	return
}

// 调用远端服务接口 异步
func (this *ClusterService) RpcGo(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) (call *rpc.MessageCall, err error) {
	return rpc.Go(ctx, servicePath, serviceMethod, args, reply)
}

// 广播调用远端服务接口
func (this *ClusterService) Broadcast(ctx context.Context, servicePath, serviceMethod string, args interface{}) (err error) {
	return rpc.Broadcast(ctx, servicePath, serviceMethod, args)
}
