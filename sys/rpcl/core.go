package rpcl

import (
	"context"

	"github.com/liwei1dao/lego/sys/rpcl/core"
)

type (
	ISys interface {
		Register(rcvr interface{}) error
		RegisterFunction(fn interface{}) error
		RegisterFunctionName(name string, fn interface{}) error
		UnRegister(name string)
		Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)                                      //同步调用 等待结果
		GoNR(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error)                                                         //异步调用 无返回
		Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *core.Call) (call *core.Call, err error) //异步调用 异步返回
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}
