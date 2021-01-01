package rpc

import (
	"github.com/liwei1dao/lego/core"
)

type IRpcServer interface {
	RpcId() string
	GetRpcInfo() (rfs []core.Rpc_Key)
	Register(id string, f interface{})
	RegisterGO(id string, f interface{})
	UnRegister(id string, f interface{})
	Done() (err error)
}

type IRpcClient interface {
	Done() (err error)
	Call(_func string, params ...interface{}) (interface{}, error)
	CallNR(_func string, params ...interface{}) (err error)
}
