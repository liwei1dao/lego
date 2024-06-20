package grpc

import (
	"net"

	"google.golang.org/grpc"
)

func newService(options *Options) (sys *Service, err error) {
	sys = &Service{
		options: options,
		s:       grpc.NewServer(),
	}

	return
}

type Service struct {
	options *Options
	s       *grpc.Server
}

func (this *Service) Start() (err error) {
	lis, err := net.Listen("tcp", this.options.ServiceAddr)
	if err = this.s.Serve(lis); err != nil {
		this.options.Log.Errorf("Service Start err:%s", err.Error())
		return
	}
	return
}
