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

func decoderOfArray(ctx *core.Ctx, typ reflect2.Type) core.IDecoder {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	decoder := DecoderOfType(ctx.Append("[arrayElem]"), arrayType.Elem())
	return &arrayDecoder{ctx.ICodec, arrayType, decoder}
}

func encoderOfArray(ctx *core.Ctx, typ reflect2.Type) core.IEncoder {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	if arrayType.Len() == 0 {
		return &emptyArrayEncoder{}
	}
	encoder := EncoderOfType(ctx.Append("[arrayElem]"), arrayType.Elem())
	return &arrayEncoder{ctx.ICodec, arrayType, encoder}
}

//array-------------------------------------------------------------------------------------------------------------------------------
type arrayEncoder struct {
	codec       core.ICodec
	arrayType   *reflect2.UnsafeArrayType
	elemEncoder core.IEncoder
}

func (codec *arrayEncoder) GetType() reflect.Kind {
	return reflect.Array
}
func (this *arrayEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	stream.WriteArrayStart()
	elemPtr := unsafe.Pointer(ptr)
	this.elemEncoder.Encode(elemPtr, stream)
	for i := 1; i < this.arrayType.Len(); i++ {
		stream.WriteMemberSplit()
		elemPtr = this.arrayType.UnsafeGetIndex(ptr, i)
		this.elemEncoder.Encode(elemPtr, stream)
	}
	stream.WriteArrayEnd()
	if stream.Error() != nil && stream.Error() != io.EOF {
		stream.SetErr(fmt.Errorf("%v: %s", this.arrayType, stream.Error().Error()))
	}
}

func (this *arrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return false
}

//编码对象到json数组
func (this *arrayEncoder) EncodeToSliceJson(ptr unsafe.Pointer) (ret []string, err error) {
	ret = make([]string, this.arrayType.Len())
	stream := this.codec.BorrowStream()
	for i := 1; i < this.arrayType.Len(); i++ {
		elemPtr := this.arrayType.UnsafeGetIndex(ptr, i)
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

type arrayDecoder struct {
	codec       core.ICodec
	arrayType   *reflect2.UnsafeArrayType
	elemDecoder core.IDecoder
}

func (codec *arrayDecoder) GetType() reflect.Kind {
	return reflect.Array
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

func (this *arrayDecoder) DecodeForSliceJson(ptr unsafe.Pointer, data []string) (err error) {
	arrayType := this.arrayType
	if data == nil {
		err = errors.New("extra is nil")
		return
	}
	extra := this.codec.BorrowExtractor([]byte{})
	arrayType.UnsafeGetIndex(ptr, len(data))
	for i, v := range data {
		elemPtr := arrayType.UnsafeGetIndex(ptr, i)
		extra.ResetBytes(StringToBytes(v))
		this.elemDecoder.Decode(elemPtr, extra)
		if extra.Error() != nil && extra.Error() != io.EOF {
			err = extra.Error()
			return
		}
	}
	return
}

type emptyArrayEncoder struct{}

func (this *emptyArrayEncoder) GetType() reflect.Kind {
	return reflect.Array
}

func (this *emptyArrayEncoder) Encode(ptr unsafe.Pointer, stream core.IStream) {
	stream.WriteEmptyArray()
}

func (this *emptyArrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return true
}
