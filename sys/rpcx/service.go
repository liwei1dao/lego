package rpcx

import (
	lgcore "github.com/liwei1dao/lego/core"
	"github.com/smallnest/rpcx/server"
)

func newService(options Options) (s *Service) {
	s = &Service{
		server:  server.NewServer(),
		options: options,
	}
	return
}

type Service struct {
	server  *server.Server
	options Options
}

func (this *Service) Start() (err error) {
	err = this.server.Serve("tcp", this.options.Addr)
	return
}

func (this *Service) Stop() (err error) {
	err = this.server.Close()
	return
}

func (this *Service) Register(id lgcore.Rpc_Key, fn interface{}) (err error) {
	err = this.server.RegisterFunction(string(id), fn, "")
	return
}
func (this *Service) UnRegister(id lgcore.Rpc_Key) (err error) {
	err = this.server.Plugins.DoUnregister(string(id))
	return
}
