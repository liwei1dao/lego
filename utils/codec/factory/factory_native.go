package factory

import (
	"reflect"
	"strconv"
	"unsafe"

	"github.com/liwei1dao/lego/utils/codec/codecore"

	"github.com/modern-go/reflect2"
)

const ptrSize = 32 << uintptr(^uintptr(0)>>63) //计算int的大小 是32 还是64

func createDecoderOfNative(ctx codecore.ICtx, typ reflect2.Type) codecore.IDecoder {
	typeName := typ.String()
	switch typ.Kind() {
	case reflect.String:
		if typeName != "string" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*string)(nil)).Elem())
		}
		return &stringCodec{}
	case reflect.Int:
		if typeName != "int" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*int)(nil)).Elem())
		}
		if strconv.IntSize == 32 {
			return &int32Codec{}
		}
		return &int64Codec{}
	case reflect.Int8:
		if typeName != "int8" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*int8)(nil)).Elem())
		}
		return &int8Codec{}
	case reflect.Int16:
		if typeName != "int16" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*int16)(nil)).Elem())
		}
		return &int16Codec{}
	case reflect.Int32:
		if typeName != "int32" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*int32)(nil)).Elem())
		}
		return &int32Codec{}
	case reflect.Int64:
		if typeName != "int64" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*int64)(nil)).Elem())
		}
		return &int64Codec{}
	case reflect.Uint:
		if typeName != "uint" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*uint)(nil)).Elem())
		}
		if strconv.IntSize == 32 {
			return &uint32Codec{}
		}
		return &uint64Codec{}
	case reflect.Uint8:
		if typeName != "uint8" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*uint8)(nil)).Elem())
		}
		return &uint8Codec{}
	case reflect.Uint16:
		if typeName != "uint16" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*uint16)(nil)).Elem())
		}
		return &uint16Codec{}
	case reflect.Uint32:
		if typeName != "uint32" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*uint32)(nil)).Elem())
		}
		return &uint32Codec{}
	case reflect.Uintptr:
		if typeName != "uintptr" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*uintptr)(nil)).Elem())
		}
		if ptrSize == 32 {
			return &uint32Codec{}
		}
		return &uint64Codec{}
	case reflect.Uint64:
		if typeName != "uint64" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*uint64)(nil)).Elem())
		}
		return &uint64Codec{}
	case reflect.Float32:
		if typeName != "float32" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*float32)(nil)).Elem())
		}
		return &float32Codec{}
	case reflect.Float64:
		if typeName != "float64" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*float64)(nil)).Elem())
		}
		return &float64Codec{}
	case reflect.Bool:
		if typeName != "bool" {
			return DecoderOfType(ctx, reflect2.TypeOfPtr((*bool)(nil)).Elem())
		}
		return &boolCodec{}
	}
	return nil
}

func createEncoderOfNative(ctx codecore.ICtx, typ reflect2.Type) codecore.IEncoder {
	typeName := typ.String()
	kind := typ.Kind()
	switch kind {
	case reflect.String:
		if typeName != "string" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*string)(nil)).Elem())
		}
		return &stringCodec{}
	case reflect.Int:
		if typeName != "int" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*int)(nil)).Elem())
		}
		if strconv.IntSize == 32 {
			return &int32Codec{}
		}
		return &int64Codec{}
	case reflect.Int8:
		if typeName != "int8" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*int8)(nil)).Elem())
		}
		return &int8Codec{}
	case reflect.Int16:
		if typeName != "int16" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*int16)(nil)).Elem())
		}
		return &int16Codec{}
	case reflect.Int32:
		if typeName != "int32" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*int32)(nil)).Elem())
		}
		return &int32Codec{}
	case reflect.Int64:
		if typeName != "int64" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*int64)(nil)).Elem())
		}
		return &int64Codec{}
	case reflect.Uint:
		if typeName != "uint" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*uint)(nil)).Elem())
		}
		if strconv.IntSize == 32 {
			return &uint32Codec{}
		}
		return &uint64Codec{}
	case reflect.Uint8:
		if typeName != "uint8" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*uint8)(nil)).Elem())
		}
		return &uint8Codec{}
	case reflect.Uint16:
		if typeName != "uint16" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*uint16)(nil)).Elem())
		}
		return &uint16Codec{}
	case reflect.Uint32:
		if typeName != "uint32" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*uint32)(nil)).Elem())
		}
		return &uint32Codec{}
	case reflect.Uintptr:
		if typeName != "uintptr" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*uintptr)(nil)).Elem())
		}
		if ptrSize == 32 {
			return &uint32Codec{}
		}
		return &uint64Codec{}
	case reflect.Uint64:
		if typeName != "uint64" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*uint64)(nil)).Elem())
		}
		return &uint64Codec{}
	case reflect.Float32:
		if typeName != "float32" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*float32)(nil)).Elem())
		}
		return &float32Codec{}
	case reflect.Float64:
		if typeName != "float64" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*float64)(nil)).Elem())
		}
		return &float64Codec{}
	case reflect.Bool:
		if typeName != "bool" {
			return EncoderOfType(ctx, reflect2.TypeOfPtr((*bool)(nil)).Elem())
		}
		return &boolCodec{}
	}
	return nil
}

type stringCodec struct {
}

func (codec *stringCodec) GetType() reflect.Kind {
	return reflect.String
}
func (codec *stringCodec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	*((*string)(ptr)) = extra.ReadString()
}

func (codec *stringCodec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	str := *((*string)(ptr))
	stream.WriteString(str)
}

func (codec *stringCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*string)(ptr)) == ""
}

type int8Codec struct {
}

func (codec *int8Codec) GetType() reflect.Kind {
	return reflect.Int8
}

func (codec *int8Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*int8)(ptr)) = extra.ReadInt8()
	}
}

func (codec *int8Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteInt8(*((*int8)(ptr)))
}

func (codec *int8Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int8)(ptr)) == 0
}

type int16Codec struct {
}

func (codec *int16Codec) GetType() reflect.Kind {
	return reflect.Int16
}
func (codec *int16Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*int16)(ptr)) = extra.ReadInt16()
	}
}

func (codec *int16Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteInt16(*((*int16)(ptr)))
}

func (codec *int16Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int16)(ptr)) == 0
}

type int32Codec struct {
}

func (codec *int32Codec) GetType() reflect.Kind {
	return reflect.Int32
}
func (codec *int32Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*int32)(ptr)) = extra.ReadInt32()
	}
}

func (codec *int32Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteInt32(*((*int32)(ptr)))
}
func (codec *int32Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int32)(ptr)) == 0
}

type int64Codec struct {
}

func (codec *int64Codec) GetType() reflect.Kind {
	return reflect.Int64
}
func (codec *int64Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*int64)(ptr)) = extra.ReadInt64()
	}
}

func (codec *int64Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteInt64(*((*int64)(ptr)))
}

func (codec *int64Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int64)(ptr)) == 0
}

type uint8Codec struct {
}

func (codec *uint8Codec) GetType() reflect.Kind {
	return reflect.Uint8
}
func (codec *uint8Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*uint8)(ptr)) = extra.ReadUint8()
	}
}

func (codec *uint8Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteUint8(*((*uint8)(ptr)))
}

func (codec *uint8Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uint8)(ptr)) == 0
}

type uint16Codec struct {
}

func (codec *uint16Codec) GetType() reflect.Kind {
	return reflect.Uint16
}
func (codec *uint16Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*uint16)(ptr)) = extra.ReadUint16()
	}
}

func (codec *uint16Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteUint16(*((*uint16)(ptr)))
}

func (codec *uint16Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uint16)(ptr)) == 0
}

type uint32Codec struct {
}

func (codec *uint32Codec) GetType() reflect.Kind {
	return reflect.Uint32
}
func (codec *uint32Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*uint32)(ptr)) = extra.ReadUint32()
	}
}

func (codec *uint32Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteUint32(*((*uint32)(ptr)))
}
func (codec *uint32Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uint32)(ptr)) == 0
}

type uint64Codec struct {
}

func (codec *uint64Codec) GetType() reflect.Kind {
	return reflect.Uint64
}
func (codec *uint64Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*uint64)(ptr)) = extra.ReadUint64()
	}
}

func (codec *uint64Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteUint64(*((*uint64)(ptr)))
}

func (codec *uint64Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*uint64)(ptr)) == 0
}

type float32Codec struct {
}

func (codec *float32Codec) GetType() reflect.Kind {
	return reflect.Float32
}
func (codec *float32Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*float32)(ptr)) = extra.ReadFloat32()
	}
}

func (codec *float32Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteFloat32(*((*float32)(ptr)))
}

func (codec *float32Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*float32)(ptr)) == 0
}

type float64Codec struct {
}

func (codec *float64Codec) GetType() reflect.Kind {
	return reflect.Float64
}
func (codec *float64Codec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*float64)(ptr)) = extra.ReadFloat64()
	}
}

func (codec *float64Codec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteFloat64(*((*float64)(ptr)))
}

func (codec *float64Codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*float64)(ptr)) == 0
}

type boolCodec struct {
}

func (codec *boolCodec) GetType() reflect.Kind {
	return reflect.Bool
}
func (codec *boolCodec) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadNil() {
		*((*bool)(ptr)) = extra.ReadBool()
	}
}

func (codec *boolCodec) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteBool(*((*bool)(ptr)))
}

func (codec *boolCodec) IsEmpty(ptr unsafe.Pointer) bool {
	return !(*((*bool)(ptr)))
}
