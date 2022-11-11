package factory

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/liwei1dao/lego/utils/codec/codecore"
	"github.com/modern-go/reflect2"
)

var marshalerType = reflect2.TypeOfPtr((*json.Marshaler)(nil)).Elem()
var unmarshalerType = reflect2.TypeOfPtr((*json.Unmarshaler)(nil)).Elem()
var textMarshalerType = reflect2.TypeOfPtr((*encoding.TextMarshaler)(nil)).Elem()
var textUnmarshalerType = reflect2.TypeOfPtr((*encoding.TextUnmarshaler)(nil)).Elem()

func createDecoderOfMarshaler(ctx codecore.ICtx, typ reflect2.Type) codecore.IDecoder {
	ptrType := reflect2.PtrTo(typ)
	if ptrType.Implements(unmarshalerType) {
		return &referenceDecoder{
			&unmarshalerDecoder{ptrType},
		}
	}
	if ptrType.Implements(textUnmarshalerType) {
		return &referenceDecoder{
			&textUnmarshalerDecoder{ptrType},
		}
	}
	return nil
}

func createEncoderOfMarshaler(ctx codecore.ICtx, typ reflect2.Type) codecore.IEncoder {
	if typ == marshalerType {
		checkIsEmpty := createCheckIsEmpty(ctx, typ)
		var encoder codecore.IEncoder = &directMarshalerEncoder{
			checkIsEmpty: checkIsEmpty,
		}
		return encoder
	}
	if typ.Implements(marshalerType) {
		checkIsEmpty := createCheckIsEmpty(ctx, typ)
		var encoder codecore.IEncoder = &marshalerEncoder{
			valType:      typ,
			checkIsEmpty: checkIsEmpty,
		}
		return encoder
	}
	ptrType := reflect2.PtrTo(typ)
	if ctx.Prefix() != "" && ptrType.Implements(marshalerType) {
		checkIsEmpty := createCheckIsEmpty(ctx, ptrType)
		var encoder codecore.IEncoder = &marshalerEncoder{
			valType:      ptrType,
			checkIsEmpty: checkIsEmpty,
		}
		return &referenceEncoder{encoder}
	}
	if typ == textMarshalerType {
		checkIsEmpty := createCheckIsEmpty(ctx, typ)
		var encoder codecore.IEncoder = &directTextMarshalerEncoder{
			checkIsEmpty:  checkIsEmpty,
			stringEncoder: ctx.EncoderOf(reflect2.TypeOf("")),
		}
		return encoder
	}
	if typ.Implements(textMarshalerType) {
		checkIsEmpty := createCheckIsEmpty(ctx, typ)
		var encoder codecore.IEncoder = &textMarshalerEncoder{
			valType:       typ,
			stringEncoder: ctx.EncoderOf(reflect2.TypeOf("")),
			checkIsEmpty:  checkIsEmpty,
		}
		return encoder
	}
	// if prefix is empty, the type is the root type
	if ctx.Prefix() != "" && ptrType.Implements(textMarshalerType) {
		checkIsEmpty := createCheckIsEmpty(ctx, ptrType)
		var encoder codecore.IEncoder = &textMarshalerEncoder{
			valType:       ptrType,
			stringEncoder: &stringCodec{},
			checkIsEmpty:  checkIsEmpty,
		}
		return &referenceEncoder{encoder}
	}
	return nil
}

type unmarshalerDecoder struct {
	valType reflect2.Type
}

func (codec *unmarshalerDecoder) GetType() reflect.Kind {
	return reflect.Ptr
}

func (decoder *unmarshalerDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	valType := decoder.valType
	obj := valType.UnsafeIndirect(ptr)
	unmarshaler := obj.(json.Unmarshaler)
	bytes := extra.SkipAndReturnBytes()
	err := unmarshaler.UnmarshalJSON(bytes)
	if err != nil {
		extra.SetErr(fmt.Errorf("unmarshalerDecoder %s", err.Error()))
	}
}

type textUnmarshalerDecoder struct {
	valType reflect2.Type
}

func (codec *textUnmarshalerDecoder) GetType() reflect.Kind {
	return reflect.Ptr
}

func (decoder *textUnmarshalerDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	valType := decoder.valType
	obj := valType.UnsafeIndirect(ptr)
	if reflect2.IsNil(obj) {
		ptrType := valType.(*reflect2.UnsafePtrType)
		elemType := ptrType.Elem()
		elem := elemType.UnsafeNew()
		ptrType.UnsafeSet(ptr, unsafe.Pointer(&elem))
		obj = valType.UnsafeIndirect(ptr)
	}
	unmarshaler := (obj).(encoding.TextUnmarshaler)
	str := extra.ReadString()
	err := unmarshaler.UnmarshalText([]byte(str))
	if err != nil {
		extra.SetErr(fmt.Errorf("textUnmarshalerDecoder %s", err.Error()))
	}
}

//Encoder --------------------------------------------------------------------------------------------------------------------------------
func createCheckIsEmpty(ctx codecore.ICtx, typ reflect2.Type) codecore.IEncoder {
	encoder := createEncoderOfNative(ctx, typ)
	if encoder != nil {
		return encoder
	}
	kind := typ.Kind()
	switch kind {
	case reflect.Interface:
		return &dynamicEncoder{typ}
	case reflect.Struct:
		return &structEncoder{typ: typ}
	case reflect.Array:
		return &arrayEncoder{}
	case reflect.Slice:
		return &sliceEncoder{}
	case reflect.Map:
		return encoderOfMap(ctx, typ)
	case reflect.Ptr:
		return &OptionalEncoder{}
	default:
		return &lazyErrorEncoder{err: fmt.Errorf("unsupported type: %v", typ)}
	}
}

type directMarshalerEncoder struct {
	checkIsEmpty codecore.IEncoder
}

func (codec *directMarshalerEncoder) GetType() reflect.Kind {
	return reflect.Ptr
}
func (encoder *directMarshalerEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	marshaler := *(*json.Marshaler)(ptr)
	if marshaler == nil {
		stream.WriteNil()
		return
	}
	bytes, err := marshaler.MarshalJSON()
	if err != nil {
		stream.SetErr(err)
	} else {
		stream.WriteBytes(bytes)
	}
}

func (encoder *directMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type marshalerEncoder struct {
	checkIsEmpty codecore.IEncoder
	valType      reflect2.Type
}

func (codec *marshalerEncoder) GetType() reflect.Kind {
	return reflect.Ptr
}
func (encoder *marshalerEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	if encoder.valType.IsNullable() && reflect2.IsNil(obj) {
		stream.WriteNil()
		return
	}
	marshaler := obj.(json.Marshaler)
	bytes, err := marshaler.MarshalJSON()
	if err != nil {
		stream.SetErr(err)
	} else {
		// html escape was already done by jsoniter
		// but the extra '\n' should be trimed
		l := len(bytes)
		if l > 0 && bytes[l-1] == '\n' {
			bytes = bytes[:l-1]
		}
		stream.WriteBytes(bytes)
	}
}

func (encoder *marshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type directTextMarshalerEncoder struct {
	stringEncoder codecore.IEncoder
	checkIsEmpty  codecore.IEncoder
}

func (codec *directTextMarshalerEncoder) GetType() reflect.Kind {
	return reflect.Ptr
}
func (encoder *directTextMarshalerEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	marshaler := *(*encoding.TextMarshaler)(ptr)
	if marshaler == nil {
		stream.WriteNil()
		return
	}
	bytes, err := marshaler.MarshalText()
	if err != nil {
		stream.SetErr(err)
	} else {
		str := string(bytes)
		encoder.stringEncoder.Encode(unsafe.Pointer(&str), stream)
	}
}

func (encoder *directTextMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type textMarshalerEncoder struct {
	valType       reflect2.Type
	stringEncoder codecore.IEncoder
	checkIsEmpty  codecore.IEncoder
}

func (codec *textMarshalerEncoder) GetType() reflect.Kind {
	return reflect.Ptr
}
func (encoder *textMarshalerEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	if encoder.valType.IsNullable() && reflect2.IsNil(obj) {
		stream.WriteNil()
		return
	}
	marshaler := (obj).(encoding.TextMarshaler)
	bytes, err := marshaler.MarshalText()
	if err != nil {
		stream.SetErr(err)
	} else {
		str := string(bytes)
		encoder.stringEncoder.Encode(unsafe.Pointer(&str), stream)
	}
}

func (encoder *textMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}
