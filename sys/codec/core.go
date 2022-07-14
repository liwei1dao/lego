package codec

import (
	"github.com/liwei1dao/lego/sys/codec/core"

	"github.com/modern-go/reflect2"
)

type (
	ISys interface {
		DecoderOf(typ reflect2.Type) core.IDecoder
		EncoderOf(typ reflect2.Type) core.IEncoder
		MarshalJson(v interface{}, option ...core.ExecuteOption) ([]byte, error)
		UnmarshalJson(data []byte, v interface{}, option ...core.ExecuteOption) error
		MarshalMapJson(val interface{}, option ...core.ExecuteOption) (ret map[string]string, err error)
		UnmarshalMapJson(data map[string]string, val interface{}, option ...core.ExecuteOption) (err error)
		MarshalSliceJson(val interface{}, option ...core.ExecuteOption) (ret []string, err error)
		UnmarshalSliceJson(data []string, val interface{}, option ...core.ExecuteOption) (err error)
	}
)

var defsys ISys

func OnInit(config map[string]interface{}, opt ...core.Option) (err error) {
	var option *core.Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...core.Option) (sys ISys, err error) {
	var option *core.Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}
func DecoderOf(typ reflect2.Type) core.IDecoder {
	return defsys.DecoderOf(typ)
}
func EncoderOf(typ reflect2.Type) core.IEncoder {
	return defsys.EncoderOf(typ)
}
func MarshalJson(v interface{}, option ...core.ExecuteOption) ([]byte, error) {
	return defsys.MarshalJson(v, option...)
}
func UnmarshalJson(data []byte, v interface{}, option ...core.ExecuteOption) error {
	return defsys.UnmarshalJson(data, v, option...)
}
func MarshalMapJson(val interface{}, option ...core.ExecuteOption) (ret map[string]string, err error) {
	return defsys.MarshalMapJson(val, option...)
}
func UnmarshalMapJson(data map[string]string, val interface{}, option ...core.ExecuteOption) (err error) {
	return defsys.UnmarshalMapJson(data, val, option...)
}
func MarshalSliceJson(val interface{}, option ...core.ExecuteOption) (ret []string, err error) {
	return defsys.MarshalSliceJson(val, option...)
}
func UnmarshalSliceJson(data []string, val interface{}, option ...core.ExecuteOption) (err error) {
	return defsys.UnmarshalSliceJson(data, val)
}
