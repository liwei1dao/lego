package factory

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"

	"github.com/modern-go/reflect2"
)

func decoderOfOptional(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	ptrType := typ.(*reflect2.UnsafePtrType)
	elemType := ptrType.Elem()
	decoder := DecoderOfType(ctx, elemType)
	return &OptionalDecoder{elemType, decoder}
}

func encoderOfOptional(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	ptrType := typ.(*reflect2.UnsafePtrType)
	elemType := ptrType.Elem()
	elemEncoder := EncoderOfType(ctx, elemType)
	encoder := &OptionalEncoder{elemEncoder}
	return encoder
}

//Optional--------------------------------------------------------------------------------------------------------------------
type OptionalDecoder struct {
	ValueType    reflect2.Type
	ValueDecoder core.IDecoder
}

func (this *OptionalDecoder) GetType() reflect.Kind {
	return this.ValueDecoder.GetType()
}
func (this *OptionalDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	if extra.ReadNil() {
		*((*unsafe.Pointer)(ptr)) = nil
	} else {
		if *((*unsafe.Pointer)(ptr)) == nil {
			newPtr := this.ValueType.UnsafeNew()
			this.ValueDecoder.Decode(newPtr, extra)
			*((*unsafe.Pointer)(ptr)) = newPtr
		} else {
			this.ValueDecoder.Decode(*((*unsafe.Pointer)(ptr)), extra)
		}
	}
}

func (this *OptionalDecoder) DecodeForMapJson(ptr unsafe.Pointer, extra map[string]string) (err error) {
	if decoderMapJson, ok := this.ValueDecoder.(core.IDecoderMapJson); !ok {
		err = fmt.Errorf("encoder %T not support EncodeToMapJson", this.ValueDecoder)
		return
	} else {
		return decoderMapJson.DecodeForMapJson(ptr, extra)
	}
}

type OptionalEncoder struct {
	ValueEncoder core.IEncoder
}

func (this *OptionalEncoder) GetType() reflect.Kind {
	return this.ValueEncoder.GetType()
}
func (this *OptionalEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
	} else {
		this.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), stream)
	}
}

func (this *OptionalEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*unsafe.Pointer)(ptr)) == nil
}

func (this *OptionalEncoder) EncodeToMapJson(ptr unsafe.Pointer) (ret map[string]string, err error) {
	if encoderMapJson, ok := this.ValueEncoder.(core.IEncoderMapJson); !ok {
		err = fmt.Errorf("encoder %T not support EncodeToMapJson", this.ValueEncoder)
		return
	} else {
		if *((*unsafe.Pointer)(ptr)) == nil {
			err = fmt.Errorf("encoder ptr is nil")
			return
		} else {
			return encoderMapJson.EncodeToMapJson(*((*unsafe.Pointer)(ptr)))
		}
	}
}

//reference--------------------------------------------------------------------------------------------------------------------
type referenceEncoder struct {
	encoder core.IEncoder
}

func (this *referenceEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	this.encoder.Encode(unsafe.Pointer(&ptr), stream)
}

func (this *referenceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return this.encoder.IsEmpty(unsafe.Pointer(&ptr))
}

type referenceDecoder struct {
	decoder core.IDecoder
}

func (this *referenceDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	this.decoder.Decode(unsafe.Pointer(&ptr), extra)
}

//dereference--------------------------------------------------------------------------------------------------------------------
type dereferenceDecoder struct {
	valueType    reflect2.Type
	valueDecoder core.IDecoder
}

func (this *dereferenceDecoder) GetType() reflect.Kind {
	return this.valueDecoder.GetType()
}
func (this *dereferenceDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		newPtr := this.valueType.UnsafeNew()
		this.valueDecoder.Decode(newPtr, extra)
		*((*unsafe.Pointer)(ptr)) = newPtr
	} else {
		this.valueDecoder.Decode(*((*unsafe.Pointer)(ptr)), extra)
	}
}

type dereferenceEncoder struct {
	ValueEncoder core.IEncoder
}

func (this *dereferenceEncoder) GetType() reflect.Kind {
	return this.ValueEncoder.GetType()
}
func (this *dereferenceEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
	} else {
		this.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), stream)
	}
}

func (this *dereferenceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	dePtr := *((*unsafe.Pointer)(ptr))
	if dePtr == nil {
		return true
	}
	return this.ValueEncoder.IsEmpty(dePtr)
}

func (this *dereferenceEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	deReferenced := *((*unsafe.Pointer)(ptr))
	if deReferenced == nil {
		return true
	}
	isEmbeddedPtrNil, converted := this.ValueEncoder.(core.IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := unsafe.Pointer(deReferenced)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}
