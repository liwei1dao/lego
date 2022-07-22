package codec

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/liwei1dao/lego/sys/codec/factory"
	"github.com/liwei1dao/lego/sys/codec/render"

	"github.com/modern-go/reflect2"
)

func newSys(options core.Options) (sys *Codec, err error) {
	sys = &Codec{
		options:      &options,
		decoderCache: new(sync.Map),
		encoderCache: new(sync.Map),
		streamPool: &sync.Pool{
			New: func() interface{} {
				return render.NewStream(sys, 512)
			},
		},
		extraPool: &sync.Pool{
			New: func() interface{} {
				return render.NewExtractor(sys)
			},
		},
	}
	return
}

type Codec struct {
	options      *core.Options
	decoderCache *sync.Map
	encoderCache *sync.Map
	streamPool   *sync.Pool
	extraPool    *sync.Pool
}

func (this *Codec) Options() *core.Options {
	return this.options
}

func (this *Codec) addDecoderToCache(cacheKey uintptr, decoder core.IDecoder) {
	this.decoderCache.Store(cacheKey, decoder)
}
func (this *Codec) addEncoderToCache(cacheKey uintptr, encoder core.IEncoder) {
	this.encoderCache.Store(cacheKey, encoder)
}

func (this *Codec) GetEncoderFromCache(cacheKey uintptr) core.IEncoder {
	encoder, found := this.encoderCache.Load(cacheKey)
	if found {
		return encoder.(core.IEncoder)
	}
	return nil
}

func (this *Codec) GetDecoderFromCache(cacheKey uintptr) core.IDecoder {
	decoder, found := this.decoderCache.Load(cacheKey)
	if found {
		return decoder.(core.IDecoder)
	}
	return nil
}

func (this *Codec) EncoderOf(typ reflect2.Type) core.IEncoder {
	cacheKey := typ.RType()
	encoder := this.GetEncoderFromCache(cacheKey)
	if encoder != nil {
		return encoder
	}
	ctx := &core.Ctx{
		ICodec:   this,
		Prefix:   "",
		Decoders: map[reflect2.Type]core.IDecoder{},
		Encoders: map[reflect2.Type]core.IEncoder{},
	}
	encoder = factory.EncoderOfType(ctx, typ)
	if typ.LikePtr() {
		encoder = factory.NewonePtrEncoder(encoder)
	}
	this.addEncoderToCache(cacheKey, encoder)
	return encoder
}

func (this *Codec) DecoderOf(typ reflect2.Type) core.IDecoder {
	cacheKey := typ.RType()
	decoder := this.GetDecoderFromCache(cacheKey)
	if decoder != nil {
		return decoder
	}
	ctx := &core.Ctx{
		ICodec:   this,
		Prefix:   "",
		Decoders: map[reflect2.Type]core.IDecoder{},
		Encoders: map[reflect2.Type]core.IEncoder{},
	}
	ptrType := typ.(*reflect2.UnsafePtrType)
	decoder = factory.DecoderOfType(ctx, ptrType.Elem())
	this.addDecoderToCache(cacheKey, decoder)
	return decoder
}

func (this *Codec) BorrowStream() core.IStream {
	stream := this.streamPool.Get().(core.IStream)
	return stream
}

func (this *Codec) ReturnStream(stream core.IStream) {
	this.streamPool.Put(stream)
}

func (this *Codec) BorrowExtractor(buf []byte) core.IExtractor {
	extra := this.extraPool.Get().(core.IExtractor)
	extra.ResetBytes(buf)
	return extra
}

func (this *Codec) ReturnExtractor(extra core.IExtractor) {
	this.extraPool.Put(extra)
}

//编码对象到json
func (this *Codec) MarshalJson(val interface{}, option ...core.ExecuteOption) (buf []byte, err error) {
	stream := this.BorrowStream()
	defer this.ReturnStream(stream)
	stream.WriteVal(val)
	if stream.Error() != nil {
		return nil, stream.Error()
	}
	result := stream.Buffer()
	copied := make([]byte, len(result))
	copy(copied, result)
	return copied, nil
}

//解码json到对象
func (this *Codec) UnmarshalJson(data []byte, v interface{}, option ...core.ExecuteOption) error {
	extra := this.BorrowExtractor(data)
	defer this.ReturnExtractor(extra)
	extra.ReadVal(v)
	return extra.Error()
}

//编码对象到mapjson
func (this *Codec) MarshalMapJson(val interface{}, option ...core.ExecuteOption) (ret map[string]string, err error) {
	if nil == val {
		err = errors.New("val is null")
		return
	}
	cacheKey := reflect2.RTypeOf(val)
	encoder := this.GetEncoderFromCache(cacheKey)
	if encoder == nil {
		typ := reflect2.TypeOf(val)
		encoder = this.EncoderOf(typ)
	}
	if encoderMapJson, ok := encoder.(core.IEncoderMapJson); !ok {
		err = fmt.Errorf("val type:%T not support MarshalMapJson", val)
	} else {
		ret, err = encoderMapJson.EncodeToMapJson(reflect2.PtrOf(val))
	}
	return
}

//解码mapjson到对象
func (this *Codec) UnmarshalMapJson(data map[string]string, val interface{}, option ...core.ExecuteOption) (err error) {
	cacheKey := reflect2.RTypeOf(val)
	decoder := this.GetDecoderFromCache(cacheKey)
	if decoder == nil {
		typ := reflect2.TypeOf(val)
		if typ == nil || typ.Kind() != reflect.Ptr {
			err = errors.New("can only unmarshal into pointer")
			return
		}
		decoder = this.DecoderOf(typ)
	}
	ptr := reflect2.PtrOf(val)
	if ptr == nil {
		err = errors.New("can not read into nil pointer")
		return
	}
	if decoderMapJson, ok := decoder.(core.IDecoderMapJson); !ok {
		err = fmt.Errorf("val type:%T not support MarshalMapJson", val)
	} else {
		err = decoderMapJson.DecodeForMapJson(ptr, data)
	}
	return
}

//编码对象到sliceJson
func (this *Codec) MarshalSliceJson(val interface{}, option ...core.ExecuteOption) (ret []string, err error) {
	if nil == val {
		err = errors.New("val is null")
		return
	}
	cacheKey := reflect2.RTypeOf(val)
	encoder := this.GetEncoderFromCache(cacheKey)
	if encoder == nil {
		typ := reflect2.TypeOf(val)
		encoder = this.EncoderOf(typ)
	}
	if encoderMapJson, ok := encoder.(core.IEncoderSliceJson); !ok {
		err = fmt.Errorf("val type:%T not support MarshalMapJson", val)
	} else {
		ret, err = encoderMapJson.EncodeToSliceJson(reflect2.PtrOf(val))
	}
	return
}

//解码sliceJson到对象
func (this *Codec) UnmarshalSliceJson(data []string, val interface{}, option ...core.ExecuteOption) (err error) {
	cacheKey := reflect2.RTypeOf(val)
	decoder := this.GetDecoderFromCache(cacheKey)
	if decoder == nil {
		typ := reflect2.TypeOf(val)
		if typ == nil || typ.Kind() != reflect.Ptr {
			err = errors.New("can only unmarshal into pointer")
			return
		}
		decoder = this.DecoderOf(typ)
	}
	ptr := reflect2.PtrOf(val)
	if ptr == nil {
		err = errors.New("can not read into nil pointer")
		return
	}
	if decoderMapJson, ok := decoder.(core.IDecoderSliceJson); !ok {
		err = fmt.Errorf("val type:%T not support UnmarshalSliceJson", val)
	} else {
		err = decoderMapJson.DecodeForSliceJson(ptr, data)
	}
	return
}

///日志***********************************************************************
func (this *Codec) Debug() bool {
	return this.options.Debug
}

func (this *Codec) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Debugf("[SYS Codec] "+format, a...)
	}
}
func (this *Codec) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Infof("[SYS Codec] "+format, a...)
	}
}
func (this *Codec) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Warnf("[SYS Codec] "+format, a...)
	}
}
func (this *Codec) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Errorf("[SYS Codec] "+format, a...)
	}
}
func (this *Codec) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Panicf("[SYS Codec] "+format, a...)
	}
}
func (this *Codec) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Fatalf("[SYS Codec] "+format, a...)
	}
}
