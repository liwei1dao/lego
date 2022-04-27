package livego

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/container/id"
)

const (
	maxQueueNum           = 1024
	SAVE_STATICS_INTERVAL = 5000
)

func NewWriter(conn core.StreamReadWriteCloser) *Writer {
	ret := &Writer{
		Uid:         id.NewXId(),
		conn:        conn,
		packetQueue: make(chan *core.Packet, maxQueueNum),
	}
	go ret.Check()
	go func() {
		err := ret.SendPacket()
		if err != nil {
			log.Warnf("%v", err)
		}
	}()
	return ret
}

type StaticsBW struct {
	StreamId               uint32
	VideoDatainBytes       uint64
	LastVideoDatainBytes   uint64
	VideoSpeedInBytesperMS uint64

	AudioDatainBytes       uint64
	LastAudioDatainBytes   uint64
	AudioSpeedInBytesperMS uint64

	LastTimestamp int64
}

type Writer struct {
	core.RWBaser
	Uid         string
	conn        core.StreamReadWriteCloser
	packetQueue chan *core.Packet
	WriteBWInfo StaticsBW
	closed      bool
}

func (v *Writer) Info() (ret core.Info) {
	ret.UID = v.Uid
	_, _, URL := v.conn.GetInfo()
	ret.URL = URL
	_url, err := url.Parse(URL)
	if err != nil {
		log.Warnf("[SYS LiveGo] err:%v", err)
	}
	ret.Key = strings.TrimLeft(_url.Path, "/")
	ret.Inter = true
	return
}

func (this *Writer) Check() {
	var c core.ChunkStream
	for {
		if err := this.conn.Read(&c); err != nil {
			this.Close(err)
			return
		}
	}
}

func (this *Writer) SendPacket() error {
	Flush := reflect.ValueOf(this.conn).MethodByName("Flush")
	var cs core.ChunkStream
	for {
		p, ok := <-this.packetQueue
		if ok {
			cs.Data = p.Data
			cs.Length = uint32(len(p.Data))
			cs.StreamID = p.StreamID
			cs.Timestamp = p.TimeStamp
			cs.Timestamp += this.BaseTimeStamp()

			if p.IsVideo {
				cs.TypeID = core.TAG_VIDEO
			} else {
				if p.IsMetadata {
					cs.TypeID = core.TAG_SCRIPTDATAAMF0
				} else {
					cs.TypeID = core.TAG_AUDIO
				}
			}

			this.SaveStatics(p.StreamID, uint64(cs.Length), p.IsVideo)
			this.SetPreTime()
			this.RecTimeStamp(cs.Timestamp, cs.TypeID)
			err := this.conn.Write(cs)
			if err != nil {
				this.closed = true
				return err
			}
			Flush.Call(nil)
		} else {
			return fmt.Errorf("closed")
		}

	}
}

func (this *Writer) SaveStatics(streamid uint32, length uint64, isVideoFlag bool) {
	nowInMS := int64(time.Now().UnixNano() / 1e6)
	this.WriteBWInfo.StreamId = streamid
	if isVideoFlag {
		this.WriteBWInfo.VideoDatainBytes = this.WriteBWInfo.VideoDatainBytes + length
	} else {
		this.WriteBWInfo.AudioDatainBytes = this.WriteBWInfo.AudioDatainBytes + length
	}

	if this.WriteBWInfo.LastTimestamp == 0 {
		this.WriteBWInfo.LastTimestamp = nowInMS
	} else if (nowInMS - this.WriteBWInfo.LastTimestamp) >= SAVE_STATICS_INTERVAL {
		diffTimestamp := (nowInMS - this.WriteBWInfo.LastTimestamp) / 1000

		this.WriteBWInfo.VideoSpeedInBytesperMS = (this.WriteBWInfo.VideoDatainBytes - this.WriteBWInfo.LastVideoDatainBytes) * 8 / uint64(diffTimestamp) / 1000
		this.WriteBWInfo.AudioSpeedInBytesperMS = (this.WriteBWInfo.AudioDatainBytes - this.WriteBWInfo.LastAudioDatainBytes) * 8 / uint64(diffTimestamp) / 1000

		this.WriteBWInfo.LastVideoDatainBytes = this.WriteBWInfo.VideoDatainBytes
		this.WriteBWInfo.LastAudioDatainBytes = this.WriteBWInfo.AudioDatainBytes
		this.WriteBWInfo.LastTimestamp = nowInMS
	}
}

func (this *Writer) Close(err error) {
	log.Debugf("[SYS LiveGo] player ", this.Info(), "closed: "+err.Error())
	if !this.closed {
		close(this.packetQueue)
	}
	this.closed = true
	this.conn.Close(err)
}
