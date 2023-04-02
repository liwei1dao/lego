package pay

import (
	"fmt"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	AppId        string //app id
	AliPublicKey string //阿里公钥
	PrivateKey   string //私钥
	NotifyURL    string //通知地址
	ReturnURL    string //返回地址
}

func SetAppId(v string) Option {
	return func(o *Options) {
		o.AppId = v
	}
}

func KeyAliPublicKey(v string) Option {
	return func(o *Options) {
		o.AliPublicKey = v
	}
}

func SetPrivateKey(v string) Option {
	return func(o *Options) {
		o.PrivateKey = v
	}
}
func SetNotifyURL(v string) Option {
	return func(o *Options) {
		o.NotifyURL = v
	}
}
func SetReturnURL(v string) Option {
	return func(o *Options) {
		o.ReturnURL = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) (options *Options, err error) {
	options = &Options{}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(options)
	}
	if options.AppId == "" || options.AliPublicKey == "" || options.PrivateKey == "" {
		err = fmt.Errorf("Missing configuration : AppId:%s AliPublicKey:%s PrivateKey:%s", options.AppId, options.AliPublicKey, options.PrivateKey)
	}
	return
}

func newOptionsByOption(opts ...Option) (options *Options, err error) {
	options = &Options{}
	for _, o := range opts {
		o(options)
	}
	if options.AppId == "" || options.AliPublicKey == "" || options.PrivateKey == "" {
		err = fmt.Errorf("Missing configuration : AppId:%s AliPublicKey:%s PrivateKey:%s", options.AppId, options.AliPublicKey, options.PrivateKey)
	}
	return
}
