package rpcx

import (
	"fmt"

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
	err = this.server.Serve("tcp", fmt.Sprintf(":%d", this.options.Port))
	return
}

func (this *Service) Stop() (err error) {
	err = this.server.Close()
	return
}

func (this *Service) Register(rcvr interface{}) (err error) {
	err = this.server.RegisterName(this.options.ServiceId, rcvr, "")
	return
}
func (this *Service) RegisterFunction(fn interface{}) (err error) {
	err = this.server.RegisterFunction(this.options.ServiceId, fn, "")
	return
}
func (this *Service) RegisterFunctionName(name string, fn interface{}) (err error) {
	err = this.server.RegisterFunctionName(this.options.ServiceId, name, fn, "")
	return
}
func (this *Service) UnregisterAll() (err error) {
	err = this.server.UnregisterAll()
	return
}
