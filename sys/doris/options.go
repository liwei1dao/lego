package doris

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	IP        string
	Port      int32
	User      string //ftp 登录用户名
	Password  string //ftp 登录密码
	DBname    string
	Separator string
}

func SetIP(v string) Option {
	return func(o *Options) {
		o.IP = v
	}
}
func SetPort(v int32) Option {
	return func(o *Options) {
		o.Port = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		IP:        "127.0.0.1",
		Port:      21,
		Separator: ",",
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		IP:        "127.0.0.1",
		Port:      21,
		Separator: ",",
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
