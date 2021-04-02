package live

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/lib/modules/live/container/flv"
	"github.com/liwei1dao/lego/lib/modules/live/liveconn"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/uid"
)

const (
	maxQueueNum           = 1024
	SAVE_STATICS_INTERVAL = 5000
)

type (
	GetInFo interface {
		GetInfo() (string, string, string)
	}
	StreamReadWriteCloser interface {
		GetInFo
		Close(error)
		Write(liveconn.ChunkStream) error
		Read(c *liveconn.ChunkStream) error
	}
	StaticsBW struct {
		StreamId               uint32
		VideoDatainBytes       uint64
		LastVideoDatainBytes   uint64
		VideoSpeedInBytesperMS uint64

		AudioDatainBytes       uint64
		LastAudioDatainBytes   uint64
		AudioSpeedInBytesperMS uint64

		LastTimestamp int64
	}
)

//-------------------------------------------------------------------------------------------------------------------------------------------------------------
//
//                            _ooOoo_
//                           o8888888o
//                           88" . "88
//                           (| -_- |)
//                           O\  =  /O
//                        ____/`---'\____
//                      .'  \\|     |//  `.
//                     /  \\|||  :  |||//  \
//                    /  _||||| -:- |||||-  \
//                    |   | \\\  -  /// |   |
//                    | \_|  ''\---/''  |   |
//                    \  .-\__  `-`  ___/-. /
//                  ___`. .'  /--.--\  `. . __
//               ."" '<  `.___\_<|>_/___.'  >'"".
//              | | :  `- \`.;`\ _ /`;.`/ - ` : | |
//              \  \ `-.   \_ __\ /__ _/   .-` /  /
//         ======`-.____`-.___\_____/___.-`____.-'======
//                            `=---='
//        ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//                      Buddha Bless, No Bug !
//VirReader----------------------------------------------------------------------------------------------------------------------------------------------------------
func NewVirReader(conn StreamReadWriteCloser, writeTimeout int) *VirReader {
	return &VirReader{
		Uid:        uid.NewId(),
		conn:       conn,
		RWBaser:    av.NewRWBaser(time.Second * time.Duration(writeTimeout)),
		demuxer:    flv.NewDemuxer(),
		ReadBWInfo: StaticsBW{0, 0, 0, 0, 0, 0, 0, 0},
	}
}

type VirReader struct {
	Uid string
	av.RWBaser
	demuxer    *flv.Demuxer
	conn       StreamReadWriteCloser
	ReadBWInfo StaticsBW
}

func (v *VirReader) SaveStatics(streamid uint32, length uint64, isVideoFlag bool) {
	nowInMS := int64(time.Now().UnixNano() / 1e6)

	v.ReadBWInfo.StreamId = streamid
	if isVideoFlag {
		v.ReadBWInfo.VideoDatainBytes = v.ReadBWInfo.VideoDatainBytes + length
	} else {
		v.ReadBWInfo.AudioDatainBytes = v.ReadBWInfo.AudioDatainBytes + length
	}

	if v.ReadBWInfo.LastTimestamp == 0 {
		v.ReadBWInfo.LastTimestamp = nowInMS
	} else if (nowInMS - v.ReadBWInfo.LastTimestamp) >= SAVE_STATICS_INTERVAL {
		diffTimestamp := (nowInMS - v.ReadBWInfo.LastTimestamp) / 1000

		//log.Printf("now=%d, last=%d, diff=%d", nowInMS, v.ReadBWInfo.LastTimestamp, diffTimestamp)
		v.ReadBWInfo.VideoSpeedInBytesperMS = (v.ReadBWInfo.VideoDatainBytes - v.ReadBWInfo.LastVideoDatainBytes) * 8 / uint64(diffTimestamp) / 1000
		v.ReadBWInfo.AudioSpeedInBytesperMS = (v.ReadBWInfo.AudioDatainBytes - v.ReadBWInfo.LastAudioDatainBytes) * 8 / uint64(diffTimestamp) / 1000

		v.ReadBWInfo.LastVideoDatainBytes = v.ReadBWInfo.VideoDatainBytes
		v.ReadBWInfo.LastAudioDatainBytes = v.ReadBWInfo.AudioDatainBytes
		v.ReadBWInfo.LastTimestamp = nowInMS
	}
}

func (v *VirReader) Read(p *av.Packet) (err error) {
	defer cbase.Recover()

	v.SetPreTime()
	var cs liveconn.ChunkStream
	for {
		err = v.conn.Read(&cs)
		if err != nil {
			return err
		}
		if cs.TypeID == av.TAG_AUDIO ||
			cs.TypeID == av.TAG_VIDEO ||
			cs.TypeID == av.TAG_SCRIPTDATAAMF0 ||
			cs.TypeID == av.TAG_SCRIPTDATAAMF3 {
			break
		}
	}

	p.IsAudio = cs.TypeID == av.TAG_AUDIO
	p.IsVideo = cs.TypeID == av.TAG_VIDEO
	p.IsMetadata = cs.TypeID == av.TAG_SCRIPTDATAAMF0 || cs.TypeID == av.TAG_SCRIPTDATAAMF3
	p.StreamID = cs.StreamID
	p.Data = cs.Data
	p.TimeStamp = cs.Timestamp

	v.SaveStatics(p.StreamID, uint64(len(p.Data)), p.IsVideo)
	v.demuxer.DemuxH(p)
	return err
}

func (v *VirReader) Info() (ret av.Info) {
	ret.UID = v.Uid
	_, _, URL := v.conn.GetInfo()
	ret.URL = URL
	_url, err := url.Parse(URL)
	if err != nil {
		log.Warnf("err:%v", err)
	}
	ret.Key = strings.TrimLeft(_url.Path, "/")
	return
}

func (v *VirReader) Close(err error) {
	log.Debugf("publisher ", v.Info(), "closed: "+err.Error())
	v.conn.Close(err)
}

//-------------------------------------------------------------------------------------------------------------------------------------------------------------
//
//                            _ooOoo_
//                           o8888888o
//                           88" . "88
//                           (| -_- |)
//                           O\  =  /O
//                        ____/`---'\____
//                      .'  \\|     |//  `.
//                     /  \\|||  :  |||//  \
//                    /  _||||| -:- |||||-  \
//                    |   | \\\  -  /// |   |
//                    | \_|  ''\---/''  |   |
//                    \  .-\__  `-`  ___/-. /
//                  ___`. .'  /--.--\  `. . __
//               ."" '<  `.___\_<|>_/___.'  >'"".
//              | | :  `- \`.;`\ _ /`;.`/ - ` : | |
//              \  \ `-.   \_ __\ /__ _/   .-` /  /
//         ======`-.____`-.___\_____/___.-`____.-'======
//                            `=---='
//        ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
//                      Buddha Bless, No Bug !
//VirWriter----------------------------------------------------------------------------------------------------------------------------------------------------------
func NewVirWriter(writeTimeout int, conn StreamReadWriteCloser) *VirWriter {
	ret := &VirWriter{
		Uid:         uid.NewId(),
		conn:        conn,
		RWBaser:     av.NewRWBaser(time.Second * time.Duration(writeTimeout)),
		packetQueue: make(chan *av.Packet, maxQueueNum),
		WriteBWInfo: StaticsBW{0, 0, 0, 0, 0, 0, 0, 0},
	}

	go ret.Check()
	go func() {
		err := ret.SendPacket()
		if err != nil {
			log.Warnf("err:%v", err)
		}
	}()
	return ret
}

type VirWriter struct {
	Uid    string
	closed bool
	conn   StreamReadWriteCloser
	av.RWBaser
	packetQueue chan *av.Packet
	WriteBWInfo StaticsBW
}

func (v *VirWriter) SaveStatics(streamid uint32, length uint64, isVideoFlag bool) {
	nowInMS := int64(time.Now().UnixNano() / 1e6)

	v.WriteBWInfo.StreamId = streamid
	if isVideoFlag {
		v.WriteBWInfo.VideoDatainBytes = v.WriteBWInfo.VideoDatainBytes + length
	} else {
		v.WriteBWInfo.AudioDatainBytes = v.WriteBWInfo.AudioDatainBytes + length
	}

	if v.WriteBWInfo.LastTimestamp == 0 {
		v.WriteBWInfo.LastTimestamp = nowInMS
	} else if (nowInMS - v.WriteBWInfo.LastTimestamp) >= SAVE_STATICS_INTERVAL {
		diffTimestamp := (nowInMS - v.WriteBWInfo.LastTimestamp) / 1000

		v.WriteBWInfo.VideoSpeedInBytesperMS = (v.WriteBWInfo.VideoDatainBytes - v.WriteBWInfo.LastVideoDatainBytes) * 8 / uint64(diffTimestamp) / 1000
		v.WriteBWInfo.AudioSpeedInBytesperMS = (v.WriteBWInfo.AudioDatainBytes - v.WriteBWInfo.LastAudioDatainBytes) * 8 / uint64(diffTimestamp) / 1000

		v.WriteBWInfo.LastVideoDatainBytes = v.WriteBWInfo.VideoDatainBytes
		v.WriteBWInfo.LastAudioDatainBytes = v.WriteBWInfo.AudioDatainBytes
		v.WriteBWInfo.LastTimestamp = nowInMS
	}
}

func (v *VirWriter) DropPacket(pktQue chan *av.Packet, info av.Info) {
	log.Warnf("[%v] packet queue max!!!", info)
	for i := 0; i < maxQueueNum-84; i++ {
		tmpPkt, ok := <-pktQue
		// try to don't drop audio
		if ok && tmpPkt.IsAudio {
			if len(pktQue) > maxQueueNum-2 {
				log.Debug("drop audio pkt")
				<-pktQue
			} else {
				pktQue <- tmpPkt
			}

		}

		if ok && tmpPkt.IsVideo {
			videoPkt, ok := tmpPkt.Header.(av.VideoPacketHeader)
			// dont't drop sps config and dont't drop key frame
			if ok && (videoPkt.IsSeq() || videoPkt.IsKeyFrame()) {
				pktQue <- tmpPkt
			}
			if len(pktQue) > maxQueueNum-10 {
				log.Debug("drop video pkt")
				<-pktQue
			}
		}

	}
	log.Debugf("packet queue len: ", len(pktQue))
}

func (v *VirWriter) SendPacket() error {
	Flush := reflect.ValueOf(v.conn).MethodByName("Flush")
	var cs liveconn.ChunkStream
	for {
		p, ok := <-v.packetQueue
		if ok {
			cs.Data = p.Data
			cs.Length = uint32(len(p.Data))
			cs.StreamID = p.StreamID
			cs.Timestamp = p.TimeStamp
			cs.Timestamp += v.BaseTimeStamp()

			if p.IsVideo {
				cs.TypeID = av.TAG_VIDEO
			} else {
				if p.IsMetadata {
					cs.TypeID = av.TAG_SCRIPTDATAAMF0
				} else {
					cs.TypeID = av.TAG_AUDIO
				}
			}

			v.SaveStatics(p.StreamID, uint64(cs.Length), p.IsVideo)
			v.SetPreTime()
			v.RecTimeStamp(cs.Timestamp, cs.TypeID)
			err := v.conn.Write(cs)
			if err != nil {
				v.closed = true
				return err
			}
			Flush.Call(nil)
		} else {
			return fmt.Errorf("closed")
		}

	}
}

func (v *VirWriter) Write(p *av.Packet) (err error) {
	err = nil

	if v.closed {
		err = fmt.Errorf("VirWriter closed")
		return
	}
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("VirWriter has already been closed:%v", e)
		}
	}()
	if len(v.packetQueue) >= maxQueueNum-24 {
		v.DropPacket(v.packetQueue, v.Info())
	} else {
		v.packetQueue <- p
	}

	return
}

func (v *VirWriter) Info() (ret av.Info) {
	ret.UID = v.Uid
	_, _, URL := v.conn.GetInfo()
	ret.URL = URL
	_url, err := url.Parse(URL)
	if err != nil {
		log.Warnf("err:%v", err)
	}
	ret.Key = strings.TrimLeft(_url.Path, "/")
	ret.Inter = true
	return
}

func (v *VirWriter) Check() {
	var c liveconn.ChunkStream
	for {
		if err := v.conn.Read(&c); err != nil {
			v.Close(err)
			return
		}
	}
}

func (v *VirWriter) Close(err error) {
	log.Warnf("player ", v.Info(), "closed: "+err.Error())
	if !v.closed {
		close(v.packetQueue)
	}
	v.closed = true
	v.conn.Close(err)
}
