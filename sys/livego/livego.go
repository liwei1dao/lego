package livego

import (
	"fmt"
	"sync"

	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/serves/api"
	"github.com/liwei1dao/lego/sys/livego/serves/hls"
	"github.com/liwei1dao/lego/sys/livego/serves/rtmp"
	"github.com/liwei1dao/lego/utils/container/id"
)

func newSys(options Options) (sys *LiveGo, err error) {
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
	options Options
	///流频道管理
	keyLock     *sync.RWMutex
	keys        map[string]string
	channelLock *sync.RWMutex
	channels    map[string]string

	///静态推送管理
	mapLock       *sync.RWMutex
	staticPushMap map[string](*core.StaticPush)

	rtmpServer core.IRtmpServer
	apiServer  core.IApiServer
	hlsServer  core.IHlsServer
}

func (this *LiveGo) init() (err error) {
	if this.options.Api {
		if this.apiServer, err = api.NewServer(this); err != nil {
			return
		}
	}
	if this.options.Hls {
		if this.hlsServer, err = hls.NewServer(this); err != nil {
			return
		}
	}
	if this.rtmpServer, err = rtmp.NewServer(this, this.hlsServer); err != nil {
		return
	}
	return
}

func (this *LiveGo) GetRtmpServer() core.IRtmpServer {
	return this.rtmpServer
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
func (this *LiveGo) GetHlsServerCrt() string {
	return this.options.HlsServerCrt
}
func (this *LiveGo) GetHlsServerKey() string {
	return this.options.HlsServerKey
}
func (this *LiveGo) GetFlv() bool {
	return this.options.Flv
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

///日志***********************************************************************
func (this *LiveGo) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Debugf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Infof("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Warnf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Errorf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Panicf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Fatalf("[SYS LiveGo] "+format, a)
	}
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
	this.Debugf("GetAndCreateStaticPushObject: %s, return %v", rtmpurl, ok)
	if !ok {
		this.mapLock.RUnlock()
		newStaticpush := core.NewStaticPush(rtmpurl)
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
		this.Debugf("ReleaseStaticPushObject %s ok", rtmpurl)
		this.mapLock.Lock()
		delete(this.staticPushMap, rtmpurl)
		this.mapLock.Unlock()
	} else {
		this.mapLock.RUnlock()
		this.Debugf("ReleaseStaticPushObject: not find %s", rtmpurl)
	}
	return
}
