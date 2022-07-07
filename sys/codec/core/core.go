package core

import (
	"unsafe"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/modern-go/reflect2"
)

//数据类型
type DataType int8

const (
	Json  DataType = iota //json 格式数据
	Bytes                 //byte 字节流数据
	Proto                 //proto google protobuff 数据
)

//结构值类型
type ValueType int

const (
	// InvalidValue invalid JSON element
	InvalidValue ValueType = iota
	// StringValue JSON element "string"
	StringValue
	// NumberValue JSON element 100 or 0.10
	NumberValue
	// NilValue JSON element null
	NilValue
	// BoolValue JSON element true or false
	BoolValue
	// ArrayValue JSON element []
	ArrayValue
	// ObjectValue JSON element {}
	ObjectValue
)

type (
	//编解码控制器
	ICodec interface {
		log.Ilogf
		Options() *Options
		GetEncoderFromCache(cacheKey uintptr) IEncoder
		EncoderOf(typ reflect2.Type) IEncoder
		BorrowExtractor() IExtractor      //借 提取器
		ReturnExtractor(extra IExtractor) //还 提取器
		BorrowStream() IStream            //借 输出对象
		ReturnStream(stream IStream)      //换 输出对象
	}
	//提取器
	IExtractor interface {
		ReadObjectStart() bool
		IncrementDepth() (success bool) //深度增加
		DecrementDepth() (success bool) //深度减少
		ReadNil() (ret bool)
		ReadChar() (ret byte)
		ReadBool() (ret bool)
		ReadInt8() (ret int8)
		ReadInt16() (ret int16)
		ReadInt32() (ret int32)
		ReadInt64() (ret int64)
		ReadUint8() (ret uint8)
		ReadUint16() (ret uint16)
		ReadUint32() (ret uint32)
		ReadUint64() (ret uint64)
		ReadFloat32() (ret float32)
		ReadFloat64() (ret float64)
		ReadString() (ret string)
		ReadStringAsSlice() (ret []byte)
		WhatIsNext() ValueType
		Skip()
		SkipBytes(chars []byte)
		SkipAndReturnBytes() []byte
		UnreadChar()
		NextToken() byte
		ResetBytes(input []byte)
		Read() interface{}
		ReadVal(obj interface{})
		ReportError(operation string, msg string)
		Error() error
		SetError(err error)
	}
	//输出器
	IStream interface {
		WriteVal(val interface{}, opt *ExecuteOptions)
		WriteObjectStart()
		WriteObjectField(field string)
		WriteMore()
		WriteObjectEnd()
		WriteEmptyObject()
		WriteArrayStart()
		WriteArrayEnd()
		WriteEmptyArray()
		Indention() int     //缩进
		SetIndention(v int) //设置缩进
		WriteNil()
		WriteChar(d byte)
		WriteBool(val bool)
		WriteInt8(val int8)
		WriteInt16(val int16)
		WriteInt32(val int32)
		WriteInt64(val int64)
		WriteUint8(val uint8)
		WriteUint16(val uint16)
		WriteUint32(val uint32)
		WriteUint64(val uint64)
		WriteFloat32(val float32)
		WriteFloat64(val float64)
		WriteBytes(d []byte)
		WriteString(d string)
		ToBuffer() []byte //输出字节缓存区
		Buffered() int    //返回已写入当前缓冲区的字节数。
		Error() error
		SetError(err error)
	}
	//编码器
	IEncoder interface {
		IsEmpty(ptr unsafe.Pointer) bool
		Encode(ptr unsafe.Pointer, stream IStream, opt *ExecuteOptions)
	}
	//解码器
	IDecoder interface {
		Decode(ptr unsafe.Pointer, extra IExtractor, opt *ExecuteOptions)
	}
	//空校验
	CheckIsEmpty interface {
		IsEmpty(ptr unsafe.Pointer) bool
	}
	//内嵌指针是否是空
	IsEmbeddedPtrNil interface {
		IsEmbeddedPtrNil(ptr unsafe.Pointer) bool
	}
)
type Ctx struct {
	ICodec
	Prefix   string
	Encoders map[reflect2.Type]IEncoder
	Decoders map[reflect2.Type]IDecoder
}

func (this *Ctx) Append(prefix string) *Ctx {
	return &Ctx{
		ICodec:   this.ICodec,
		Prefix:   this.Prefix + " " + prefix,
		Encoders: this.Encoders,
		Decoders: this.Decoders,
	}
}
