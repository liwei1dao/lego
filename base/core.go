package base

import (
	"context"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc"
	"github.com/smallnest/rpcx/client"
)

type Result struct {
	Index  string
	Result interface{}
	Err    error
}

type ISingleService interface {
	core.IService
}

type IClusterService interface {
	core.IService
	GetTag() string //获取集群标签
	Register(rcvr interface{}) (err error)
	RegisterFunction(fn interface{}) (err error)
	RegisterFunctionName(name string, fn interface{}) (err error)
	RpcCall(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) (err error)
	RpcGo(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) (call *rpc.MessageCall, err error)
	Broadcast(ctx context.Context, servicePath, serviceMethod string, args interface{}) (err error)
}

type IRPCXService interface {
	core.IService
	GetTag() string //获取集群标签
	RegisterFunction(fn interface{}) (err error)
	RegisterFunctionName(name string, fn interface{}) (err error)
	RpcCall(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
	RpcBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
	RpcGo(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *client.Call, err error)
	AcrossClusterRpcCall(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
	AcrossClusterRpcGo(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *client.Call, err error)
	AcrossClusterBroadcast(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
	ClusterBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
}
