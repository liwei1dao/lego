package rtmp

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/liwei1dao/lego/lib/modules/live/utils/uid"
	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/packet"
	"github.com/liwei1dao/lego/sys/log"
)

type StaticsBW struct {
	StreamId               uint32
	VideoDatainBytes       uint64
	LastVideoDatainBytes   uint64
	VideoSpeedInBytesperMS uint64
	AudioDatainBytes       uint64
	LastAudioDatainBytes   uint64
	AudioSpeedInBytesperMS uint64
	LastTimestamp          int64
}

func NewVirWriter(conn StreamReadWriteCloser, writeTimeout int) *VirWriter {
	ret := &VirWriter{
		Uid:         uid.NewId(),
		conn:        conn,
		RWBaser:     packet.NewRWBaser(time.Second * time.Duration(writeTimeout)),
		packetQueue: make(chan *packet.Packet, maxQueueNum),
		WriteBWInfo: StaticsBW{0, 0, 0, 0, 0, 0, 0, 0},
	}

	go ret.Check()
	go func() {
		err := ret.SendPacket()
		if err != nil {
			log.Warnf("[SYS LiveGo] %v", err)
		}
	}()
	return ret
}

type VirWriter struct {
	Uid    string
	closed bool
	packet.RWBaser
	conn        StreamReadWriteCloser
	packetQueue chan *packet.Packet
	WriteBWInfo StaticsBW
}

func (this *VirWriter) Write(p *packet.Packet) (err error) {
	err = nil

	if this.closed {
		err = fmt.Errorf("VirWriter closed")
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("VirWriter has already been closed:%v", e)
		}
	}()
	if len(this.packetQueue) >= maxQueueNum-24 {
		this.DropPacket(this.packetQueue, this.Info())
	} else {
		this.packetQueue <- p
	}

	return
}

func (v *VirWriter) DropPacket(pktQue chan *packet.Packet, info packet.Info) {
	log.Warnf("[SYS LiveGo] [%v] packet queue max!!!", info)
	for i := 0; i < maxQueueNum-84; i++ {
		tmpPkt, ok := <-pktQue
		// try to don't drop audio
		if ok && tmpPkt.IsAudio {
			if len(pktQue) > maxQueueNum-2 {
				log.Debugf("[SYS LiveGo] drop audio pkt")
				<-pktQue
			} else {
				pktQue <- tmpPkt
			}

		}

		if ok && tmpPkt.IsVideo {
			videoPkt, ok := tmpPkt.Header.(packet.VideoPacketHeader)
			// dont't drop sps config and dont't drop key frame
			if ok && (videoPkt.IsSeq() || videoPkt.IsKeyFrame()) {
				pktQue <- tmpPkt
			}
			if len(pktQue) > maxQueueNum-10 {
				log.Debug("[SYS LiveGo] drop video pkt")
				<-pktQue
			}
		}

	}
	log.Debugf("[SYS LiveGo] packet queue len: ", len(pktQue))
}

func (this *VirWriter) Info() (ret packet.Info) {
	ret.UID = this.Uid
	_, _, URL := this.conn.GetInfo()
	ret.URL = URL
	_url, err := url.Parse(URL)
	if err != nil {
		log.Warnf("[SYS LiveGo] err:%v", err)
	}
	ret.Key = strings.TrimLeft(_url.Path, "/")
	ret.Inter = true
	return
}

func (this *VirWriter) Check() {
	var c core.ChunkStream
	for {
		if err := this.conn.Read(&c); err != nil {
			this.Close(err)
			return
		}
	}
}

func (this *VirWriter) Close(err error) {
	log.Warnf("[SYS LiveGo] player:%v closed:%v", this.Info(), err)
	if !this.closed {
		close(this.packetQueue)
	}
	this.closed = true
	this.conn.Close(err)
}

func (this *VirWriter) SendPacket() error {
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
				cs.TypeID = packet.TAG_VIDEO
			} else {
				if p.IsMetadata {
					cs.TypeID = packet.TAG_SCRIPTDATAAMF0
				} else {
					cs.TypeID = packet.TAG_AUDIO
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

func (this *VirWriter) SaveStatics(streamid uint32, length uint64, isVideoFlag bool) {
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
