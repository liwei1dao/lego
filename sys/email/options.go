package email

import (
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type EmailType uint8

const (
	QQEmail    EmailType = iota //QQ邮箱
	GmailEmail                  //Google邮箱
)

type Option func(*Options)
type Options struct {
	EmailType EmailType
	FromEmail string
	Password  string
}

func Set_EmailType(v EmailType) Option {
	return func(o *Options) {
		o.EmailType = v
	}
}

func Set_FromEmail(v string) Option {
	return func(o *Options) {
		o.FromEmail = v
	}
}

func Set_Password(v string) Option {
	return func(o *Options) {
		o.Password = v
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
	if options.FromEmail == "" || options.Password == "" {
		log.Panicf("start email Missing necessary configuration : Serverhost:%s Fromemail:%s Fompasswd:%s", options.FromEmail, options.EmailType)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	if options.FromEmail == "" || options.Password == "" {
		log.Panicf("start email Missing necessary configuration : Serverhost:%s Fromemail:%s Fompasswd:%s", options.FromEmail, options.EmailType)
	}
	return options
}
