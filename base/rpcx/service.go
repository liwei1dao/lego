package rpcx

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
	"github.com/liwei1dao/lego/sys/rpcx"
	"github.com/smallnest/rpcx/client"
)

/*
XRPC 服务对象
*/
type RPCXService struct {
	cbase.ServiceBase
	option         *Options
	serviceNode    *core.ServiceNode
	clusterService base.IClusterService
}

func (this *RPCXService) GetTag() string {
	return this.option.Setting.Tag
}
func (this *RPCXService) GetId() string {
	return this.option.Setting.Id
}
func (this *RPCXService) GetType() string {
	return this.option.Setting.Type
}
func (this *RPCXService) GetVersion() string {
	return this.option.Version
}
func (this *RPCXService) GetSettings() core.ServiceSttings {
	return this.option.Setting
}
func (this *RPCXService) Options() *Options {
	return this.option
}
func (this *RPCXService) Configure(option ...Option) {
	this.option = newOptions(option...)
	this.serviceNode = &core.ServiceNode{
		Tag:     this.option.Setting.Tag,
		Id:      this.option.Setting.Id,
		Type:    this.option.Setting.Type,
		Version: this.option.Version,
		Meta:    make(map[string]string),
	}
}

func (this *RPCXService) Init(service core.IService) (err error) {
	this.clusterService = service.(base.IClusterService)
	return this.ServiceBase.Init(service)
}

func (this *RPCXService) InitSys() {
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

func (this *RPCXService) Destroy() (err error) {
	if err = rpc.Close(); err != nil {
		return
	}
	cron.Close()
	err = this.ServiceBase.Destroy()
	return
}

// 注册服务方法
func (this *RPCXService) RegisterFunction(fn interface{}) (err error) {
	err = rpcx.RegisterFunction(fn)
	return
}

// 注册服务方法 自定义方法名称
func (this *RPCXService) RegisterFunctionName(name string, fn interface{}) (err error) {
	err = rpcx.RegisterFunctionName(name, fn)
	return
}

// /同步 执行目标远程服务方法
// /servicePath = worker  表示采用负载的方式调用 worker类型服务执行rpc方法
// /servicePath = worker/worker_1  表示寻找目标服务节点调用rpc方法
// /servicePath = worker/!worker_1  表示选择非worker_1的节点随机选择节点执行rpc方法
// /servicePath = worker/[worker_1,worker_2]  表示随机选择[]里面的服务节点执行rpc方法
// /servicePath = worker/![worker_1,worker_2]  表示随机选择非[]里面的服务节点执行rpc方法
func (this *RPCXService) RpcCall(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return rpcx.Call(ctx, servicePath, serviceMethod, args, reply)
}

// /广播 执行目标远程服务方法
// /servicePath = worker  表示采用负载的方式调用 worker类型服务执行rpc方法
func (this *RPCXService) RpcBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return rpcx.Broadcast(ctx, servicePath, serviceMethod, args, reply)
}

// /异步 执行目标远程服务方法
// /servicePath = worker  表示采用负载的方式调用 worker类型服务执行rpc方法
// /servicePath = worker/worker_1  表示寻找目标服务节点调用rpc方法
// /servicePath = worker/!worker_1  表示选择非worker_1的节点随机选择节点执行rpc方法
// /servicePath = worker/[worker_1,worker_2]  表示随机选择[]里面的服务节点执行rpc方法
// /servicePath = worker/![worker_1,worker_2]  表示随机选择非[]里面的服务节点执行rpc方法
func (this *RPCXService) RpcGo(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *client.Call, err error) {
	return rpcx.Go(ctx, servicePath, serviceMethod, args, reply, nil)
}

// /跨集群 同步 执行目标远程服务方法
// clusterTag 集群标签
// /servicePath = worker  表示采用负载的方式调用 worker类型服务执行rpc方法
// /servicePath = worker/worker_1  表示寻找目标服务节点调用rpc方法
// /servicePath = worker/!worker_1  表示选择非worker_1的节点随机选择节点执行rpc方法
// /servicePath = worker/[worker_1,worker_2]  表示随机选择[]里面的服务节点执行rpc方法
// /servicePath = worker/![worker_1,worker_2]  表示随机选择非[]里面的服务节点执行rpc方法
func (this *RPCXService) AcrossClusterRpcCall(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return rpcx.AcrossClusterCall(ctx, clusterTag, servicePath, serviceMethod, args, reply)
}

// /跨集群 异步 执行目标远程服务方法
// clusterTag 集群标签
// /servicePath = worker  表示采用负载的方式调用 worker类型服务执行rpc方法
// /servicePath = worker/worker_1  表示寻找目标服务节点调用rpc方法
// /servicePath = worker/!worker_1  表示选择非worker_1的节点随机选择节点执行rpc方法
// /servicePath = worker/[worker_1,worker_2]  表示随机选择[]里面的服务节点执行rpc方法
// /servicePath = worker/![worker_1,worker_2]  表示随机选择非[]里面的服务节点执行rpc方法
func (this *RPCXService) AcrossClusterRpcGo(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *client.Call, err error) {
	return rpcx.AcrossClusterGo(ctx, clusterTag, servicePath, serviceMethod, args, reply, nil)
}

// 目标集群广播消息
// clusterTag 集群标签
// /servicePath = worker  表示采用负载的方式调用 worker类型服务执行rpc方法
// /servicePath = worker/worker_1  表示寻找目标服务节点调用rpc方法
// /servicePath = worker/!worker_1  表示选择非worker_1的节点随机选择节点执行rpc方法
// /servicePath = worker/[worker_1,worker_2]  表示随机选择[]里面的服务节点执行rpc方法
// /servicePath = worker/![worker_1,worker_2]  表示随机选择非[]里面的服务节点执行rpc方法
func (this *RPCXService) AcrossClusterBroadcast(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return rpcx.AcrossClusterBroadcast(ctx, clusterTag, servicePath, serviceMethod, args, reply)
}

// /全集群广播  执行目标远程服务方法
// /servicePath = worker  表示采用负载的方式调用 worker类型服务执行rpc方法
func (this *RPCXService) ClusterBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return rpcx.ClusterBroadcast(ctx, servicePath, serviceMethod, args, reply)
}
