package reflect

import (
	"encoding"
	"encoding/json"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

var marshalerType = reflect2.TypeOfPtr((*json.Marshaler)(nil)).Elem()
var unmarshalerType = reflect2.TypeOfPtr((*json.Unmarshaler)(nil)).Elem()
var textMarshalerType = reflect2.TypeOfPtr((*encoding.TextMarshaler)(nil)).Elem()
var textUnmarshalerType = reflect2.TypeOfPtr((*encoding.TextUnmarshaler)(nil)).Elem()

func createDecoderOfMarshaler(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
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

func createEncoderOfMarshaler(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	if typ == marshalerType {
		checkIsEmpty := createCheckIsEmpty(ctx, typ)
		var encoder core.IEncoder = &directMarshalerEncoder{
			checkIsEmpty: checkIsEmpty,
		}
		return encoder
	}
	if typ.Implements(marshalerType) {
		checkIsEmpty := createCheckIsEmpty(ctx, typ)
		var encoder core.IEncoder = &marshalerEncoder{
			valType:      typ,
			checkIsEmpty: checkIsEmpty,
		}
		return encoder
	}
	ptrType := reflect2.PtrTo(typ)
	if ctx.Prefix != "" && ptrType.Implements(marshalerType) {
		checkIsEmpty := createCheckIsEmpty(ctx, ptrType)
		var encoder core.IEncoder = &marshalerEncoder{
			valType:      ptrType,
			checkIsEmpty: checkIsEmpty,
		}
		return &referenceEncoder{encoder}
	}
	if typ == textMarshalerType {
		checkIsEmpty := createCheckIsEmpty(ctx, typ)
		var encoder core.IEncoder = &directTextMarshalerEncoder{
			checkIsEmpty:  checkIsEmpty,
			stringEncoder: ctx.EncoderOf(reflect2.TypeOf("")),
		}
		return encoder
	}
	if typ.Implements(textMarshalerType) {
		checkIsEmpty := createCheckIsEmpty(ctx, typ)
		var encoder core.IEncoder = &textMarshalerEncoder{
			valType:       typ,
			stringEncoder: ctx.EncoderOf(reflect2.TypeOf("")),
			checkIsEmpty:  checkIsEmpty,
		}
		return encoder
	}
	// if prefix is empty, the type is the root type
	if ctx.Prefix != "" && ptrType.Implements(textMarshalerType) {
		checkIsEmpty := createCheckIsEmpty(ctx, ptrType)
		var encoder core.IEncoder = &textMarshalerEncoder{
			valType:       ptrType,
			stringEncoder: ctx.EncoderOf(reflect2.TypeOf("")),
			checkIsEmpty:  checkIsEmpty,
		}
		return &referenceEncoder{encoder}
	}
	return nil
}

type unmarshalerDecoder struct {
	valType reflect2.Type
}

func (decoder *unmarshalerDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
	valType := decoder.valType
	obj := valType.UnsafeIndirect(ptr)
	unmarshaler := obj.(json.Unmarshaler)
	extra.NextToken()
	extra.UnreadChar() // skip spaces
	bytes := extra.SkipAndReturnBytes()
	err := unmarshaler.UnmarshalJSON(bytes)
	if err != nil {
		extra.ReportError("unmarshalerDecoder", err.Error())
	}
}

type marshalerEncoder struct {
	checkIsEmpty core.CheckIsEmpty
	valType      reflect2.Type
}

func (encoder *marshalerEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	if encoder.valType.IsNullable() && reflect2.IsNil(obj) {
		stream.WriteNil()
		return
	}
	marshaler := obj.(json.Marshaler)
	bytes, err := marshaler.MarshalJSON()
	if err != nil {
		stream.SetError(err)
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

type textUnmarshalerDecoder struct {
	valType reflect2.Type
}

func (decoder *textUnmarshalerDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
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
		extra.ReportError("textUnmarshalerDecoder", err.Error())
	}
}

type textMarshalerEncoder struct {
	valType       reflect2.Type
	stringEncoder core.IEncoder
	checkIsEmpty  core.CheckIsEmpty
}

func (encoder *textMarshalerEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	if encoder.valType.IsNullable() && reflect2.IsNil(obj) {
		stream.WriteNil()
		return
	}
	marshaler := (obj).(encoding.TextMarshaler)
	bytes, err := marshaler.MarshalText()
	if err != nil {
		stream.SetError(err)
	} else {
		str := string(bytes)
		encoder.stringEncoder.Encode(unsafe.Pointer(&str), stream, opt)
	}
}

func (encoder *textMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type directTextMarshalerEncoder struct {
	stringEncoder core.IEncoder
	checkIsEmpty  core.CheckIsEmpty
}

func (encoder *directTextMarshalerEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	marshaler := *(*encoding.TextMarshaler)(ptr)
	if marshaler == nil {
		stream.WriteNil()
		return
	}
	bytes, err := marshaler.MarshalText()
	if err != nil {
		stream.SetError(err)
	} else {
		str := string(bytes)
		encoder.stringEncoder.Encode(unsafe.Pointer(&str), stream, opt)
	}
}

func (encoder *directTextMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}

type directMarshalerEncoder struct {
	checkIsEmpty core.CheckIsEmpty
}

func (encoder *directMarshalerEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	marshaler := *(*json.Marshaler)(ptr)
	if marshaler == nil {
		stream.WriteNil()
		return
	}
	bytes, err := marshaler.MarshalJSON()
	if err != nil {
		stream.SetError(err)
	} else {
		stream.WriteBytes(bytes)
	}
}

func (encoder *directMarshalerEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.checkIsEmpty.IsEmpty(ptr)
}
