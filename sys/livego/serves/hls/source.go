package hls

import (
	"bytes"
	"fmt"
	"time"

	"github.com/liwei1dao/lego/sys/livego/container/flv"
	"github.com/liwei1dao/lego/sys/livego/container/ts"
	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/parser"
	"github.com/liwei1dao/lego/sys/log"
)

const (
	duration = 3000
)
const (
	videoHZ      = 90000
	aacSampleLen = 1024
	maxQueueNum  = 512

	h264_default_hz uint64 = 90
)

func NewSource(sys core.ISys, log log.ILogger, info core.Info) *Source {
	info.Inter = true
	s := &Source{
		sys:         sys,
		log:         log,
		info:        info,
		align:       &align{},
		stat:        newStatus(),
		RWBaser:     core.NewRWBaser(time.Second * 10),
		cache:       newAudioCache(),
		demuxer:     flv.NewDemuxer(),
		muxer:       ts.NewMuxer(),
		tsCache:     NewTSCacheItem(info.Key),
		tsparser:    parser.NewCodecParser(),
		bwriter:     bytes.NewBuffer(make([]byte, 100*1024)),
		packetQueue: make(chan *core.Packet, maxQueueNum),
	}
	go func() {
		err := s.SendPacket()
		if err != nil {
			s.log.Debugf("send packet error:%v", err)
			s.closed = true
		}
	}()
	return s
}

type Source struct {
	core.RWBaser
	sys         core.ISys
	log         log.ILogger
	seq         int
	info        core.Info
	stat        *status
	pts, dts    uint64
	align       *align
	bwriter     *bytes.Buffer
	btswriter   *bytes.Buffer
	demuxer     *flv.Demuxer
	muxer       *ts.Muxer
	cache       *audioCache
	tsCache     *TSCacheItem
	tsparser    *parser.CodecParser
	packetQueue chan *core.Packet
	closed      bool
}

func (this *Source) DropPacket(pktQue chan *core.Packet, info core.Info) {
	this.log.Warnf("[%v] packet queue max!!!", info)
	for i := 0; i < maxQueueNum-84; i++ {
		tmpPkt, ok := <-pktQue
		// try to don't drop audio
		if ok && tmpPkt.IsAudio {
			if len(pktQue) > maxQueueNum-2 {
				<-pktQue
			} else {
				pktQue <- tmpPkt
			}
		}

		if ok && tmpPkt.IsVideo {
			videoPkt, ok := tmpPkt.Header.(core.VideoPacketHeader)
			// dont't drop sps config and dont't drop key frame
			if ok && (videoPkt.IsSeq() || videoPkt.IsKeyFrame()) {
				pktQue <- tmpPkt
			}
			if len(pktQue) > maxQueueNum-10 {
				<-pktQue
			}
		}

	}
	this.log.Debugf("packet queue len:%d", len(pktQue))
}
func (this *Source) Write(p *core.Packet) (err error) {
	err = nil
	if this.closed {
		err = fmt.Errorf("hls source closed")
		return
	}
	this.SetPreTime()
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("hls source has already been closed:%v", e)
		}
	}()
	if len(this.packetQueue) >= maxQueueNum-24 {
		this.DropPacket(this.packetQueue, this.info)
	} else {
		if !this.closed {
			this.packetQueue <- p
		}
	}
	return
}
func (this *Source) Info() (ret core.Info) {
	return this.info
}
func (this *Source) Close(err error) {
	this.log.Debugf("hls source closed:%v", this.info)
	if !this.closed && !this.sys.GetHLSKeepAfterEnd() {
		this.cleanup()
	}
	this.closed = true
}

func (this *Source) cleanup() {
	close(this.packetQueue)
	this.bwriter = nil
	this.btswriter = nil
	this.cache = nil
	this.tsCache = nil
}

func (this *Source) tsMux(p *core.Packet) error {
	if p.IsVideo {
		return this.muxer.Mux(p, this.btswriter)
	} else {
		this.cache.Cache(p.Data, this.pts)
		return this.muxAudio(cache_max_frames)
	}
}
func (this *Source) muxAudio(limit byte) error {
	if this.cache.CacheNum() < limit {
		return nil
	}
	var p core.Packet
	_, pts, buf := this.cache.GetFrame()
	p.Data = buf
	p.TimeStamp = uint32(pts / h264_default_hz)
	return this.muxer.Mux(&p, this.btswriter)
}

func (this *Source) SendPacket() error {
	defer func() {
		this.log.Debugf("[%v] hls sender stop", this.info)
		if r := recover(); r != nil {
			this.log.Warnf("hls SendPacket panic: ", r)
		}
	}()

	this.log.Debugf("[%v] hls sender start", this.info)
	for {
		if this.closed {
			return fmt.Errorf("closed")
		}

		p, ok := <-this.packetQueue
		if ok {
			if p.IsMetadata {
				continue
			}

			err := this.demuxer.Demux(p)
			if err == flv.ErrAvcEndSEQ {
				this.log.Warnf("hls err:%v", err)
				continue
			} else {
				if err != nil {
					this.log.Warnf("hls err:%v", err)
					return err
				}
			}
			compositionTime, isSeq, err := this.parse(p)
			if err != nil {
				this.log.Warnf("hls err:%v", err)
			}
			if err != nil || isSeq {
				continue
			}
			if this.btswriter != nil {
				this.stat.update(p.IsVideo, p.TimeStamp)
				this.calcPtsDts(p.IsVideo, p.TimeStamp, uint32(compositionTime))
				this.tsMux(p)
			}
		} else {
			return fmt.Errorf("closed")
		}
	}
}

func (this *Source) GetCacheInc() *TSCacheItem {
	return this.tsCache
}

func (this *Source) parse(p *core.Packet) (int32, bool, error) {
	var compositionTime int32
	var ah core.AudioPacketHeader
	var vh core.VideoPacketHeader
	if p.IsVideo {
		vh = p.Header.(core.VideoPacketHeader)
		if vh.CodecID() != core.VIDEO_H264 {
			return compositionTime, false, ErrNoSupportVideoCodec
		}
		compositionTime = vh.CompositionTime()
		if vh.IsKeyFrame() && vh.IsSeq() {
			return compositionTime, true, this.tsparser.Parse(p, this.bwriter)
		}
	} else {
		ah = p.Header.(core.AudioPacketHeader)
		if ah.SoundFormat() != core.SOUND_AAC {
			return compositionTime, false, ErrNoSupportAudioCodec
		}
		if ah.AACPacketType() == core.AAC_SEQHDR {
			return compositionTime, true, this.tsparser.Parse(p, this.bwriter)
		}
	}
	this.bwriter.Reset()
	if err := this.tsparser.Parse(p, this.bwriter); err != nil {
		return compositionTime, false, err
	}
	p.Data = this.bwriter.Bytes()

	if p.IsVideo && vh.IsKeyFrame() {
		this.cut()
	}
	return compositionTime, false, nil
}

func (this *Source) calcPtsDts(isVideo bool, ts, compositionTs uint32) {
	this.dts = uint64(ts) * h264_default_hz
	if isVideo {
		this.pts = this.dts + uint64(compositionTs)*h264_default_hz
	} else {
		sampleRate, _ := this.tsparser.SampleRate()
		this.align.align(&this.dts, uint32(videoHZ*aacSampleLen/sampleRate))
		this.pts = this.dts
	}
}

func (this *Source) cut() {
	newf := true
	if this.btswriter == nil {
		this.btswriter = bytes.NewBuffer(nil)
	} else if this.btswriter != nil && this.stat.durationMs() >= duration {
		this.flushAudio()

		this.seq++
		filename := fmt.Sprintf("/%s/%d.ts", this.info.Key, time.Now().Unix())
		item := NewTSItem(filename, int(this.stat.durationMs()), this.seq, this.btswriter.Bytes())
		this.tsCache.SetItem(filename, item)

		this.btswriter.Reset()
		this.stat.resetAndNew()
	} else {
		newf = false
	}
	if newf {
		this.btswriter.Write(this.muxer.PAT())
		this.btswriter.Write(this.muxer.PMT(core.SOUND_AAC, true))
	}
}

func (this *Source) flushAudio() error {
	return this.muxAudio(1)
}
