package codec

import (
	"sync"

	"github.com/liwei1dao/lego/utils/codec/codecore"
	"github.com/liwei1dao/lego/utils/codec/factory"

	"github.com/modern-go/reflect2"
)

var decoderCache = new(sync.Map)
var encoderCache = new(sync.Map)

func AddDecoder(cacheKey uintptr, decoder codecore.IDecoder) {
	decoderCache.Store(cacheKey, decoder)
}
func AddEncoder(cacheKey uintptr, encoder codecore.IEncoder) {
	encoderCache.Store(cacheKey, encoder)
}
func GetEncoder(cacheKey uintptr) codecore.IEncoder {
	encoder, found := encoderCache.Load(cacheKey)
	if found {
		return encoder.(codecore.IEncoder)
	}
	return nil
}
func GetDecoder(cacheKey uintptr) codecore.IDecoder {
	decoder, found := decoderCache.Load(cacheKey)
	if found {
		return decoder.(codecore.IDecoder)
	}
	return nil
}
func EncoderOf(typ reflect2.Type, config *codecore.Config) codecore.IEncoder {
	cacheKey := typ.RType()
	encoder := GetEncoder(cacheKey)
	if encoder != nil {
		return encoder
	}
	ctx := &Ctx{
		config:   config,
		prefix:   "",
		decoders: map[reflect2.Type]codecore.IDecoder{},
		encoders: map[reflect2.Type]codecore.IEncoder{},
	}
	encoder = factory.EncoderOfType(ctx, typ)
	if typ.LikePtr() {
		encoder = factory.NewonePtrEncoder(encoder)
	}
	AddEncoder(cacheKey, encoder)
	return encoder
}
func DecoderOf(typ reflect2.Type, config *codecore.Config) codecore.IDecoder {
	cacheKey := typ.RType()
	decoder := GetDecoder(cacheKey)
	if decoder != nil {
		return decoder
	}
	ctx := &Ctx{
		config:   config,
		prefix:   "",
		decoders: map[reflect2.Type]codecore.IDecoder{},
		encoders: map[reflect2.Type]codecore.IEncoder{},
	}
	ptrType := typ.(*reflect2.UnsafePtrType)
	decoder = factory.DecoderOfType(ctx, ptrType.Elem())
	AddDecoder(cacheKey, decoder)
	return decoder
}

type Ctx struct {
	config   *codecore.Config
	prefix   string
	encoders map[reflect2.Type]codecore.IEncoder
	decoders map[reflect2.Type]codecore.IDecoder
}

func (this *Ctx) Config() *codecore.Config {
	return this.config
}

func (this *Ctx) Prefix() string {
	return this.prefix
}
func (this *Ctx) GetEncoder(rtype reflect2.Type) codecore.IEncoder {
	return this.encoders[rtype]
}
func (this *Ctx) SetEncoder(rtype reflect2.Type, encoder codecore.IEncoder) {
	this.encoders[rtype] = encoder
}
func (this *Ctx) GetDecoder(rtype reflect2.Type) codecore.IDecoder {
	return this.decoders[rtype]
}
func (this *Ctx) SetDecoder(rtype reflect2.Type, decoder codecore.IDecoder) {
	this.decoders[rtype] = decoder
}
func (this *Ctx) Append(prefix string) codecore.ICtx {
	return &Ctx{
		config:   this.config,
		prefix:   this.prefix + " " + prefix,
		encoders: this.encoders,
		decoders: this.decoders,
	}
}

func (this *Ctx) EncoderOf(typ reflect2.Type) codecore.IEncoder {
	return EncoderOf(typ, this.config)
}
func (this *Ctx) DecoderOf(typ reflect2.Type) codecore.IDecoder {
	return DecoderOf(typ, this.config)
}
