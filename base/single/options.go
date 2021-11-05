package single

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/utils/ip"
)

type Option func(*Options)

type Options struct {
	ConfPath string
	Setting  core.ServiceSttings
}

func SetConfPath(v string) Option {
	return func(o *Options) {
		o.ConfPath = v
	}
}

func SetSetting(v core.ServiceSttings) Option {
	return func(o *Options) {
		o.Setting = v
	}
}
func newOptions(option ...Option) *Options {
	options := &Options{
		ConfPath: "conf/cluster.toml",
	}
	for _, o := range option {
		o(options)
	}
	// confpath := fmt.Sprintf("conf/%s.toml", options.Id)
	_, err := toml.DecodeFile(options.ConfPath, &options.Setting)
	if err != nil {
		panic(fmt.Sprintf("读取服务配置【%s】文件失败err:%v:", options.ConfPath, err))
	}
	if len(options.Setting.Id) == 0 || len(options.Setting.Type) == 0 || len(options.Setting.Tag) == 0 {
		panic(fmt.Sprintf("[%s] 配置缺少必要配置: %+v", options.ConfPath, options))
	}
	if len(options.Setting.Ip) == 0 {
		if ipinfo := ip.GetEthernetInfo(); ipinfo != nil { //获取以太网Ip地址
			options.Setting.Ip = ipinfo.IP
		} else {
			options.Setting.Ip = ip.GetOutboundIP() //局域网ip
		}
	}
	return options
}
