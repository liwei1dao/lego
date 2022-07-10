package factory

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

func decoderOfSlice(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	sliceType := typ.(*reflect2.UnsafeSliceType)
	decoder := DecoderOfType(ctx.Append("[sliceElem]"), sliceType.Elem())
	return &sliceDecoder{sliceType, decoder}
}
func encoderOfSlice(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	sliceType := typ.(*reflect2.UnsafeSliceType)
	encoder := EncoderOfType(ctx.Append("[sliceElem]"), sliceType.Elem())
	return &sliceEncoder{sliceType, encoder}
}

type sliceEncoder struct {
	sliceType   *reflect2.UnsafeSliceType
	elemEncoder core.IEncoder
}

func (encoder *sliceEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	if encoder.sliceType.UnsafeIsNil(ptr) {
		stream.WriteNil()
		return
	}
	length := encoder.sliceType.UnsafeLengthOf(ptr)
	if length == 0 {
		stream.WriteEmptyArray()
		return
	}
	stream.WriteArrayStart()
	encoder.elemEncoder.Encode(encoder.sliceType.UnsafeGetIndex(ptr, 0), stream)
	for i := 1; i < length; i++ {
		stream.WriteMemberSplit()
		elemPtr := encoder.sliceType.UnsafeGetIndex(ptr, i)
		encoder.elemEncoder.Encode(elemPtr, stream)
	}
	stream.WriteArrayEnd()
	if stream.Error() != nil && stream.Error() != io.EOF {
		stream.SetErr(fmt.Errorf("%v: %s", encoder.sliceType, stream.Error().Error()))
	}
}

func (encoder *sliceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return encoder.sliceType.UnsafeLengthOf(ptr) == 0
}

type sliceDecoder struct {
	sliceType   *reflect2.UnsafeSliceType
	elemDecoder core.IDecoder
}

func (decoder *sliceDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	sliceType := decoder.sliceType
	if extra.ReadNil() {
		sliceType.UnsafeSetNil(ptr)
		return
	}
	if !extra.ReadArrayStart() {
		return
	}
	if extra.CheckNextIsArrayEnd() {
		sliceType.UnsafeSet(ptr, sliceType.UnsafeMakeSlice(0, 0))
		return
	}
	sliceType.UnsafeGrow(ptr, 1)
	elemPtr := sliceType.UnsafeGetIndex(ptr, 0)
	decoder.elemDecoder.Decode(elemPtr, extra)
	length := 1
	for extra.ReadMemberSplit() {
		idx := length
		length += 1
		sliceType.UnsafeGrow(ptr, length)
		elemPtr = sliceType.UnsafeGetIndex(ptr, idx)
		decoder.elemDecoder.Decode(elemPtr, extra)
	}
	if extra.ReadArrayEnd() {
		return
	}
	if extra.Error() != nil && extra.Error() != io.EOF {
		extra.SetErr(fmt.Errorf("%v: %s", decoder.sliceType, extra.Error().Error()))
	}
}
