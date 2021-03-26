package live

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib"
)

type Live struct {
	cbase.ModuleBase
}

func (this *Live) GetType() core.M_Modules {
	return lib.SM_MonitorModule
}

func (this *Live) NewOptions() (options core.IModuleOptions) {
	return new(Options)
}
