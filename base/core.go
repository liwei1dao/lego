package base

import (
	"context"
	"time"

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

type IClusterServiceBase interface {
	core.IService
	GetTag() string               //获取集群标签
	GetCategory() core.S_Category //服务类别 例如游戏服
	GetRpcId() string             //获取rpc通信id
	GetPreWeight() float64        //集群服务负载值 暂时可以不用理会
	SetPreWeight(weight int32)    //设置服务器权重
}

type IClusterService interface {
	IClusterServiceBase
	Register(rcvr interface{}) (err error)
	RegisterFunction(fn interface{}) (err error)
	RegisterFunctionName(name string, fn interface{}) (err error)
	RpcCall(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) (err error)
	RpcGo(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) (call *rpc.MessageCall, err error)
	Broadcast(ctx context.Context, servicePath, serviceMethod string, args interface{}) (err error)
}

type IRPCXService interface {
	IClusterServiceBase
	GetOpentime() time.Time
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
