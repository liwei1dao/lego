package reflect

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

func decoderOfMap(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	mapType := typ.(*reflect2.UnsafeMapType)
	keyDecoder := decoderOfMapKey(ctx.Append("[mapKey]"), mapType.Key())
	elemDecoder := decoderOfType(ctx.Append("[mapElem]"), mapType.Elem())
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
	if ctx.Options().SortMapKeys {
		return &sortKeysMapEncoder{
			codec:       ctx.ICodec,
			mapType:     mapType,
			keyEncoder:  encoderOfMapKey(ctx.Append("[mapKey]"), mapType.Key()),
			elemEncoder: EncoderOfType(ctx.Append("[mapElem]"), mapType.Elem()),
		}
	}
	return &mapEncoder{
		mapType:     mapType,
		keyEncoder:  encoderOfMapKey(ctx.Append("[mapKey]"), mapType.Key()),
		elemEncoder: EncoderOfType(ctx.Append("[mapElem]"), mapType.Elem()),
	}
}

func decoderOfMapKey(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	ptrType := reflect2.PtrTo(typ)
	if ptrType.Implements(unmarshalerType) {
		return &referenceDecoder{
			&unmarshalerDecoder{
				valType: ptrType,
			},
		}
	}
	if typ.Implements(unmarshalerType) {
		return &unmarshalerDecoder{
			valType: typ,
		}
	}
	if ptrType.Implements(textUnmarshalerType) {
		return &referenceDecoder{
			&textUnmarshalerDecoder{
				valType: ptrType,
			},
		}
	}
	if typ.Implements(textUnmarshalerType) {
		return &textUnmarshalerDecoder{
			valType: typ,
		}
	}

	switch typ.Kind() {
	case reflect.String:
		return decoderOfType(ctx, reflect2.DefaultTypeOfKind(reflect.String))
	case reflect.Bool,
		reflect.Uint8, reflect.Int8,
		reflect.Uint16, reflect.Int16,
		reflect.Uint32, reflect.Int32,
		reflect.Uint64, reflect.Int64,
		reflect.Uint, reflect.Int,
		reflect.Float32, reflect.Float64,
		reflect.Uintptr:
		typ = reflect2.DefaultTypeOfKind(typ.Kind())
		return &numericMapKeyDecoder{decoderOfType(ctx, typ)}
	default:
		return &lazyErrorDecoder{err: fmt.Errorf("unsupported map key type: %v", typ)}
	}
}

func encoderOfMapKey(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	if typ == textMarshalerType {
		return &directTextMarshalerEncoder{
			stringEncoder: ctx.EncoderOf(reflect2.TypeOf("")),
		}
	}
	if typ.Implements(textMarshalerType) {
		return &textMarshalerEncoder{
			valType:       typ,
			stringEncoder: ctx.EncoderOf(reflect2.TypeOf("")),
		}
	}

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

func (encoder *mapEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	if *(*unsafe.Pointer)(ptr) == nil {
		stream.WriteNil()
		return
	}
	stream.WriteObjectStart()
	iter := encoder.mapType.UnsafeIterate(ptr)
	for i := 0; iter.HasNext(); i++ {
		if i != 0 {
			stream.WriteMore()
		}
		key, elem := iter.UnsafeNext()
		encoder.keyEncoder.Encode(key, stream, opt)
		if stream.Indention() > 0 {
			stream.WriteBytes([]byte{':', ' '})
		} else {
			stream.WriteChar(':')
		}
		encoder.elemEncoder.Encode(elem, stream, opt)
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

func (decoder *mapDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
	mapType := decoder.mapType
	c := extra.NextToken()
	if c == 'n' {
		extra.SkipBytes([]byte{'u', 'l', 'l'})
		*(*unsafe.Pointer)(ptr) = nil
		mapType.UnsafeSet(ptr, mapType.UnsafeNew())
		return
	}
	if mapType.UnsafeIsNil(ptr) {
		mapType.UnsafeSet(ptr, mapType.UnsafeMakeMap(0))
	}
	if c != '{' {
		extra.ReportError("ReadMapCB", `expect { or n, but found `+string([]byte{c}))
		return
	}
	c = extra.NextToken()
	if c == '}' {
		return
	}
	extra.UnreadChar()
	key := decoder.keyType.UnsafeNew()
	decoder.keyDecoder.Decode(key, extra, opt)
	c = extra.NextToken()
	if c != ':' {
		extra.ReportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
		return
	}
	elem := decoder.elemType.UnsafeNew()
	decoder.elemDecoder.Decode(elem, extra, opt)
	decoder.mapType.UnsafeSetIndex(ptr, key, elem)
	for c = extra.NextToken(); c == ','; c = extra.NextToken() {
		key := decoder.keyType.UnsafeNew()
		decoder.keyDecoder.Decode(key, extra, opt)
		c = extra.NextToken()
		if c != ':' {
			extra.ReportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
			return
		}
		elem := decoder.elemType.UnsafeNew()
		decoder.elemDecoder.Decode(elem, extra, opt)
		decoder.mapType.UnsafeSetIndex(ptr, key, elem)
	}
	if c != '}' {
		extra.ReportError("ReadMapCB", `expect }, but found `+string([]byte{c}))
	}
}

//SortKeysMap-------------------------------------------------------------------------------------------------------------------------------
type sortKeysMapEncoder struct {
	codec       core.ICodec
	mapType     *reflect2.UnsafeMapType
	keyEncoder  core.IEncoder
	elemEncoder core.IEncoder
}

func (this *sortKeysMapEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	if *(*unsafe.Pointer)(ptr) == nil {
		stream.WriteNil()
		return
	}
	stream.WriteObjectStart()
	mapIter := this.mapType.UnsafeIterate(ptr)
	subStream := this.codec.BorrowStream()
	subExtra := this.codec.BorrowExtractor()
	keyValues := encodedKeyValues{}
	for mapIter.HasNext() {
		key, elem := mapIter.UnsafeNext()
		subStreamIndex := subStream.Buffered()
		this.keyEncoder.Encode(key, subStream, opt)
		if subStream.Error() != nil && subStream.Error() != io.EOF && stream.Error() == nil {
			stream.SetError(subStream.Error())
		}
		encodedKey := subStream.ToBuffer()[subStreamIndex:]
		subExtra.ResetBytes(encodedKey)
		decodedKey := subExtra.ReadString()
		if stream.Indention() > 0 {
			subStream.WriteBytes([]byte{':', ' '})
		} else {
			subStream.WriteChar(':')
		}
		this.elemEncoder.Encode(elem, subStream, opt)
		keyValues = append(keyValues, encodedKV{
			key:      decodedKey,
			keyValue: subStream.ToBuffer()[subStreamIndex:],
		})
	}
	sort.Sort(keyValues)
	for i, keyValue := range keyValues {
		if i != 0 {
			stream.WriteMore()
		}
		stream.WriteBytes(keyValue.keyValue)
	}
	if subStream.Error() != nil && stream.Error() == nil {
		stream.SetError(subStream.Error())
	}
	stream.WriteObjectEnd()
	this.codec.ReturnStream(subStream)
	this.codec.ReturnExtractor(subExtra)
}

func (encoder *sortKeysMapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	iter := encoder.mapType.UnsafeIterate(ptr)
	return !iter.HasNext()
}

type encodedKeyValues []encodedKV

type encodedKV struct {
	key      string
	keyValue []byte
}

func (sv encodedKeyValues) Len() int           { return len(sv) }
func (sv encodedKeyValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv encodedKeyValues) Less(i, j int) bool { return sv[i].key < sv[j].key }

//NumericMap-------------------------------------------------------------------------------------------------------------------------------
type numericMapKeyDecoder struct {
	decoder core.IDecoder
}

func (decoder *numericMapKeyDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
	c := extra.NextToken()
	if c != '"' {
		extra.ReportError("ReadMapCB", `expect ", but found `+string([]byte{c}))
		return
	}
	decoder.decoder.Decode(ptr, extra, opt)
	c = extra.NextToken()
	if c != '"' {
		extra.ReportError("ReadMapCB", `expect ", but found `+string([]byte{c}))
		return
	}
}

type numericMapKeyEncoder struct {
	encoder core.IEncoder
}

func (encoder *numericMapKeyEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	stream.WriteChar('"')
	encoder.encoder.Encode(ptr, stream, opt)
	stream.WriteChar('"')
}

func (encoder *numericMapKeyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

//------------------------------------------------------------------------------------------------------------------
type dynamicMapKeyEncoder struct {
	ctx     *core.Ctx
	valType reflect2.Type
}

func (encoder *dynamicMapKeyEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	obj := encoder.valType.UnsafeIndirect(ptr)
	encoderOfMapKey(encoder.ctx, reflect2.TypeOf(obj)).Encode(reflect2.PtrOf(obj), stream, opt)
}

func (encoder *dynamicMapKeyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	obj := encoder.valType.UnsafeIndirect(ptr)
	return encoderOfMapKey(encoder.ctx, reflect2.TypeOf(obj)).IsEmpty(reflect2.PtrOf(obj))
}
