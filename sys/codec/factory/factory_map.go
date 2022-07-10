package factory

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

func decoderOfMap(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	mapType := typ.(*reflect2.UnsafeMapType)
	keyDecoder := decoderOfMapKey(ctx.Append("[mapKey]"), mapType.Key())
	elemDecoder := DecoderOfType(ctx.Append("[mapElem]"), mapType.Elem())
	return &mapDecoder{
		mapType:     mapType,
		keyType:     mapType.Key(),
		elemType:    mapType.Elem(),
		keyDecoder:  keyDecoder,
		elemDecoder: elemDecoder,
	}
}

func encoderOfMap(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	mapType := typ.(*reflect2.UnsafeMapType)
	return &mapEncoder{
		mapType:     mapType,
		keyEncoder:  encoderOfMapKey(ctx.Append("[mapKey]"), mapType.Key()),
		elemEncoder: EncoderOfType(ctx.Append("[mapElem]"), mapType.Elem()),
	}
}

func decoderOfMapKey(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	switch typ.Kind() {
	case reflect.String:
		return DecoderOfType(ctx, reflect2.DefaultTypeOfKind(reflect.String))
	case reflect.Bool,
		reflect.Uint8, reflect.Int8,
		reflect.Uint16, reflect.Int16,
		reflect.Uint32, reflect.Int32,
		reflect.Uint64, reflect.Int64,
		reflect.Uint, reflect.Int,
		reflect.Float32, reflect.Float64,
		reflect.Uintptr:
		typ = reflect2.DefaultTypeOfKind(typ.Kind())
		return &numericMapKeyDecoder{DecoderOfType(ctx, typ)}
	default:
		return &lazyErrorDecoder{err: fmt.Errorf("unsupported map key type: %v", typ)}
	}
}

func encoderOfMapKey(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	switch typ.Kind() {
	case reflect.String:
		return EncoderOfType(ctx, reflect2.DefaultTypeOfKind(reflect.String))
	case reflect.Bool,
		reflect.Uint8, reflect.Int8,
		reflect.Uint16, reflect.Int16,
		reflect.Uint32, reflect.Int32,
		reflect.Uint64, reflect.Int64,
		reflect.Uint, reflect.Int,
		reflect.Float32, reflect.Float64,
		reflect.Uintptr:
		typ = reflect2.DefaultTypeOfKind(typ.Kind())
		return &numericMapKeyEncoder{EncoderOfType(ctx, typ)}
	default:
		if typ.Kind() == reflect.Interface {
			return &dynamicMapKeyEncoder{ctx, typ}
		}
		return &lazyErrorEncoder{err: fmt.Errorf("unsupported map key type: %v", typ)}
	}
}

//Map--------------------------------------------------------------------------------------------------------------------------------------
type mapEncoder struct {
	mapType     *reflect2.UnsafeMapType
	keyEncoder  core.IEncoder
	elemEncoder core.IEncoder
}

func (encoder *mapEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	if *(*unsafe.Pointer)(ptr) == nil {
		stream.WriteNil()
		return
	}
	stream.WriteObjectStart()
	iter := encoder.mapType.UnsafeIterate(ptr)
	for i := 0; iter.HasNext(); i++ {
		if i != 0 {
			stream.WriteMemberSplit()
		}
		key, elem := iter.UnsafeNext()
		encoder.keyEncoder.Encode(key, stream)
		stream.WriteKVSplit()
		encoder.elemEncoder.Encode(elem, stream)
	}
	stream.WriteObjectEnd()
}

func (encoder *mapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	iter := encoder.mapType.UnsafeIterate(ptr)
	return !iter.HasNext()
}

type mapDecoder struct {
	mapType     *reflect2.UnsafeMapType
	keyType     reflect2.Type
	elemType    reflect2.Type
	keyDecoder  core.IDecoder
	elemDecoder core.IDecoder
}

func (decoder *mapDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	mapType := decoder.mapType
	if extra.ReadNil() {
		*(*unsafe.Pointer)(ptr) = nil
		mapType.UnsafeSet(ptr, mapType.UnsafeNew())
		return
	}
	if mapType.UnsafeIsNil(ptr) {
		mapType.UnsafeSet(ptr, mapType.UnsafeMakeMap(0))
	}
	if extra.ReadObjectStart() {
		return
	}
	if extra.CheckNextIsObjectEnd() {
		return
	}
	key := decoder.keyType.UnsafeNew()
	decoder.keyDecoder.Decode(key, extra)
	if extra.ReadKVSplit() {
		return
	}
	elem := decoder.elemType.UnsafeNew()
	decoder.elemDecoder.Decode(elem, extra)
	decoder.mapType.UnsafeSetIndex(ptr, key, elem)
	for extra.ReadMemberSplit() {
		key := decoder.keyType.UnsafeNew()
		decoder.keyDecoder.Decode(key, extra)
		if extra.ReadMemberSplit() {
			return
		}
		elem := decoder.elemType.UnsafeNew()
		decoder.elemDecoder.Decode(elem, extra)
		decoder.mapType.UnsafeSetIndex(ptr, key, elem)
	}
	extra.ReadObjectEnd()
}

//NumericMap-------------------------------------------------------------------------------------------------------------------------------
type numericMapKeyDecoder struct {
	decoder core.IDecoder
}

func (decoder *numericMapKeyDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	if extra.ReadKeyStart() {
		return
	}
	decoder.decoder.Decode(ptr, extra)
	if extra.ReadKeyEnd() {
		return
	}
}

type numericMapKeyEncoder struct {
	encoder core.IEncoder
}

func (encoder *numericMapKeyEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	stream.WriteKeyStart()
	encoder.encoder.Encode(ptr, stream)
	stream.WriteKeyEnd()
}

func (encoder *numericMapKeyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

//------------------------------------------------------------------------------------------------------------------
type dynamicMapKeyEncoder struct {
	ctx     *core.Ctx
	valType reflect2.Type
}

func (encoder *dynamicMapKeyEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	encoderOfMapKey(encoder.ctx, reflect2.TypeOf(obj)).Encode(reflect2.PtrOf(obj), stream)
}

func (encoder *dynamicMapKeyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	obj := encoder.valType.UnsafeIndirect(ptr)
	return encoderOfMapKey(encoder.ctx, reflect2.TypeOf(obj)).IsEmpty(reflect2.PtrOf(obj))
}
