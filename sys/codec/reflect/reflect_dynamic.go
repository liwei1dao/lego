package reflect

import (
	"reflect"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

type dynamicEncoder struct {
	valType reflect2.Type
}

func (encoder *dynamicEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	stream.WriteVal(obj, opt)
}

func (encoder *dynamicEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.valType.UnsafeIndirect(ptr) == nil
}

type efaceDecoder struct {
}

func (decoder *efaceDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
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
	if extra.WhatIsNext() == core.NilValue {
		if ptrElemType.Kind() != reflect.Ptr {
			extra.SkipBytes([]byte{'n', 'u', 'l', 'l'})
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

func (decoder *ifaceDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
	if extra.ReadNil() {
		decoder.valType.UnsafeSet(ptr, decoder.valType.UnsafeNew())
		return
	}
	obj := decoder.valType.UnsafeIndirect(ptr)
	if reflect2.IsNil(obj) {
		extra.ReportError("decode non empty interface", "can not unmarshal into nil")
		return
	}
	extra.ReadVal(obj)
}
