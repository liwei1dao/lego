package reflect

import (
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

func decoderOfOptional(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	ptrType := typ.(*reflect2.UnsafePtrType)
	elemType := ptrType.Elem()
	decoder := decoderOfType(ctx, elemType)
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

func (decoder *OptionalDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
	if extra.ReadNil() {
		*((*unsafe.Pointer)(ptr)) = nil
	} else {
		if *((*unsafe.Pointer)(ptr)) == nil {
			//pointer to null, we have to allocate memory to hold the value
			newPtr := decoder.ValueType.UnsafeNew()
			decoder.ValueDecoder.Decode(newPtr, extra, opt)
			*((*unsafe.Pointer)(ptr)) = newPtr
		} else {
			//reuse existing instance
			decoder.ValueDecoder.Decode(*((*unsafe.Pointer)(ptr)), extra, opt)
		}
	}
}

type OptionalEncoder struct {
	ValueEncoder core.IEncoder
}

func (encoder *OptionalEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
	} else {
		encoder.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), stream, opt)
	}
}

func (encoder *OptionalEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*unsafe.Pointer)(ptr)) == nil
}

//reference--------------------------------------------------------------------------------------------------------------------
type referenceEncoder struct {
	encoder core.IEncoder
}

func (encoder *referenceEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	encoder.encoder.Encode(unsafe.Pointer(&ptr), stream, opt)
}

func (encoder *referenceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.encoder.IsEmpty(unsafe.Pointer(&ptr))
}

type referenceDecoder struct {
	decoder core.IDecoder
}

func (decoder *referenceDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
	decoder.decoder.Decode(unsafe.Pointer(&ptr), extra, opt)
}

//dereference--------------------------------------------------------------------------------------------------------------------
type dereferenceDecoder struct {
	valueType    reflect2.Type
	valueDecoder core.IDecoder
}

func (this *dereferenceDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		newPtr := this.valueType.UnsafeNew()
		this.valueDecoder.Decode(newPtr, extra, opt)
		*((*unsafe.Pointer)(ptr)) = newPtr
	} else {
		this.valueDecoder.Decode(*((*unsafe.Pointer)(ptr)), extra, opt)
	}
}

type dereferenceEncoder struct {
	ValueEncoder core.IEncoder
}

func (encoder *dereferenceEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	if *((*unsafe.Pointer)(ptr)) == nil {
		stream.WriteNil()
	} else {
		encoder.ValueEncoder.Encode(*((*unsafe.Pointer)(ptr)), stream, opt)
	}
}

func (encoder *dereferenceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	dePtr := *((*unsafe.Pointer)(ptr))
	if dePtr == nil {
		return true
	}
	return encoder.ValueEncoder.IsEmpty(dePtr)
}

func (encoder *dereferenceEncoder) IsEmbeddedPtrNil(ptr unsafe.Pointer) bool {
	deReferenced := *((*unsafe.Pointer)(ptr))
	if deReferenced == nil {
		return true
	}
	isEmbeddedPtrNil, converted := encoder.ValueEncoder.(core.IsEmbeddedPtrNil)
	if !converted {
		return false
	}
	fieldPtr := unsafe.Pointer(deReferenced)
	return isEmbeddedPtrNil.IsEmbeddedPtrNil(fieldPtr)
}
