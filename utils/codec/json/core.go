package json

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/liwei1dao/lego/utils/codec"
	"github.com/liwei1dao/lego/utils/codec/codecore"

	"github.com/modern-go/reflect2"
)

const defsize = 512

var writerPool = &sync.Pool{
	New: func() interface{} {
		return &JsonWriter{
			config: &codecore.Config{
				SortMapKeys:           true,
				IndentionStep:         0,
				OnlyTaggedField:       false,
				DisallowUnknownFields: false,
				CaseSensitive:         false,
				TagKey:                "json",
			},
			buf:       make([]byte, 0, defsize),
			err:       nil,
			indention: 0,
		}
	},
}

func GetReader(buf []byte, r io.Reader) codecore.IReader {
	reader := readerPool.Get().(codecore.IReader)
	reader.ResetBytes(buf, r)
	return reader
}

func PutReader(r codecore.IReader) {
	readerPool.Put(r)
}

var readerPool = &sync.Pool{
	New: func() interface{} {
		return &JsonReader{
			config: &codecore.Config{
				SortMapKeys:           true,
				IndentionStep:         0,
				OnlyTaggedField:       false,
				DisallowUnknownFields: false,
				CaseSensitive:         false,
				TagKey:                "json",
			},
			buf:    nil,
			reader: nil,
			head:   0,
			tail:   0,
			depth:  0,
			err:    nil,
		}
	},
}

func GetWriter(w io.Writer) codecore.IWriter {
	writer := writerPool.Get().(codecore.IWriter)
	writer.Reset(w)
	return writer
}

func PutWriter(w codecore.IWriter) {
	w.Reset(nil)
	writerPool.Put(w)
}

const maxDepth = 10000

var valueTypes []codecore.ValueType

func init() {
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
	SortMapKeys:           true,
	IndentionStep:         1,
	OnlyTaggedField:       false,
	DisallowUnknownFields: false,
	CaseSensitive:         false,
	TagKey:                "json",
}

func Marshal(val interface{}) (buf []byte, err error) {
	writer := GetWriter(nil)
	defer PutWriter(writer)
	writer.WriteVal(val)
	if writer.Error() != nil {
		return nil, writer.Error()
	}
	result := writer.Buffer()
	copied := make([]byte, len(result))
	copy(copied, result)
	return copied, nil
}

func MarshalIndent(val interface{}, prefix, indent string) ([]byte, error) {
	writer := GetWriter(nil)

	defer func() {
		writer.Config().IndentionStep = 0
		PutWriter(writer)
	}()
	writer.Config().IndentionStep = len(indent)
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
	extra := GetReader(data, nil)
	defer PutReader(extra)
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
		w := GetWriter(nil)
		ret, err = encoderMapJson.EncodeToMapJson(reflect2.PtrOf(val), w)
		PutWriter(w)
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
		r := GetReader([]byte{}, nil)
		err = decoderMapJson.DecodeForMapJson(ptr, r, data)
		PutReader(r)
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
		w := GetWriter(nil)
		ret, err = encoderMapJson.EncodeToSliceJson(reflect2.PtrOf(val), w)
		w.PutWriter(w)
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
		r := GetReader([]byte{}, nil)
		err = decoderMapJson.DecodeForSliceJson(ptr, r, data)
		PutReader(r)
	}
	return
}
