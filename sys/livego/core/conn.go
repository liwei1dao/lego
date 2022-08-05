package core

import (
	"encoding/binary"
	"net"
	"time"

	"github.com/liwei1dao/lego/sys/livego/utils/pio"
	"github.com/liwei1dao/lego/sys/livego/utils/pool"
	"github.com/liwei1dao/lego/sys/log"
)

const (
	_ = iota
	idSetChunkSize
	idAbortMessage
	idAck
	idUserControlMessages
	idWindowAckSize
	idSetPeerBandwidth
)

const (
	streamBegin      uint32 = 0
	streamEOF        uint32 = 1
	streamDry        uint32 = 2
	setBufferLen     uint32 = 3
	streamIsRecorded uint32 = 4
	pingRequest      uint32 = 6
	pingResponse     uint32 = 7
)

func NewConn(c net.Conn, sys ISys, log log.ILogger) *Conn {
	return &Conn{
		sys:             sys,
		log:             log,
		Conn:            c,
		timeout:         time.Second * time.Duration(sys.GetTimeout()),
		chunkSize:       128,
		remoteChunkSize: 128,
		pool:            pool.NewPool(),
		rw:              NewReadWriter(c, sys.GetConnBuffSzie()),
	}
}

type Conn struct {
	net.Conn
	sys                 ISys
	log                 log.ILogger
	rw                  *ReadWriter
	timeout             time.Duration
	chunkSize           uint32
	remoteChunkSize     uint32
	remoteWindowAckSize uint32
	received            uint32
	ackReceived         uint32
	pool                *pool.Pool
	chunks              map[uint32]ChunkStream
}

func (this *Conn) Sys() ISys {
	return this.sys
}
func (this *Conn) Log() log.ILogger {
	return this.log
}

func (conn *Conn) NewAck(size uint32) ChunkStream {
	return initControlMsg(idAck, 4, size)
}

func (this *Conn) Read(c *ChunkStream) error {
	for {
		h, _ := this.rw.ReadUintBE(1)
		format := h >> 6
		csid := h & 0x3f
		cs, ok := this.chunks[csid]
		if !ok {
			cs = ChunkStream{}
			this.chunks[csid] = cs
		}
		cs.tmpFromat = format
		cs.CSID = csid
		err := cs.readChunk(this.rw, this.remoteChunkSize, this.pool)
		if err != nil {
			return err
		}
		this.chunks[csid] = cs
		if cs.full() {
			*c = cs
			break
		}
	}
	this.handleControlMsg(c)
	this.ack(c.Length)
	return nil
}

func (this *Conn) Write(c *ChunkStream) error {
	if c.TypeID == idSetChunkSize {
		this.chunkSize = binary.BigEndian.Uint32(c.Data)
	}
	return c.writeChunk(this.rw, int(this.chunkSize))
}

func (this *Conn) Flush() error {
	return this.rw.Flush()
}

func (this *Conn) ack(size uint32) {
	this.received += uint32(size)
	this.ackReceived += uint32(size)
	if this.received >= 0xf0000000 {
		this.received = 0
	}
	if this.ackReceived >= this.remoteWindowAckSize {
		cs := this.NewAck(this.ackReceived)
		cs.writeChunk(this.rw, int(this.chunkSize))
		this.ackReceived = 0
	}
}

func (this *Conn) handleControlMsg(c *ChunkStream) {
	if c.TypeID == idSetChunkSize {
		this.remoteChunkSize = binary.BigEndian.Uint32(c.Data)
	} else if c.TypeID == idWindowAckSize {
		this.remoteWindowAckSize = binary.BigEndian.Uint32(c.Data)
	}
}

func (this *Conn) NewWindowAckSize(size uint32) ChunkStream {
	return initControlMsg(idWindowAckSize, 4, size)
}

func (this *Conn) NewSetPeerBandwidth(size uint32) ChunkStream {
	ret := initControlMsg(idSetPeerBandwidth, 5, size)
	ret.Data[4] = 2
	return ret
}

func (this *Conn) NewSetChunkSize(size uint32) ChunkStream {
	return initControlMsg(idSetChunkSize, 4, size)
}

func (this *Conn) SetRecorded() {
	ret := this.userControlMsg(streamIsRecorded, 4)
	for i := 0; i < 4; i++ {
		ret.Data[2+i] = byte(1 >> uint32((3-i)*8) & 0xff)
	}
	this.Write(&ret)
}

func (this *Conn) SetBegin() {
	ret := this.userControlMsg(streamBegin, 4)
	for i := 0; i < 4; i++ {
		ret.Data[2+i] = byte(1 >> uint32((3-i)*8) & 0xff)
	}
	this.Write(&ret)
}

func (this *Conn) userControlMsg(eventType, buflen uint32) ChunkStream {
	var ret ChunkStream
	buflen += 2
	ret = ChunkStream{
		Format:   0,
		CSID:     2,
		TypeID:   4,
		StreamID: 1,
		Length:   buflen,
		Data:     make([]byte, buflen),
	}
	ret.Data[0] = byte(eventType >> 8 & 0xff)
	ret.Data[1] = byte(eventType & 0xff)
	return ret
}

func initControlMsg(id, size, value uint32) ChunkStream {
	ret := ChunkStream{
		Format:   0,
		CSID:     2,
		TypeID:   id,
		StreamID: 0,
		Length:   size,
		Data:     make([]byte, size),
	}
	pio.PutU32BE(ret.Data[:size], value)
	return ret
}
