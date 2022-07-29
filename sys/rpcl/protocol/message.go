package protocol

import (
	"encoding/binary"

	"github.com/liwei1dao/lego/utils/codec"
	"github.com/smallnest/rpcx/util"
	"github.com/valyala/bytebufferpool"
)

var bufferPool = util.NewLimitedPool(512, 4096)

var Compressors = map[CompressType]Compressor{
	CompressNone: &RawDataCompressor{},
	CompressGzip: &GzipCompressor{},
}

var (
	zeroHeaderArray Header
	zeroHeader      = zeroHeaderArray[1:]
)

type Message struct {
	*Header
	ServicePath   string
	ServiceMethod string
	Metadata      map[string]string
	Payload       []byte
}

func (this *Message) Reset() {
	resetHeader(this.Header)
	this.Metadata = nil
	this.Payload = []byte{}
	this.ServiceMethod = ""
}
func (this Message) Clone() *Message {
	header := *this.Header
	c := GetPooledMsg()
	header.SetCompressType(CompressNone)
	c.Header = &header
	c.ServiceMethod = this.ServiceMethod
	return c
}

func (m Message) EncodeSlicePointer() *[]byte {
	bb := bytebufferpool.Get()
	encodeMetadata(m.Metadata, bb)
	meta := bb.Bytes()

	spL := len(m.ServicePath)
	smL := len(m.ServiceMethod)

	var err error
	payload := m.Payload
	if m.CompressType() != CompressNone {
		compressor := Compressors[m.CompressType()]
		if compressor == nil {
			m.SetCompressType(CompressNone)
		} else {
			payload, err = compressor.Zip(m.Payload)
			if err != nil {
				m.SetCompressType(CompressNone)
				payload = m.Payload
			}
		}
	}

	totalL := (4 + spL) + (4 + smL) + (4 + len(meta)) + (4 + len(payload))

	// header + dataLen + spLen + sp + smLen + sm + metaL + meta + payloadLen + payload
	metaStart := 12 + 4 + (4 + spL) + (4 + smL)

	payLoadStart := metaStart + (4 + len(meta))
	l := 12 + 4 + totalL

	data := bufferPool.Get(l)
	copy(*data, m.Header[:])

	// totalLen
	binary.BigEndian.PutUint32((*data)[12:16], uint32(totalL))

	binary.BigEndian.PutUint32((*data)[16:20], uint32(spL))
	copy((*data)[20:20+spL], util.StringToSliceByte(m.ServicePath))

	binary.BigEndian.PutUint32((*data)[20+spL:24+spL], uint32(smL))
	copy((*data)[24+spL:metaStart], util.StringToSliceByte(m.ServiceMethod))

	binary.BigEndian.PutUint32((*data)[metaStart:metaStart+4], uint32(len(meta)))
	copy((*data)[metaStart+4:], meta)

	bytebufferpool.Put(bb)

	binary.BigEndian.PutUint32((*data)[payLoadStart:payLoadStart+4], uint32(len(payload)))
	copy((*data)[payLoadStart+4:], payload)

	return data
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

func resetHeader(h *Header) {
	copy(h[1:], zeroHeader)
}
