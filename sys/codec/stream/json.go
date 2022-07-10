package stream

import (
	"github.com/liwei1dao/lego/sys/codec/core"
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
