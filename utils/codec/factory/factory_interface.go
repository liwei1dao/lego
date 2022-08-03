package factory

import (
	"errors"
	"reflect"
	"unsafe"

	"github.com/liwei1dao/lego/utils/codec/codecore"

	"github.com/modern-go/reflect2"
)

type dynamicEncoder struct {
	valType reflect2.Type
}

func (codec *dynamicEncoder) GetType() reflect.Kind {
	return reflect.Interface
}
func (encoder *dynamicEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	stream.WriteVal(obj)
}

func (encoder *dynamicEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.valType.UnsafeIndirect(ptr) == nil
}

type efaceDecoder struct {
}

func (codec *efaceDecoder) GetType() reflect.Kind {
	return reflect.Interface
}
func (decoder *efaceDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	pObj := (*interface{})(ptr)
	obj := *pObj
	if obj == nil {
		*pObj = extra.Read()
		return
	}
	typ := reflect2.TypeOf(obj)
	if typ.Kind() != reflect.Ptr {
		*pObj = extra.Read()
		return
	}
	ptrType := typ.(*reflect2.UnsafePtrType)
	ptrElemType := ptrType.Elem()
	if extra.WhatIsNext() == codecore.NilValue {
		if ptrElemType.Kind() != reflect.Ptr {
			extra.ReadNil()
			*pObj = nil
			return
		}
	}
	if reflect2.IsNil(obj) {
		obj := ptrElemType.New()
		extra.ReadVal(obj)
		*pObj = obj
		return
	}
	extra.ReadVal(obj)
}

type ifaceDecoder struct {
	valType *reflect2.UnsafeIFaceType
}

func (codec *ifaceDecoder) GetType() reflect.Kind {
	return reflect.Interface
}
func (decoder *ifaceDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if extra.ReadNil() {
		decoder.valType.UnsafeSet(ptr, decoder.valType.UnsafeNew())
		return
	}
	obj := decoder.valType.UnsafeIndirect(ptr)
	if reflect2.IsNil(obj) {
		extra.SetErr(errors.New("decode non empty interface can not unmarshal into nil"))
		return
	}
	extra.ReadVal(obj)
}
