package cbase

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type ServiceCompOptions struct {
	core.ICompOptions
}

func (this *ServiceCompOptions) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		mapstructure.Decode(settings, &this)
	}
	return
}
