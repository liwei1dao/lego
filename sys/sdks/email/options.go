package email

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

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

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	if options.Serverhost == "" || options.Fromemail == "" || options.Fompasswd == "" {
		log.Panicf("start email Missing necessary configuration : Serverhost:%s Fromemail:%s Fompasswd:%s", options.Serverhost, options.Fromemail, options.Fompasswd)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	if options.Serverhost == "" || options.Fromemail == "" || options.Fompasswd == "" {
		log.Panicf("start email Missing necessary configuration : Serverhost:%s Fromemail:%s Fompasswd:%s", options.Serverhost, options.Fromemail, options.Fompasswd)
	}
	return options
}
