package factory

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/liwei1dao/lego/sys/codec/core"

	"github.com/modern-go/reflect2"
)

func decoderOfSlice(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	sliceType := typ.(*reflect2.UnsafeSliceType)
	decoder := DecoderOfType(ctx.Append("[sliceElem]"), sliceType.Elem())
	return &sliceDecoder{ctx.ICodec, sliceType, decoder}
}
func encoderOfSlice(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	sliceType := typ.(*reflect2.UnsafeSliceType)
	encoder := EncoderOfType(ctx.Append("[sliceElem]"), sliceType.Elem())
	return &sliceEncoder{ctx.ICodec, sliceType, encoder}
}

type sliceEncoder struct {
	codec       core.ICodec
	sliceType   *reflect2.UnsafeSliceType
	elemEncoder core.IEncoder
}

func (codec *sliceEncoder) GetType() reflect.Kind {
	return reflect.Slice
}
func (this *sliceEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	if this.sliceType.UnsafeIsNil(ptr) {
		stream.WriteNil()
		return
	}
	length := this.sliceType.UnsafeLengthOf(ptr)
	if length == 0 {
		stream.WriteEmptyArray()
		return
	}
	stream.WriteArrayStart()
	this.elemEncoder.Encode(this.sliceType.UnsafeGetIndex(ptr, 0), stream)
	for i := 1; i < length; i++ {
		stream.WriteMemberSplit()
		elemPtr := this.sliceType.UnsafeGetIndex(ptr, i)
		this.elemEncoder.Encode(elemPtr, stream)
	}
	stream.WriteArrayEnd()
	if stream.Error() != nil && stream.Error() != io.EOF {
		stream.SetErr(fmt.Errorf("%v: %s", this.sliceType, stream.Error().Error()))
	}
}

func (this *sliceEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return this.sliceType.UnsafeLengthOf(ptr) == 0
}

//编码对象到json数组
func (this *sliceEncoder) EncodeToSliceJson(ptr unsafe.Pointer) (ret []string, err error) {
	if this.sliceType.UnsafeIsNil(ptr) {
		err = errors.New("val is nil")
		return
	}
	length := this.sliceType.UnsafeLengthOf(ptr)
	ret = make([]string, length)
	if length == 0 {
		return
	}
	stream := this.codec.BorrowStream()
	for i := 1; i < length; i++ {
		elemPtr := this.sliceType.UnsafeGetIndex(ptr, i)
		this.elemEncoder.Encode(elemPtr, stream)
		if stream.Error() != nil && stream.Error() != io.EOF {
			err = stream.Error()
			return
		}
		ret[i] = BytesToString(stream.Buffer())
		stream.Reset(512)
	}
	return
}

type sliceDecoder struct {
	codec       core.ICodec
	sliceType   *reflect2.UnsafeSliceType
	elemDecoder core.IDecoder
}

func (codec *sliceDecoder) GetType() reflect.Kind {
	return reflect.Slice
}
func (this *sliceDecoder) Decode(ptr unsafe.Pointer, extra core.IExtractor) {
	sliceType := this.sliceType
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
	this.elemDecoder.Decode(elemPtr, extra)
	length := 1
	for extra.ReadMemberSplit() {
		idx := length
		length += 1
		sliceType.UnsafeGrow(ptr, length)
		elemPtr = sliceType.UnsafeGetIndex(ptr, idx)
		this.elemDecoder.Decode(elemPtr, extra)
	}
	if extra.ReadArrayEnd() {
		return
	}
	if extra.Error() != nil && extra.Error() != io.EOF {
		extra.SetErr(fmt.Errorf("%v: %s", this.sliceType, extra.Error().Error()))
	}
}

func (this *sliceDecoder) DecodeForSliceJson(ptr unsafe.Pointer, data []string) (err error) {
	sliceType := this.sliceType
	if data == nil {
		err = errors.New("extra is nil")
		return
	}
	extra := this.codec.BorrowExtractor([]byte{})
	sliceType.UnsafeGrow(ptr, len(data))
	for i, v := range data {
		elemPtr := sliceType.UnsafeGetIndex(ptr, i)
		extra.ResetBytes(StringToBytes(v))
		this.elemDecoder.Decode(elemPtr, extra)
		if extra.Error() != nil && extra.Error() != io.EOF {
			err = extra.Error()
			return
		}
	}
	return
}
