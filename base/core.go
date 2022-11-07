package base

import (
	"context"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc"
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
	RpcNoCall(ctx context.Context, servicePath, serviceMethod string, args interface{}) (err error)
	Broadcast(ctx context.Context, servicePath, serviceMethod string, args interface{}) (err error)
}
