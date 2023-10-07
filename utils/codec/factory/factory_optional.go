package factory

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/liwei1dao/lego/utils/codec/codecore"

	"github.com/modern-go/reflect2"
)

func decoderOfOptional(ctx codecore.ICtx, typ reflect2.Type) codecore.IDecoder {
	ptrType := typ.(*reflect2.UnsafePtrType)
	elemType := ptrType.Elem()
	decoder := DecoderOfType(ctx, elemType)
	return &OptionalDecoder{elemType, decoder}
}

func encoderOfOptional(ctx codecore.ICtx, typ reflect2.Type) codecore.IEncoder {
	ptrType := typ.(*reflect2.UnsafePtrType)
	elemType := ptrType.Elem()
	elemEncoder := EncoderOfType(ctx, elemType)
	encoder := &OptionalEncoder{elemEncoder}
	return encoder
}

//Optional--------------------------------------------------------------------------------------------------------------------
type OptionalDecoder struct {
	ValueType    reflect2.Type
	ValueDecoder codecore.IDecoder
}

func (this *OptionalDecoder) GetType() reflect.Kind {
	return this.ValueDecoder.GetType()
}
func (this *OptionalDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
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

func (this *OptionalDecoder) DecodeForMapJson(ptr unsafe.Pointer, r codecore.IReader, extra map[string]string) (err error) {
	if decoderMapJson, ok := this.ValueDecoder.(codecore.IDecoderMapJson); !ok {
		err = fmt.Errorf("encoder %T not support EncodeToMapJson", this.ValueDecoder)
		return
	} else {
		return decoderMapJson.DecodeForMapJson(ptr, r, extra)
	}
}

type OptionalEncoder struct {
	ValueEncoder codecore.IEncoder
}

func (this *OptionalEncoder) GetType() reflect.Kind {
	return this.ValueEncoder.GetType()
}
func (this *OptionalEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
	} else {
		this.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), stream)
	}
}

func (this *OptionalEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*unsafe.Pointer)(ptr)) == nil
}

func (this *OptionalEncoder) EncodeToMapJson(ptr unsafe.Pointer, w codecore.IWriter) (ret map[string]string, err error) {
	if encoderMapJson, ok := this.ValueEncoder.(codecore.IEncoderMapJson); !ok {
		err = fmt.Errorf("encoder %T not support EncodeToMapJson", this.ValueEncoder)
		return
	} else {
		if *((*unsafe.Pointer)(ptr)) == nil {
			err = fmt.Errorf("encoder ptr is nil")
			return
		} else {
			return encoderMapJson.EncodeToMapJson(*((*unsafe.Pointer)(ptr)), w)
		}
	}
}

//reference--------------------------------------------------------------------------------------------------------------------
type referenceEncoder struct {
	encoder codecore.IEncoder
}

func (this *referenceEncoder) GetType() reflect.Kind {
	return this.encoder.GetType()
}
func (this *referenceEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	this.encoder.Encode(unsafe.Pointer(&ptr), stream)
}

func (this *referenceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return this.encoder.IsEmpty(unsafe.Pointer(&ptr))
}

type referenceDecoder struct {
	decoder codecore.IDecoder
}

func (this *referenceDecoder) GetType() reflect.Kind {
	return this.decoder.GetType()
}

func (this *referenceDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	this.decoder.Decode(unsafe.Pointer(&ptr), extra)
}

//dereference--------------------------------------------------------------------------------------------------------------------
type dereferenceDecoder struct {
	valueType    reflect2.Type
	valueDecoder codecore.IDecoder
}

func (this *dereferenceDecoder) GetType() reflect.Kind {
	return this.valueDecoder.GetType()
}
func (this *dereferenceDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		newPtr := this.valueType.UnsafeNew()
		this.valueDecoder.Decode(newPtr, extra)
		*((*unsafe.Pointer)(ptr)) = newPtr
	} else {
		this.valueDecoder.Decode(*((*unsafe.Pointer)(ptr)), extra)
	}
}

type dereferenceEncoder struct {
	ValueEncoder codecore.IEncoder
}

func (this *dereferenceEncoder) GetType() reflect.Kind {
	return this.ValueEncoder.GetType()
}
func (this *dereferenceEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
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
	isEmbeddedPtrNil, converted := this.ValueEncoder.(codecore.IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := unsafe.Pointer(deReferenced)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}
