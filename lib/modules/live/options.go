package live

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
		Get_GopNum() int
	}
	Options struct {
		GopNum int
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

func (this *Options) Get_GopNum() int {
	return this.GopNum
}
