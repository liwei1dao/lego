package factory

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"unsafe"

	"github.com/liwei1dao/lego/utils/codec/codecore"

	"github.com/modern-go/reflect2"
)

func decoderOfArray(ctx codecore.ICtx, typ reflect2.Type) codecore.IDecoder {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	decoder := DecoderOfType(ctx.Append("[arrayElem]"), arrayType.Elem())
	return &arrayDecoder{arrayType, decoder}
}

func encoderOfArray(ctx codecore.ICtx, typ reflect2.Type) codecore.IEncoder {
	arrayType := typ.(*reflect2.UnsafeArrayType)
	if arrayType.Len() == 0 {
		return &emptyArrayEncoder{}
	}
	encoder := EncoderOfType(ctx.Append("[arrayElem]"), arrayType.Elem())
	return &arrayEncoder{arrayType, encoder}
}

//array-------------------------------------------------------------------------------------------------------------------------------
type arrayEncoder struct {
	arrayType   *reflect2.UnsafeArrayType
	elemEncoder codecore.IEncoder
}

func (codec *arrayEncoder) GetType() reflect.Kind {
	return reflect.Array
}
func (this *arrayEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
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
func (this *arrayEncoder) EncodeToSliceJson(ptr unsafe.Pointer, w codecore.IWriter) (ret []string, err error) {
	ret = make([]string, this.arrayType.Len())
	for i := 1; i < this.arrayType.Len(); i++ {
		elemPtr := this.arrayType.UnsafeGetIndex(ptr, i)
		this.elemEncoder.Encode(elemPtr, w)
		if w.Error() != nil && w.Error() != io.EOF {
			err = w.Error()
			return
		}
		ret[i] = string(w.Buffer())
		w.Reset(nil)
	}
	return
}

type arrayDecoder struct {
	arrayType   *reflect2.UnsafeArrayType
	elemDecoder codecore.IDecoder
}

func (codec *arrayDecoder) GetType() reflect.Kind {
	return reflect.Array
}
func (this *arrayDecoder) Decode(ptr unsafe.Pointer, extra codecore.IReader) {
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

func (this *arrayDecoder) DecodeForSliceJson(ptr unsafe.Pointer, r codecore.IReader, data []string) (err error) {
	arrayType := this.arrayType
	if data == nil {
		err = errors.New("extra is nil")
		return
	}
	arrayType.UnsafeGetIndex(ptr, len(data))
	for i, v := range data {
		elemPtr := arrayType.UnsafeGetIndex(ptr, i)
		r.ResetBytes(StringToBytes(v), nil)
		this.elemDecoder.Decode(elemPtr, r)
		if r.Error() != nil && r.Error() != io.EOF {
			err = r.Error()
			return
		}
	}
	return
}

type emptyArrayEncoder struct{}

func (this *emptyArrayEncoder) GetType() reflect.Kind {
	return reflect.Array
}

func (this *emptyArrayEncoder) Encode(ptr unsafe.Pointer, stream codecore.IWriter) {
	stream.WriteEmptyArray()
}

func (this *emptyArrayEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return true
}
