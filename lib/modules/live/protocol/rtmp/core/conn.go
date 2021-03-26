package core

import (
	"net"

	"github.com/liwei1dao/lego/utils/pool"
)

type Conn struct {
	net.Conn
	chunkSize           uint32
	remoteChunkSize     uint32
	windowAckSize       uint32
	remoteWindowAckSize uint32
	received            uint32
	ackReceived         uint32
	rw                  *ReadWriter
	pool                *pool.Pool
	chunks              map[uint32]ChunkStream
}
