package live

import (
	"bytes"

	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/lib/modules/live/container/flv"
	"github.com/liwei1dao/lego/lib/modules/live/container/ts"
)

type Source struct {
	av.RWBaser
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
