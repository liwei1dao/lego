package live

import (
	"bytes"
	"fmt"
	"time"

	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/lib/modules/live/container/flv"
	"github.com/liwei1dao/lego/lib/modules/live/container/ts"
	"github.com/liwei1dao/lego/lib/modules/live/parser"
	"github.com/liwei1dao/lego/sys/log"
)

const (
	duration = 3000
)

const (
	videoHZ           = 90000
	aacSampleLen      = 1024
	sourcemaxQueueNum = 512

	h264_default_hz uint64 = 90
)

func NewSource(options IOptions, info av.Info) *Source {
	info.Inter = true
	s := &Source{
		options:     options,
		info:        info,
		align:       &align{},
		stat:        newStatus(),
		RWBaser:     av.NewRWBaser(time.Second * 10),
		cache:       newAudioCache(),
		demuxer:     flv.NewDemuxer(),
		muxer:       ts.NewMuxer(),
		tsCache:     NewTSCacheItem(info.Key),
		tsparser:    parser.NewCodecParser(),
		bwriter:     bytes.NewBuffer(make([]byte, 100*1024)),
		packetQueue: make(chan *av.Packet, sourcemaxQueueNum),
	}
	go func() {
		err := s.SendPacket()
		if err != nil {
			log.Debugf("send packet error: %v", err)
			s.closed = true
		}
	}()
	return s
}

type Source struct {
	av.RWBaser
	options     IOptions
	seq         int
	info        av.Info
	bwriter     *bytes.Buffer
	btswriter   *bytes.Buffer
	demuxer     *flv.Demuxer
	muxer       *ts.Muxer
	pts, dts    uint64
	stat        *status
	align       *align
	cache       *audioCache
	tsCache     *TSCacheItem
	tsparser    *parser.CodecParser
	closed      bool
	packetQueue chan *av.Packet
}

func (source *Source) GetCacheInc() *TSCacheItem {
	return source.tsCache
}
func (source *Source) DropPacket(pktQue chan *av.Packet, info av.Info) {
	log.Warnf("[%v] packet queue max!!!", info)
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
			videoPkt, ok := tmpPkt.Header.(av.VideoPacketHeader)
			// dont't drop sps config and dont't drop key frame
			if ok && (videoPkt.IsSeq() || videoPkt.IsKeyFrame()) {
				pktQue <- tmpPkt
			}
			if len(pktQue) > maxQueueNum-10 {
				<-pktQue
			}
		}

	}
	log.Debugf("packet queue len: %d", len(pktQue))
}
func (source *Source) Write(p *av.Packet) (err error) {
	err = nil
	if source.closed {
		err = fmt.Errorf("hls source closed")
		return
	}
	source.SetPreTime()
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("hls source has already been closed:%v", e)
		}
	}()
	if len(source.packetQueue) >= maxQueueNum-24 {
		source.DropPacket(source.packetQueue, source.info)
	} else {
		if !source.closed {
			source.packetQueue <- p
		}
	}
	return
}

func (source *Source) Info() (ret av.Info) {
	return source.info
}

func (source *Source) cleanup() {
	close(source.packetQueue)
	source.bwriter = nil
	source.btswriter = nil
	source.cache = nil
	source.tsCache = nil
}

func (source *Source) Close(err error) {
	log.Debugf("hls source closed: ", source.info)
	if !source.closed && !source.options.GetHlsKeepAfterEnd() {
		source.cleanup()
	}
	source.closed = true
}

func (source *Source) SendPacket() error {
	defer func() {
		log.Debugf("[%v] hls sender stop", source.info)
		if r := recover(); r != nil {
			log.Warnf("hls SendPacket panic: %v", r)
		}
	}()

	log.Debugf("[%v] hls sender start", source.info)
	for {
		if source.closed {
			return fmt.Errorf("closed")
		}

		p, ok := <-source.packetQueue
		if ok {
			if p.IsMetadata {
				continue
			}

			err := source.demuxer.Demux(p)
			if err == flv.ErrAvcEndSEQ {
				log.Warnf("err:%v", err)
				continue
			} else {
				if err != nil {
					log.Warnf("err:%v", err)
					return err
				}
			}
			compositionTime, isSeq, err := source.parse(p)
			if err != nil {
				log.Warnf("err:%v", err)
			}
			if err != nil || isSeq {
				continue
			}
			if source.btswriter != nil {
				source.stat.update(p.IsVideo, p.TimeStamp)
				source.calcPtsDts(p.IsVideo, p.TimeStamp, uint32(compositionTime))
				source.tsMux(p)
			}
		} else {
			return fmt.Errorf("closed")
		}
	}
}

func (source *Source) parse(p *av.Packet) (int32, bool, error) {
	var compositionTime int32
	var ah av.AudioPacketHeader
	var vh av.VideoPacketHeader
	if p.IsVideo {
		vh = p.Header.(av.VideoPacketHeader)
		if vh.CodecID() != av.VIDEO_H264 {
			return compositionTime, false, ErrNoSupportVideoCodec
		}
		compositionTime = vh.CompositionTime()
		if vh.IsKeyFrame() && vh.IsSeq() {
			return compositionTime, true, source.tsparser.Parse(p, source.bwriter)
		}
	} else {
		ah = p.Header.(av.AudioPacketHeader)
		if ah.SoundFormat() != av.SOUND_AAC {
			return compositionTime, false, ErrNoSupportAudioCodec
		}
		if ah.AACPacketType() == av.AAC_SEQHDR {
			return compositionTime, true, source.tsparser.Parse(p, source.bwriter)
		}
	}
	source.bwriter.Reset()
	if err := source.tsparser.Parse(p, source.bwriter); err != nil {
		return compositionTime, false, err
	}
	p.Data = source.bwriter.Bytes()

	if p.IsVideo && vh.IsKeyFrame() {
		source.cut()
	}
	return compositionTime, false, nil
}

func (source *Source) calcPtsDts(isVideo bool, ts, compositionTs uint32) {
	source.dts = uint64(ts) * h264_default_hz
	if isVideo {
		source.pts = source.dts + uint64(compositionTs)*h264_default_hz
	} else {
		sampleRate, _ := source.tsparser.SampleRate()
		source.align.align(&source.dts, uint32(videoHZ*aacSampleLen/sampleRate))
		source.pts = source.dts
	}
}

func (source *Source) cut() {
	newf := true
	if source.btswriter == nil {
		source.btswriter = bytes.NewBuffer(nil)
	} else if source.btswriter != nil && source.stat.durationMs() >= duration {
		source.flushAudio()

		source.seq++
		filename := fmt.Sprintf("/%s/%d.ts", source.info.Key, time.Now().Unix())
		item := NewTSItem(filename, int(source.stat.durationMs()), source.seq, source.btswriter.Bytes())
		source.tsCache.SetItem(filename, item)

		source.btswriter.Reset()
		source.stat.resetAndNew()
	} else {
		newf = false
	}
	if newf {
		source.btswriter.Write(source.muxer.PAT())
		source.btswriter.Write(source.muxer.PMT(av.SOUND_AAC, true))
	}
}

func (source *Source) tsMux(p *av.Packet) error {
	if p.IsVideo {
		return source.muxer.Mux(p, source.btswriter)
	} else {
		source.cache.Cache(p.Data, source.pts)
		return source.muxAudio(cache_max_frames)
	}
}

func (source *Source) muxAudio(limit byte) error {
	if source.cache.CacheNum() < limit {
		return nil
	}
	var p av.Packet
	_, pts, buf := source.cache.GetFrame()
	p.Data = buf
	p.TimeStamp = uint32(pts / h264_default_hz)
	return source.muxer.Mux(&p, source.btswriter)
}

func (source *Source) flushAudio() error {
	return source.muxAudio(1)
}
