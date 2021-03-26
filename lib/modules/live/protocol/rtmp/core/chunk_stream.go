package core

import (
	"fmt"

	"github.com/liwei1dao/lego/lib/modules/live/av"
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

func (chunkStream *ChunkStream) writeHeader(w *ReadWriter) error {
	//Chunk Basic Header
	h := chunkStream.Format << 6
	switch {
	case chunkStream.CSID < 64:
		h |= chunkStream.CSID
		w.WriteUintBE(h, 1)
	case chunkStream.CSID-64 < 256:
		h |= 0
		w.WriteUintBE(h, 1)
		w.WriteUintLE(chunkStream.CSID-64, 1)
	case chunkStream.CSID-64 < 65536:
		h |= 1
		w.WriteUintBE(h, 1)
		w.WriteUintLE(chunkStream.CSID-64, 2)
	}
	//Chunk Message Header
	ts := chunkStream.Timestamp
	if chunkStream.Format == 3 {
		goto END
	}
	if chunkStream.Timestamp > 0xffffff {
		ts = 0xffffff
	}
	w.WriteUintBE(ts, 3)
	if chunkStream.Format == 2 {
		goto END
	}
	if chunkStream.Length > 0xffffff {
		return fmt.Errorf("length=%d", chunkStream.Length)
	}
	w.WriteUintBE(chunkStream.Length, 3)
	w.WriteUintBE(chunkStream.TypeID, 1)
	if chunkStream.Format == 1 {
		goto END
	}
	w.WriteUintLE(chunkStream.StreamID, 4)
END:
	//Extended Timestamp
	if ts >= 0xffffff {
		w.WriteUintBE(chunkStream.Timestamp, 4)
	}
	return w.WriteError()
}

func (chunkStream *ChunkStream) writeChunk(w *ReadWriter, chunkSize int) error {
	if chunkStream.TypeID == av.TAG_AUDIO {
		chunkStream.CSID = 4
	} else if chunkStream.TypeID == av.TAG_VIDEO ||
		chunkStream.TypeID == av.TAG_SCRIPTDATAAMF0 ||
		chunkStream.TypeID == av.TAG_SCRIPTDATAAMF3 {
		chunkStream.CSID = 6
	}

	totalLen := uint32(0)
	numChunks := (chunkStream.Length / uint32(chunkSize))
	for i := uint32(0); i <= numChunks; i++ {
		if totalLen == chunkStream.Length {
			break
		}
		if i == 0 {
			chunkStream.Format = uint32(0)
		} else {
			chunkStream.Format = uint32(3)
		}
		if err := chunkStream.writeHeader(w); err != nil {
			return err
		}
		inc := uint32(chunkSize)
		start := uint32(i) * uint32(chunkSize)
		if uint32(len(chunkStream.Data))-start <= inc {
			inc = uint32(len(chunkStream.Data)) - start
		}
		totalLen += inc
		end := start + inc
		buf := chunkStream.Data[start:end]
		if _, err := w.Write(buf); err != nil {
			return err
		}
	}

	return nil
}
