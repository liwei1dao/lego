package single

import (
	"fmt"
	"io/ioutil"

	"github.com/liwei1dao/lego/core"
	"gopkg.in/yaml.v2"
)

type Option func(*Options)

type Options struct {
	ConfPath string
	Version  string //服务版本
	Setting  core.ServiceSttings
}

func SetConfPath(v string) Option {
	return func(o *Options) {
		o.ConfPath = v
	}
}

func SetVersion(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

func newOptions(option ...Option) *Options {
	options := &Options{
		ConfPath: "conf/cluster.yaml",
	}
	for _, o := range option {
		o(options)
	}
	// confpath := fmt.Sprintf("conf/%s.toml", options.Id)
	yamlFile, err := ioutil.ReadFile(options.ConfPath)
	if err != nil {
		panic(fmt.Sprintf("读取服务配置【%s】文件失败err:%v:", options.ConfPath, err))
	}
	err = yaml.Unmarshal(yamlFile, &options.Setting)
	if err != nil {
		panic(fmt.Sprintf("读取服务配置【%s】文件失败err:%v:", options.ConfPath, err))
	}
	if len(options.Setting.Id) == 0 || len(options.Setting.Type) == 0 || len(options.Setting.Tag) == 0 {
		panic(fmt.Sprintf("[%s] 配置缺少必要配置: %+v", options.ConfPath, options))
	}
	return options
}
