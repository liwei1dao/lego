package email

type Option func(*Options)
type Options struct {
	Serverhost string
	Fromemail  string
	Fompasswd  string
	Serverport int
}

func SetServerhost(v string) Option {
	return func(o *Options) {
		o.Serverhost = v
	}
}

func SetFromemail(v string) Option {
	return func(o *Options) {
		o.Fromemail = v
	}
}

func SetFompasswd(v string) Option {
	return func(o *Options) {
		o.Fompasswd = v
	}
}

func SetServerport(v int) Option {
	return func(o *Options) {
		o.Serverport = v
	}
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
