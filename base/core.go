package base

import (
	"context"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpcl"
)

type Result struct {
	Index  string
	Result interface{}
	Err    error
}

type ISingleService interface {
	core.IService
}

type IServiceSession interface {
	GetId() string
	GetIp() string
	GetRpcId() string
	GetType() string
	GetVersion() string
	SetVersion(v string)
	GetPreWeight() float64
	SetPreWeight(p float64)
	Done()
}

type IClusterServiceBase interface {
	core.IService
	GetTag() string            //获取集群标签
	SetPreWeight(weight int32) //设置服务器权重
}

type IClusterService interface {
	IClusterServiceBase
	Register(rcvr interface{}) (err error)
	RegisterFunction(fn interface{}) (err error)
	RegisterFunctionName(name string, fn interface{}) (err error)
	RpcCall(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) (err error)
	RpcGo(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *rpcl.MessageCall) (call *rpcl.MessageCall, err error)
	RpcNoCall(ctx context.Context, servicePath, serviceMethod string, args interface{}) (err error)
	Broadcast(ctx context.Context, servicePath, serviceMethod string, args interface{}) (err error)
}
