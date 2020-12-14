package cbase

import (
	"github.com/liwei1dao/lego/core"
)

type ModuleBase struct {
	comps []core.IModuleComp
}

func (this *ModuleBase) NewOptions() (options core.IModuleOptions) {
	return new(ModuleOptions)
}

func (this *ModuleBase) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	module.OnInstallComp()
	for _, v := range this.comps {
		err = v.Init(service, module, v, options)
		if err != nil {
			return
		}
	}
	return
}

func (this *ModuleBase) OnInstallComp() {

}

func (this *ModuleBase) Start() (err error) {
	for _, v := range this.comps {
		err = v.Start()
		if err != nil {
			return
		}
	}
	return
}

func (this *ModuleBase) Run(closeSig chan bool) (err error) {

	return
}

func (this *ModuleBase) Destroy() (err error) {
	for _, v := range this.comps {
		err = v.Destroy()
		if err != nil {
			return
		}
	}
	return
}

func (this *ModuleBase) RegisterComp(comp core.IModuleComp) interface{} {
	this.comps = append(this.comps, comp)
	return comp
}
