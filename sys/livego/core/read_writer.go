package core

import (
	"bufio"
	"io"
)

func NewReadWriter(rw io.ReadWriter, bufSize int) *ReadWriter {
	return &ReadWriter{
		ReadWriter: bufio.NewReadWriter(bufio.NewReaderSize(rw, bufSize), bufio.NewWriterSize(rw, bufSize)),
	}
}

type ReadWriter struct {
	*bufio.ReadWriter
	readError  error
	writeError error
}

func (this *ReadWriter) Read(p []byte) (int, error) {
	if this.readError != nil {
		return 0, this.readError
	}
	n, err := io.ReadAtLeast(this.ReadWriter, p, len(p))
	this.readError = err
	return n, err
}

func (this *ReadWriter) ReadError() error {
	return this.readError
}

func (this *ReadWriter) ReadUintBE(n int) (uint32, error) {
	if this.readError != nil {
		return 0, this.readError
	}
	ret := uint32(0)
	for i := 0; i < n; i++ {
		b, err := this.ReadByte()
		if err != nil {
			this.readError = err
			return 0, err
		}
		ret = ret<<8 + uint32(b)
	}
	return ret, nil
}

func (this *ReadWriter) ReadUintLE(n int) (uint32, error) {
	if this.readError != nil {
		return 0, this.readError
	}
	ret := uint32(0)
	for i := 0; i < n; i++ {
		b, err := this.ReadByte()
		if err != nil {
			this.readError = err
			return 0, err
		}
		ret += uint32(b) << uint32(i*8)
	}
	return ret, nil
}

func (this *ReadWriter) Flush() error {
	if this.writeError != nil {
		return this.writeError
	}

	if this.ReadWriter.Writer.Buffered() == 0 {
		return nil
	}
	return this.ReadWriter.Flush()
}

func (this *ReadWriter) Write(p []byte) (int, error) {
	if this.writeError != nil {
		return 0, this.writeError
	}
	return this.ReadWriter.Write(p)
}

func (this *ReadWriter) WriteError() error {
	return this.writeError
}

func (this *ReadWriter) WriteUintBE(v uint32, n int) error {
	if this.writeError != nil {
		return this.writeError
	}
	for i := 0; i < n; i++ {
		b := byte(v>>uint32((n-i-1)<<3)) & 0xff
		if err := this.WriteByte(b); err != nil {
			this.writeError = err
			return err
		}
	}
	return nil
}

func (this *ReadWriter) WriteUintLE(v uint32, n int) error {
	if this.writeError != nil {
		return this.writeError
	}
	for i := 0; i < n; i++ {
		b := byte(v) & 0xff
		if err := this.WriteByte(b); err != nil {
			this.writeError = err
			return err
		}
		v = v >> 8
	}
	return nil
}
