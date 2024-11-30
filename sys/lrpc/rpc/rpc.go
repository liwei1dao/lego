package rpc

func newsys(options *Options) (cpool *Ppc, err error) {
	cpool = &Ppc{
		options: options,
	}
	return
}

type Ppc struct {
	options *Options
}
