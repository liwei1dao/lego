package rpcx

import (
	"context"
	"strings"
	"time"

	"github.com/smallnest/rpcx/client"
)

func newSys(options *Options) (sys ISys, err error) {
	// if options.RpcxStartType == RpcxStartByService { //创建RPCX 服务端
	// 	sys, err = newService(options)
	// 	return
	// }

	// if options.RpcxStartType == RpcxStartByClient { //创建RPCX 客户端
	// 	sys, err = newClient(options)
	// 	return
	// }
	var (
		service ISys
		client  ISys
	)
	service, err = newService(options)
	client, err = newClient(options)
	sys = &RPCX{
		options: options,
		service: service,
		client:  client,
	}
	return
}

type RPCX struct {
	options *Options
	service ISys
	client  ISys
}

func (this *RPCX) Start() (err error) {
	this.service.Start()
	this.client.Start()

	return
}

func (this *RPCX) Stop() (err error) {
	this.service.Stop()
	this.client.Stop()
	return
}

// 获取服务集群列表
func (this *RPCX) GetServiceTags() []string {
	return this.GetServiceTags()
}
func (this *RPCX) RegisterFunction(fn interface{}) (err error) {
	this.service.RegisterFunction(fn)
	this.client.RegisterFunction(fn)
	return
}

func (this *RPCX) RegisterFunctionName(name string, fn interface{}) (err error) {
	this.service.RegisterFunctionName(name, fn)
	this.client.RegisterFunctionName(name, fn)
	return
}

func (this *RPCX) UnregisterAll() (err error) {
	err = this.service.UnregisterAll()
	err = this.client.UnregisterAll()
	return
}

// 同步调用
func (this *RPCX) Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		cancel context.CancelFunc
	)
	if this.options.OutTime > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(this.options.OutTime)*time.Second)
		defer cancel()
	}

	if this.options.RpcxStartType == RpcxStartByService {
		err = this.service.Call(ctx, servicePath, serviceMethod, args, reply)
		return
	}

	if this.options.RpcxStartType == RpcxStartByClient {
		err = this.client.Call(ctx, servicePath, serviceMethod, args, reply)
		return
	}

	//先排查下 服务端是否存在连接对象 不存在 在使用客户端对象连接
	err = this.service.Call(ctx, servicePath, serviceMethod, args, reply)
	if err != nil && strings.Contains(err.Error(), "on found") {
		return this.client.Call(ctx, servicePath, serviceMethod, args, reply)
	}
	return
}

// 广播调用
func (this *RPCX) Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		cancel context.CancelFunc
	)
	if this.options.OutTime > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(this.options.OutTime)*time.Second)
		defer cancel()
	}

	if this.options.RpcxStartType == RpcxStartByService {
		err = this.service.Broadcast(ctx, servicePath, serviceMethod, args, reply)
		return
	}
	if this.options.RpcxStartType == RpcxStartByClient {
		err = this.client.Broadcast(ctx, servicePath, serviceMethod, args, reply)
		return
	}

	err = this.service.Broadcast(ctx, servicePath, serviceMethod, args, reply)
	if err != nil && strings.Contains(err.Error(), "on found") {
		return this.client.Broadcast(ctx, servicePath, serviceMethod, args, reply)
	}
	return
}

// 异步调用
func (this *RPCX) Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (call *client.Call, err error) {

	if this.options.RpcxStartType == RpcxStartByService {
		call, err = this.service.Go(ctx, servicePath, serviceMethod, args, reply, done)
		return
	}
	if this.options.RpcxStartType == RpcxStartByClient {
		call, err = this.client.Go(ctx, servicePath, serviceMethod, args, reply, done)
		return
	}

	call, err = this.service.Go(ctx, servicePath, serviceMethod, args, reply, done)
	if err != nil && strings.Contains(err.Error(), "on found") {
		return this.client.Go(ctx, servicePath, serviceMethod, args, reply, done)
	}
	return
}

// 跨服同步调用
func (this *RPCX) AcrossClusterCall(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		cancel context.CancelFunc
	)
	if this.options.OutTime > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(this.options.OutTime)*time.Second)
		defer cancel()
	}
	if this.options.RpcxStartType == RpcxStartByService {
		err = this.service.AcrossClusterCall(ctx, clusterTag, servicePath, serviceMethod, args, reply)
		return
	}
	if this.options.RpcxStartType == RpcxStartByClient {
		err = this.client.AcrossClusterCall(ctx, clusterTag, servicePath, serviceMethod, args, reply)
		return
	}
	err = this.service.AcrossClusterCall(ctx, clusterTag, servicePath, serviceMethod, args, reply)
	if err != nil && strings.Contains(err.Error(), "on found") {
		return this.client.AcrossClusterCall(ctx, clusterTag, servicePath, serviceMethod, args, reply)
	}
	return
}

// 跨集群 广播
func (this *RPCX) AcrossClusterBroadcast(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		cancel context.CancelFunc
	)
	if this.options.OutTime > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(this.options.OutTime)*time.Second)
		defer cancel()
	}
	if this.options.RpcxStartType == RpcxStartByService {
		err = this.service.AcrossClusterBroadcast(ctx, clusterTag, servicePath, serviceMethod, args, reply)
		return
	}
	if this.options.RpcxStartType == RpcxStartByClient {
		err = this.client.AcrossClusterBroadcast(ctx, clusterTag, servicePath, serviceMethod, args, reply)
		return
	}
	err = this.service.AcrossClusterBroadcast(ctx, clusterTag, servicePath, serviceMethod, args, reply)
	if err != nil && strings.Contains(err.Error(), "on found") {
		return this.client.AcrossClusterBroadcast(ctx, clusterTag, servicePath, serviceMethod, args, reply)
	}
	return
}

// 跨服异步调用
func (this *RPCX) AcrossClusterGo(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (call *client.Call, err error) {
	if this.options.RpcxStartType == RpcxStartByService {
		call, err = this.service.AcrossClusterGo(ctx, clusterTag, servicePath, serviceMethod, args, reply, done)
		return
	}
	if this.options.RpcxStartType == RpcxStartByClient {
		call, err = this.client.AcrossClusterGo(ctx, clusterTag, servicePath, serviceMethod, args, reply, done)
		return
	}

	call, err = this.service.AcrossClusterGo(ctx, clusterTag, servicePath, serviceMethod, args, reply, done)
	if err != nil && strings.Contains(err.Error(), "on found") {
		return this.client.AcrossClusterGo(ctx, clusterTag, servicePath, serviceMethod, args, reply, done)
	}
	return
}

// 全集群广播
func (this *RPCX) ClusterBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	var (
		cancel context.CancelFunc
	)
	if this.options.OutTime > 0 {
		ctx, cancel = context.WithTimeout(ctx, time.Duration(this.options.OutTime)*time.Second)
		defer cancel()
	}
	if this.options.RpcxStartType == RpcxStartByService {
		err = this.service.ClusterBroadcast(ctx, servicePath, serviceMethod, args, reply)
		return
	}
	if this.options.RpcxStartType == RpcxStartByClient {
		err = this.client.ClusterBroadcast(ctx, servicePath, serviceMethod, args, reply)
		return
	}
	err = this.service.ClusterBroadcast(ctx, servicePath, serviceMethod, args, reply)
	if err != nil && strings.Contains(err.Error(), "on found") {
		return this.client.ClusterBroadcast(ctx, servicePath, serviceMethod, args, reply)
	}
	return
}
