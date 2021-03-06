package cluster

import (
	"fmt"

	"github.com/liwei1dao/lego/core"

	"github.com/BurntSushi/toml"
)

type Option func(*Options)

type Options struct {
	Id      string
	Setting core.ServiceSttings
}

func SetId(v string) Option {
	return func(o *Options) {
		o.Id = v
	}
}

func SetSetting(v core.ServiceSttings) Option {
	return func(o *Options) {
		o.Setting = v
	}
}

func newOptions(option ...Option) *Options {
	options := &Options{
		Id: "cluster_1",
	}
	for _, o := range option {
		o(options)
	}
	confpath := fmt.Sprintf("conf/%s.toml", options.Id)
	_, err := toml.DecodeFile(confpath, &options.Setting)
	if err != nil {
		panic(fmt.Sprintf("读取服务配置【%s】文件失败err:%v:", confpath, err))
	}
	if options.Setting.Id == "" || options.Setting.Type == "" || options.Setting.Tag == "" {
		panic(fmt.Sprintf("服务[%s] 配置缺少必要配置: %+v", options.Id, options))
	}
	return options
}
