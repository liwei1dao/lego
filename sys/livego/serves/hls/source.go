package hls

import (
	"bytes"

	"github.com/liwei1dao/lego/sys/livego/container/flv"
	"github.com/liwei1dao/lego/sys/livego/container/ts"
	"github.com/liwei1dao/lego/sys/livego/packet"
	"github.com/liwei1dao/lego/sys/livego/parser"
)

const (
	syncms       = 2 // ms
	videoHZ      = 90000
	aacSampleLen = 1024
	maxQueueNum  = 512

	h264_default_hz uint64 = 90
)

type Source struct {
	packet.RWBaser
	seq         int
	info        packet.Info
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
	packetQueue chan *packet.Packet
}

func (this *Source) Info() (ret packet.Info) {
	return this.info
}

func (this *Source) GetCacheInc() *TSCacheItem {
	return this.tsCache
}
