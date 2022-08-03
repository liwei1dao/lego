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
	ctx := &codecore.Ctx{
		Config:   config,
		Prefix:   "",
		Decoders: map[reflect2.Type]codecore.IDecoder{},
		Encoders: map[reflect2.Type]codecore.IEncoder{},
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
	ctx := &codecore.Ctx{
		Config:   config,
		Prefix:   "",
		Decoders: map[reflect2.Type]codecore.IDecoder{},
		Encoders: map[reflect2.Type]codecore.IEncoder{},
	}
	ptrType := typ.(*reflect2.UnsafePtrType)
	decoder = factory.DecoderOfType(ctx, ptrType.Elem())
	AddDecoder(cacheKey, decoder)
	return decoder
}
