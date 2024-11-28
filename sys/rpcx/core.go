package rpcx

import (
	"context"

	"github.com/smallnest/rpcx/client"
)

const (
	ServiceClusterTag = "ctag"
	ServiceMetaKey    = "smeta"
	ServiceAddrKey    = "addr"
	CallRoutRulesKey  = "callrules"
)

const RpcX_ShakeHands = "RpcX_ShakeHands" //握手

type (
	ISys interface {
		Start() (err error)
		Stop() (err error)
		GetServiceTags() []string
		RegisterFunction(fn interface{}) (err error)
		RegisterFunctionName(name string, fn interface{}) (err error)
		UnregisterAll() (err error)
		Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
		Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
		Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (call *client.Call, err error)
		AcrossClusterCall(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
		AcrossClusterBroadcast(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
		AcrossClusterGo(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (call *client.Call, err error)
		ClusterBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error)
	}
	ISelector interface {
		client.Selector
		Find(ctx context.Context, servicePath, serviceMethod string, args interface{}) []string
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var options *Options
	if options, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(options)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var options *Options
	if options, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(options)
	return
}

func Start() (err error) {
	return defsys.Start()
}

func Stop() (err error) {
	return defsys.Stop()
}

func RegisterFunction(fn interface{}) (err error) {
	return defsys.RegisterFunction(fn)
}
func RegisterFunctionName(name string, fn interface{}) (err error) {
	return defsys.RegisterFunctionName(name, fn)
}

func UnregisterAll() (err error) {
	return defsys.UnregisterAll()
}

func Call(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return defsys.Call(ctx, servicePath, serviceMethod, args, reply)
}

func Broadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return defsys.Broadcast(ctx, servicePath, serviceMethod, args, reply)
}

func Go(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (call *client.Call, err error) {
	return defsys.Go(ctx, servicePath, serviceMethod, args, reply, done)
}
func AcrossClusterCall(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return defsys.AcrossClusterCall(ctx, clusterTag, servicePath, serviceMethod, args, reply)
}
func AcrossClusterGo(ctx context.Context, clusterTag, servicePath string, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (_call *client.Call, err error) {
	return defsys.AcrossClusterGo(ctx, clusterTag, servicePath, serviceMethod, args, reply, done)
}
func AcrossClusterBroadcast(ctx context.Context, clusterTag string, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return defsys.AcrossClusterBroadcast(ctx, clusterTag, servicePath, serviceMethod, args, reply)
}
func ClusterBroadcast(ctx context.Context, servicePath string, serviceMethod string, args interface{}, reply interface{}) (err error) {
	return defsys.ClusterBroadcast(ctx, servicePath, serviceMethod, args, reply)
}
