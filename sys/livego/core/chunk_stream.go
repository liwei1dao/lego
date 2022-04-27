package core

import (
	"encoding/binary"
	"fmt"

	"github.com/liwei1dao/lego/sys/livego/utils/pool"
)

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

func (this *ChunkStream) new(pool *pool.Pool) {
	this.got = false
	this.index = 0
	this.remain = this.Length
	this.Data = pool.Get(int(this.Length))
}

func (this *ChunkStream) full() bool {
	return this.got
}

func (this *ChunkStream) readChunk(r *ReadWriter, chunkSize uint32, pool *pool.Pool) error {
	if this.remain != 0 && this.tmpFromat != 3 {
		return fmt.Errorf("invalid remain = %d", this.remain)
	}
	switch this.CSID {
	case 0:
		id, _ := r.ReadUintLE(1)
		this.CSID = id + 64
	case 1:
		id, _ := r.ReadUintLE(2)
		this.CSID = id + 64
	}

	switch this.tmpFromat {
	case 0:
		this.Format = this.tmpFromat
		this.Timestamp, _ = r.ReadUintBE(3)
		this.Length, _ = r.ReadUintBE(3)
		this.TypeID, _ = r.ReadUintBE(1)
		this.StreamID, _ = r.ReadUintLE(4)
		if this.Timestamp == 0xffffff {
			this.Timestamp, _ = r.ReadUintBE(4)
			this.exted = true
		} else {
			this.exted = false
		}
		this.new(pool)
	case 1:
		this.Format = this.tmpFromat
		timeStamp, _ := r.ReadUintBE(3)
		this.Length, _ = r.ReadUintBE(3)
		this.TypeID, _ = r.ReadUintBE(1)
		if timeStamp == 0xffffff {
			timeStamp, _ = r.ReadUintBE(4)
			this.exted = true
		} else {
			this.exted = false
		}
		this.timeDelta = timeStamp
		this.Timestamp += timeStamp
		this.new(pool)
	case 2:
		this.Format = this.tmpFromat
		timeStamp, _ := r.ReadUintBE(3)
		if timeStamp == 0xffffff {
			timeStamp, _ = r.ReadUintBE(4)
			this.exted = true
		} else {
			this.exted = false
		}
		this.timeDelta = timeStamp
		this.Timestamp += timeStamp
		this.new(pool)
	case 3:
		if this.remain == 0 {
			switch this.Format {
			case 0:
				if this.exted {
					timestamp, _ := r.ReadUintBE(4)
					this.Timestamp = timestamp
				}
			case 1, 2:
				var timedet uint32
				if this.exted {
					timedet, _ = r.ReadUintBE(4)
				} else {
					timedet = this.timeDelta
				}
				this.Timestamp += timedet
			}
			this.new(pool)
		} else {
			if this.exted {
				b, err := r.Peek(4)
				if err != nil {
					return err
				}
				tmpts := binary.BigEndian.Uint32(b)
				if tmpts == this.Timestamp {
					r.Discard(4)
				}
			}
		}
	default:
		return fmt.Errorf("invalid format=%d", this.Format)
	}
	size := int(this.remain)
	if size > int(chunkSize) {
		size = int(chunkSize)
	}

	buf := this.Data[this.index : this.index+uint32(size)]
	if _, err := r.Read(buf); err != nil {
		return err
	}
	this.index += uint32(size)
	this.remain -= uint32(size)
	if this.remain == 0 {
		this.got = true
	}

	return r.readError
}

func (this *ChunkStream) writeHeader(w *ReadWriter) error {
	//Chunk Basic Header
	h := this.Format << 6
	switch {
	case this.CSID < 64:
		h |= this.CSID
		w.WriteUintBE(h, 1)
	case this.CSID-64 < 256:
		h |= 0
		w.WriteUintBE(h, 1)
		w.WriteUintLE(this.CSID-64, 1)
	case this.CSID-64 < 65536:
		h |= 1
		w.WriteUintBE(h, 1)
		w.WriteUintLE(this.CSID-64, 2)
	}
	//Chunk Message Header
	ts := this.Timestamp
	if this.Format == 3 {
		goto END
	}
	if this.Timestamp > 0xffffff {
		ts = 0xffffff
	}
	w.WriteUintBE(ts, 3)
	if this.Format == 2 {
		goto END
	}
	if this.Length > 0xffffff {
		return fmt.Errorf("length=%d", this.Length)
	}
	w.WriteUintBE(this.Length, 3)
	w.WriteUintBE(this.TypeID, 1)
	if this.Format == 1 {
		goto END
	}
	w.WriteUintLE(this.StreamID, 4)
END:
	//Extended Timestamp
	if ts >= 0xffffff {
		w.WriteUintBE(this.Timestamp, 4)
	}
	return w.WriteError()
}

func (this *ChunkStream) writeChunk(w *ReadWriter, chunkSize int) error {
	if this.TypeID == TAG_AUDIO {
		this.CSID = 4
	} else if this.TypeID == TAG_VIDEO ||
		this.TypeID == TAG_SCRIPTDATAAMF0 ||
		this.TypeID == TAG_SCRIPTDATAAMF3 {
		this.CSID = 6
	}
	totalLen := uint32(0)
	numChunks := (this.Length / uint32(chunkSize))
	for i := uint32(0); i <= numChunks; i++ {
		if totalLen == this.Length {
			break
		}
		if i == 0 {
			this.Format = uint32(0)
		} else {
			this.Format = uint32(3)
		}
		if err := this.writeHeader(w); err != nil {
			return err
		}
		inc := uint32(chunkSize)
		start := uint32(i) * uint32(chunkSize)
		if uint32(len(this.Data))-start <= inc {
			inc = uint32(len(this.Data)) - start
		}
		totalLen += inc
		end := start + inc
		buf := this.Data[start:end]
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}
	return nil
}
