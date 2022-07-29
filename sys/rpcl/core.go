package rpcl

import (
	"context"
)

type (
	ISys interface {
		Start() (err error)
		Close() (err error)
		Register(rcvr interface{}) error
		RegisterFunction(fn interface{}) error
		RegisterFunctionName(name string, fn interface{}) (err error)
		UnRegister(name string)
		Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)                  //同步调用 等待结果
		GoNR(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error)                                     //异步调用 无返回
		Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *MessageCall, err error) //异步调用 异步返回
		Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

func Start() (err error) {
	return defsys.Start()
}
func Close() (err error) {
	return defsys.Close()
}

func Register(rcvr interface{}) error {
	return defsys.Register(rcvr)
}
func RegisterFunction(fn interface{}) error {
	return defsys.RegisterFunction(fn)
}
func RegisterFunctionName(name string, fn interface{}) error {
	return defsys.RegisterFunctionName(name, fn)
}
func UnRegister(name string) {
	defsys.UnRegister(name)
}
func Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return defsys.Call(ctx, servicePath, serviceMethod, args, reply)
}
func GoNR(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error) {
	return defsys.GoNR(ctx, servicePath, serviceMethod, args)
}
func Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (call *MessageCall, err error) {
	return defsys.Go(ctx, servicePath, serviceMethod, args, reply)
}
func Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}) (err error) {
	return defsys.Broadcast(ctx, servicePath, serviceMethod, args)
}
