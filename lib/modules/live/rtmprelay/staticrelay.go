package rtmprelay

import (
	"fmt"
	"sync"

	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/lib/modules/live/liveconn"
	"github.com/liwei1dao/lego/sys/log"
)

var G_StaticPushMap = make(map[string](*StaticPush))
var g_MapLock = new(sync.RWMutex)
var G_PushUrlList []string = nil

var (
	STATIC_RELAY_STOP_CTRL = "STATIC_RTMPRELAY_STOP"
)

func GetStaticPushList(pushurlList []string, ok bool) ([]string, error) {
	if G_PushUrlList == nil {
		// Do not unmarshel the config every time, lots of reflect works -gs
		// pushurlList, ok := configure.GetStaticPushUrlList(appname)
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

func GetStaticPushObject(rtmpurl string) (*StaticPush, error) {
	g_MapLock.RLock()
	if staticpush, ok := G_StaticPushMap[rtmpurl]; ok {
		g_MapLock.RUnlock()
		return staticpush, nil
	}
	g_MapLock.RUnlock()

	return nil, fmt.Errorf("G_StaticPushMap[%s] not exist....", rtmpurl)
}

func ReleaseStaticPushObject(rtmpurl string) {
	g_MapLock.RLock()
	if _, ok := G_StaticPushMap[rtmpurl]; ok {
		g_MapLock.RUnlock()

		log.Debugf("ReleaseStaticPushObject %s ok", rtmpurl)
		g_MapLock.Lock()
		delete(G_StaticPushMap, rtmpurl)
		g_MapLock.Unlock()
	} else {
		g_MapLock.RUnlock()
		log.Debugf("ReleaseStaticPushObject: not find %s", rtmpurl)
	}
}

type StaticPush struct {
	RtmpUrl       string
	packet_chan   chan *av.Packet
	sndctrl_chan  chan string
	connectClient *liveconn.ConnClient
	startflag     bool
}

func (self *StaticPush) Stop() {
	if !self.startflag {
		return
	}

	log.Debugf("StaticPush Stop: %s", self.RtmpUrl)
	self.sndctrl_chan <- STATIC_RELAY_STOP_CTRL
	self.startflag = false
}

func (self *StaticPush) WriteAvPacket(packet *av.Packet) {
	if !self.startflag {
		return
	}

	self.packet_chan <- packet
}
