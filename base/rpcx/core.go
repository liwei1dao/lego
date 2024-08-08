package rpcx

import (
	"context"

	"github.com/liwei1dao/lego/core"
	"github.com/smallnest/rpcx/client"
)

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
