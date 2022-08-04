package pools

import (
	"io"
	"strconv"
	"sync"
	"time"
)

/*
	bytes buffer 缓存池
*/

const _size = 1024

var bufferPool = &BufferPool{p: &sync.Pool{
	New: func() interface{} {
		return &Buffer{bs: make([]byte, 0, _size)}
	},
}}

func BufferPoolGet() *Buffer {
	return bufferPool.get()
}
func BufferPoolPut(buf *Buffer) {
	bufferPool.put(buf)
}

type BufferPool struct {
	p *sync.Pool
}

func (p BufferPool) get() *Buffer {
	buf := p.p.Get().(*Buffer)
	buf.reset()
	buf.pool = p
	return buf
}

func (p BufferPool) put(buf *Buffer) {
	p.p.Put(buf)
}

type Buffer struct {
	bs   []byte
	pool BufferPool
}

func (b *Buffer) reset() {
	b.bs = b.bs[:0]
}
func (b *Buffer) Len() int {
	return len(b.bs)
}
func (b *Buffer) Cap() int {
	return cap(b.bs)
}
func (b *Buffer) Bytes() []byte {
	return b.bs
}
func (b *Buffer) String() string {
	return string(b.bs)
}
func (b *Buffer) ReadFrom(r io.Reader) (int64, error) {
	p := b.bs
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
	if nMax == 0 {
		nMax = 64
		p = make([]byte, nMax)
	} else {
		p = p[:nMax]
	}
	for {
		if n == nMax {
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		}
		nn, err := r.Read(p[n:])
		n += int64(nn)
		if err != nil {
			b.bs = p[:n]
			n -= nStart
			if err == io.EOF {
				return n, nil
			}
			return n, err
		}
	}
}
func (b *Buffer) AppendByte(v byte) {
	b.bs = append(b.bs, v)
}
func (b *Buffer) AppendInt(i int64) {
	b.bs = strconv.AppendInt(b.bs, i, 10)
}
func (b *Buffer) AppendUint(i uint64) {
	b.bs = strconv.AppendUint(b.bs, i, 10)
}
func (b *Buffer) AppendFloat(f float64, bitSize int) {
	b.bs = strconv.AppendFloat(b.bs, f, 'f', -1, bitSize)
}
func (b *Buffer) AppendBool(v bool) {
	b.bs = strconv.AppendBool(b.bs, v)
}
func (b *Buffer) AppendTime(t time.Time, layout string) {
	b.bs = t.AppendFormat(b.bs, layout)
}
func (b *Buffer) AppendString(s string) {
	b.bs = append(b.bs, s...)
}
func (b *Buffer) AppendBytes(s []byte) {
	b.bs = append(b.bs, s...)
}
func (b *Buffer) Write(p []byte) (int, error) {
	b.bs = append(b.bs, p...)
	return len(p), nil
}
func (b *Buffer) WriteString(s string) (int, error) {
	b.AppendString(s)
	return len(s), nil
}
func (b *Buffer) WriteByte(v byte) error {
	b.AppendByte(v)
	return nil
}
func (b *Buffer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.bs)
	return int64(n), err
}
func (b *Buffer) Free() {
	b.pool.put(b)
}
