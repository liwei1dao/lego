package paypal

func newSys(options *Options) (sys *Pay, err error) {
	sys = &Pay{options: options}

	return
}

type Pay struct {
	options *Options
}
