package single

import (
	"fmt"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
)

type SingleService struct {
	cbase.ServiceBase
	opts *Options
}

func (this *SingleService) GetId() string {
	return this.opts.Id
}
func (this *SingleService) GetType() string {
	return this.opts.Type
}
func (this *SingleService) GetVersion() int32 {
	return this.opts.Version
}
func (this *SingleService) GetSettings() core.ServiceSttings {
	return this.opts.Setting
}
func (this *SingleService) GetWorkPath() string {
	return this.opts.WorkPath
}
func (this *SingleService) Configure(opts ...Option) {
	this.opts = newOptions(opts...)
}

//初始化服务系统模块
func (this *SingleService) InitSys() {
	if err := log.OnInit(this.opts.Setting.Sys["log"]); err != nil {
		panic(fmt.Sprintf("初始化log系统失败 %s", err.Error()))
	}
}
