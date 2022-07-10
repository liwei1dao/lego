package factory

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

func decoderOfArray(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	decoder := DecoderOfType(ctx.Append("[arrayElem]"), arrayType.Elem())
	return &arrayDecoder{arrayType, decoder}
}

func encoderOfArray(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	if arrayType.Len() == 0 {
		return emptyArrayEncoder{}
	}
	encoder := EncoderOfType(ctx.Append("[arrayElem]"), arrayType.Elem())
	return &arrayEncoder{arrayType, encoder}
}

//array-------------------------------------------------------------------------------------------------------------------------------
type arrayEncoder struct {
	arrayType   *reflect2.UnsafeArrayType
	elemEncoder core.IEncoder
}

func (encoder *arrayEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	stream.WriteArrayStart()
	elemPtr := unsafe.Pointer(ptr)
	encoder.elemEncoder.Encode(elemPtr, stream)
	for i := 1; i < encoder.arrayType.Len(); i++ {
		stream.WriteMemberSplit()
		elemPtr = encoder.arrayType.UnsafeGetIndex(ptr, i)
		encoder.elemEncoder.Encode(elemPtr, stream)
	}
	stream.WriteArrayEnd()
	if stream.Error() != nil && stream.Error() != io.EOF {
		stream.SetErr(fmt.Errorf("%v: %s", encoder.arrayType, stream.Error().Error()))
	}
}

func (encoder *arrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type arrayDecoder struct {
	arrayType   *reflect2.UnsafeArrayType
	elemDecoder core.IDecoder
}

func (this *arrayDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	arrayType := this.arrayType
	if extra.ReadNil() {
		return
	}
	if extra.ReadArrayStart() {
		return
	}
	if extra.CheckNextIsArrayEnd() {
		return
	}
	elemPtr := arrayType.UnsafeGetIndex(ptr, 0)
	this.elemDecoder.Decode(elemPtr, extra)
	length := 1
	for extra.ReadMemberSplit() {
		idx := length
		length += 1
		elemPtr = arrayType.UnsafeGetIndex(ptr, idx)
		this.elemDecoder.Decode(elemPtr, extra)
	}
	if extra.ReadArrayEnd() {
		return
	}
	if extra.Error() != nil && extra.Error() != io.EOF {
		extra.SetErr(fmt.Errorf("%v: %s", this.arrayType, extra.Error().Error()))
	}
}

type emptyArrayEncoder struct{}

func (this emptyArrayEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	stream.WriteEmptyArray()
}

func (this emptyArrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return true
}
