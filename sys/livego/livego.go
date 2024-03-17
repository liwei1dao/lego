package livego

import (
	"fmt"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/serves/api"
	"github.com/liwei1dao/lego/sys/livego/serves/hls"
	"github.com/liwei1dao/lego/sys/livego/serves/httpflv"
	"github.com/liwei1dao/lego/sys/livego/serves/rtmp"
	"github.com/liwei1dao/lego/utils/container/id"
)

func newSys(options *Options) (sys *LiveGo, err error) {
	sys = &LiveGo{
		options:       options,
		keyLock:       new(sync.RWMutex),
		keys:          make(map[string]string),
		channelLock:   new(sync.RWMutex),
		channels:      make(map[string]string),
		mapLock:       new(sync.RWMutex),
		staticPushMap: make(map[string]*core.StaticPush),
	}
	err = sys.init()
	return
}

type LiveGo struct {
	options *Options
	///流管理
	streams *sync.Map
	///流频道管理
	keyLock     *sync.RWMutex
	keys        map[string]string
	channelLock *sync.RWMutex
	channels    map[string]string
	///静态推送管理
	mapLock       *sync.RWMutex
	staticPushMap map[string](*core.StaticPush)
	rtmpServer    core.IRtmpServer
	apiServer     core.IApiServer
	hlsServer     core.IHlsServer
	httpflvServer core.IHttpFlvServer
}

func (this *LiveGo) init() (err error) {
	if this.options.Hls {
		if this.hlsServer, err = hls.NewServer(this, this.options.Log); err != nil {
			return
		}
	}
	if this.options.Flv {
		if this.httpflvServer, err = httpflv.NewServer(this, this.options.Log); err != nil {
			return
		}
	}
	if this.options.Api {
		if this.apiServer, err = api.NewServer(this, this.options.Log); err != nil {
			return
		}
	}
	if this.rtmpServer, err = rtmp.NewServer(this, this.options.Log, this.hlsServer); err != nil {
		return
	}
	return
}

func (this *LiveGo) GetRtmpServer() core.IRtmpServer {
	return this.rtmpServer
}

///流管理***********************************************************************
func (this *LiveGo) GetStreams() *sync.Map {
	return this.streams
}

//监测存活
func (this *LiveGo) CheckAlive() {
	for {
		<-time.After(5 * time.Second)
		this.streams.Range(func(key, val interface{}) bool {
			v := val.(*core.Stream)
			if v.CheckAlive() == 0 {
				this.streams.Delete(key)
			}
			return true
		})
	}
}

//读处理
func (this *LiveGo) HandleReader(r core.ReadCloser) {
	info := r.Info()
	this.options.Log.Debugf("HandleReader: info[%v]", info)
	var stm *core.Stream
	i, ok := this.streams.Load(info.Key)
	if stm, ok = i.(*core.Stream); ok {
		stm.TransStop()
		id := stm.ID()
		if id != core.EmptyID && id != info.UID {
			ns := core.NewStream(this, this.options.Log)
			stm.Copy(ns)
			stm = ns
			this.streams.Store(info.Key, ns)
		}
	} else {
		stm = core.NewStream(this, this.options.Log)
		this.streams.Store(info.Key, stm)
		stm.SetInfo(info)
	}
	stm.AddReader(r)
}

//写处理
func (this *LiveGo) HandleWriter(w core.WriteCloser) {
	info := w.Info()
	this.options.Log.Debugf("HandleWriter: info[%v]", info)

	var s *core.Stream
	item, ok := this.streams.Load(info.Key)
	if !ok {
		this.options.Log.Debugf("HandleWriter: not found create new info[%v]", info)
		s = core.NewStream(this, this.options.Log)
		this.streams.Store(info.Key, s)
		s.SetInfo(info)
	} else {
		s = item.(*core.Stream)
		s.AddWriter(w)
	}
}

///参数列表***********************************************************************
func (this *LiveGo) GetAppname() string {
	return this.options.Appname
}
func (this *LiveGo) GetHls() bool {
	return this.options.Hls
}
func (this *LiveGo) GetHLSAddr() string {
	return this.options.HLSAddr
}
func (this *LiveGo) GetUseHlsHttps() bool {
	return this.options.UseHlsHttps
}

func (this *LiveGo) GetHLSKeepAfterEnd() bool {
	return this.options.UseHlsHttps
}
func (this *LiveGo) GetHlsServerCrt() string {
	return this.options.HLSKeepAfterEnd
}
func (this *LiveGo) GetHlsServerKey() string {
	return this.options.HlsServerKey
}
func (this *LiveGo) GetFlv() bool {
	return this.options.Flv
}

func (this *LiveGo) GetHTTPFLVAddr() string {
	return this.options.HTTPFLVAddr
}

func (this *LiveGo) GetApi() bool {
	return this.options.Api
}
func (this *LiveGo) GetApiAddr() string {
	return this.options.APIAddr
}
func (this *LiveGo) GetJWTSecret() string {
	return this.options.JWTSecret
}
func (this *LiveGo) GetJWTAlgorithm() string {
	return this.options.JWTAlgorithm
}
func (this *LiveGo) GetStaticPush() []string {
	return this.options.StaticPush
}
func (this *LiveGo) GetRTMPAddr() string {
	return this.options.RTMPAddr
}

func (this *LiveGo) GetRTMPNoAuth() bool {
	return this.options.RTMPNoAuth
}
func (this *LiveGo) GetFLVDir() string {
	return this.options.FLVDir
}
func (this *LiveGo) GetFLVArchive() bool {
	return this.options.FLVArchive
}
func (this *LiveGo) GetEnableTLSVerify() bool {
	return this.options.EnableTLSVerify
}
func (this *LiveGo) GetTimeout() int {
	return this.options.Timeout
}
func (this *LiveGo) GetConnBuffSzie() int {
	return this.options.ConnBuffSzie
}
func (this *LiveGo) GetDebug() bool {
	return this.options.Debug
}

///房间***********************************************************************
func (this *LiveGo) SetKey(channel string) (key string, err error) {
	key = id.NewXId()
	this.keyLock.Lock()
	this.keys[key] = channel
	this.keyLock.Unlock()
	this.channelLock.Lock()
	this.channels[channel] = key
	this.channelLock.Unlock()
	return
}
func (this *LiveGo) GetKey(channel string) (key string, err error) {
	var ok bool
	this.channelLock.RLock()
	key, ok = this.channels[channel]
	this.channelLock.RUnlock()
	if !ok {
		err = fmt.Errorf("no channel:%s", channel)
	}
	return
}
func (this *LiveGo) GetChannel(key string) (channel string, err error) {
	var ok bool
	this.channelLock.RLock()
	channel, ok = this.keys[key]
	this.channelLock.RUnlock()
	if !ok {
		err = fmt.Errorf("no key:%s", channel)
	}
	return
}
func (this *LiveGo) DeleteChannel(channel string) (ok bool) {
	this.channelLock.Lock()
	delete(this.channels, channel)
	this.channelLock.Unlock()
	ok = true
	return
}
func (this *LiveGo) DeleteKey(key string) (ok bool) {
	this.keyLock.Lock()
	delete(this.keys, key)
	this.keyLock.Unlock()
	ok = true
	return
}

///静态推送管理***********************************************************************
func (this *LiveGo) GetAndCreateStaticPushObject(rtmpurl string) (pushobj *core.StaticPush) {
	this.mapLock.RLock()
	staticpush, ok := this.staticPushMap[rtmpurl]
	this.options.Log.Debugf("GetAndCreateStaticPushObject: %s, return %v", rtmpurl, ok)
	if !ok {
		this.mapLock.RUnlock()
		newStaticpush := core.NewStaticPush(this, this.options.Log, rtmpurl)
		this.mapLock.Lock()
		this.staticPushMap[rtmpurl] = newStaticpush
		this.mapLock.Unlock()
		return newStaticpush
	}
	this.mapLock.RUnlock()
	return staticpush
}

func (this *LiveGo) GetStaticPushObject(rtmpurl string) (pushobj *core.StaticPush, err error) {
	this.mapLock.RLock()
	if staticpush, ok := this.staticPushMap[rtmpurl]; ok {
		this.mapLock.RUnlock()
		return staticpush, nil
	}
	this.mapLock.RUnlock()

	return nil, fmt.Errorf("G_StaticPushMap[%s] not exist....", rtmpurl)
}

func (this *LiveGo) ReleaseStaticPushObject(rtmpurl string) {
	this.mapLock.RLock()
	if _, ok := this.staticPushMap[rtmpurl]; ok {
		this.mapLock.RUnlock()
		this.options.Log.Debugf("ReleaseStaticPushObject %s ok", rtmpurl)
		this.mapLock.Lock()
		delete(this.staticPushMap, rtmpurl)
		this.mapLock.Unlock()
	} else {
		this.mapLock.RUnlock()
		this.options.Log.Debugf("ReleaseStaticPushObject: not find %s", rtmpurl)
	}
	return
}
