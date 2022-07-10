package codec

import (
	"sync"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/liwei1dao/lego/sys/codec/factory"
	"github.com/liwei1dao/lego/sys/codec/stream"
	"github.com/modern-go/reflect2"
)

func newSys(options core.Options) (sys *codec, err error) {
	sys = &codec{
		options:      &options,
		decoderCache: new(sync.Map),
		encoderCache: new(sync.Map),
		streamPool: &sync.Pool{
			New: func() interface{} {
				return stream.NewStream(sys, 512)
			},
		},
	}
	return
}

type codec struct {
	options      *core.Options
	decoderCache *sync.Map
	encoderCache *sync.Map
	streamPool   *sync.Pool
	extraPool    *sync.Pool
}

func (this *codec) Options() *core.Options {
	return this.options
}

//序列化Josn
func (this *codec) MarshalJson(val interface{}, option ...core.ExecuteOption) (buf []byte, err error) {
	stream := this.BorrowStream()
	defer this.ReturnStream(stream)
	stream.WriteVal(val)
	if stream.Error != nil {
		return nil, stream.Error()
	}
	result := stream.Buffer()
	copied := make([]byte, len(result))
	copy(copied, result)
	return copied, nil
}

func (this *codec) Unmarshal(data []byte, v interface{}) error {
	return nil
}

func (this *codec) addEncoderToCache(cacheKey uintptr, encoder core.IEncoder) {
	this.encoderCache.Store(cacheKey, encoder)
}

func (this *codec) GetEncoderFromCache(cacheKey uintptr) core.IEncoder {
	encoder, found := this.encoderCache.Load(cacheKey)
	if found {
		return encoder.(core.IEncoder)
	}
	return nil
}

func (this *codec) EncoderOf(typ reflect2.Type) core.IEncoder {
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

func (this *codec) BorrowStream() core.IStream {
	stream := this.streamPool.Get().(core.IStream)
	return stream
}

func (this *codec) ReturnStream(stream core.IStream) {
	this.streamPool.Put(stream)
}

func (this *codec) BorrowExtractor() core.IExtractor {
	return this.extraPool.Get().(core.IExtractor)
}

func (this *codec) ReturnExtractor(extra core.IExtractor) {
	this.extraPool.Put(extra)
}

///日志***********************************************************************
func (this *codec) Debug() bool {
	return this.options.Debug
}

func (this *codec) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Debugf("[SYS Gin] "+format, a...)
	}
}
func (this *codec) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Infof("[SYS Gin] "+format, a...)
	}
}
func (this *codec) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Warnf("[SYS Gin] "+format, a...)
	}
}
func (this *codec) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Errorf("[SYS Gin] "+format, a...)
	}
}
func (this *codec) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Panicf("[SYS Gin] "+format, a...)
	}
}
func (this *codec) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Fatalf("[SYS Gin] "+format, a...)
	}
}
