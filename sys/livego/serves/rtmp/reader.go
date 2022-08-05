package rtmp

import (
	"net/url"
	"strings"
	"time"

	"github.com/liwei1dao/lego/sys/livego/container/flv"
	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/utils/container/id"
)

func NewReader(conn core.StreamReadWriteCloser) *Reader {
	return &Reader{
		Uid:  id.NewXId(),
		conn: conn,
	}
}

type Reader struct {
	core.RWBaser
	Uid        string
	demuxer    *flv.Demuxer
	conn       core.StreamReadWriteCloser
	ReadBWInfo StaticsBW
}

func (this *Reader) Read(p *core.Packet) (err error) {
	defer func() {
		if r := recover(); r != nil {
			this.conn.Log().Warnf("rtmp read packet panic:%v", r)
		}
	}()

	this.SetPreTime()
	var cs core.ChunkStream
	for {
		err = this.conn.Read(&cs)
		if err != nil {
			return err
		}
		if cs.TypeID == core.TAG_AUDIO ||
			cs.TypeID == core.TAG_VIDEO ||
			cs.TypeID == core.TAG_SCRIPTDATAAMF0 ||
			cs.TypeID == core.TAG_SCRIPTDATAAMF3 {
			break
		}
	}

	p.IsAudio = cs.TypeID == core.TAG_AUDIO
	p.IsVideo = cs.TypeID == core.TAG_VIDEO
	p.IsMetadata = cs.TypeID == core.TAG_SCRIPTDATAAMF0 || cs.TypeID == core.TAG_SCRIPTDATAAMF3
	p.StreamID = cs.StreamID
	p.Data = cs.Data
	p.TimeStamp = cs.Timestamp

	this.SaveStatics(p.StreamID, uint64(len(p.Data)), p.IsVideo)
	this.demuxer.DemuxH(p)
	return err
}

func (this *Reader) Info() (ret core.Info) {
	ret.UID = this.Uid
	_, _, URL := this.conn.GetInfo()
	ret.URL = URL
	_url, err := url.Parse(URL)
	if err != nil {
		this.conn.Log().Warnf("err:%v", err)
	}
	ret.Key = strings.TrimLeft(_url.Path, "/")
	return
}

func (this *Reader) SaveStatics(streamid uint32, length uint64, isVideoFlag bool) {
	nowInMS := int64(time.Now().UnixNano() / 1e6)

	this.ReadBWInfo.StreamId = streamid
	if isVideoFlag {
		this.ReadBWInfo.VideoDatainBytes = this.ReadBWInfo.VideoDatainBytes + length
	} else {
		this.ReadBWInfo.AudioDatainBytes = this.ReadBWInfo.AudioDatainBytes + length
	}

	if this.ReadBWInfo.LastTimestamp == 0 {
		this.ReadBWInfo.LastTimestamp = nowInMS
	} else if (nowInMS - this.ReadBWInfo.LastTimestamp) >= SAVE_STATICS_INTERVAL {
		diffTimestamp := (nowInMS - this.ReadBWInfo.LastTimestamp) / 1000

		//log.Printf("now=%d, last=%d, diff=%d", nowInMS, v.ReadBWInfo.LastTimestamp, diffTimestamp)
		this.ReadBWInfo.VideoSpeedInBytesperMS = (this.ReadBWInfo.VideoDatainBytes - this.ReadBWInfo.LastVideoDatainBytes) * 8 / uint64(diffTimestamp) / 1000
		this.ReadBWInfo.AudioSpeedInBytesperMS = (this.ReadBWInfo.AudioDatainBytes - this.ReadBWInfo.LastAudioDatainBytes) * 8 / uint64(diffTimestamp) / 1000

		this.ReadBWInfo.LastVideoDatainBytes = this.ReadBWInfo.VideoDatainBytes
		this.ReadBWInfo.LastAudioDatainBytes = this.ReadBWInfo.AudioDatainBytes
		this.ReadBWInfo.LastTimestamp = nowInMS
	}
}

func (this *Reader) Close(err error) {
	this.conn.Log().Debugf("publisher %v  closed: ", this.Info(), err)
	this.conn.Close(err)
}
