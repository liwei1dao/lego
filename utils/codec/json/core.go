package json

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/liwei1dao/lego/utils/codec"
	"github.com/liwei1dao/lego/utils/codec/codecore"
	"github.com/modern-go/reflect2"
)

const maxDepth = 10000

var hexDigits []byte
var valueTypes []codecore.ValueType

func init() {
	hexDigits = make([]byte, 256)
	for i := 0; i < len(hexDigits); i++ {
		hexDigits[i] = 255
	}
	for i := '0'; i <= '9'; i++ {
		hexDigits[i] = byte(i - '0')
	}
	for i := 'a'; i <= 'f'; i++ {
		hexDigits[i] = byte((i - 'a') + 10)
	}
	for i := 'A'; i <= 'F'; i++ {
		hexDigits[i] = byte((i - 'A') + 10)
	}
	valueTypes = make([]codecore.ValueType, 256)
	for i := 0; i < len(valueTypes); i++ {
		valueTypes[i] = codecore.InvalidValue
	}
	valueTypes['"'] = codecore.StringValue
	valueTypes['-'] = codecore.NumberValue
	valueTypes['0'] = codecore.NumberValue
	valueTypes['1'] = codecore.NumberValue
	valueTypes['2'] = codecore.NumberValue
	valueTypes['3'] = codecore.NumberValue
	valueTypes['4'] = codecore.NumberValue
	valueTypes['5'] = codecore.NumberValue
	valueTypes['6'] = codecore.NumberValue
	valueTypes['7'] = codecore.NumberValue
	valueTypes['8'] = codecore.NumberValue
	valueTypes['9'] = codecore.NumberValue
	valueTypes['t'] = codecore.BoolValue
	valueTypes['f'] = codecore.BoolValue
	valueTypes['n'] = codecore.NilValue
	valueTypes['['] = codecore.ArrayValue
	valueTypes['{'] = codecore.ObjectValue
}

var defconf = &codecore.Config{
	IndentionStep:         1,
	OnlyTaggedField:       false,
	DisallowUnknownFields: false,
	CaseSensitive:         false,
	TagKey:                "json",
}

func Marshal(val interface{}) (buf []byte, err error) {
	writer := BorrowWriter()
	defer writer.Free()
	writer.WriteVal(val)
	if writer.Error() != nil {
		return nil, writer.Error()
	}
	result := writer.Buffer()
	copied := make([]byte, len(result))
	copy(copied, result)
	return copied, nil
}

func Unmarshal(data []byte, v interface{}) error {
	extra := BorrowReader(data)
	defer ReturnReader(extra)
	extra.ReadVal(v)
	return extra.Error()
}

func MarshalMap(val interface{}) (ret map[string]string, err error) {
	if nil == val {
		err = errors.New("val is null")
		return
	}
	cacheKey := reflect2.RTypeOf(val)
	encoder := codec.GetEncoder(cacheKey)
	if encoder == nil {
		typ := reflect2.TypeOf(val)
		encoder = codec.EncoderOf(typ, defconf)
	}
	if encoderMapJson, ok := encoder.(codecore.IEncoderMapJson); !ok {
		err = fmt.Errorf("val type:%T not support MarshalMapJson", val)
	} else {
		w := BorrowWriter()
		ret, err = encoderMapJson.EncodeToMapJson(reflect2.PtrOf(val), w)
		w.Free()
	}
	return
}

func UnmarshalMap(data map[string]string, val interface{}) (err error) {
	cacheKey := reflect2.RTypeOf(val)
	decoder := codec.GetDecoder(cacheKey)
	if decoder == nil {
		typ := reflect2.TypeOf(val)
		if typ == nil || typ.Kind() != reflect.Ptr {
			err = errors.New("can only unmarshal into pointer")
			return
		}
		decoder = codec.DecoderOf(typ, defconf)
	}
	ptr := reflect2.PtrOf(val)
	if ptr == nil {
		err = errors.New("can not read into nil pointer")
		return
	}
	if decoderMapJson, ok := decoder.(codecore.IDecoderMapJson); !ok {
		err = fmt.Errorf("val type:%T not support MarshalMapJson", val)
	} else {
		r := BorrowReader([]byte{})
		err = decoderMapJson.DecodeForMapJson(ptr, r, data)
		r.Free()
	}
	return
}

func MarshalSlice(val interface{}) (ret []string, err error) {
	if nil == val {
		err = errors.New("val is null")
		return
	}
	cacheKey := reflect2.RTypeOf(val)
	encoder := codec.GetEncoder(cacheKey)
	if encoder == nil {
		typ := reflect2.TypeOf(val)
		encoder = codec.EncoderOf(typ, defconf)
	}
	if encoderMapJson, ok := encoder.(codecore.IEncoderSliceJson); !ok {
		err = fmt.Errorf("val type:%T not support MarshalMapJson", val)
	} else {
		w := BorrowWriter()
		ret, err = encoderMapJson.EncodeToSliceJson(reflect2.PtrOf(val), w)
		w.Free()
	}
	return
}

func UnmarshalSlice(data []string, val interface{}) (err error) {
	cacheKey := reflect2.RTypeOf(val)
	decoder := codec.GetDecoder(cacheKey)
	if decoder == nil {
		typ := reflect2.TypeOf(val)
		if typ == nil || typ.Kind() != reflect.Ptr {
			err = errors.New("can only unmarshal into pointer")
			return
		}
		decoder = codec.DecoderOf(typ, defconf)
	}
	ptr := reflect2.PtrOf(val)
	if ptr == nil {
		err = errors.New("can not read into nil pointer")
		return
	}
	if decoderMapJson, ok := decoder.(codecore.IDecoderSliceJson); !ok {
		err = fmt.Errorf("val type:%T not support UnmarshalSliceJson", val)
	} else {
		r := BorrowReader([]byte{})
		err = decoderMapJson.DecodeForSliceJson(ptr, r, data)
		r.Free()
	}
	return
}
