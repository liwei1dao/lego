package colly

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib"
)

type Colly struct {
	cbase.ModuleBase
}

func (this *Colly) GetType() core.M_Modules {
	return lib.SM_CollyModule
}

func (this *Colly) NewOptions() (options core.IModuleOptions) {
	return new(Options)
}
func (this *Colly) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	err = this.ModuleBase.Init(service, module, options)
	return
}
