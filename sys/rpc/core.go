package rpc

import (
	"context"

	"github.com/liwei1dao/lego/core"
)

/*
系统描述:rpc服务通信系统，参考xrpc设计思路
*/
//系统事件
const (
	///发现新的节点
	Event_RpcDiscoverNewNodes core.Event_Key = "Event_RpcDiscoverNewNodes"
	///丢失节点
	Event_RpcLoseNodes core.Event_Key = "Event_RpcLoseNodes"
	///节点属性变更
	Event_RpChangeNodes core.Event_Key = "Event_RpChangeNodes"
)

type (
	ISys interface {
		Start() (err error)
		Close() (err error)
		Register(name string, fn interface{}) (err error)
		UnRegister() (err error)
		Call(ctx context.Context, service string, args interface{}, reply interface{}) (err error)                  //同步调用 等待结果
		Go(ctx context.Context, service string, args interface{}, reply interface{}) (call *MessageCall, err error) //异步调用 异步返回
		Broadcast(ctx context.Context, service string, args interface{}) (err error)
	}

	IClient interface {
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

func Register(name string, fn interface{}) error {
	return defsys.Register(name, fn)
}
func UnRegister() {
	defsys.UnRegister()
}
func Call(ctx context.Context, service string, args interface{}, reply interface{}) (err error) {
	return defsys.Call(ctx, service, args, reply)
}
func Go(ctx context.Context, service string, args interface{}, reply interface{}) (call *MessageCall, err error) {
	return defsys.Go(ctx, service, args, reply)
}
func Broadcast(ctx context.Context, service string, args interface{}) (err error) {
	return defsys.Broadcast(ctx, service, args)
}
