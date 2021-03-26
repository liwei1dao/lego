package core

type ChunkStream struct {
	Format    uint32
	CSID      uint32
	Timestamp uint32
	Length    uint32
	TypeID    uint32
	StreamID  uint32
	timeDelta uint32
	exted     bool
	index     uint32
	remain    uint32
	got       bool
	tmpFromat uint32
	Data      []byte
}
