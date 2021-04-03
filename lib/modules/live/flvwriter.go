package live

import (
	"fmt"
	"net/http"
	"time"

	"github.com/liwei1dao/lego/lib/modules/live/amf"
	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/pio"
	"github.com/liwei1dao/lego/utils/uid"
)

const (
	headerLen = 11
	// maxQueueNum = 1024
)

func NewFLVWriter(app, title, url string, ctx http.ResponseWriter) *FLVWriter {
	ret := &FLVWriter{
		Uid:         uid.NewId(),
		app:         app,
		title:       title,
		url:         url,
		ctx:         ctx,
		RWBaser:     av.NewRWBaser(time.Second * 10),
		closedChan:  make(chan struct{}),
		buf:         make([]byte, headerLen),
		packetQueue: make(chan *av.Packet, maxQueueNum),
	}

	if _, err := ret.ctx.Write([]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09}); err != nil {
		log.Errorf("Error on response writer")
		ret.closed = true
	}
	pio.PutI32BE(ret.buf[:4], 0)
	if _, err := ret.ctx.Write(ret.buf[:4]); err != nil {
		log.Errorf("Error on response writer")
		ret.closed = true
	}
	go func() {
		err := ret.SendPacket()
		if err != nil {
			log.Debugf("SendPacket error: %v", err)
			ret.closed = true
		}

	}()
	return ret
}

type FLVWriter struct {
	Uid string
	av.RWBaser
	app, title, url string
	buf             []byte
	closed          bool
	closedChan      chan struct{}
	ctx             http.ResponseWriter
	packetQueue     chan *av.Packet
}

func (flvWriter *FLVWriter) DropPacket(pktQue chan *av.Packet, info av.Info) {
	log.Warnf("[%v] packet queue max!!!", info)
	for i := 0; i < maxQueueNum-84; i++ {
		tmpPkt, ok := <-pktQue
		if ok && tmpPkt.IsVideo {
			videoPkt, ok := tmpPkt.Header.(av.VideoPacketHeader)
			// dont't drop sps config and dont't drop key frame
			if ok && (videoPkt.IsSeq() || videoPkt.IsKeyFrame()) {
				log.Debug("insert keyframe to queue")
				pktQue <- tmpPkt
			}

			if len(pktQue) > maxQueueNum-10 {
				<-pktQue
			}
			// drop other packet
			<-pktQue
		}
		// try to don't drop audio
		if ok && tmpPkt.IsAudio {
			log.Debug("insert audio to queue")
			pktQue <- tmpPkt
		}
	}
	log.Debugf("packet queue len: %d", len(pktQue))
}

func (flvWriter *FLVWriter) Write(p *av.Packet) (err error) {
	err = nil
	if flvWriter.closed {
		err = fmt.Errorf("flvwrite source closed")
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("FLVWriter has already been closed:%v", e)
		}
	}()

	if len(flvWriter.packetQueue) >= maxQueueNum-24 {
		flvWriter.DropPacket(flvWriter.packetQueue, flvWriter.Info())
	} else {
		flvWriter.packetQueue <- p
	}

	return
}

func (flvWriter *FLVWriter) SendPacket() error {
	for {
		p, ok := <-flvWriter.packetQueue
		if ok {
			flvWriter.RWBaser.SetPreTime()
			h := flvWriter.buf[:headerLen]
			typeID := av.TAG_VIDEO
			if !p.IsVideo {
				if p.IsMetadata {
					var err error
					typeID = av.TAG_SCRIPTDATAAMF0
					p.Data, err = amf.MetaDataReform(p.Data, amf.DEL)
					if err != nil {
						return err
					}
				} else {
					typeID = av.TAG_AUDIO
				}
			}
			dataLen := len(p.Data)
			timestamp := p.TimeStamp
			timestamp += flvWriter.BaseTimeStamp()
			flvWriter.RWBaser.RecTimeStamp(timestamp, uint32(typeID))

			preDataLen := dataLen + headerLen
			timestampbase := timestamp & 0xffffff
			timestampExt := timestamp >> 24 & 0xff

			pio.PutU8(h[0:1], uint8(typeID))
			pio.PutI24BE(h[1:4], int32(dataLen))
			pio.PutI24BE(h[4:7], int32(timestampbase))
			pio.PutU8(h[7:8], uint8(timestampExt))

			if _, err := flvWriter.ctx.Write(h); err != nil {
				return err
			}

			if _, err := flvWriter.ctx.Write(p.Data); err != nil {
				return err
			}

			pio.PutI32BE(h[:4], int32(preDataLen))
			if _, err := flvWriter.ctx.Write(h[:4]); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("closed")
		}

	}
}

func (flvWriter *FLVWriter) Wait() {
	select {
	case <-flvWriter.closedChan:
		return
	}
}

func (flvWriter *FLVWriter) Close(error) {
	log.Debug("http flv closed")
	if !flvWriter.closed {
		close(flvWriter.packetQueue)
		close(flvWriter.closedChan)
	}
	flvWriter.closed = true
}

func (flvWriter *FLVWriter) Info() (ret av.Info) {
	ret.UID = flvWriter.Uid
	ret.URL = flvWriter.url
	ret.Key = flvWriter.app + "/" + flvWriter.title
	ret.Inter = true
	return
}
