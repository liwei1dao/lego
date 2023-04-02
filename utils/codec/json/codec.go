package json

import (
	"io"

	"github.com/liwei1dao/lego/utils/codec/codecore"
)

func NewDecoder(r io.Reader) *Decoder {
	iter := GetReader(make([]byte, 512), r)
	return &Decoder{iter}
}

type Decoder struct {
	iter codecore.IReader
}

func (this *Decoder) Decode(obj interface{}) error {
	this.iter.ReadVal(obj)
	err := this.iter.Error()
	if err == io.EOF {
		return nil
	}
	return this.iter.Error()
}

func (this *Decoder) UseNumber() {
	this.iter.Config().UseNumber = true
}

func (this *Decoder) DisallowUnknownFields() {
	this.iter.Config().DisallowUnknownFields = true
}

func NewEncoder(w io.Writer) *Encoder {
	stream := GetWriter(w)
	return &Encoder{stream}
}

type Encoder struct {
	stream codecore.IWriter
}

func (this *Encoder) Encode(val interface{}) error {
	this.stream.WriteVal(val)
	this.stream.WriteString("\n")
	this.stream.Flush()
	return this.stream.Error()
}
func (this *Encoder) SetEscapeHTML(escapeHTML bool) {
	this.stream.Config().EscapeHTML = true
}
