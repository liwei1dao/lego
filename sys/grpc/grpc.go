package grpc

func newSys(options *Options) (sys ISys, err error) {
	var (
		service ISys
	)
	service, err = newService(options)
	sys = &RPCX{
		options: options,
		service: service,
	}
	return
}

type RPCX struct {
	options *Options
	service ISys
	client  ISys
}
