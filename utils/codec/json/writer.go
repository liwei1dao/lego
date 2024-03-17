package json

import (
	"io"

	"github.com/liwei1dao/lego/utils/codec"
	"github.com/liwei1dao/lego/utils/codec/codecore"
	"github.com/liwei1dao/lego/utils/codec/utils"

	"github.com/modern-go/reflect2"
)

type JsonWriter struct {
	config    *codecore.Config
	out       io.Writer
	buf       []byte
	indention int
	err       error
}

func (this *JsonWriter) Config() *codecore.Config {
	return this.config
}
func (this *JsonWriter) GetReader(buf []byte, r io.Reader) codecore.IReader {
	return GetReader(buf, r)
}
func (this *JsonWriter) PutReader(r codecore.IReader) {
	PutReader(r)
}
func (this *JsonWriter) GetWriter() codecore.IWriter {
	return GetWriter(nil)
}
func (this *JsonWriter) PutWriter(w codecore.IWriter) {
	PutWriter(w)
}

//写入对象
func (this *JsonWriter) WriteVal(val interface{}) {
	if nil == val {
		this.WriteNil()
		return
	}
	cacheKey := reflect2.RTypeOf(val)
	encoder := codec.GetEncoder(cacheKey)
	if encoder == nil {
		typ := reflect2.TypeOf(val)
		encoder = codec.EncoderOf(typ, defconf)
	}
	encoder.Encode(reflect2.PtrOf(val), this)
}

func (this *JsonWriter) WriteNil() {
	this.writeFourBytes('n', 'u', 'l', 'l')
}
func (this *JsonWriter) WriteEmptyArray() {
	this.writeTwoBytes('[', ']')
}
func (this *JsonWriter) WriteArrayStart() {
	this.indention += this.config.IndentionStep
	this.writeByte('[')
	this.writeIndention(0)
}
func (this *JsonWriter) WriteArrayEnd() {
	this.writeIndention(this.config.IndentionStep)
	this.indention -= this.config.IndentionStep
	this.writeByte(']')
}
func (this *JsonWriter) WriteEmptyObject() {
	this.writeTwoBytes('{', '}')
}
func (this *JsonWriter) WriteObjectStart() {
	this.indention += this.config.IndentionStep
	this.writeByte('{')
	this.writeIndention(0)
}
func (this *JsonWriter) WriteObjectEnd() {
	this.writeIndention(this.config.IndentionStep)
	this.indention -= this.config.IndentionStep
	this.writeByte('}')
}
func (this *JsonWriter) WriteMemberSplit() {
	this.writeByte(',')
	this.writeIndention(0)
}
func (this *JsonWriter) WriteKVSplit() {
	this.writeByte(':')
}
func (this *JsonWriter) WriteKeyStart() {
	this.writeByte('"')
}
func (this *JsonWriter) WriteKeyEnd() {
	this.writeByte('"')
}
func (this *JsonWriter) WriteObjectFieldName(val string) {
	this.WriteString(val)
	if this.indention > 0 {
		this.writeTwoBytes(':', ' ')
	} else {
		this.writeByte(':')
	}
}
func (this *JsonWriter) WriteBool(val bool) {
	if val {
		this.writeTrue()
	} else {
		this.writeFalse()
	}
}
func (this *JsonWriter) WriteInt8(val int8) {
	utils.WriteInt8ToString(&this.buf, val)
}
func (this *JsonWriter) WriteInt16(val int16) {
	utils.WriteInt16ToString(&this.buf, val)
}
func (this *JsonWriter) WriteInt32(val int32) {
	utils.WriteInt32ToString(&this.buf, val)
}
func (this *JsonWriter) WriteInt64(val int64) {
	utils.WriteInt64ToString(&this.buf, val)
}
func (this *JsonWriter) WriteUint8(val uint8) {
	utils.WriteUint8ToString(&this.buf, val)
}
func (this *JsonWriter) WriteUint16(val uint16) {
	utils.WriteUint16ToString(&this.buf, val)
}
func (this *JsonWriter) WriteUint32(val uint32) {
	utils.WriteUint32ToString(&this.buf, val)
}
func (this *JsonWriter) WriteUint64(val uint64) {
	utils.WriteUint64ToString(&this.buf, val)
}
func (this *JsonWriter) WriteFloat32(val float32) {
	utils.WriteFloat32ToString(&this.buf, val)
}
func (this *JsonWriter) WriteFloat64(val float64) {
	utils.WriteFloat64ToString(&this.buf, val)
}
func (this *JsonWriter) WriteString(val string) {
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

func (this *JsonWriter) Write(p []byte) (nn int, err error) {
	this.buf = append(this.buf, p...)
	if this.out != nil {
		nn, err = this.out.Write(this.buf)
		this.buf = this.buf[nn:]
		return
	}
	return len(p), nil
}

func (this *JsonWriter) Flush() error {
	if this.out == nil {
		return nil
	}
	if this.err != nil {
		return this.err
	}
	_, err := this.out.Write(this.buf)
	if err != nil {
		if this.err == nil {
			this.err = err
		}
		return err
	}
	this.buf = this.buf[:0]
	return nil
}

func (this *JsonWriter) WriteBytes(val []byte) {

	this.buf = append(this.buf, val...)
}
func (this *JsonWriter) Reset(w io.Writer) {
	this.buf = []byte{}
	this.err = nil
	this.indention = 0
	return
}
func (this *JsonWriter) Buffer() []byte {
	return this.buf
}
func (this *JsonWriter) Buffered() int {
	return len(this.buf)
}

func (this *JsonWriter) Error() error {
	return this.err
}
func (this *JsonWriter) SetErr(err error) {
	this.err = err
}

//-------------------------------------------------------------------------------------------------------------------------------------------------------------------

// WriteTrue write true to stream
func (stream *JsonWriter) writeTrue() {
	stream.writeFourBytes('t', 'r', 'u', 'e')
}

// WriteFalse write false to stream
func (stream *JsonWriter) writeFalse() {
	stream.writeFiveBytes('f', 'a', 'l', 's', 'e')
}

func (this *JsonWriter) writeByte(c byte) {
	this.buf = append(this.buf, c)
}

func (this *JsonWriter) writeTwoBytes(c1 byte, c2 byte) {
	this.buf = append(this.buf, c1, c2)
}

func (this *JsonWriter) writeThreeBytes(c1 byte, c2 byte, c3 byte) {
	this.buf = append(this.buf, c1, c2, c3)
}

func (this *JsonWriter) writeFourBytes(c1 byte, c2 byte, c3 byte, c4 byte) {
	this.buf = append(this.buf, c1, c2, c3, c4)
}

func (this *JsonWriter) writeFiveBytes(c1 byte, c2 byte, c3 byte, c4 byte, c5 byte) {
	this.buf = append(this.buf, c1, c2, c3, c4, c5)
}

func (stream *JsonWriter) writeIndention(delta int) {
	if stream.indention == 0 {
		return
	}
	stream.writeByte('\n')
	toWrite := stream.indention - delta
	for i := 0; i < toWrite; i++ {
		stream.buf = append(stream.buf, ' ')
	}
}
