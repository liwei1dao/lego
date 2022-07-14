package core

import (
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
)

//数据类型
type CodecType int8

const (
	Json  CodecType = iota //json 格式数据
	Bytes                  //byte 字节流数据
	Proto                  //proto google protobuff 数据
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
	ICodec interface {
		Options() *Options
		GetEncoderFromCache(cacheKey uintptr) IEncoder
		GetDecoderFromCache(cacheKey uintptr) IDecoder
		EncoderOf(typ reflect2.Type) IEncoder
		DecoderOf(typ reflect2.Type) IDecoder
		BorrowExtractor(buf []byte) IExtractor //借 提取器
		ReturnExtractor(extra IExtractor)      //还 提取器
		BorrowStream() IStream                 //借 输出对象
		ReturnStream(stream IStream)           //换 输出对象
	}
	IExtractor interface {
		ReadVal(obj interface{}) //读取指定类型对象
		WhatIsNext() ValueType
		Read() interface{}
		ReadNil() (ret bool)              //读取空 null
		ReadArrayStart() (ret bool)       //读数组开始 [
		CheckNextIsArrayEnd() (ret bool)  //下一个是否是数组结束符 不影响数据游标位置
		ReadArrayEnd() (ret bool)         //读数组结束 ]
		ReadObjectStart() (ret bool)      //读对象或Map开始 {
		CheckNextIsObjectEnd() (ret bool) //下一个是否是数对象结束符 不影响数据游标位置
		ReadObjectEnd() (ret bool)        //读对象或Map结束 }
		ReadMemberSplit() (ret bool)      //读成员分割符 ,
		ReadKVSplit() (ret bool)          //读取key value 分割符 :
		ReadKeyStart() (ret bool)         //读取字段名 开始
		ReadKeyEnd() (ret bool)           //读取字段名 结束
		Skip()                            //跳过一个数据单元 容错处理
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
		ResetBytes(d []byte)
		Error() error
		SetErr(err error)
	}
	IStream interface {
		WriteVal(val interface{})        //写入一个对象
		WriteNil()                       //写空 null
		WriteEmptyArray()                //写空数组 []
		WriteArrayStart()                //写数组开始 [
		WriteArrayEnd()                  //写数组结束 ]
		WriteEmptyObject()               //写空对象 {}
		WriteObjectStart()               //写对象或Map开始 {
		WriteObjectEnd()                 //写对象或Map结束 }
		WriteMemberSplit()               //写成员分割符 ,
		WriteKVSplit()                   //写key value 分割符 :
		WriteKeyStart()                  //写入字段名 开始
		WriteKeyEnd()                    //写入字段名 结束
		WriteObjectFieldName(val string) //写入对象字段名
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
		WriteString(val string)
		WriteBytes(val []byte)
		Reset(bufSize int)
		Buffer() []byte //返回缓存区数据
		Error() error
		SetErr(err error)
	}
	//Json 编码器
	IEncoder interface {
		GetType() reflect.Kind
		IsEmpty(ptr unsafe.Pointer) bool
		Encode(ptr unsafe.Pointer, stream IStream)
	}
	//MapJson 编码器
	IEncoderMapJson interface {
		EncodeToMapJson(ptr unsafe.Pointer) (ret map[string]string, err error)
	}
	//SliceJson 编码器
	IEncoderSliceJson interface {
		EncodeToSliceJson(ptr unsafe.Pointer) (ret []string, err error)
	}

	//Json 解码器
	IDecoder interface {
		GetType() reflect.Kind
		Decode(ptr unsafe.Pointer, extra IExtractor)
	}
	//MapJson 解码器
	IDecoderMapJson interface {
		DecodeForMapJson(ptr unsafe.Pointer, extra map[string]string) (err error)
	}
	//MapJson 解码器
	IDecoderSliceJson interface {
		DecodeForSliceJson(ptr unsafe.Pointer, extra []string) (err error)
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
