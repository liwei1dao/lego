package rpcx

import (
	"context"

	lgcore "github.com/liwei1dao/lego/core"
	"github.com/smallnest/rpcx/client"
)

type (
	ISys interface {
		IRPCXServer
		NewRpcClient(addr string) (clent IRPCXClient, err error)
	}

	IRPCXServer interface {
		Start() (err error)
		Stop() (err error)
		Register(id lgcore.Rpc_Key, f interface{}) (err error)
		UnRegister(id lgcore.Rpc_Key) (err error)
	}

	IRPCXClient interface {
		Stop() (err error)
		Call(ctx context.Context, serviceMethod lgcore.Rpc_Key, args interface{}, reply interface{}) error
		Go(ctx context.Context, serviceMethod lgcore.Rpc_Key, args interface{}, reply interface{}, done chan *client.Call) (*client.Call, error)
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

func Register(id lgcore.Rpc_Key, f interface{}) (err error) {
	return defsys.Register(id, f)
}
func UnRegister(id lgcore.Rpc_Key) (err error) {
	return defsys.UnRegister(id)
}

func NewRpcClient(addr string) (clent IRPCXClient, err error) {
	return defsys.NewRpcClient(addr)
}
