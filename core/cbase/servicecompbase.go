package cbase

import (
	"github.com/liwei1dao/lego/core"
)

type ServiceCompBase struct {
}

func (this *ServiceCompBase) NewOptions() (options core.ICompOptions) {
	return new(ServiceCompOptions)
}

func (this *ServiceCompBase) Init(service core.IService, comp core.IServiceComp, options core.ICompOptions) (err error) {
	return
}

func (this *ServiceCompBase) Start() (err error) {
	return
}

func (this *ServiceCompBase) Destroy() (err error) {
	return
}
