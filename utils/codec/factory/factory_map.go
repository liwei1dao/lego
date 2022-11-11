package factory

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"unsafe"

	"github.com/liwei1dao/lego/utils/codec/codecore"

	"github.com/modern-go/reflect2"
)

func decoderOfMap(ctx codecore.ICtx, typ reflect2.Type) codecore.IDecoder {
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

func encoderOfMap(ctx codecore.ICtx, typ reflect2.Type) codecore.IEncoder {
	mapType := typ.(*reflect2.UnsafeMapType)
	if ctx.Config().SortMapKeys {
		return &sortKeysMapEncoder{
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

func decoderOfMapKey(ctx codecore.ICtx, typ reflect2.Type) codecore.IDecoder {
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

func encoderOfMapKey(ctx codecore.ICtx, typ reflect2.Type) codecore.IEncoder {
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
type sortKeysMapEncoder struct {
	mapType     *reflect2.UnsafeMapType
	keyEncoder  codecore.IEncoder
	elemEncoder codecore.IEncoder
}

func (this *sortKeysMapEncoder) GetType() reflect.Kind {
	return reflect.Map
}
func (encoder *sortKeysMapEncoder) Encode(ptr unsafe.Pointer, w codecore.IWriter) {
	if *(*unsafe.Pointer)(ptr) == nil {
		w.WriteNil()
		return
	}
	w.WriteObjectStart()
	mapIter := encoder.mapType.UnsafeIterate(ptr)
	subStream := w.GetWriter()
	subIter := w.GetReader(nil, nil)
	keyValues := encodedKeyValues{}
	for mapIter.HasNext() {
		key, elem := mapIter.UnsafeNext()
		subStreamIndex := subStream.Buffered()
		encoder.keyEncoder.Encode(key, subStream)
		if subStream.Error() != nil && subStream.Error() != io.EOF && w.Error() == nil {
			w.SetErr(subStream.Error())
		}
		encodedKey := subStream.Buffer()[subStreamIndex:]
		subIter.ResetBytes(encodedKey, nil)
		decodedKey := subIter.ReadString()
		subStream.WriteKVSplit()
		encoder.elemEncoder.Encode(elem, subStream)
		keyValues = append(keyValues, encodedKV{
			key:      decodedKey,
			keyValue: subStream.Buffer()[subStreamIndex:],
		})
	}
	sort.Sort(keyValues)
	for i, keyValue := range keyValues {
		if i != 0 {
			w.WriteMemberSplit()
		}
		w.WriteBytes(keyValue.keyValue)
	}
	if subStream.Error() != nil && w.Error() == nil {
		w.SetErr(subStream.Error())
	}
	w.WriteObjectEnd()
	w.PutWriter(subStream)
	w.PutReader(subIter)
}

func (this *sortKeysMapEncoder) EncodeToMapJson(ptr unsafe.Pointer, w codecore.IWriter) (ret map[string]string, err error) {
	ret = make(map[string]string)
	var (
		k, v string
	)
	keystream := w.GetWriter()
	elemstream := w.GetWriter()
	defer func() {
		w.PutWriter(keystream)
		w.PutWriter(elemstream)
	}()
	iter := this.mapType.UnsafeIterate(ptr)
	for i := 0; iter.HasNext(); i++ {
		key, elem := iter.UnsafeNext()
		if this.keyEncoder.GetType() != reflect.String {
			this.keyEncoder.Encode(key, keystream)
			if keystream.Error() != nil && keystream.Error() != io.EOF {
				err = keystream.Error()
				return
			}
			k = BytesToString(keystream.Buffer())
		} else {
			k = *((*string)(key))
		}
		if this.elemEncoder.GetType() != reflect.String {
			this.elemEncoder.Encode(elem, elemstream)
			if elemstream.Error() != nil && elemstream.Error() != io.EOF {
				err = elemstream.Error()
				return
			}
			v = BytesToString(elemstream.Buffer())
		} else {
			v = *((*string)(elem))
		}
		ret[k] = v
		keystream.Reset(nil)
		elemstream.Reset(nil)
	}
	return
}

func (encoder *sortKeysMapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	iter := encoder.mapType.UnsafeIterate(ptr)
	return !iter.HasNext()
}

type mapEncoder struct {
	mapType     *reflect2.UnsafeMapType
	keyEncoder  codecore.IEncoder
	elemEncoder codecore.IEncoder
}

func (this *mapEncoder) GetType() reflect.Kind {
	return reflect.Map
}
func (this *mapEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	if *(*unsafe.Pointer)(ptr) == nil {
		stream.WriteNil()
		return
	}
	stream.WriteObjectStart()
	iter := this.mapType.UnsafeIterate(ptr)
	for i := 0; iter.HasNext(); i++ {
		if i != 0 {
			stream.WriteMemberSplit()
		}
		key, elem := iter.UnsafeNext()
		this.keyEncoder.Encode(key, stream)
		stream.WriteKVSplit()
		this.elemEncoder.Encode(elem, stream)
	}
	stream.WriteObjectEnd()
}

func (this *mapEncoder) EncodeToMapJson(ptr unsafe.Pointer, w codecore.IWriter) (ret map[string]string, err error) {
	ret = make(map[string]string)
	var (
		k, v string
	)
	keystream := w.GetWriter()
	elemstream := w.GetWriter()
	defer func() {
		w.PutWriter(keystream)
		w.PutWriter(elemstream)
	}()
	iter := this.mapType.UnsafeIterate(ptr)
	for i := 0; iter.HasNext(); i++ {
		key, elem := iter.UnsafeNext()
		if this.keyEncoder.GetType() != reflect.String {
			this.keyEncoder.Encode(key, keystream)
			if keystream.Error() != nil && keystream.Error() != io.EOF {
				err = keystream.Error()
				return
			}
			k = BytesToString(keystream.Buffer())
		} else {
			k = *((*string)(key))
		}
		if this.elemEncoder.GetType() != reflect.String {
			this.elemEncoder.Encode(elem, elemstream)
			if elemstream.Error() != nil && elemstream.Error() != io.EOF {
				err = elemstream.Error()
				return
			}
			v = BytesToString(elemstream.Buffer())
		} else {
			v = *((*string)(elem))
		}
		ret[k] = v
		keystream.Reset(nil)
		elemstream.Reset(nil)
	}
	return
}

func (this *mapEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	iter := this.mapType.UnsafeIterate(ptr)
	return !iter.HasNext()
}

type mapDecoder struct {
	mapType     *reflect2.UnsafeMapType
	keyType     reflect2.Type
	elemType    reflect2.Type
	keyDecoder  codecore.IDecoder
	elemDecoder codecore.IDecoder
}

func (codec *mapDecoder) GetType() reflect.Kind {
	return reflect.Map
}
func (this *mapDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	mapType := this.mapType
	if extra.ReadNil() {
		*(*unsafe.Pointer)(ptr) = nil
		mapType.UnsafeSet(ptr, mapType.UnsafeNew())
		return
	}
	if mapType.UnsafeIsNil(ptr) {
		mapType.UnsafeSet(ptr, mapType.UnsafeMakeMap(0))
	}
	if !extra.ReadObjectStart() {
		return
	}
	if extra.CheckNextIsObjectEnd() {
		// extra.ReadObjectEnd()
		return
	}
	key := this.keyType.UnsafeNew()
	this.keyDecoder.Decode(key, extra)
	if !extra.ReadKVSplit() {
		return
	}
	elem := this.elemType.UnsafeNew()
	this.elemDecoder.Decode(elem, extra)
	this.mapType.UnsafeSetIndex(ptr, key, elem)
	for extra.ReadMemberSplit() {
		key := this.keyType.UnsafeNew()
		this.keyDecoder.Decode(key, extra)
		if !extra.ReadKVSplit() {
			return
		}
		elem := this.elemType.UnsafeNew()
		this.elemDecoder.Decode(elem, extra)
		this.mapType.UnsafeSetIndex(ptr, key, elem)
	}
	extra.ReadObjectEnd()
}

//解码对象从MapJson 中
func (this *mapDecoder) DecodeForMapJson(ptr unsafe.Pointer, r codecore.IReader, extra map[string]string) (err error) {
	keyext := r.GetReader([]byte{}, nil)
	elemext := r.GetReader([]byte{}, nil)
	defer func() {
		r.PutReader(keyext)
		r.PutReader(keyext)
	}()
	for k, v := range extra {
		key := this.keyType.UnsafeNew()
		if this.keyDecoder.GetType() != reflect.String {
			keyext.ResetBytes(StringToBytes(k), nil)
			this.keyDecoder.Decode(key, keyext)
			if keyext.Error() != nil && keyext.Error() != io.EOF {
				err = keyext.Error()
				return
			}
		} else {
			*((*string)(key)) = k
		}
		elem := this.elemType.UnsafeNew()
		if this.elemDecoder.GetType() != reflect.String {
			elemext.ResetBytes(StringToBytes(v), nil)
			this.elemDecoder.Decode(elem, elemext)
			this.mapType.UnsafeSetIndex(ptr, key, elem)
			if elemext.Error() != nil && elemext.Error() != io.EOF {
				err = elemext.Error()
				return
			}
		} else {
			*((*string)(elem)) = v
		}
	}
	return
}

//NumericMap-------------------------------------------------------------------------------------------------------------------------------
type numericMapKeyDecoder struct {
	decoder codecore.IDecoder
}

func (this *numericMapKeyDecoder) GetType() reflect.Kind {
	return this.decoder.GetType()
}
func (this *numericMapKeyDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
	if !extra.ReadKeyStart() {
		return
	}
	this.decoder.Decode(ptr, extra)
	if !extra.ReadKeyEnd() {
		return
	}
}

type numericMapKeyEncoder struct {
	encoder codecore.IEncoder
}

func (this *numericMapKeyEncoder) GetType() reflect.Kind {
	return this.encoder.GetType()
}
func (this *numericMapKeyEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteKeyStart()
	this.encoder.Encode(ptr, stream)
	stream.WriteKeyEnd()
}

func (encoder *numericMapKeyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

//------------------------------------------------------------------------------------------------------------------
type dynamicMapKeyEncoder struct {
	ctx     codecore.ICtx
	valType reflect2.Type
}

func (this *dynamicMapKeyEncoder) GetType() reflect.Kind {
	return reflect.Interface
}

func (this *dynamicMapKeyEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	obj := this.valType.UnsafeIndirect(ptr)
	encoderOfMapKey(this.ctx, reflect2.TypeOf(obj)).Encode(reflect2.PtrOf(obj), stream)
}

func (this *dynamicMapKeyEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	obj := this.valType.UnsafeIndirect(ptr)
	return encoderOfMapKey(this.ctx, reflect2.TypeOf(obj)).IsEmpty(reflect2.PtrOf(obj))
}

type encodedKeyValues []encodedKV

type encodedKV struct {
	key      string
	keyValue []byte
}

func (sv encodedKeyValues) Len() int           { return len(sv) }
func (sv encodedKeyValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv encodedKeyValues) Less(i, j int) bool { return sv[i].key < sv[j].key }
