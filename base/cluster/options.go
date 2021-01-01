package cluster

import (
	"fmt"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils"

	"github.com/BurntSushi/toml"
)

type Option func(*Options)

type Options struct {
	Tag       string
	Id        string
	Type      string
	Category  core.S_Category
	Version   int32
	WorkPath  string
	Setting   core.ServiceSttings
	Debugmode bool
	LogLvel   log.Loglevel
}

func SetTag(v string) Option {
	return func(o *Options) {
		o.Tag = v
	}
}
func SetId(v string) Option {
	return func(o *Options) {
		o.Id = v
	}
}
func SetType(v string) Option {
	return func(o *Options) {
		o.Type = v
	}
}

func SetSetting(v core.ServiceSttings) Option {
	return func(o *Options) {
		o.Setting = v
	}
}

func SetCategory(v core.S_Category) Option {
	return func(o *Options) {
		o.Category = v
	}
}
func SetVersion(v int32) Option {
	return func(o *Options) {
		o.Version = v
	}
}

func SetWorkPath(v string) Option {
	return func(o *Options) {
		o.WorkPath = v
	}
}

func SetDebugMode(v bool) Option {
	return func(o *Options) {
		o.Debugmode = v
	}
}

func SetLogLvel(v log.Loglevel) Option {
	return func(o *Options) {
		o.LogLvel = v
	}
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		Tag:       "liwie1dao",
		Id:        "cluster_1",
		Type:      "cluster",
		Category:  core.S_Category_BusinessService,
		Version:   1,
		WorkPath:  utils.GetApplicationDir(),
		LogLvel:   log.InfoLevel,
		Debugmode: false,
		Setting: core.ServiceSttings{Settings: map[string]interface{}{
			"ConsulAddr": "127.0.0.1:8500",
			"NatsAddr":   "127.0.0.1:4222",
		}},
	}
	for _, o := range opts {
		o(opt)
	}
	confpath := fmt.Sprintf("conf/%s.toml", opt.Id)
	_, err := toml.DecodeFile(confpath, &opt.Setting)
	if err != nil {
		fmt.Printf("警告 读取服务配置【%s】文件失败err=%s:\n", confpath, err.Error())
	}
	return opt
}
