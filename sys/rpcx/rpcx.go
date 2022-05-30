package rpcx

func newSys(options Options) (sys *RPCX, err error) {
	sys = &RPCX{
		options: options,
		service: newService(options),
	}
	return
}

type RPCX struct {
	options Options
	service IRPCXServer
}

func (this *RPCX) Start() (err error) {
	this.service.Start()
	return
}

func (this *RPCX) Stop() (err error) {
	err = this.service.Stop()
	return
}

func (this *RPCX) Register(rcvr interface{}) (err error) {
	err = this.service.Register(rcvr)
	return
}

func (this *RPCX) RegisterFunction(fn interface{}) (err error) {
	err = this.service.RegisterFunction(fn)
	return
}

func (this *RPCX) RegisterFunctionName(name string, fn interface{}) (err error) {
	err = this.service.RegisterFunctionName(name, fn)
	return
}

func (this *RPCX) UnregisterAll() (err error) {
	return this.service.UnregisterAll()
}

func (this *RPCX) NewRpcClient(addr, sId string) (clent IRPCXClient, err error) {
	return newClient(addr, sId)
}
