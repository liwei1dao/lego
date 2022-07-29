package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/gogo/protobuf/proto"
	"github.com/liwei1dao/lego/core"
	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/liwei1dao/lego/utils/codec"
	"github.com/smallnest/rpcx/util"
	"github.com/valyala/bytebufferpool"
)

var bufferPool = util.NewLimitedPool(512, 4096)

var Compressors = map[lcore.CompressType]Compressor{
	lcore.CompressNone: &RawDataCompressor{},
	lcore.CompressGzip: &GzipCompressor{},
}

var (
	zeroHeaderArray Header
	zeroHeader      = zeroHeaderArray[1:]
)

func NewMessage() *Message {
	header := Header([12]byte{})
	header[0] = magicNumber

	return &Message{
		Header: &header,
	}
}

type Message struct {
	*Header
	serviceMethod string
	from          *core.ServiceNode
	metadata      map[string]string
	payload       []byte
	data          []byte
}

func (this *Message) ServiceMethod() string {
	return this.serviceMethod
}
func (this *Message) SetServiceMethod(v string) {
	this.serviceMethod = v
}

func (this *Message) From() *core.ServiceNode {
	return this.from
}
func (this *Message) SetFrom(v *core.ServiceNode) {
	*this.from = *v
}
func (this *Message) Payload() []byte {
	return this.payload
}
func (this *Message) SetPayload(m []byte) {
	this.payload = m
}
func (this *Message) Metadata() map[string]string {
	return this.metadata
}
func (this *Message) SetMetadata(m map[string]string) {
	this.metadata = m
}
func Read(r io.Reader) (*Message, error) {
	msg := NewMessage()
	err := msg.Decode(r)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// Decode decodes a message from reader.
func (this *Message) Decode(r io.Reader) error {
	// validate rest length for each step?

	// parse header
	_, err := io.ReadFull(r, this.Header[:1])
	if err != nil {
		return err
	}
	if !this.Header.CheckMagicNumber() {
		return fmt.Errorf("wrong magic number: %v", this.Header[0])
	}

	_, err = io.ReadFull(r, this.Header[1:])
	if err != nil {
		return err
	}

	// total
	lenData := make([]byte, 4)
	_, err = io.ReadFull(r, lenData)
	if err != nil {
		return err
	}
	l := binary.BigEndian.Uint32(lenData)

	totalL := int(l)
	if cap(this.data) >= totalL { // reuse data
		this.data = this.data[:totalL]
	} else {
		this.data = make([]byte, totalL)
	}
	data := this.data
	_, err = io.ReadFull(r, data)
	if err != nil {
		return err
	}

	n := 0

	// parse serviceMethod
	l = binary.BigEndian.Uint32(data[n : n+4])
	n = n + 4
	nEnd := n + int(l)
	this.serviceMethod = util.SliceByteToString(data[n:nEnd])
	n = nEnd

	// parse serviceMethod
	l = binary.BigEndian.Uint32(data[n : n+4])
	n = n + 4
	nEnd = n + int(l)
	err = proto.Unmarshal(data[n:nEnd], this.from)
	if err != nil {
		return err
	}
	n = nEnd

	// parse meta
	l = binary.BigEndian.Uint32(data[n : n+4])
	n = n + 4
	nEnd = n + int(l)

	if l > 0 {
		this.metadata, err = decodeMetadata(l, data[n:nEnd])
		if err != nil {
			return err
		}
	}
	n = nEnd

	// parse payload
	l = binary.BigEndian.Uint32(data[n : n+4])
	_ = l
	n = n + 4
	this.payload = data[n:]

	if this.CompressType() != lcore.CompressNone {
		compressor := Compressors[this.CompressType()]
		if compressor == nil {
			return lcore.ErrUnsupportedCompressor
		}
		this.payload, err = compressor.Unzip(this.payload)
		if err != nil {
			return err
		}
	}

	return err
}
func (this *Message) Reset() {
	resetHeader(this.Header)
	this.metadata = nil
	this.payload = []byte{}
	this.serviceMethod = ""
}
func (this Message) Clone() lcore.IMessage {
	header := *this.Header
	c := GetPooledMsg()
	header.SetCompressType(lcore.CompressNone)
	c.Header = &header
	c.serviceMethod = this.serviceMethod
	return c
}

func (m Message) EncodeSlicePointer() *[]byte {
	bb := bytebufferpool.Get()
	encodeMetadata(m.metadata, bb)
	fdata, _ := proto.Marshal(m.from)
	meta := bb.Bytes()
	smL := len(m.serviceMethod)
	fml := len(fdata)
	var err error
	payload := m.payload
	if m.CompressType() != lcore.CompressNone {
		compressor := Compressors[m.CompressType()]
		if compressor == nil {
			m.SetCompressType(lcore.CompressNone)
		} else {
			payload, err = compressor.Zip(m.payload)
			if err != nil {
				m.SetCompressType(lcore.CompressNone)
				payload = m.payload
			}
		}
	}

	totalL := (4 + smL) + (4 + fml) + (4 + len(meta)) + (4 + len(payload))

	// header + dataLen + spLen + sp + smLen + sm + metaL + meta + payloadLen + payload
	metaStart := 12 + 4 + (4 + smL) + (4 + fml)

	payLoadStart := metaStart + (4 + len(meta))
	l := 12 + 4 + totalL

	data := bufferPool.Get(l)
	copy(*data, m.Header[:])

	// totalLen
	binary.BigEndian.PutUint32((*data)[12:16], uint32(totalL))

	binary.BigEndian.PutUint32((*data)[16:20], uint32(smL))
	copy((*data)[20:20+smL], util.StringToSliceByte(m.serviceMethod))

	binary.BigEndian.PutUint32((*data)[20+smL:24+smL], uint32(fml))
	copy((*data)[24+smL:metaStart], fdata)

	binary.BigEndian.PutUint32((*data)[metaStart:metaStart+4], uint32(len(meta)))
	copy((*data)[metaStart+4:], meta)

	bytebufferpool.Put(bb)

	binary.BigEndian.PutUint32((*data)[payLoadStart:payLoadStart+4], uint32(len(payload)))
	copy((*data)[payLoadStart+4:], payload)

	return data
}

// WriteTo writes message to writers.
func (m Message) WriteTo(w io.Writer) (int64, error) {
	fdata, _ := proto.Marshal(m.from)
	nn, err := w.Write(m.Header[:])
	n := int64(nn)
	if err != nil {
		return n, err
	}

	bb := bytebufferpool.Get()
	encodeMetadata(m.metadata, bb)
	meta := bb.Bytes()

	smL := len(m.serviceMethod)
	fml := len(fdata)

	payload := m.payload
	if m.CompressType() != lcore.CompressNone {
		compressor := Compressors[m.CompressType()]
		if compressor == nil {
			return n, lcore.ErrUnsupportedCompressor
		}
		payload, err = compressor.Zip(m.payload)
		if err != nil {
			return n, err
		}
	}

	totalL := (4 + smL) + (4 + fml) + (4 + len(meta)) + (4 + len(payload))
	err = binary.Write(w, binary.BigEndian, uint32(totalL))
	if err != nil {
		return n, err
	}

	// write servicePath and serviceMethod
	err = binary.Write(w, binary.BigEndian, uint32(len(m.serviceMethod)))
	if err != nil {
		return n, err
	}
	_, err = w.Write(util.StringToSliceByte(m.serviceMethod))
	if err != nil {
		return n, err
	}
	// write servicePath and serviceMethod
	err = binary.Write(w, binary.BigEndian, uint32(fml))
	if err != nil {
		return n, err
	}
	_, err = w.Write(fdata)
	if err != nil {
		return n, err
	}
	// write meta
	err = binary.Write(w, binary.BigEndian, uint32(len(meta)))
	if err != nil {
		return n, err
	}
	_, err = w.Write(meta)
	if err != nil {
		return n, err
	}

	bytebufferpool.Put(bb)

	// write payload
	err = binary.Write(w, binary.BigEndian, uint32(len(payload)))
	if err != nil {
		return n, err
	}

	nn, err = w.Write(payload)
	return int64(nn), err
}

func PutData(data *[]byte) {
	bufferPool.Put(data)
}

func encodeMetadata(m map[string]string, bb *bytebufferpool.ByteBuffer) {
	if len(m) == 0 {
		return
	}
	d := poolUint32Data.Get().(*[]byte)
	for k, v := range m {
		binary.BigEndian.PutUint32(*d, uint32(len(k)))
		bb.Write(*d)
		bb.Write(codec.StringToBytes(k))
		binary.BigEndian.PutUint32(*d, uint32(len(v)))
		bb.Write(*d)
		bb.Write(codec.StringToBytes(v))
	}
}
func decodeMetadata(l uint32, data []byte) (map[string]string, error) {
	m := make(map[string]string, 10)
	n := uint32(0)
	for n < l {
		// parse one key and value
		// key
		sl := binary.BigEndian.Uint32(data[n : n+4])
		n = n + 4
		if n+sl > l-4 {
			return m, lcore.ErrMetaKVMissing
		}
		k := string(data[n : n+sl])
		n = n + sl

		// value
		sl = binary.BigEndian.Uint32(data[n : n+4])
		n = n + 4
		if n+sl > l {
			return m, lcore.ErrMetaKVMissing
		}
		v := string(data[n : n+sl])
		n = n + sl
		m[k] = v
	}

	return m, nil
}
func resetHeader(h *Header) {
	copy(h[1:], zeroHeader)
}
