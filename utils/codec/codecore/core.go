package codecore

import (
	"io"
	"reflect"
	"unsafe"

	"github.com/modern-go/reflect2"
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
	IPools interface {
		GetReader(buf []byte, r io.Reader) IReader
		PutReader(w IReader)
		GetWriter() IWriter
		PutWriter(w IWriter)
	}
	IReader interface {
		IPools
		Config() *Config         //解析配置
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
		SkipAndReturnBytes() []byte       //跳过一个数据单元并返回中间的数据
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
		ResetBytes(d []byte, r io.Reader)
		Error() error
		SetErr(err error)
	}
	IWriter interface {
		IPools
		Config() *Config
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
		Write(p []byte) (nn int, err error)
		Flush() error
		Reset(w io.Writer)
		Buffer() []byte //返回缓存区数据
		Buffered() int
		Error() error
		SetErr(err error)
	}
	IDecoder interface {
		GetType() reflect.Kind
		Decode(ptr unsafe.Pointer, r IReader)
	}
	//MapJson 解码器
	IDecoderMapJson interface {
		DecodeForMapJson(ptr unsafe.Pointer, r IReader, extra map[string]string) (err error)
	}
	//MapJson 解码器
	IDecoderSliceJson interface {
		DecodeForSliceJson(ptr unsafe.Pointer, r IReader, extra []string) (err error)
	}
	IEncoder interface {
		GetType() reflect.Kind
		IsEmpty(ptr unsafe.Pointer) bool
		Encode(ptr unsafe.Pointer, w IWriter)
	}
	//MapJson 编码器
	IEncoderMapJson interface {
		EncodeToMapJson(ptr unsafe.Pointer, w IWriter) (ret map[string]string, err error)
	}
	//SliceJson 编码器
	IEncoderSliceJson interface {
		EncodeToSliceJson(ptr unsafe.Pointer, w IWriter) (ret []string, err error)
	}
	//内嵌指针是否是空
	IsEmbeddedPtrNil interface {
		IsEmbeddedPtrNil(ptr unsafe.Pointer) bool
	}
	ICtx interface {
		Config() *Config
		Prefix() string
		Append(prefix string) ICtx
		GetEncoder(rtype reflect2.Type) IEncoder
		SetEncoder(rtype reflect2.Type, encoder IEncoder)
		GetDecoder(rtype reflect2.Type) IDecoder
		SetDecoder(rtype reflect2.Type, decoder IDecoder)
		EncoderOf(typ reflect2.Type) IEncoder
		DecoderOf(typ reflect2.Type) IDecoder
	}
)

//序列化配置
type Config struct {
	EscapeHTML            bool
	SortMapKeys           bool //排序mapkey
	UseNumber             bool
	IndentionStep         int    //缩进步骤
	OnlyTaggedField       bool   //仅仅处理标签字段
	DisallowUnknownFields bool   //禁止未知字段
	CaseSensitive         bool   //是否区分大小写
	TagKey                string //标签
}
