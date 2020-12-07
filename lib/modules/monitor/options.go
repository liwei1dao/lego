package monitor

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IMonitorOptions interface {
	}
	MonitorOptions struct {
	}
)

func (this *MonitorOptions) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		mapstructure.Decode(settings, &this)
	}

	return
}
