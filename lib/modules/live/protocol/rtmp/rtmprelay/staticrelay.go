package rtmprelay

import (
	"fmt"
	"sync"

	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/lib/modules/live/protocol/rtmp/core"
	"github.com/liwei1dao/lego/sys/log"
)

var G_PushUrlList []string = nil
var g_MapLock = new(sync.RWMutex)
var G_StaticPushMap = make(map[string](*StaticPush))

func GetStaticPushList(appname string) ([]string, error) {
	if G_PushUrlList == nil {
		// Do not unmarshel the config every time, lots of reflect works -gs
		// pushurlList, ok := configure.GetStaticPushUrlList(appname)
		pushurlList, ok := []string{}, false
		if !ok {
			G_PushUrlList = []string{}
		} else {
			G_PushUrlList = pushurlList
		}
	}

	if len(G_PushUrlList) == 0 {
		return nil, fmt.Errorf("no static push url")
	}

	return G_PushUrlList, nil
}

func GetAndCreateStaticPushObject(rtmpurl string) *StaticPush {
	g_MapLock.RLock()
	staticpush, ok := G_StaticPushMap[rtmpurl]
	log.Debugf("GetAndCreateStaticPushObject: %s, return %v", rtmpurl, ok)
	if !ok {
		g_MapLock.RUnlock()
		newStaticpush := NewStaticPush(rtmpurl)

		g_MapLock.Lock()
		G_StaticPushMap[rtmpurl] = newStaticpush
		g_MapLock.Unlock()

		return newStaticpush
	}
	g_MapLock.RUnlock()

	return staticpush
}

func GetStaticPushObject(rtmpurl string) (*StaticPush, error) {
	g_MapLock.RLock()
	if staticpush, ok := G_StaticPushMap[rtmpurl]; ok {
		g_MapLock.RUnlock()
		return staticpush, nil
	}
	g_MapLock.RUnlock()

	return nil, fmt.Errorf("G_StaticPushMap[%s] not exist....", rtmpurl)
}

func NewStaticPush(rtmpurl string) *StaticPush {
	return &StaticPush{
		RtmpUrl:       rtmpurl,
		packet_chan:   make(chan *av.Packet, 500),
		sndctrl_chan:  make(chan string),
		connectClient: nil,
		startflag:     false,
	}
}

type StaticPush struct {
	RtmpUrl       string
	packet_chan   chan *av.Packet
	sndctrl_chan  chan string
	connectClient *core.ConnClient
	startflag     bool
}

func (self *StaticPush) WriteAvPacket(packet *av.Packet) {
	if !self.startflag {
		return
	}

	self.packet_chan <- packet
}
