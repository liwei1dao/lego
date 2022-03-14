package ftp

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type FtpServerType int8

const (
	FTP FtpServerType = iota
	SFTP
)

type Option func(*Options)
type Options struct {
	SType    FtpServerType //服务类型
	IP       string        //ftp 地址
	Port     int32         //ftp 访问端口		ftp:21 sftp:22
	User     string        //ftp 登录用户名
	Password string        //ftp 登录密码
}

func SetSType(v FtpServerType) Option {
	return func(o *Options) {
		o.SType = v
	}
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
func SetUser(v string) Option {
	return func(o *Options) {
		o.User = v
	}
}
func SetPassword(v string) Option {
	return func(o *Options) {
		o.Password = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		SType: FTP,
		IP:    "127.0.0.1",
		Port:  21,
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
		IP:   "127.0.0.1",
		Port: 21,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
