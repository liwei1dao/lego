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
	"github.com/liwei1dao/lego/sys/rpcl"
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
	}
}

func (this *ClusterService) Init(service core.IService) (err error) {
	this.clusterService = service.(base.IClusterService)
	return this.ServiceBase.Init(service)
}

func (this *ClusterService) InitSys() {
	if err := log.OnInit(this.option.Setting.Sys["log"]); err != nil {
		panic(fmt.Sprintf("初始化log系统失败 err:%v", err))
	} else {
		log.Infof("Sys log Init success !")
	}
	if err := event.OnInit(this.option.Setting.Sys["event"]); err != nil {
		log.Panicf(fmt.Sprintf("初始化event系统失败 err:%v", err))
	} else {
		log.Infof("Sys event Init success !")
	}
	if err := rpcl.OnInit(this.option.Setting.Sys["rpcl"], rpcl.SetServiceNode(this.serviceNode)); err != nil {
		log.Panicf(fmt.Sprintf("初始化rpc系统 err:%v", err))
	} else {
		log.Infof("Sys rpc Init success !")
	}
	event.Register(core.Event_ServiceStartEnd, func() { //阻塞 先注册服务集群 保证其他服务能及时发现
		if err := rpcl.Start(); err != nil {
			log.Panicf(fmt.Sprintf("启动RPC失败 err:%v", err))
		}
	})
}

func (this *ClusterService) Destroy() (err error) {
	if err = rpcl.Close(); err != nil {
		return
	}
	cron.Close()
	err = this.ServiceBase.Destroy()
	return
}

//注册服务对象
func (this *ClusterService) Register(rcvr interface{}) (err error) {
	return
}

//注册服务方法
func (this *ClusterService) RegisterFunction(fn interface{}) (err error) {
	return
}

//注册服务方法 自定义服务名
func (this *ClusterService) RegisterFunctionName(name string, fn interface{}) (err error) {
	return
}

//调用远端服务接口 同步
func (this *ClusterService) RpcCall(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return
}

//调用远端服务接口 异步
func (this *ClusterService) RpcGo(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *rpcl.MessageCall) (call *rpcl.MessageCall, err error) {
	return
}

//调用远端服务接口 无回应
func (this *ClusterService) RpcNoCall(ctx context.Context, servicePath, serviceMethod string, args interface{}) (err error) {
	return
}

//广播调用远端服务接口
func (this *ClusterService) Broadcast(ctx context.Context, servicePath, serviceMethod string, args interface{}) (err error) {
	return
}
