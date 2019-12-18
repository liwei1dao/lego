package single

import (
	"lego/core"
	"lego/core/cbase"
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
