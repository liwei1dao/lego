package monitor

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
	}
	Options struct {
	}
)

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		mapstructure.Decode(settings, &this)
	}

	return
}
