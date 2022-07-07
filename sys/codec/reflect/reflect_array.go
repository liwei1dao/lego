package reflect

import (
	"fmt"
	"io"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

func decoderOfArray(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	decoder := decoderOfType(ctx.Append("[arrayElem]"), arrayType.Elem())
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

func (encoder *arrayEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	stream.WriteArrayStart()
	elemPtr := unsafe.Pointer(ptr)
	encoder.elemEncoder.Encode(elemPtr, stream, opt)
	for i := 1; i < encoder.arrayType.Len(); i++ {
		stream.WriteMore()
		elemPtr = encoder.arrayType.UnsafeGetIndex(ptr, i)
		encoder.elemEncoder.Encode(elemPtr, stream, opt)
	}
	stream.WriteArrayEnd()
	if stream.Error() != nil && stream.Error() != io.EOF {
		stream.SetError(fmt.Errorf("%v: %s", encoder.arrayType, stream.Error().Error()))
	}
}

func (encoder *arrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

type arrayDecoder struct {
	arrayType   *reflect2.UnsafeArrayType
	elemDecoder core.IDecoder
}

func (this *arrayDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
	this.doDecode(ptr, extra, opt)
	if extra.Error() != nil && extra.Error() != io.EOF {
		extra.SetError(fmt.Errorf("%v: %s", this.arrayType, extra.Error().Error()))
	}
}
func (this *arrayDecoder) doDecode(ptr unsafe.Pointer, extra core.IExtractor, opt *core.ExecuteOptions) {
	c := extra.NextToken()
	arrayType := this.arrayType
	if c == 'n' {
		extra.SkipBytes([]byte{'u', 'l', 'l'})
		return
	}
	if c != '[' {
		extra.ReportError("decode array", "expect [ or n, but found "+string([]byte{c}))
		return
	}
	c = extra.NextToken()
	if c == ']' {
		return
	}
	extra.UnreadChar()
	elemPtr := arrayType.UnsafeGetIndex(ptr, 0)
	this.elemDecoder.Decode(elemPtr, extra, opt)
	length := 1
	for c = extra.NextToken(); c == ','; c = extra.NextToken() {
		if length >= arrayType.Len() {
			extra.Skip()
			continue
		}
		idx := length
		length += 1
		elemPtr = arrayType.UnsafeGetIndex(ptr, idx)
		this.elemDecoder.Decode(elemPtr, extra, opt)
	}
	if c != ']' {
		extra.ReportError("decode array", "expect ], but found "+string([]byte{c}))
		return
	}
}

type emptyArrayEncoder struct{}

func (this emptyArrayEncoder) Encode(ptr unsafe.Pointer, stream core.IStream, opt *core.ExecuteOptions) {
	stream.WriteEmptyArray()
}

func (this emptyArrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return true
}
