package codec

import (
	"github.com/liwei1dao/lego/sys/codec/core"
)

type (
	ISys interface {
		Marshal(v interface{}) ([]byte, error)
		Unmarshal(data []byte, v interface{}) error
	}
)

var defsys ISys

func OnInit(config map[string]interface{}, option ...core.Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...core.Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}
