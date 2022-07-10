package factory

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

type sortableBindings []*Binding

func (this sortableBindings) Len() int {
	return len(this)
}

func (this sortableBindings) Less(i, j int) bool {
	left := this[i].levels
	right := this[j].levels
	k := 0
	for {
		if left[k] < right[k] {
			return true
		} else if left[k] > right[k] {
			return false
		}
		k++
	}
}

func (this sortableBindings) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

//----------------------------------------------------------------------------------------------------------------------------------------------------------------
func EncoderOfType(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	encoder := ctx.Encoders[typ]
	if encoder != nil {
		return encoder
	}
	root := &rootEncoder{}
	ctx.Encoders[typ] = root
	encoder = _createEncoderOfType(ctx, typ)
	root.encoder = encoder
	return encoder
}

func _createEncoderOfType(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	var encoder core.IEncoder
	encoder = createEncoderOfNative(ctx, typ)
	if encoder != nil {
		return encoder
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Interface:
		return &dynamicEncoder{typ}
	case reflect.Struct:
		return encoderOfStruct(ctx, typ)
	case reflect.Array:
		return encoderOfArray(ctx, typ)
	case reflect.Slice:
		return encoderOfSlice(ctx, typ)
	case reflect.Map:
		return encoderOfMap(ctx, typ)
	case reflect.Ptr:
		return encoderOfOptional(ctx, typ)
	default:
		return &lazyErrorEncoder{err: fmt.Errorf("%s %s is unsupported type", ctx.Prefix, typ.String())}
	}
}

func DecoderOfType(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	decoder := ctx.Decoders[typ]
	if decoder != nil {
		return decoder
	}
	root := &rootDecoder{}
	ctx.Decoders[typ] = root
	decoder = _createDecoderOfType(ctx, typ)
	root.decoder = decoder
	return decoder
}

func _createDecoderOfType(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	var decoder core.IDecoder

	decoder = createDecoderOfNative(ctx, typ)
	if decoder != nil {
		return decoder
	}
	switch typ.Kind() {
	case reflect.Interface:
		ifaceType, isIFace := typ.(*reflect2.UnsafeIFaceType)
		if isIFace {
			return &ifaceDecoder{valType: ifaceType}
		}
		return &efaceDecoder{}
	case reflect.Struct:
		return decoderOfStruct(ctx, typ)
	case reflect.Array:
		return decoderOfArray(ctx, typ)
	case reflect.Slice:
		return decoderOfSlice(ctx, typ)
	case reflect.Map:
		return decoderOfMap(ctx, typ)
	case reflect.Ptr:
		return decoderOfOptional(ctx, typ)
	default:
		return &lazyErrorDecoder{err: fmt.Errorf("%s %s is unsupported type", ctx.Prefix, typ.String())}
	}
}

//根节点 -------------------------------------------------------------------
type rootDecoder struct {
	decoder core.IDecoder
}

func (this *rootDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	this.decoder.Decode(ptr, extra)
}

type rootEncoder struct {
	encoder core.IEncoder
}

func (this *rootEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	this.encoder.Encode(ptr, stream)
}

func (this *rootEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return this.encoder.IsEmpty(ptr)
}

//onePtrEncoder---------------------------------------------------------------
func NewonePtrEncoder(encoder core.IEncoder) core.IEncoder {
	return &onePtrEncoder{encoder}
}

type onePtrEncoder struct {
	encoder core.IEncoder
}

func (encoder *onePtrEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.encoder.IsEmpty(unsafe.Pointer(&ptr))
}

func (encoder *onePtrEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	encoder.encoder.Encode(unsafe.Pointer(&ptr), stream)
}

//错误节点 ------------------------------------------------------------------
type lazyErrorDecoder struct {
	err error
}

func (this *lazyErrorDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	if extra.Error() == nil {
		extra.SetErr(this.err)
	}
}

type lazyErrorEncoder struct {
	err error
}

func (this *lazyErrorEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	if ptr == nil {
		stream.WriteNil()
	} else if stream.Error() == nil {
		stream.SetErr(this.err)
	}
}

func (this *lazyErrorEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}
