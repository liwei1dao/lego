package httpflv

import (
	"fmt"
	"net/http"
	"time"

	"github.com/liwei1dao/lego/sys/livego/codec"
	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/utils/pio"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/container/id"
)

const (
	headerLen   = 11
	maxQueueNum = 1024
)

func NewFLVWriter(sys core.ISys, log log.ILogger, app, title, url string, ctx http.ResponseWriter) *FLVWriter {
	ret := &FLVWriter{
		sys:         sys,
		log:         log,
		Uid:         id.NewXId(),
		app:         app,
		title:       title,
		url:         url,
		ctx:         ctx,
		RWBaser:     core.NewRWBaser(time.Second * 10),
		closedChan:  make(chan struct{}),
		buf:         make([]byte, headerLen),
		packetQueue: make(chan *core.Packet, maxQueueNum),
	}

	if _, err := ret.ctx.Write([]byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09}); err != nil {
		ret.log.Errorf("Error on response writer")
		ret.closed = true
	}
	pio.PutI32BE(ret.buf[:4], 0)
	if _, err := ret.ctx.Write(ret.buf[:4]); err != nil {
		ret.log.Errorf("Error on response writer")
		ret.closed = true
	}
	go func() {
		err := ret.SendPacket()
		if err != nil {
			ret.log.Debugf("SendPacket error:%v", err)
			ret.closed = true
		}

	}()
	return ret
}

type FLVWriter struct {
	sys core.ISys
	log log.ILogger
	Uid string
	core.RWBaser
	app, title, url string
	buf             []byte
	closed          bool
	closedChan      chan struct{}
	ctx             http.ResponseWriter
	packetQueue     chan *core.Packet
}

func (this *FLVWriter) DropPacket(pktQue chan *core.Packet, info core.Info) {
	this.log.Warnf("[%v] packet queue max!!!", info)
	for i := 0; i < maxQueueNum-84; i++ {
		tmpPkt, ok := <-pktQue
		if ok && tmpPkt.IsVideo {
			videoPkt, ok := tmpPkt.Header.(core.VideoPacketHeader)
			// dont't drop sps config and dont't drop key frame
			if ok && (videoPkt.IsSeq() || videoPkt.IsKeyFrame()) {
				this.log.Debugf("insert keyframe to queue")
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
			this.log.Debugf("insert audio to queue")
			pktQue <- tmpPkt
		}
	}
	this.log.Debugf("packet queue len: ", len(pktQue))
}
func (this *FLVWriter) SendPacket() error {
	for {
		p, ok := <-this.packetQueue
		if ok {
			this.RWBaser.SetPreTime()
			h := this.buf[:headerLen]
			typeID := core.TAG_VIDEO
			if !p.IsVideo {
				if p.IsMetadata {
					var err error
					typeID = core.TAG_SCRIPTDATAAMF0
					p.Data, err = codec.MetaDataReform(p.Data, codec.DEL)
					if err != nil {
						return err
					}
				} else {
					typeID = core.TAG_AUDIO
				}
			}
			dataLen := len(p.Data)
			timestamp := p.TimeStamp
			timestamp += this.BaseTimeStamp()
			this.RWBaser.RecTimeStamp(timestamp, uint32(typeID))

			preDataLen := dataLen + headerLen
			timestampbase := timestamp & 0xffffff
			timestampExt := timestamp >> 24 & 0xff

			pio.PutU8(h[0:1], uint8(typeID))
			pio.PutI24BE(h[1:4], int32(dataLen))
			pio.PutI24BE(h[4:7], int32(timestampbase))
			pio.PutU8(h[7:8], uint8(timestampExt))

			if _, err := this.ctx.Write(h); err != nil {
				return err
			}

			if _, err := this.ctx.Write(p.Data); err != nil {
				return err
			}

			pio.PutI32BE(h[:4], int32(preDataLen))
			if _, err := this.ctx.Write(h[:4]); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("closed")
		}
	}
}

func (this *FLVWriter) Write(p *core.Packet) (err error) {
	err = nil
	if this.closed {
		err = fmt.Errorf("flvwrite source closed")
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("FLVWriter has already been closed:%v", e)
		}
	}()

	if len(this.packetQueue) >= maxQueueNum-24 {
		this.DropPacket(this.packetQueue, this.Info())
	} else {
		this.packetQueue <- p
	}

	return
}

func (this *FLVWriter) Wait() {
	select {
	case <-this.closedChan:
		return
	}
}

func (this *FLVWriter) Close(error) {
	this.log.Debugf("http flv closed")
	if !this.closed {
		close(this.packetQueue)
		close(this.closedChan)
	}
	this.closed = true
}

func (flvWriter *FLVWriter) Info() (ret core.Info) {
	ret.UID = flvWriter.Uid
	ret.URL = flvWriter.url
	ret.Key = flvWriter.app + "/" + flvWriter.title
	ret.Inter = true
	return
}
