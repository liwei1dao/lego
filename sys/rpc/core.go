package rpc

type IRpcServer interface {
	RpcId() string
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
