package stream

import (
	"github.com/liwei1dao/lego/sys/codec/core"
	"github.com/modern-go/reflect2"
)

func NewStream(codec core.ICodec, bufSize int) *Stream {
	return &Stream{
		codec:     codec,
		buf:       make([]byte, 0, bufSize),
		err:       nil,
		indention: 0,
	}
}

type Stream struct {
	codec     core.ICodec
	err       error
	buf       []byte
	indention int
}

func (stream *Stream) WriteVal(val interface{}, opt *core.ExecuteOptions) {
	if nil == val {
		stream.WriteNil()
		return
	}
	cacheKey := reflect2.RTypeOf(val)
	encoder := stream.codec.GetEncoderFromCache(cacheKey)
	if encoder == nil {
		typ := reflect2.TypeOf(val)
		encoder = stream.codec.EncoderOf(typ)
	}
	encoder.Encode(reflect2.PtrOf(val), stream, opt)
}

func (stream *Stream) WriteNil() {
	stream.writeFourBytes('n', 'u', 'l', 'l')
}

func (stream *Stream) WriteObjectStart() {
	stream.indention += stream.codec.Options().IndentionStep
	stream.WriteChar('{')
	stream.writeIndention(0)
}

func (stream *Stream) WriteObjectField(field string) {
	stream.WriteString(field)
	if stream.indention > 0 {
		stream.writeTwoBytes(':', ' ')
	} else {
		stream.WriteChar(':')
	}
}

func (stream *Stream) WriteObjectEnd() {
	stream.writeIndention(stream.codec.Options().IndentionStep)
	stream.indention -= stream.codec.Options().IndentionStep
	stream.WriteChar('}')
}
func (stream *Stream) WriteEmptyObject() {
	stream.WriteChar('{')
	stream.WriteChar('}')
}
func (stream *Stream) WriteMore() {
	stream.WriteChar(',')
	stream.writeIndention(0)
}
func (stream *Stream) WriteArrayStart() {
	stream.indention += stream.codec.Options().IndentionStep
	stream.WriteChar('[')
	stream.writeIndention(0)
}
func (stream *Stream) WriteArrayEnd() {
	stream.writeIndention(stream.codec.Options().IndentionStep)
	stream.indention -= stream.codec.Options().IndentionStep
	stream.WriteChar(']')
}
func (stream *Stream) WriteEmptyArray() {
	stream.writeTwoBytes('[', ']')
}
func (stream *Stream) WriteChar(c byte) {
	stream.buf = append(stream.buf, c)
}

func (stream *Stream) WriteBytes(d []byte) {
	stream.buf = append(stream.buf, d...)
}
func (stream *Stream) WriteBool(val bool) {
	if val {
		stream.WriteTrue()
	} else {
		stream.WriteFalse()
	}
}
func (stream *Stream) WriteRaw(s string) {
	stream.buf = append(stream.buf, s...)
}
func (stream *Stream) Indention() int {
	return stream.indention
}
func (stream *Stream) SetIndention(v int) {
	stream.indention = v
}
func (stream *Stream) Buffered() int {
	return len(stream.buf)
}
func (stream *Stream) ToBuffer() []byte {
	return stream.buf
}
func (stream *Stream) Error() error {
	return stream.err
}
func (stream *Stream) SetError(err error) {
	stream.err = err
}

func (stream *Stream) WriteTrue() {
	stream.writeFourBytes('t', 'r', 'u', 'e')
}

// WriteFalse write false to stream
func (stream *Stream) WriteFalse() {
	stream.writeFiveBytes('f', 'a', 'l', 's', 'e')
}

func (stream *Stream) writeTwoBytes(c1 byte, c2 byte) {
	stream.buf = append(stream.buf, c1, c2)
}

func (stream *Stream) writeThreeBytes(c1 byte, c2 byte, c3 byte) {
	stream.buf = append(stream.buf, c1, c2, c3)
}

func (stream *Stream) writeFourBytes(c1 byte, c2 byte, c3 byte, c4 byte) {
	stream.buf = append(stream.buf, c1, c2, c3, c4)
}

func (stream *Stream) writeFiveBytes(c1 byte, c2 byte, c3 byte, c4 byte, c5 byte) {
	stream.buf = append(stream.buf, c1, c2, c3, c4, c5)
}
func (stream *Stream) writeIndention(delta int) {
	if stream.indention == 0 {
		return
	}
	stream.WriteChar('\n')
	toWrite := stream.indention - delta
	for i := 0; i < toWrite; i++ {
		stream.buf = append(stream.buf, ' ')
	}
}
