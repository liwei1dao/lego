package rpcx

import (
	lgcore "github.com/liwei1dao/lego/core"
)

func newSys(options Options) (sys *RPCX, err error) {
	var (
		service *Service
	)
	sys = &RPCX{
		options: options,
	}
	service = &Service{}
	sys.service = service
	return
}

type RPCX struct {
	options Options
	service IRPCXServer
}

func (this *RPCX) Start() (err error) {
	return
}

func (this *RPCX) Stop() (err error) {
	return
}

func (this *RPCX) NewRpcClient(addr string) (clent IRPCXClient, err error) {
	clent, err = newClient(addr)
	return
}

func (this *RPCX) Register(id lgcore.Rpc_Key, fn interface{}) (err error) {
	err = this.service.Register(id, fn)
	return
}

func (this *RPCX) UnRegister(id lgcore.Rpc_Key) (err error) {
	return this.service.UnRegister(id)
}
