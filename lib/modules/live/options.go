package live

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
		GetRtmpAddr() string
	}
	Options struct {
		RtmpAddr string
	}
)

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		if err = mapstructure.Decode(settings, this); err != nil {
			return
		}
	}
	return
}
