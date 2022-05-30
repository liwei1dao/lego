package rpcx

import (
	"context"

	"github.com/smallnest/rpcx/client"
)

type (
	ISys interface {
		IRPCXServer
		NewRpcClient(addr, sId string) (clent IRPCXClient, err error)
	}

	IRPCXServer interface {
		Start() (err error)
		Stop() (err error)
		Register(rcvr interface{}) (err error)
		RegisterFunction(fn interface{}) (err error)
		RegisterFunctionName(name string, fn interface{}) (err error)
		UnregisterAll() (err error)
	}

	IRPCXClient interface {
		Stop() (err error)
		Call(ctx context.Context, serviceMethod string, args interface{}, reply interface{}) error
		Go(ctx context.Context, serviceMethod string, args interface{}, reply interface{}, done chan *client.Call) (*client.Call, error)
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

func Start() (err error) {
	return defsys.Start()
}

func Stop() (err error) {
	return defsys.Stop()
}

func Register(rcvr interface{}) (err error) {
	return defsys.Register(rcvr)
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

func NewRpcClient(addr, sId string) (clent IRPCXClient, err error) {
	return defsys.NewRpcClient(addr, sId)
}
