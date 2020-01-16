package cbase

import (
	"github.com/liwei1dao/lego/core"
)

type ModuleCompBase struct{}

func (this *ModuleCompBase) Init(service core.IService, module core.IModule, comp core.IModuleComp, setting map[string]interface{}) (err error) {
	return
}

func (this *ModuleCompBase) Start() (err error) {
	return
}

func (this *ModuleCompBase) Destroy() (err error) {
	return
}
