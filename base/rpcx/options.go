package rpcx

import (
	"fmt"
	"os"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/utils/container/id"
	"github.com/liwei1dao/lego/utils/mapstructure"

	"gopkg.in/yaml.v2"
)

type RPCXConfig struct {
	Port           int32    //监听短裤
	Addr           string   //对外地址
	ETCDServers    []string //etcd 地址
	UpdateInterval int32    //保活间隔
}

type Option func(*Options)

type Options struct {
	ConfPath   string
	Version    string //服务版本
	Setting    core.ServiceSttings
	RPCXConfig RPCXConfig
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
	yamlFile, err := os.ReadFile(options.ConfPath)
	if err != nil {
		panic(fmt.Sprintf("读取服务配置【%s】文件失败err:%v:", options.ConfPath, err))
	}
	err = yaml.Unmarshal(yamlFile, &options.Setting)
	if err != nil {
		panic(fmt.Sprintf("读取服务配置【%s】文件失败err:%v:", options.ConfPath, err))
	}
	if len(options.Setting.Type) == 0 || len(options.Setting.Tag) == 0 {
		panic(fmt.Sprintf("[%s] 配置缺少必要配置: %+v", options.ConfPath, options))
	}
	if len(options.Setting.Id) == 0 {
		options.Setting.Id = fmt.Sprintf("%s-%s", options.Setting.Type, id.NewXId())
	}
	if err = mapstructure.Decode(options.Setting.Exte, &options.RPCXConfig); err != nil {
		panic(fmt.Sprintf("[%s] 配置缺少必要配置 Exte: %+v", options.ConfPath, options))
	}
	return options
}
