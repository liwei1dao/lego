package single

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/utils"
)

type Option func(*Options)

type Options struct {
	Id       string
	Type     string
	Version  int32
	WorkPath string
	Setting  core.ServiceSttings
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

func newOptions(opts ...Option) *Options {
	opt := &Options{
		Id:       "cluster_1",
		Type:     "cluster",
		Version:  1,
		WorkPath: utils.GetApplicationDir(),
	}
	for _, o := range opts {
		o(opt)
	}
	confpath := fmt.Sprintf(opt.WorkPath+"conf/%s.toml", opt.Id)
	_, err := toml.DecodeFile(confpath, &opt.Setting)
	if err != nil {
		panic(fmt.Sprintf("读取服务配置【%s】文件失败err=%s:", confpath, err.Error()))
	}
	return opt
}
