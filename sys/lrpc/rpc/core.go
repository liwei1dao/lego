package rpc

import (
	"context"

	"github.com/liwei1dao/lego/core"
)

type (
	//选择器
	ISelector interface {
		Select(ctx context.Context, servicePath string) []core.IServiceNode
		UpdateServer(servers []core.IServiceNode)
	}

	ISys interface {
		Start() (err error)
		Close() (err error)
		Register(rcvr interface{}) error
		RegisterFunction(fn interface{}) error
		RegisterFunctionName(name string, fn interface{}) (err error)
		UnRegister(name string)
		Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)           //同步调用 等待结果
		Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *Call, err error) //异步调用 异步返回
		Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error)
	}
)

var (
	defsys ISys
)
