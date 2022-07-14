package render

import (
	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/liwei1dao/lego/sys/codec/utils"

	"github.com/modern-go/reflect2"
)

func NewStream(codec core.ICodec, bufSize int) *JsonStream {
	return &JsonStream{
		codec:     codec,
		buf:       make([]byte, 0, bufSize),
		err:       nil,
		indention: 0,
	}
}

type JsonStream struct {
	codec     core.ICodec
	err       error
	buf       []byte
	indention int
}

//写入对象
func (this *JsonStream) WriteVal(val interface{}) {
	if nil == val {
		this.WriteNil()
		return
	}
	cacheKey := reflect2.RTypeOf(val)
	encoder := this.codec.GetEncoderFromCache(cacheKey)
	if encoder == nil {
		typ := reflect2.TypeOf(val)
		encoder = this.codec.EncoderOf(typ)
	}
	encoder.Encode(reflect2.PtrOf(val), this)
}

func (this *JsonStream) WriteNil() {
	this.writeFourBytes('n', 'u', 'l', 'l')
}
func (this *JsonStream) WriteEmptyArray() {
	this.writeTwoBytes('[', ']')
}
func (this *JsonStream) WriteArrayStart() {
	this.indention += this.codec.Options().IndentionStep
	this.writeByte('[')
	this.writeIndention(0)
}
func (this *JsonStream) WriteArrayEnd() {
	this.writeIndention(this.codec.Options().IndentionStep)
	this.indention -= this.codec.Options().IndentionStep
	this.writeByte(']')
}
func (this *JsonStream) WriteEmptyObject() {
	this.writeTwoBytes('{', '}')
}
func (this *JsonStream) WriteObjectStart() {
	this.indention += this.codec.Options().IndentionStep
	this.writeByte('{')
	this.writeIndention(0)
}
func (this *JsonStream) WriteObjectEnd() {
	this.writeIndention(this.codec.Options().IndentionStep)
	this.indention -= this.codec.Options().IndentionStep
	this.writeByte('}')
}
func (this *JsonStream) WriteMemberSplit() {
	this.writeByte(',')
	this.writeIndention(0)
}
func (this *JsonStream) WriteKVSplit() {
	this.writeByte(':')
}
func (this *JsonStream) WriteKeyStart() {
	this.writeByte('"')
}
func (this *JsonStream) WriteKeyEnd() {
	this.writeByte('"')
}
func (this *JsonStream) WriteObjectFieldName(val string) {
	this.WriteString(val)
	if this.indention > 0 {
		this.writeTwoBytes(':', ' ')
	} else {
		this.writeByte(':')
	}
}
func (this *JsonStream) WriteBool(val bool) {
	if val {
		this.writeTrue()
	} else {
		this.writeFalse()
	}
}
func (this *JsonStream) WriteInt8(val int8) {
	utils.WriteInt8ToString(&this.buf, val)
}
func (this *JsonStream) WriteInt16(val int16) {
	utils.WriteInt16ToString(&this.buf, val)
}
func (this *JsonStream) WriteInt32(val int32) {
	utils.WriteInt32ToString(&this.buf, val)
}
func (this *JsonStream) WriteInt64(val int64) {
	utils.WriteInt64ToString(&this.buf, val)
}
func (this *JsonStream) WriteUint8(val uint8) {
	utils.WriteUint8ToString(&this.buf, val)
}
func (this *JsonStream) WriteUint16(val uint16) {
	utils.WriteUint16ToString(&this.buf, val)
}
func (this *JsonStream) WriteUint32(val uint32) {
	utils.WriteUint32ToString(&this.buf, val)
}
func (this *JsonStream) WriteUint64(val uint64) {
	utils.WriteUint64ToString(&this.buf, val)
}
func (this *JsonStream) WriteFloat32(val float32) {
	utils.WriteFloat32ToString(&this.buf, val)
}
func (this *JsonStream) WriteFloat64(val float64) {
	utils.WriteFloat64ToString(&this.buf, val)
}
func (this *JsonStream) WriteString(val string) {
	valLen := len(val)
	this.buf = append(this.buf, '"')
	i := 0
	for ; i < valLen; i++ {
		c := val[i]
		if c > 31 && c != '"' && c != '\\' {
			this.buf = append(this.buf, c)
		} else {
			break
		}
	}
	if i == valLen {
		this.buf = append(this.buf, '"')
		return
	}
	utils.WriteStringSlowPath(&this.buf, i, val, valLen)
}
func (this *JsonStream) WriteBytes(val []byte) {
	this.buf = append(this.buf, val...)
}
func (this *JsonStream) Reset(bufSize int) {
	this.buf = make([]byte, 0, bufSize)
	this.err = nil
	this.indention = 0
	return
}
func (this *JsonStream) Buffer() []byte {
	return this.buf
}
func (this *JsonStream) Error() error {
	return this.err
}
func (this *JsonStream) SetErr(err error) {
	this.err = err
}

//-------------------------------------------------------------------------------------------------------------------------------------------------------------------

// WriteTrue write true to stream
func (stream *JsonStream) writeTrue() {
	stream.writeFourBytes('t', 'r', 'u', 'e')
}

// WriteFalse write false to stream
func (stream *JsonStream) writeFalse() {
	stream.writeFiveBytes('f', 'a', 'l', 's', 'e')
}

func (this *JsonStream) writeByte(c byte) {
	this.buf = append(this.buf, c)
}

func (this *JsonStream) writeTwoBytes(c1 byte, c2 byte) {
	this.buf = append(this.buf, c1, c2)
}

func (this *JsonStream) writeThreeBytes(c1 byte, c2 byte, c3 byte) {
	this.buf = append(this.buf, c1, c2, c3)
}

func (this *JsonStream) writeFourBytes(c1 byte, c2 byte, c3 byte, c4 byte) {
	this.buf = append(this.buf, c1, c2, c3, c4)
}

func (this *JsonStream) writeFiveBytes(c1 byte, c2 byte, c3 byte, c4 byte, c5 byte) {
	this.buf = append(this.buf, c1, c2, c3, c4, c5)
}

func (stream *JsonStream) writeIndention(delta int) {
	if stream.indention == 0 {
		return
	}
	stream.writeByte('\n')
	toWrite := stream.indention - delta
	for i := 0; i < toWrite; i++ {
		stream.buf = append(stream.buf, ' ')
	}
}
